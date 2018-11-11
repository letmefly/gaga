package main

import (
	"context"
	"log"
	"sync"
)
import (
	"pb"
	"services"

	"google.golang.org/grpc/metadata"
)

type ClientMsg struct {
	//network int
	Name string
	Data []byte
}

type Server struct {
	users   sync.Map
	msgPool sync.Pool
}

type ServerInterface interface {
}

func (s *Server) Init(ctx context.Context) {
	s.msgPool = sync.Pool{
		New: func() interface{} {
			return &ClientMsg{}
		},
	}
}

func (s *Server) CreateUser(userId string) *User {
	user, ok := s.GetUser(userId)
	if !ok {
		user = NewUser(userId, s)
		s.users.Store(userId, user)
	}
	return user.(*User)
}

func (s *Server) GetUser(userId string) (*User, bool) {
	user, ok := s.users.Load(userId)
	if !ok {
		return nil, false
	}
	return user.(*User), true
}

// grpc api
func (s *Server) CreateStream(stream pb.Stream_CreateStreamServer) error {
	log.Println("new stream is coming..")
	meta, _ := metadata.FromIncomingContext(stream.Context())
	userId := meta["user-id"][0]
	user := s.CreateUser(userId)

	defer func() {}()

	// receive loop
	for {
		frame, err := stream.Recv()
		if err != nil {
			log.Println(err)
			return err
		}
		msgName := services.ToMsgName(uint32(frame.MsgId))
		msgData := frame.MsgData
		codec := frame.Codec
		//clientMsg := &ClientMsg{}
		//clientMsg.Name = msgName
		//clientMsg.Data = frame.MsgData
		user.HandleStreamMessage(codec, msgName, msgData)

		switch frame.Type {
		case pb.StreamFrameType_Message:
		case pb.StreamFrameType_Ping:
		case pb.StreamFrameType_Kick:
		}
	}
	return nil
}

func (s *Server) GetProtoUseList() []string {
	protoUseList := []string{
		"TemplateMsgTest",
		"TemplateMsgTestAck",
		"TemplateMsgTestNtf",
	}
	return protoUseList
}
