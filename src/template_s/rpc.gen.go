// DOT NOT EDIT THIS FILE, AUTO GEN!!
package main

import (
	"context"
)

import (
	"pb"
	. "pb/template"
	"google.golang.org/grpc"
)


func (s *server) Test(ctx context.Context, param *TestParam) (*TestRet, error) {
	ret, err := s.logic.handleMessage("TestParam", param)
	return ret.(*TestRet), err
}


func (s *server) registerPbServers(grpcServer *grpc.Server) {
	pb.RegisterStreamServer(grpcServer, s)
	RegisterTemplateServer(grpcServer, s)
}
