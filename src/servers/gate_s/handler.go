package main

import (
	"context"
)
import (
	"pb/auth"
	"pb/gate"
	"services"
)

func gateHandler(agent *agent, msgName string, msgInterface interface{}) error {
	switch msgName {
	case "LoginReq":
		msg := gate.ToLoginReq(msgInterface)

		serviceClient, _ := services.GetServiceClient("auth")
		authClient := auth.NewAuthClient(serviceClient.Conn)
		loginRet, _ := authClient.Login(context.Background(), &auth.LoginParam{Account: msg.Account, Password: msg.Password})

		userId := loginRet.UserId
		sess, err := getSessionManager().createSession(userId)
		if err != nil {
			return err
		}
		// bindng each other
		sess.bindAgent(agent)
		msgData, _ := gate.EncodeMessage("protobuf", gate.LoginAck{Error: "ok", UserId: userId})
		agent.toClient(services.ToMsgId("LoginAck"), msgData)
	}
	return nil
}
