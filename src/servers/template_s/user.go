package main

import (
	"log"
)

import (
	"actors"
	"pb/template"
)

type User struct {
	actors.BaseActorHost
	userId      string
	out         chan interface{}
	streamCodec string
}

func NewUser(userId string, out chan interface{}) *User {
	user := &User{
		userId: userId,
		out:    out,
	}
	actors.NewActor(user)
	return user
}

func (u *User) SetOutChan(out chan interface{}) {
	u.out = out
}

func (u *User) HandleMessages(name string, msg interface{}) {
	switch name {
	case "TemplateMsgTest":
		msgTemplateMsgTest := template.ToTemplateMsgTest(msg)
		if msgTemplateMsgTest == nil {
			log.Fatalln("type error")
		}
	}
}

func (u *User) SendMessage(message interface{}) {
	if u.out == nil {
		return
	}
	u.out <- message
	/*
		msgId := services.ToMsgId(name)
		data, _ := EncodeMessage(u.streamCodec, name, message)
		u.stream.Send(&pb.StreamFrame{
			Type:    pb.StreamFrameType_Message,
			Codec:   u.streamCodec,
			MsgId:   int32(msgId),
			MsgData: data,
		})
	*/
}
