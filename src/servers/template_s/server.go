package main

import (
	"context"
	"log"
	"reflect"
	"strings"
	"sync"
)
import (
	"actors"
	"pb"
	"pb/template"
	"services"

	"google.golang.org/grpc/metadata"
)

type Server struct {
	template.GrpcServer
	users sync.Map
	//msgPool sync.Pool
}

func (s *Server) Init(ctx context.Context) {
	/*
		s.msgPool = sync.Pool{
			New: func() interface{} {
				return &ClientMsg{}
			},
		}
	*/
}

func (s *Server) CreateUser(userId string, out chan interface{}) *User {
	user, ok := s.GetUser(userId)
	if !ok {
		user = NewUser(userId, out)
		s.users.Store(userId, user)
	}
	return user
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
	codec := meta["codec"][0]
	out := make(chan interface{}, 1000)
	user := s.CreateUser(userId, out)

	go func() {
		for {
			select {
			case outMsg := <-out:
				msgName := reflect.TypeOf(outMsg).String()
				msgName = strings.Replace(msgName, "*", "", 1)
				msgId := services.ToMsgId(msgName)
				log.Println("send msgName:", msgName, msgId)
				data, _ := template.EncodeMessage(codec, outMsg)
				stream.Send(&pb.StreamFrame{
					Type:    pb.StreamFrameType_Message,
					Codec:   codec,
					MsgId:   int32(msgId),
					MsgData: data,
				})
			case <-stream.Context().Done():
				return
			}
		}
	}()

	defer func() {
		close(out)
	}()

	// receive loop
	for {
		frame, err := stream.Recv()
		if err != nil {
			log.Println(err)
			return err
		}
		codec := frame.Codec
		msgName := services.ToMsgName(uint32(frame.MsgId))
		log.Println("stream recv msgName", msgName)
		msgData := frame.MsgData
		msg, decodeErr := template.DecodeMessage(codec, msgName, msgData)
		if decodeErr != nil {
			log.Println(decodeErr)
		} else {
			actors.AsynCall(user.ActorId(), (*User).HandleMessages, msgName, msg)
		}

		switch frame.Type {
		case pb.StreamFrameType_Message:
		case pb.StreamFrameType_Ping:
		case pb.StreamFrameType_Kick:
		}
	}
	return nil
}
