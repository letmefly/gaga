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

func NewActor(host ActorHost) *Actor {
	actor := _actorManager.createActor(host)
	return actor
}

func FreeActor(actorId int32) {
	_actorManager.freeActor(actorId)
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

func AsynCall(actorId int32, function interface{}, params ...interface{}) error {
	currActor, ok := _actorManager.getActor(actorId)
	if !ok {
		log.Println("ERR: no actor find for ", actorId)
		return errors.New(fmt.Sprintf("ERR: no actor find for %d", actorId))
	}
	err := currActor.putAsynCall(nil, function, params...)
	if err != nil {
		return err
	}
	return nil
}

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

func (m *actorManager) createActor(host ActorHost) *Actor {
	actorId := m.assignId()
	currActor := &Actor{
		actorId: actorId,
		host:    host,
	}
	currActor.init(actorId, host)
	m.actors.Store(actorId, currActor)
	return currActor
}

func (m *actorManager) getActor(actorId int32) (*Actor, bool) {
	currActor, ok := m.actors.Load(actorId)
	if ok {
		return currActor.(*Actor), true
	}
	return nil, false
}

func (m *actorManager) freeActor(actorId int32) {
	currActor, ok := m.actors.Load(actorId)
	if ok {
		currActor.(*Actor).exit()
		m.actors.Delete(actorId)
	}
}
