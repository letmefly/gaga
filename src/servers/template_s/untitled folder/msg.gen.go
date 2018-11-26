package main

import (
	"errors"
	"reflect"

	"github.com/golang/protobuf/proto"
)

func decode_TemplateMsgTest(codec string, msgName string, msgData []byte) (*TemplateMsgTest, error) {
	if codec == "protobuf" {
		msg := &TemplateMsgTest{}
		err := proto.Unmarshal(msgData, msg)
		return msg, err
	}
	return nil, errors.New("not support codec " + codec)
}

func encode_TemplateMsgTest(codec string, msgName string, msg *TemplateMsgTest) ([]byte, error) {
	if codec == "protobuf" {
		data, err := proto.Marshal(msg)
		return data, err
	}
	return nil, errors.New("not support codec " + codec)
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
		}
	} else if codec == "json" {

	}
}

func ToTemplateMsgTest(msg interface{}) (*TemplateMsgTest, error) {
	if reflect.TypeOf(msg).String() != "TemplateMsgTest" {
		return nil, errors.New("msg type error")
	}
	return msg.(*TemplateMsgTest), nil
}

func EncodeMessage(codec string, msg interface{}) ([]byte, error) {
	if codec == "protobuf" {
		buf, err := proto.Marshal(msg)
		if err != nil {
			return nil, err
		}
		return buf, nil
	} else if codec == "json" {
	}
	return nil, errors.New("no proto support for " + codec)
}

func GetProtoUseList() []string {
	protoUseList := []string{
		"TemplateMsgTest",
		"TemplateMsgTestAck",
		"TemplateMsgTestNtf",
	}
	return protoUseList
}
