// DOT NOT EDIT THIS FILE, AUTO GEN!!
package main

import (
	"context"
)

import (
	"pb"
	. "pb/auth"
	"types"
	"google.golang.org/grpc"
)


func (s *server) Login(ctx context.Context, param *LoginParam) (*LoginRet, error) {
	ret, err := s.logic.handleMessage(types.GRPC, "LoginParam", param)
	return ret.(*LoginRet), err
}


func (s *server) registerPbServers(grpcServer *grpc.Server) {
	pb.RegisterStreamServer(grpcServer, s)
	RegisterAuthServer(grpcServer, s)
}
