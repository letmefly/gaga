package main

import (
	"actors"
	"fmt"
	"time"
)

type Room struct {
	actors.BaseActorHost
}

func (r *Room) Add(a int, b int, out *int) error {
	fmt.Println("call:", a, b)
	*out = a + b
	return nil
}

type UserList struct {
	Users []string
}

func (r *Room) GetUserList(userId string, out *UserList) error {
	out.Users = []string{"aaaaa", "bbbbbb", "ccccc", userId}
	return nil
}

func (r *Room) Join(userId string) {
	fmt.Println("Join Room: ", userId)
}

func main() {
	actors.Init()
	room := &Room{}
	actorId := actors.NewActor(room)
	actor := room.Actor()
	fmt.Println("actorId:", actorId)
	var ret int
	err := actors.Call(actorId, (*Room).Add, 100, 200, &ret)
	if err != nil {
		fmt.Println(err)
	}

	if ret == 300 {
		fmt.Println("ret:", ret)
	}
	var userList UserList
	err = actors.Call(actorId, (*Room).GetUserList, "xxxxxx", &userList)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(userList.Users)

	err = actors.AsynCall(actorId, (*Room).Join, "chris-li")
	if err != nil {
		fmt.Println(err)
	}

	actor.SetTimeoutTimer(5000*time.Millisecond, func() {
		fmt.Println("TimeoutTimer is time out now")
	})

	actor.SetLifeTickTimer(1000*time.Millisecond, 10, func(i int32) {
		fmt.Println("Life Ticker: tick ", i)
		if i == 10 {
			fmt.Println("Life Tiker is over")
		}
	})

	actor.SetTickTimer(1000*time.Millisecond, func(i int32) {
		fmt.Println("Test Tick:", i)
		a, b := actor.GetTimerStats()
		fmt.Println("timeout timers: ", a, " tick timers: ", b)
	})

	for {
		//actors.Call(actorId, (*Room).Add, 200, 300)
	}

}
