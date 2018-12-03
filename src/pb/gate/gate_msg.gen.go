// DOT NOT EDIT THIS FILE, AUTO GEN!!
package gate

import (
	"errors"
	"reflect"
	"log"

	"github.com/golang/protobuf/proto"
)


func GetProtoUseList() []string {
	protoUseList := []string{
		"GateMsgTest",
		"GateMsgTestAck",
		"LoginReq",
		"LoginAck",
	}
	return protoUseList
}

func DecodeMessage(codec string, msgName string, msgData []byte) (interface{}, error) {
	if codec == "protobuf" {
		switch msgName {
			
		case "GateMsgTest":
			msg := &GateMsgTest{}
			err := proto.Unmarshal(msgData, msg)
			if err != nil {
				return nil, err
			} else {
				return msg, nil
			}

		case "GateMsgTestAck":
			msg := &GateMsgTestAck{}
			err := proto.Unmarshal(msgData, msg)
			if err != nil {
				return nil, err
			} else {
				return msg, nil
			}

		case "LoginReq":
			msg := &LoginReq{}
			err := proto.Unmarshal(msgData, msg)
			if err != nil {
				return nil, err
			} else {
				return msg, nil
			}

		case "LoginAck":
			msg := &LoginAck{}
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

func ToGateMsgTest(msg interface{}) *GateMsgTest {
	if reflect.TypeOf(msg).String() != "GateMsgTest" {
		log.Panicln("msg type error")
	}
	return msg.(*GateMsgTest)
}

func ToGateMsgTestAck(msg interface{}) *GateMsgTestAck {
	if reflect.TypeOf(msg).String() != "GateMsgTestAck" {
		log.Panicln("msg type error")
	}
	return msg.(*GateMsgTestAck)
}

func ToLoginReq(msg interface{}) *LoginReq {
	if reflect.TypeOf(msg).String() != "LoginReq" {
		log.Panicln("msg type error")
	}
	return msg.(*LoginReq)
}

func ToLoginAck(msg interface{}) *LoginAck {
	if reflect.TypeOf(msg).String() != "LoginAck" {
		log.Panicln("msg type error")
	}
	return msg.(*LoginAck)
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
