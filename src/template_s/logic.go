package main

import (
	. "pb/template"
)

type logic struct {
	server *server
}

func newLogic() *logic {
	return &logic{}
}

func (l *logic) init(s *server) {
	l.server = s
}

// handle all logic message including client message or grpc message
func (l *logic) handleMessage(msgName, msg interface{}) (interface{}, error) {
	switch msgName {
	case "TemplateMsgTest":
		pbMsg := msg.(*TemplateMsgTest)
		if pbMsg != nil {
		}
	}
	return nil, nil
}
