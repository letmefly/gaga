package main

import (
	"log"
)

import (
	"actors"
	"pb/auth"
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
	case "AuthMsgTest":
		msgAuthMsgTest := auth.ToAuthMsgTest(msg)
		if msgAuthMsgTest == nil {
			log.Fatalln("type error")
		}
	}
}

func (u *User) SendMessage(message interface{}) {
	if u.out == nil {
		return
	}
	u.out <- message
}
