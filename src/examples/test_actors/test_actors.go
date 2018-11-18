package main

import (
	"actors"
	"fmt"
)

type Room struct {
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

	for {
		//actors.Call(actorId, (*Room).Add, 200, 300)
	}

}
