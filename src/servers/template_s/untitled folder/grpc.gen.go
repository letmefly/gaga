package main

import (
	"pb"
)

type GrpcServer struct {
}

func (s *GrpcServer) CreateStream(stream pb.Stream_CreateStreamServer) error {
	return nil
}
