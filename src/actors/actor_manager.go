package actors

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
)

var (
	_actorManager actorManager
	_once         sync.Once
)

func Init() {
	_once.Do(func() { _actorManager.init() })
}

func NewActor(host ActorHost) int32 {
	actorId := _actorManager.createActor(host)
	return actorId
}

func FreeActor(actor *actor) {
}
func Call(actorId int32, function interface{}, params ...interface{}) error {
	currActor, ok := _actorManager.getActor(actorId)
	if !ok {
		log.Println("ERR: no actor find for ", actorId)
		return errors.New(fmt.Sprintf("ERR: no actor find for %d", actorId))
	}
	callRetChan := make(chan *actorCallRet, 0)
	err := currActor.putCall(callRetChan, function, params...)
	if err != nil {
		return err
	}
	callRet, ok := <-callRetChan
	if !ok {
		return errors.New("ret chan error")
	}
	// check results validation
	funcT := reflect.TypeOf(function)
	if len(callRet.results) != 1 {
		return errors.New("Call return must be error")
	}
	if callRet.results[0].IsNil() == false {
		if callRet.results[0].Type().AssignableTo(funcT.Out(0)) == false {
			return errors.New("dest call function first return type is wrong")
		}
	}
	return nil
}

/*
func Call(actorId int32, function interface{}, params ...interface{}) (interface{}, error) {
	currActor, ok := _actorManager.getActor(actorId)
	if !ok {
		log.Println("ERR: no actor find for ", actorId)
		return nil, errors.New(fmt.Sprintf("ERR: no actor find for %d", actorId))
	}
	callRetChan := make(chan *actorCallRet, 0)
	err := currActor.putCall(callRetChan, function, params...)
	if err != nil {
		return nil, err
	}
	callRet, ok := <-callRetChan
	if !ok {
		return nil, errors.New("ret chan error")
	}

	// check results validation
	if len(callRet.results) != 2 {
		return nil, errors.New("Call return must 2 results, one is Msg, the other is error")
	}

	if callRet.results[1].IsNil() == false {
		if callRet.results[1].Type().AssignableTo(reflect.TypeOf(error)) == false {
			return nil, errors.New("dest call function first return type is wrong)
		}
		return nil, callRet.results[1].Interface().(error)
	}
	if callRet.results[0].IsNil() == false {
		if callRet.results[0].Type().AssignableTo(funcT.Out(0)) == false {
			return nil, errors.New("dest call function first return type is wrong)
		}
	}

	return callRet.results[0].Interface(), nil
}
*/
type actorManager struct {
	count  int32
	actors sync.Map
}

func (m *actorManager) init() {
	m.count = 0
}

func (m *actorManager) assignId() int32 {
	if m.count >= (1<<31 - 10) {
		m.count = 0
	}
	m.count += 1
	return m.count
}

func (m *actorManager) createActor(host ActorHost) int32 {
	actorId := m.assignId()
	currActor := &actor{
		actorId: actorId,
		host:    host,
	}
	currActor.init(actorId, host)
	m.actors.Store(actorId, currActor)
	return actorId
}

func (m *actorManager) getActor(actorId int32) (*actor, bool) {
	currActor, ok := m.actors.Load(actorId)
	if ok {
		return currActor.(*actor), true
	}
	return nil, false
}
