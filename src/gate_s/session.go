package main

import (
	"context"
	"log"
	"sync"
	"time"
)

import (
	"pb"
	"services"
	"types"
	"utils"

	"google.golang.org/grpc/metadata"
)

// session is user state in gate, it can store user's service using state.
// it is created when user login, and removed when user logout or timeout. each
// session has only one agent. agent find service by session, service repley find
// agent by session too.
type session struct {
	userId     string
	sessId     string
	agent      *agent
	serviceMap map[string]string                 //serviceType to serviceId
	streams    map[string]pb.Stream_StreamClient //binding serviceType and streamClient
	activeTime time.Time
	cancel     context.CancelFunc
}

func (s *session) init(userId string) {
	s.userId = userId
	s.sessId = utils.CreateUUID()
	s.agent = nil
	s.serviceMap = make(map[string]string, 0)
	s.streams = make(map[string]pb.Stream_StreamClient, 0)
	s.activeTime = time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	serviceUseList := services.GetServiceUseList()
	for _, v := range serviceUseList {
		serviceType := v
		serviceId, err := services.AssignServiceId(serviceType)
		if err != nil {
			return
		}
		s.serviceMap[serviceType] = serviceId
		serviceClient, err := services.GetServiceClient(serviceId)
		if err != nil {
			log.Println(err.Error())
			return
		}
		if serviceClient.Conf.IsStream {
			cli := pb.NewStreamClient(serviceClient.Conn)
			streamCtx := metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"user-id": s.userId}))
			streamClient, err := cli.Stream(streamCtx)
			if err != nil {
				log.Println(err.Error())
				return
			}
			s.streams[serviceType] = streamClient
			go func() {
				defer func() {

				}()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						frame, err := streamClient.Recv()
						if err != nil {
							log.Println(err)
							return
						}
						if s.agent != nil {
							s.agent.toClient(uint32(frame.MsgId), frame.MsgData)
						}
					}
				}
			}()
		}
	}
}

func (s *session) distroy() {
	s.cancel()
	for _, v := range s.serviceMap {
		services.UnassignServiceId(v)
	}
}

func (s *session) bindAgent(a *agent) {
	a.sessId = s.sessId
	s.agent = a
}

func (s *session) unbindAgent(a *agent) {
	a.sessId = ""
	s.agent = nil
}

func (s *session) streamSend(serviceType string, msgId uint32, msgData []byte) error {
	s.activeTime = time.Now()
	streamClient, ok := s.streams[serviceType]
	if !ok {
		return types.ERR_NO_SERVICE
	}
	frame := &pb.StreamFrame{
		Type:    pb.StreamFrameType_Message,
		MsgId:   int32(msgId),
		MsgData: msgData,
	}
	err := streamClient.Send(frame)
	return err
}

type sessionManager struct {
	userId2sessId map[string]string
	sessions      map[string]*session
	timeout       int
	pool          sync.Pool
}

func (p *sessionManager) init() {
	p.userId2sessId = make(map[string]string, 0)
	p.sessions = make(map[string]*session, 0)
	p.timeout = 5 * 60
	p.pool = sync.Pool{
		New: func() interface{} {
			return new(session)
		},
	}
}

func (p *sessionManager) createSession(userId string) (*session, error) {
	var sess *session
	sessId, ok := p.userId2sessId[userId]
	if !ok {
		sess = p.pool.Get().(*session)
		sess.init(userId)
		p.sessions[sess.sessId] = sess
		p.userId2sessId[userId] = sess.sessId
	} else {
		sess, _ = p.sessions[sessId]
		sess.activeTime = time.Now()
	}
	return sess, nil
}

func (p *sessionManager) getSession(sessId string) (*session, bool) {
	sess, ok := p.sessions[sessId]
	return sess, ok
}

func (p *sessionManager) releaseSession(sess *session) {
	for k, v := range p.userId2sessId {
		if v == sess.sessId {
			delete(p.userId2sessId, k)
			break
		}
	}
	delete(p.sessions, sess.sessId)
	sess.distroy()
	p.pool.Put(sess)
}

// remove seesion that active time is expired
func (p *sessionManager) watchSessions() {
	go func() {
		for {
			now := time.Now()
			isDel := false
			for _, v := range p.sessions {
				if now.Second()-v.activeTime.Second() > p.timeout {
					p.releaseSession(v)
					isDel = true
					break
				}
			}
			if isDel == false {
				time.Sleep(1)
			}
		}
	}()
}

var (
	_session_manager sessionManager
	_once            sync.Once
)

func getSessionManager() *sessionManager {
	_once.Do(func() { _session_manager.init() })
	return &_session_manager
}
