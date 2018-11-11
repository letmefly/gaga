package actors

import (
	"errors"
	//"fmt"
	"reflect"
)

type ActorHost interface {
}

type actorCallRet struct {
	results []reflect.Value
}
type actorCall struct {
	method      reflect.Method
	params      []reflect.Value
	callRetChan chan *actorCallRet
}

type actor struct {
	actorId   int32
	callQueue chan *actorCall
	host      ActorHost
}

func (a *actor) init(actorId int32, host ActorHost) {
	a.actorId = actorId
	a.host = host
	a.callQueue = make(chan *actorCall, 1000)
	a.startLoop()
}

func (a *actor) putCall(callRetChan chan *actorCallRet, function interface{}, params ...interface{}) error {
	method, err := a.getHostMethod(function, params)
	if err != nil {
		return err
	}
	methodParams := make([]reflect.Value, len(params)+1)
	methodParams[0] = reflect.ValueOf(a.host)
	for i := 0; i < len(params); i++ {
		methodParams[i+1] = reflect.ValueOf(params[i])
	}
	call := &actorCall{
		method:      method,
		params:      methodParams,
		callRetChan: callRetChan,
	}
	a.callQueue <- call
	return nil
}

func (a *actor) getHostMethod(function interface{}, params []interface{}) (reflect.Method, error) {
	var method reflect.Method
	funcT := reflect.TypeOf(function)
	funcV := reflect.ValueOf(function)
	hostT := reflect.TypeOf(a.host)
	hostV := reflect.ValueOf(a.host)
	if funcT == nil {
		return method, errors.New("call function is invalid")
	}
	if funcT.Kind() != reflect.Func || funcT.NumIn() < 1 {
		return method, errors.New("call function type is not Func")
	}
	if funcT.In(0) != hostT {
		return method, errors.New("dest actor id or call function is not right")
	}
	if funcT.NumOut() != 1 {
		return method, errors.New("call function must return only 1 result")
	}
	if funcT.Out(0).Name() != "error" {
		return method, errors.New("call function's second return type must be error")
	}
	// check function params
	if len(params) != funcT.NumIn()-1 {
		return method, errors.New("params num is not right")
	}
	for i := 1; i < funcT.NumIn(); i++ {
		t := reflect.TypeOf(params[i-1])
		if t.AssignableTo(funcT.In(i)) == false {
			return method, errors.New("papram type is wrong")
		}
	}
	// check if dest actor has this function
	isFind := false
	for i := 0; i < hostV.NumMethod(); i++ {
		if hostT.Method(i).Func.Pointer() == funcV.Pointer() {
			isFind = true
			method = hostT.Method(i)
			break
		}
	}
	if isFind == false {
		return method, errors.New("dest actor has no such call function")
	}
	return method, nil
}

func (a *actor) startLoop() {
	go func() {
		for {
			select {
			case call, _ := <-a.callQueue:
				if nil != call {
					ret := call.method.Func.Call(call.params)
					call.callRetChan <- &actorCallRet{results: ret}
				}
			default:
			}
		}
	}()
}
