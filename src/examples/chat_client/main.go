package main

import (
	"fmt"
	//"utils"
	//"time"
)

func main() {
	client := NewChatClient("127.0.0.1:12345", "chris", "123456")
	client.Start()
	fmt.Println("start chat client ...")
	for i := 0; i < 50; i++ {
		//fmt.Println(utils.CreateUUID2())
		//fmt.Println(utils.CreateUUID())
	}
	/*
		var sumCost int64
		sumCost = 0
		for i := 0; i < 10; i++ {
			beginStamp := time.Now().UnixNano()
			client.LoginAuth("chris", "123456")
			endStamp := time.Now().UnixNano()
			sumCost += (endStamp - beginStamp)
			cost := (endStamp - beginStamp) / (1000 * 1000)
			fmt.Println("call cost:", cost, " ms")
			//time.Sleep(1 * time.Millisecond)
		}

		fmt.Println("Avg Cost:", sumCost/(1000*1000), "us")

		fmt.Println("start chat client ...")
	*/

	select {}
}
