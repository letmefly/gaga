package main

var server_gen_header_tpl = `// DOT NOT EDIT THIS FILE, AUTO GEN!!
package main

import (
	"log"
)

import (
	"pb"
	. "pb/#{server_name}"
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
`

var server_gen_rpc_header_tpl = `// DOT NOT EDIT THIS FILE, AUTO GEN!!
package main

import (
	"context"
)

import (
	"pb"
	. "pb/#{server_name}"
	"google.golang.org/grpc"
)

`

var server_gen_protouse_tpl = `
func (s *server) getProtoUseList() []string {
	protoUseList := []string{#{proto_use_list}
	}
	return protoUseList
}
`

var server_gen_grpc_tpl = `
func (s *server) #{rpc_name}(ctx context.Context, param *#{rpc_param}) (*#{rpc_ret}, error) {
	ret, err := s.logic.handleMessage("#{rpc_param}", param)
	return ret.(*#{rpc_ret}), err
}
`
var server_gen_pbregister_tpl = `
func (s *server) registerPbServers(grpcServer *grpc.Server) {
	pb.RegisterStreamServer(grpcServer, s)
	Register#{service_name}Server(grpcServer, s)
}
`

var server_gen_register_part1_tpl = `func (s *server) registerClientHandlers() {`
var server_gen_register_part2_tpl = `
	s.registerHandler("#{message_name}", func(data []byte) {
		msg := &#{message_name}{}
		err := proto.Unmarshal(data, msg)
		if err != nil {
			log.Println(err)
		} else {
			s.logic.handleMessage("#{message_name}", msg)
		}
	})
`
var server_gen_register_part3_tpl = `}`

var server_gen_sendfunc_tpl = `
func (s *server) send_#{message_name}(userId string, msg *#{message_name}) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClient(userId, "#{message_name}", buf)
}
`

var server_gen_ntffunc_tpl = `
func (s *server) ntf_#{message_name}(userIdList []string, msg *#{message_name}) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClients(userIdList, "#{message_name}", buf)
}
`