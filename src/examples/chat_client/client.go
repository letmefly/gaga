package main

import (
	"flag"
	"fmt"
	"net/url"
	"reflect"
	"utils"

	"pb/gate"
	"pb/template"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

type ChatClient struct {
	Conn     *websocket.Conn
	seq      int32
	done     chan bool
	addr     string
	account  string
	password string
	userId   int32
}

func NewChatClient(addr string, account string, password string) *ChatClient {
	ret := &ChatClient{
		seq:      0,
		done:     make(chan bool, 0),
		addr:     addr,
		account:  account,
		password: password,
	}
	return ret
}

func (c *ChatClient) Start() {
	var host = flag.String("addr", c.addr, "http service address")
	u := url.URL{Scheme: "ws", Host: *host, Path: "/wsGate"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	c.Conn = conn
	if err != nil {
		fmt.Printf("dial:", err)
		return
	}
	c.StartRecvLoop()

	// start login and auto chat
	// message 1
	c.LoginAuth(c.account, c.password)
}

func (c *ChatClient) StartRecvLoop() {
	go func() {
		for {
			_, message, err := c.Conn.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				return
			}
			recevSeq, msgId, msgData := utils.UnpackMsg(message)
			fmt.Printf("recv: recevSeq %d, msgId %d\n", recevSeq, msgId)
			c.HandleMessage(msgId, msgData)
		}
	}()
}

func (c *ChatClient) Send(msg interface{}) {
	fmt.Println("call start")
	msgName := reflect.TypeOf(msg).String()
	msgName = strings.Replace(msgName, "*", "", 1)
	msgId := utils.HashCode(msgName)
	//fmt.Println(msgName, "msgId", msgId)
	buf, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	c.seq += 1
	pack := utils.PackMsg(c.seq, msgId, buf)
	c.Conn.WriteMessage(websocket.BinaryMessage, pack)
}

func (c *ChatClient) HandleMessage(msgId uint32, msgData []byte) {
	switch msgId {
	// message 2
	case utils.HashCode("gate.LoginAck"):
		loginAck := &gate.LoginAck{}
		err := proto.Unmarshal(msgData, loginAck)
		if err != nil {
			fmt.Errorf(err.Error())
			return
		}
		c.userId = loginAck.UserId
		//fmt.Println("loginAck:", loginAck.Error, loginAck.UserId)
		//c.done <- true
		//fmt.Println("call end")

		c.TemplateMsgTest()
	}
}

func (c *ChatClient) Distroy() {
	defer c.Conn.Close()
}

func (c *ChatClient) LoginAuth(account string, password string) {
	c.Send(&gate.LoginReq{Account: account, Password: password})
	/*
		select {
		case <-c.done:
			return
		}
	*/
}

func (c *ChatClient) TemplateMsgTest() {
	c.Send(&template.TemplateMsgTest{I: 10009, S: "this is test msg"})
}

func (c *ChatClient) GetRoom(userId string) {

}

func (c *ChatClient) JoinRoom(roomId string) {

}
