package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Config struct {
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	go httpServer(nil)
	select {}
}

func httpServer(config *Config) {
	http.HandleFunc("/httpGate", httpGateHandler)
	http.HandleFunc("/wsGate", wsGateHandler)
	log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
	//err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	err := http.ListenAndServe(":12345", nil)
	log.Fatal(err)
	fmt.Println("exit httpServer")
}

func httpGateHandler(w http.ResponseWriter, req *http.Request) {

}

func wsGateHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
}
