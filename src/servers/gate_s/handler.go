package main

import (
	"context"
	"log"
)
import (
	"pb/auth"
	"pb/gate"
	"services"
)

func gateHandler(agent *agent, msgName string, msgInterface interface{}) error {
	switch msgName {
	case "gate.LoginReq":
		msg := gate.ToLoginReq(msgInterface)
		log.Println("gateHandler", msg.Account, msg.Password)
		serviceId, err1 := services.AssignServiceId("auth")
		if err1 != nil {
			log.Println(err1)
			return err1
		}
		defer services.UnassignServiceId(serviceId)
		serviceClient, err2 := services.GetServiceClient(serviceId)
		if err2 != nil {
			log.Println(err2)
			return err2
		}
		authClient := auth.NewAuthClient(serviceClient.Conn)
		loginRet, err3 := authClient.Login(context.Background(), &auth.LoginParam{Account: msg.Account, Password: msg.Password})
		if err3 != nil {
			log.Println(err3)
			return err3
		}
		userId := loginRet.UserId
		sess, err4 := getSessionManager().createSession(userId)
		if err4 != nil {
			log.Println(err4)
			return err4
		}
		// bindng each other
		sess.bindAgent(agent)
		msgData, _ := gate.EncodeMessage("protobuf", &gate.LoginAck{Error: "ok", UserId: userId})
		agent.toClient(services.ToMsgId("gate.LoginAck"), msgData)
	}
	return nil
}
