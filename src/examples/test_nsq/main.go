package main

import (
	"encoding/json"
	"events"
	"fmt"
)

type TestEvent struct {
	K string
	V int
}

func main() {
	events.Init("127.0.0.1:4150")

	events.Register("topic_test", "ch01", func(data []byte) {
		//fmt.Println(string(data))
		var t TestEvent
		json.Unmarshal(data, &t)
		fmt.Println(t.K, t.V)
	})

	for i := 0; i < 10; i++ {
		e := &TestEvent{K: "key", V: i}
		data, _ := json.Marshal(e)
		events.Publish("topic_test", data)
	}
	select {}
}
