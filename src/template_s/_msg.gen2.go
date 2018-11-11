// DOT NOT EDIT THIS FILE, AUTO GEN!!
package main

import (
	"log"
)

import (
	. "pb/template"

	"github.com/golang/protobuf/proto"
)

func (s *server) registerClientHandlers() {
	s.registerHandler("TemplateMsgTest", func(network int, data []byte) {
		msg := &TemplateMsgTest{}
		err := proto.Unmarshal(data, msg)
		if err != nil {
			log.Println(err)
		} else {
			s.logic.handleMessage(network, "TemplateMsgTest", msg)
		}
	})

	s.registerHandler("TemplateMsgTestAck", func(network int, data []byte) {
		msg := &TemplateMsgTestAck{}
		err := proto.Unmarshal(data, msg)
		if err != nil {
			log.Println(err)
		} else {
			s.logic.handleMessage(network, "TemplateMsgTestAck", msg)
		}
	})

	s.registerHandler("TemplateMsgTestNtf", func(network int, data []byte) {
		msg := &TemplateMsgTestNtf{}
		err := proto.Unmarshal(data, msg)
		if err != nil {
			log.Println(err)
		} else {
			s.logic.handleMessage(network, "TemplateMsgTestNtf", msg)
		}
	})
}
func (s *server) send_TemplateMsgTest(network int, userId string, msg *TemplateMsgTest) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClient(network, userId, "TemplateMsgTest", buf)
}

func (s *server) ntf_TemplateMsgTest(network int, userIdList []string, msg *TemplateMsgTest) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClients(network, userIdList, "TemplateMsgTest", buf)
}

func (s *server) send_TemplateMsgTestAck(network int, userId string, msg *TemplateMsgTestAck) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClient(network, userId, "TemplateMsgTestAck", buf)
}

func (s *server) ntf_TemplateMsgTestAck(network int, userIdList []string, msg *TemplateMsgTestAck) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClients(network, userIdList, "TemplateMsgTestAck", buf)
}

func (s *server) send_TemplateMsgTestNtf(network int, userId string, msg *TemplateMsgTestNtf) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClient(network, userId, "TemplateMsgTestNtf", buf)
}

func (s *server) ntf_TemplateMsgTestNtf(network int, userIdList []string, msg *TemplateMsgTestNtf) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClients(network, userIdList, "TemplateMsgTestNtf", buf)
}

func (s *server) getProtoUseList() []string {
	protoUseList := []string{
		"TemplateMsgTest",
		"TemplateMsgTestAck",
		"TemplateMsgTestNtf",
	}
	return protoUseList
}
