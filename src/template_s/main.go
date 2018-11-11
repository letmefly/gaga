package main

import (
	"context"
	"log"
	"net"
	"os"
)

import (
	"pb"
	"services"
	"utils"

	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	server := &Server{}
	server.Init(ctx)

	serviceAddr := "127.0.0.1:9990"
	serviceId := utils.CreateServiceId("template", serviceAddr)

	lis, err := net.Listen("tcp", serviceAddr)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	log.Println("Listen on", lis.Addr())

	s := grpc.NewServer()
	pb.RegisterStreamServer(s, server)
	services.Register(context.Background(), &services.ServiceConf{
		ServiceType:    "template",
		ServiceId:      serviceId,
		ServiceAddr:    serviceAddr,
		IsStream:       true,
		ProtoUseList:   server.GetProtoUseList(),
		ServiceUseList: []string{},
		TTL:            4,
	})
	s.Serve(lis)

	defer func() {
		cancel()
	}()
	select {}
}
