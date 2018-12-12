package main

import (
	"fmt"
)

func main() {
	var addr = "127.0.0.1:12345"
	client := NewChatClient()
	client.ConnectServer(addr)
	client.LoginAuth("chris", "123456")
	fmt.Println("start chat client ...")

	select {}
}
