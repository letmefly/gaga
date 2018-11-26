// DOT NOT EDIT THIS FILE, AUTO GEN!!
package template

import (
	"errors"
	"reflect"
	"log"

	"github.com/golang/protobuf/proto"
)


func GetProtoUseList() []string {
	protoUseList := []string{
		"TemplateMsgTest",
		"TemplateMsgTestAck",
		"TemplateMsgTestNtf",
	}
	return protoUseList
}

func DecodeMessage(codec string, msgName string, msgData []byte) (interface{}, error) {
	if codec == "protobuf" {
		switch msgName {
			
		case "TemplateMsgTest":
			msg := &TemplateMsgTest{}
			err := proto.Unmarshal(msgData, msg)
			if err != nil {
				return nil, err
			} else {
				return msg, nil
			}

		case "TemplateMsgTestAck":
			msg := &TemplateMsgTestAck{}
			err := proto.Unmarshal(msgData, msg)
			if err != nil {
				return nil, err
			} else {
				return msg, nil
			}

		case "TemplateMsgTestNtf":
			msg := &TemplateMsgTestNtf{}
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

func ToTemplateMsgTest(msg interface{}) *TemplateMsgTest {
	if reflect.TypeOf(msg).String() != "TemplateMsgTest" {
		log.Panicln("msg type error")
	}
	return msg.(*TemplateMsgTest)
}

func ToTemplateMsgTestAck(msg interface{}) *TemplateMsgTestAck {
	if reflect.TypeOf(msg).String() != "TemplateMsgTestAck" {
		log.Panicln("msg type error")
	}
	return msg.(*TemplateMsgTestAck)
}

func ToTemplateMsgTestNtf(msg interface{}) *TemplateMsgTestNtf {
	if reflect.TypeOf(msg).String() != "TemplateMsgTestNtf" {
		log.Panicln("msg type error")
	}
	return msg.(*TemplateMsgTestNtf)
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
