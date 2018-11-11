package main

import (
	"context"
	"log"
)
import (
	"pb"
	"services"

	"google.golang.org/grpc/metadata"
)

type User struct {
	userId      string
	server      ServerInterface
	stream      pb.Stream_StreamServer
	streamCodec string
}

type UserInterface interface {
}

func NewUser(userId string, server ServerInterface) *User {
	user := &User{
		userId: userId,
		server: server,
	}
	return user
}

func (u *User) SetStream(stream pb.Stream_StreamServer) {
	u.stream = stream
}

func (u *User) HandleStreamMessage(codec string, name string, data []byte) {
	u.streamCodec = codec
	message, err := DecodeMessage(codec, name, data)
	if err != nil {
		log.Println(err)
		return
	}
	// handle or forward ?

}

func (u *User) SendMessage(name string, message interface{}) {
	if u.stream == nil {
		return
	}
	msgId := services.ToMsgId(name)
	data, _ := EncodeMessage(u.streamCodec, name, message)
	u.stream.Send(&pb.StreamFrame{
		Type:    pb.StreamFrameType_Message,
		Codec:   u.streamCodec,
		MsgId:   int32(msgId),
		MsgData: data,
	})
}
