// DOT NOT EDIT THIS FILE, AUTO GEN!!
package main

import (
	"context"
)

import (
	"google.golang.org/grpc"
	"pb"
	. "pb/template"
	"types"
)

func (s *server) Test(ctx context.Context, param *TestParam) (*TestRet, error) {
	ret, err := s.logic.handleMessage(types.GRPC, "TestParam", param)
	return ret.(*TestRet), err
}

func (s *server) registerPbServers(grpcServer *grpc.Server) {
	pb.RegisterStreamServer(grpcServer, s)
	RegisterTemplateServer(grpcServer, s)
}
