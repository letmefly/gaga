package main

import (
	"context"
	"log"
	"sync"
)

import (
	"pb"
	"services"
	"utils"
)

type clientMsg struct {
	msgName string
	msgData []byte
}
type msgHandler func([]byte)

type server struct {
	logic          *logic
	clientMsgQueue chan *clientMsg
	streamMap      *utils.SafeMap
	msgPool        sync.Pool
	msghandlers    map[string]msgHandler
}

/****************************** internal server api******************************/
func (s *server) init(ctx context.Context) {
	s.clientMsgQueue = make(chan *clientMsg, 1000)
	s.streamMap = utils.NewSafeMap()
	s.msghandlers = make(map[string]msgHandler, 0)
	s.msgPool = sync.Pool{
		New: func() interface{} {
			return &clientMsg{}
		},
	}
	s.registerClientHandlers()
	s.startLogicLoop(ctx)
	s.initLogic()
}

func (s *server) sendClient(userId string, msgName string, msgData []byte) {
	stream := s.streamMap.Get(userId)
	if stream == nil {
		return
	}
	msgId := services.ToMsgId(msgName)
	stream.(pb.Stream_StreamServer).Send(&pb.StreamFrame{
		Type:    pb.StreamFrameType_Message,
		MsgId:   int32(msgId),
		MsgData: msgData,
	})
}

func (s *server) sendClients(userIdList []string, msgName string, msgData []byte) {
	for _, userId := range userIdList {
		s.sendClient(userId, msgName, msgData)
	}
	/*
		msgId := services.ToMsgId(msgName)
			for _, v := range s.streamMap.Items() {
				v.(pb.Stream_StreamServer).Send(&pb.StreamFrame{
					Type:    pb.StreamFrameType_Message,
					MsgId:   int32(msgId),
					MsgData: msgData,
				})
			}
	*/
}

func (s *server) handleClientMsg(msg *clientMsg) {
	msgName := msg.msgName
	msgData := msg.msgData
	handler, ok := s.msghandlers[msgName]
	if ok {
		handler(msgData)
	} else {
		log.Println("[server]no message handler for", msgName)
	}
}

func (s *server) startLogicLoop(ctx context.Context) {
	defer func() {
	}()
	go func() {
		for {
			select {
			case msg, ok := <-s.clientMsgQueue:
				if !ok {
					return
				}
				s.handleClientMsg(msg)
				s.msgPool.Put(msg)

			case <-ctx.Done():
			}
		}
	}()
}

func (s *server) registerHandler(msgName string, handler msgHandler) {
	s.msghandlers[msgName] = handler
}

// add your logic init code here
func (s *server) initLogic() {
	s.logic = newLogic()
	s.logic.init(s)
}
