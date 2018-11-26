package main

import (
	. "pb/template"
	t "types"
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
func (l *logic) handleMessage(network int, msgName string, msg interface{}) (interface{}, error) {
	switch network {
	case t.HTTP:
		return l.handleHttpMessage(network, msgName, msg)
	case t.WS:
		return l.handleWsMessage(network, msgName, msg)
	case t.GRPC:
		return l.handleGrpcMessage(network, msgName, msg)
	case t.TCP:
	case t.UDP:
	case t.KCP:
	}
	return nil, nil
}

func (l *logic) handleHttpMessage(network int, msgName string, msg interface{}) (interface{}, error) {
	switch msgName {
	case "TemplateMsgTest":
		pbMsg := msg.(*TemplateMsgTest)
		if pbMsg != nil {
		}
	}
	return nil, nil
}

func (l *logic) handleWsMessage(network int, msgName string, msg interface{}) (interface{}, error) {
	switch msgName {
	case "TemplateMsgTest":
		pbMsg := msg.(*TemplateMsgTest)
		if pbMsg != nil {
		}
	}
	return nil, nil
}

func (l *logic) handleGrpcMessage(network int, msgName string, msg interface{}) (interface{}, error) {
	switch msgName {
	case "TemplateMsgTest":
		pbMsg := msg.(*TemplateMsgTest)
		if pbMsg != nil {
		}
	}
	return nil, nil
}
