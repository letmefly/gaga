// DOT NOT EDIT THIS FILE, AUTO GEN!!
package main

import (
	"log"
)

import (
	"pb"
	. "pb/template"
	"services"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/metadata"
)

func (s *server) Stream(stream pb.Stream_StreamServer) error {
	log.Println("new stream is coming..")
	meta, _ := metadata.FromIncomingContext(stream.Context())
	userId := meta["user-id"][0]
	s.streamMap.Set(userId, stream)

	defer func() {
		s.streamMap.Delete(userId)
	}()

	// receive loop
	for {
		frame, err := stream.Recv()
		if err != nil {
			log.Println(err)
			return nil
		}
		msgName := services.ToMsgName(uint32(frame.MsgId))
		clientMsg := s.msgPool.Get().(*clientMsg)
		clientMsg.msgName = msgName
		clientMsg.msgData = frame.MsgData
		select {
		case s.clientMsgQueue <- clientMsg:
		}
		switch frame.Type {
		case pb.StreamFrameType_Message:
		case pb.StreamFrameType_Ping:
		case pb.StreamFrameType_Kick:
		}
	}

	return nil
}
func (s *server) registerClientHandlers() {
	s.registerHandler("TemplateMsgTest", func(data []byte) {
		msg := &TemplateMsgTest{}
		err := proto.Unmarshal(data, msg)
		if err != nil {
			log.Println(err)
		} else {
			s.logic.handleMessage("TemplateMsgTest", msg)
		}
	})

	s.registerHandler("TemplateMsgTestAck", func(data []byte) {
		msg := &TemplateMsgTestAck{}
		err := proto.Unmarshal(data, msg)
		if err != nil {
			log.Println(err)
		} else {
			s.logic.handleMessage("TemplateMsgTestAck", msg)
		}
	})

	s.registerHandler("TemplateMsgTestNtf", func(data []byte) {
		msg := &TemplateMsgTestNtf{}
		err := proto.Unmarshal(data, msg)
		if err != nil {
			log.Println(err)
		} else {
			s.logic.handleMessage("TemplateMsgTestNtf", msg)
		}
	})
}
func (s *server) send_TemplateMsgTest(userId string, msg *TemplateMsgTest) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClient(userId, "TemplateMsgTest", buf)
}

func (s *server) ntf_TemplateMsgTest(userIdList []string, msg *TemplateMsgTest) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClients(userIdList, "TemplateMsgTest", buf)
}

func (s *server) send_TemplateMsgTestAck(userId string, msg *TemplateMsgTestAck) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClient(userId, "TemplateMsgTestAck", buf)
}

func (s *server) ntf_TemplateMsgTestAck(userIdList []string, msg *TemplateMsgTestAck) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClients(userIdList, "TemplateMsgTestAck", buf)
}

func (s *server) send_TemplateMsgTestNtf(userId string, msg *TemplateMsgTestNtf) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClient(userId, "TemplateMsgTestNtf", buf)
}

func (s *server) ntf_TemplateMsgTestNtf(userIdList []string, msg *TemplateMsgTestNtf) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClients(userIdList, "TemplateMsgTestNtf", buf)
}

func (s *server) getProtoUseList() []string {
	protoUseList := []string{
		"TemplateMsgTest",
		"TemplateMsgTestAck",
		"TemplateMsgTestNtf",
	}
	return protoUseList
}
