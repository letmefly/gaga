package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

import (
	"services"
	"utils"

	"github.com/gorilla/websocket"
	//log "github.com/sirupsen/logrus"
)

type Config struct {
	Test string
	List []int
	Hash map[string]string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	/*
			testArray2 := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			testArray2[0] = 99
			log.Println(testArray2)

			var testArray []int = make([]int, 0)
			//testArray[0] = 100
			testArray = append(testArray, 10)
			testArray = append(testArray, 20)
			testArray = append(testArray, 30)
			testArray = append(testArray, 40)
			testArray = append(testArray, 5)

			testArray = append(testArray[:1], testArray[2:]...)

			for i := 0; i < len(testArray); i++ {
				log.Println("slice:", testArray[i])
			}
			for i, v := range testArray {
				log.Println("range:", i, v)
			}
			log.Println(testArray)

			var testHash map[string]int
			if testHash == nil {
				log.Println("map is nil. Going to make one.")
				testHash = make(map[string]int, 0)
			}
			testHash["test_key2"] = 678
			testHash["test_key"] = 908
			for k, v := range testHash {
				log.Println(k, v)
			}

			testHash2 := make(map[string]int, 0)
			testHash2["test_key2"] = 111
			testHash2["test_key"] = 222
			testHash2["k"] = 222
			log.Println(testHash2)
			delete(testHash2, "k")
			log.Println(testHash2)
			for k, v := range testHash2 {
				log.Println(k, v)
			}

			log.Println("test_key", testHash)

			var testInt int
			log.Println("testInt:", testInt)
			log.Println("_test_global:", _test_global)
			var c Config
			c.Test = "test c"
			log.Println("test:", c.Test)

			d := Config{}
			d.Test = "test d"
			log.Println("test:", d.Test)
			d.List = append(d.List, 456)
			log.Println("test:", d.List)

			e := &Config{Test: "test e"}
			log.Println("test:", e.Test)

			f := new(Config)
			f.Test = "test f"
			log.Println("test:", f.Test)

		uuid := utils.CreateUUID()
		log.Println("uuid:", uuid)
		b64 := utils.Base64EncodeV2(uuid)
		log.Println("base64v2:", b64)
		reb64 := utils.Base64DecodeV2(b64)
		log.Println("base64decodev2", reb64)

		reb64_2 := utils.Base64DecodeV1(b64)
		log.Println("base64Decodev1", reb64_2)
	*/
	currServiceAddr := "127.0.0.1:9999"
	currServiceId := utils.CreateServiceId("gate", currServiceAddr)

	//ctx := context.WithValue(context.Background(), "curr_service_id", currServiceId)
	//ctx = context.WithValue(ctx, "curr_service_addr", currServiceAddr)
	ctx := context.Background()
	services.Register(ctx, &services.ServiceConf{
		ServiceType:    "gate",
		ServiceId:      currServiceId,
		ServiceAddr:    currServiceAddr,
		IsStream:       false,
		ProtoUseList:   []string{"login"},
		ServiceUseList: []string{"template"},
		TTL:            4,
	})
	go httpServer(nil)
	select {}
}

func httpServer(config *Config) {
	http.HandleFunc("/httpGate", httpHandler)
	http.HandleFunc("/wsGate", wsHandler)
	log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
	//err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	err := http.ListenAndServe(":12345", nil)
	log.Fatal(err)
	fmt.Println("exit httpServer")
}

func httpHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("New Http Connection")
	vars := req.URL.Query()
	sessId := vars["sessId"][0]
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	agent := newAgent(sessId, func(data []byte) error {
		_, err := w.Write(data)
		cancel()
		return err
	})
	agent.start(ctx)
	defer freeAgent(agent)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	agent.toService(body)
	select {
	case <-ctx.Done():
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go handleWsClient(conn, nil)
}

func handleWsClient(conn *websocket.Conn, config *Config) {
	log.Println("New Websocket Connection")
	ctx, cancel := context.WithCancel(context.Background())
	sessId := ""
	agent := newAgent(sessId, func(data []byte) error {
		messageType := websocket.TextMessage
		err := conn.WriteMessage(messageType, data)
		return err
	})
	agent.start(ctx)

	// receive client message loop
	defer func() {
		conn.Close()
		cancel()
		freeAgent(agent)
		log.Printf("webscocket connection closed")
	}()

	for {
		// messageType is websocket.BinaryMessage or websocket.TextMessage.
		// message is a []byte
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		switch messageType {
		case websocket.TextMessage:
			agent.toService(message)
		case websocket.BinaryMessage:
		case websocket.CloseMessage:
		case websocket.PingMessage:
		case websocket.PongMessage:
		}
	}
}
