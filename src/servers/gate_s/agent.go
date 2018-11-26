package main

import (
	"context"
	"log"
	"sync"
	"utils"
)

import (
	"services"
	"types"
)

type clientSendFunc func([]byte) error

// agent is client's agent in gate, its role is fowarding client's message to service,
// and forwarding service message to client. It is created when socket connection or
// http request is coming, removed when connection close or http request is over.
// each agent has one session
type agent struct {
	sessId     string
	clientSend clientSendFunc
	inBuf      chan []byte
	outBuf     chan []byte
	status     int // 0 - not work, 1 - working
	mu         sync.RWMutex
	inSeq      int32
	outSeq     int32
}

var _agent_pool = sync.Pool{
	New: func() interface{} {
		return new(agent)
	},
}

func newAgent(sessId string, clientSend clientSendFunc) *agent {
	a := _agent_pool.Get().(*agent)
	a.sessId = sessId
	a.clientSend = clientSend
	a.inBuf = make(chan []byte, 100)
	a.outBuf = make(chan []byte, 100)
	a.status = 0
	a.inSeq = -1
	a.outSeq = 0
	return a
}

func freeAgent(a *agent) {
	_agent_pool.Put(a)
}

func (a *agent) start(ctx context.Context) {
	a.setStatus(1) // now agent is working
	// main loop
	go func() {
		defer func() {
			a.exit()
			log.Println("main loop end")
		}()
		for {
			select {
			case msg, ok := <-a.inBuf:
				if !ok {
					log.Fatalln("in Buf not ok")
				}
				log.Printf("recv: %s", msg)
				seq, msgId, msgData := utils.UnpackMsg(msg)
				if seq <= a.inSeq || msgId == 0 || msgData == nil {
					log.Println("invalid message format")
					return
				}
				a.inSeq = seq
				t := services.ToServiceType(msgId)
				if t == "" {
					log.Println("invalid msgId")
					return
				}
				if t == "gate" {
					a.handle(msgId, msgData)
				} else {
					err := a.forward(t, msgId, msgData)
					switch err {
					case types.ERR_NO_SESSION, types.ERR_NO_SERVICE:
					}
				}

			case msg, ok := <-a.outBuf:
				if !ok {
					return
				}
				_ = a.clientSend(msg)

			case <-ctx.Done():
				log.Print("exit ntf")
				return
			}
		}
	}()
}

func (a *agent) getStatus() int {
	status := 0
	a.mu.RLock()
	status = a.status
	a.mu.RUnlock()
	return status
}

func (a *agent) setStatus(status int) {
	a.mu.Lock()
	a.status = status
	a.mu.Unlock()
}

func (a *agent) exit() {
	a.setStatus(0)
	sess, ok := getSessionManager().getSession(a.sessId)
	if ok {
		sess.unbindAgent(a)
	}
	close(a.inBuf)
	close(a.outBuf)
}

func (a *agent) handle(msgId uint32, msgData []byte) error {
	if msgId == utils.HashCode("login") {
		userId := "100000"
		sess, err := getSessionManager().createSession(userId)
		if err != nil {
			return err
		}
		// bindng each other
		sess.bindAgent(a)
	}
	return nil
}

func (a *agent) forward(serviceType string, msgId uint32, msgData []byte) error {
	sess, ok := getSessionManager().getSession(a.sessId)
	if !ok {
		return types.ERR_NO_SESSION
	}
	return sess.streamSend(serviceType, msgId, msgData)
}

func (a *agent) toService(msg []byte) {
	if a.getStatus() == 0 {
		return
	}
	select {
	case a.inBuf <- msg:
	}
}

func (a *agent) toClient(msgId uint32, msgData []byte) {
	if a.getStatus() == 0 {
		return
	}
	a.outSeq = a.outSeq + 1
	msg := utils.PackMsg(a.outSeq, msgId, msgData)
	select {
	case a.outBuf <- msg:
	}
}
