package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
	"utils"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:12345", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	for i := 0; i < 20; i++ {
		log.Println(utils.HashCode(string("test-hash-code") + string(i)))
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/wsGate"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			recevSeq, msgId, msgData := utils.UnpackMsg(message)
			log.Printf("recv: %d %d %s", recevSeq, msgId, msgData)
		}
	}()

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	sendSeq := int32(0)
	sendSeq = sendSeq + 1
	msgId := utils.HashCode("login")
	tmpmsg := utils.PackMsg(sendSeq, msgId, []byte("fuck"))
	c.WriteMessage(websocket.TextMessage, tmpmsg)
	for {
		select {
		case <-done:
			return

		case t := <-ticker.C:
			sendSeq = sendSeq + 1
			msg := utils.PackMsg(sendSeq, utils.HashCode("join"), []byte(t.String()))

			err := c.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write:", err)
				return
			}

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
