package main

import (
	"context"
	"log"
	"net"
	"os"
)

import (
	"pb"
	"pb/auth"
	"services"
	"utils"

	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	server := &Server{}
	server.Init(ctx)

	serviceAddr := "127.0.0.1:9991"
	serviceId := utils.CreateServiceId("auth", serviceAddr)

	lis, err := net.Listen("tcp", serviceAddr)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	log.Println("Listen on", lis.Addr())

	s := grpc.NewServer()
	pb.RegisterStreamServer(s, server)
	auth.RegisterAuthServer(s, server)

	services.Register(context.Background(), &services.ServiceConf{
		ServiceType:    "auth",
		ServiceId:      serviceId,
		ServiceAddr:    serviceAddr,
		IsStream:       true,
		ProtoUseList:   auth.GetProtoUseList(),
		ServiceUseList: []string{},
		TTL:            4,
	})
	s.Serve(lis)

	defer func() {
		cancel()
	}()
	select {}
}
