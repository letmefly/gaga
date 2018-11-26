// DOT NOT EDIT THIS FILE, AUTO GEN!!
package auth

import (
	"errors"
	"reflect"
	"log"

	"github.com/golang/protobuf/proto"
)


func GetProtoUseList() []string {
	protoUseList := []string{
		"Test",
	}
	return protoUseList
}

func DecodeMessage(codec string, msgName string, msgData []byte) (interface{}, error) {
	if codec == "protobuf" {
		switch msgName {
			
		case "Test":
			msg := &Test{}
			err := proto.Unmarshal(msgData, msg)
			if err != nil {
				return nil, err
			} else {
				return msg, nil
			}

		}
	} else if codec == "json" {
	}
	return nil, errors.New("no proto support for " + codec)
}

func ToTest(msg interface{}) *Test {
	if reflect.TypeOf(msg).String() != "Test" {
		log.Panicln("msg type error")
	}
	return msg.(*Test)
}

func EncodeMessage(codec string, msg interface{}) ([]byte, error) {
	if codec == "protobuf" {
		buf, err := proto.Marshal(msg.(proto.Message))
		if err != nil {
			return nil, err
		}
		return buf, nil
	} else if codec == "json" {
	}
	return nil, errors.New("no proto support for " + codec)
}
