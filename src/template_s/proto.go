package main

import (
	"github.com/golang/protobuf/proto"
)

func DecodeMessage(codec string, msgName string, msgData []byte) (interface{}, error) {
	if codec == "protobuf" {
		switch msgName {
		case "TestMsg":
			msg := &TestMsg{}
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

func EncodeMessage(codec string, msgName string, msg interface{}) ([]byte, error) {
	if codec == "protobuf" {

	} else if codec == "json" {

	}
}
