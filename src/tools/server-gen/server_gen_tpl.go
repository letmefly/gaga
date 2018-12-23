package main

var grpc_gen_header_tpl = `// DOT NOT EDIT THIS FILE, AUTO GEN!!
package #{server_name}

import (
	"context"
	"errors"
)

type GrpcServer struct {
}
`

var msg_gen_header_tpl = `// DOT NOT EDIT THIS FILE, AUTO GEN!!
package #{server_name}

import (
	"errors"
	"reflect"
	"log"

	"github.com/golang/protobuf/proto"
)

`

var msg_gen_protouse_tpl = `
func GetProtoUseList() []string {
	protoUseList := []string{#{proto_use_list}
	}
	return protoUseList
}
`

var msg_gen_grpc_tpl = `
func (s *GrpcServer) #{rpc_name}(ctx context.Context, param *#{rpc_param}) (*#{rpc_ret}, error) {
	return nil, errors.New("this api not support")
}
`
var msg_gen_decode_tpl = `
func DecodeMessage(codec string, msgName string, msgData []byte) (interface{}, error) {
	if codec == "protobuf" {
		switch msgName {
			#{decode_case_list}
		}
	} else if codec == "json" {
	}
	return nil, errors.New("no proto support for " + codec)
}
`

var msg_gen_decode_case_tpl = `
		case "#{server_name}.#{message_name}":
			msg := &#{message_name}{}
			err := proto.Unmarshal(msgData, msg)
			if err != nil {
				return nil, err
			} else {
				return msg, nil
			}
`

var msg_gen_to_msg_tpl = `
func To#{message_name}(msg interface{}) *#{message_name} {
	if reflect.TypeOf(msg).String() != "*#{server_name}.#{message_name}" {
		log.Panicln("msg type error")
	}
	return msg.(*#{message_name})
}
`

var msg_gen_encode_tpl = `
func EncodeMessage(codec string, msg interface{}) ([]byte, error) {
	if codec == "protobuf" {
		buf, err := proto.Marshal(msg.(proto.Message))
		if err != nil {
			return nil, err
		}
		return buf, nil
	} else if codec == "json" {
	}
	return nil, errors.New("no proto support for " + codec)
}
`

/*
var msg_gen_pbregister_tpl = `
func RegisterPbServers(s *server, grpcServer *grpc.Server) {
	pb.RegisterStreamServer(grpcServer, s)
	Register#{service_name}Server(grpcServer, s)
}
`


var server_gen_register_part1_tpl = `func (s *server) registerClientHandlers() {`
var server_gen_register_part2_tpl = `
	s.registerHandler("#{message_name}", func(network int, data []byte) {
		msg := &#{message_name}{}
		err := proto.Unmarshal(data, msg)
		if err != nil {
			log.Println(err)
		} else {
			s.logic.handleMessage(network, "#{message_name}", msg)
		}
	})
`
var server_gen_register_part3_tpl = `}`

var server_gen_sendfunc_tpl = `
func (s *server) send_#{message_name}(network int, userId string, msg *#{message_name}) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClient(network, userId, "#{message_name}", buf)
}
`

var server_gen_ntffunc_tpl = `
func (s *server) ntf_#{message_name}(network int, userIdList []string, msg *#{message_name}) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	s.sendClients(network, userIdList, "#{message_name}", buf)
}
`
*/
