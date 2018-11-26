// DOT NOT EDIT THIS FILE, AUTO GEN!!
package main

import (
	"log"
)

import (
	. "pb/auth"
	
	"github.com/golang/protobuf/proto"
)
func (s *server) registerClientHandlers() {
	s.registerHandler("Test", func(network int, data []byte) {
		msg := &Test{}
		err := proto.Unmarshal(data, msg)
		if err != nil {
			log.Println(err)
		} else {
			s.logic.handleMessage(network, "Test", msg)
		}
	})
}
func (s *server) send_Test(network int, userId string, msg *Test) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClient(network, userId, "Test", buf)
}

func (s *server) ntf_Test(network int, userIdList []string, msg *Test) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClients(network, userIdList, "Test", buf)
}

func (s *server) getProtoUseList() []string {
	protoUseList := []string{
		"Test",
	}
	return protoUseList
}
