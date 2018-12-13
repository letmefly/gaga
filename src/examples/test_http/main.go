package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"pb/gate"
	"reflect"
	"utils"

	"strings"

	"github.com/golang/protobuf/proto"
)

func main() {
	response, _ := http.Get("http://127.0.0.1:12345/httpGate?sessId=xxdfdd")
	//defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	//tmp := `{"name":"junneyang", "age": 88}`
	pack := PackMsg(&gate.LoginReq{Account: "chris.li", Password: "123456"})
	postReq := bytes.NewBuffer([]byte(pack))
	bodyType := "application/json;charset=utf-8"
	resp, _ := http.Post("http://127.0.0.1:12345/httpGate?sessId=xxdfdd", bodyType, postReq)
	body, _ = ioutil.ReadAll(resp.Body)
	_, _, msgData := utils.UnpackMsg(body)
	loginAck := &gate.LoginAck{}
	proto.Unmarshal(msgData, loginAck)
	fmt.Println("loginAck:", loginAck.String())
	select {}
}

func PackMsg(msg interface{}) []byte {
	msgName := reflect.TypeOf(msg).String()
	msgName = strings.Replace(msgName, "*", "", 1)
	msgId := utils.HashCode(msgName)
	buf, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		fmt.Errorf(err.Error())
		return nil
	}
	seq := 1
	pack := utils.PackMsg(int32(seq), msgId, buf)
	return pack
}
