package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name string
}

type User2 struct {
	Name string
}

func (u *User) Hello() {
	fmt.Println("hello")
}
func (u *User) Hello2(a int, s string) string {
	fmt.Println("hello")
	return ""
}
func (u *User) ShakeHand(name string) {
	fmt.Printf("Shake hand with %s\n", name)
}
func main() {
	//user1 := &User{Name: "xxxx"}
	funcT := (*User).Hello

	fmt.Println("user1 func pointer:", reflect.ValueOf(funcT).Pointer())
	fmt.Println("f string:", reflect.TypeOf(funcT))

	user2 := &User{Name: "xxxx"}
	funcT2 := (*User).Hello
	fmt.Println("user2 func's struct:", reflect.TypeOf(funcT2).In(0))
	fmt.Println("user2 func pointer:", reflect.ValueOf(funcT2).Pointer())
	//user2 := user1
	v := reflect.ValueOf(user2)
	//t := v.Type()
	t := reflect.TypeOf(user2)
	fmt.Println("use 2 type name", t.String())
	user3 := &User2{}
	if reflect.TypeOf(funcT2).In(0) == reflect.TypeOf(user3) {
		fmt.Println("===")
	}
	for i := 0; i < v.NumMethod(); i++ {
		//fmt.Printf("method[%d]%s\n", i, t.Method(i).Name)
		fmt.Println(t.Method(i).Name, "func name:", t.Method(i).Func.Pointer())
	}
}
