package actors

import (
	"errors"
	"reflect"
	"time"
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
type timeoutTimerCb func()
type tickTimerCb func(int32)

type timeoutTimer struct {
	done      chan bool
	timeoutCb timeoutTimerCb
	count     int32
}

type tickTimer struct {
	done   chan bool
	tickCb tickTimerCb
	count  int32
}

type Actor struct {
	actorId        int32
	callQueue      chan *actorCall
	done           chan bool
	host           ActorHost
	timerCount     int32
	timeoutTimers  map[int32]*timeoutTimer
	tickTimers     map[int32]*tickTimer
	timeoutCbQueue chan *timeoutTimer
	tickCbQueue    chan *tickTimer
}

func (a *Actor) init(actorId int32, host ActorHost) {
	a.actorId = actorId
	a.host = host
	a.callQueue = make(chan *actorCall, 1000)
	a.done = make(chan bool, 1)
	a.timerCount = 0
	a.timeoutTimers = make(map[int32]*timeoutTimer, 0)
	a.tickTimers = make(map[int32]*tickTimer, 0)
	a.timeoutCbQueue = make(chan *timeoutTimer, 1000)
	a.tickCbQueue = make(chan *tickTimer, 1000)
	a.startLoop()
}

func (a *Actor) exit() {
	a.done <- true
}

func (a *Actor) GetTimerStats() (int, int) {
	return len(a.timeoutTimers), len(a.tickTimers)
}

func (a *Actor) ActorId() int32 {
	return a.actorId
}

func (a *Actor) assignTimerId() int32 {
	if a.timerCount >= (1<<31 - 10) {
		a.timerCount = 0
	}
	a.timerCount += 1
	return a.timerCount
}

func (a *Actor) SetTimeoutTimer(timeout time.Duration, cb timeoutTimerCb) int32 {
	timerId := a.assignTimerId()
	after := time.After(timeout)
	timer := &timeoutTimer{
		done:      make(chan bool, 1),
		timeoutCb: cb,
		count:     0,
	}
	a.timeoutTimers[timerId] = timer
	go func() {
		for {
			select {
			case <-after:
				a.timeoutCbQueue <- timer
				a.UnsetTimeoutTimer(timerId)
			case <-timer.done:
				return
			}
		}
	}()
	return timerId
}

func (a *Actor) UnsetTimeoutTimer(timerId int32) {
	timer, ok := a.timeoutTimers[timerId]
	if ok {
		timer.done <- true
		delete(a.timeoutTimers, timerId)
	}
}

func (a *Actor) SetTickTimer(du time.Duration, cb tickTimerCb) int32 {
	timerId := a.assignTimerId()
	ticker := time.NewTicker(du)
	timer := &tickTimer{
		done:   make(chan bool, 1),
		tickCb: cb,
		count:  0,
	}
	a.tickTimers[timerId] = timer
	go func() {
		for {
			select {
			case <-ticker.C:
				timer.count += 1
				a.tickCbQueue <- timer
			case <-timer.done:
				ticker.Stop()
				return
			}
		}
	}()
	return timerId
}

func (a *Actor) UnsetTickTimer(timerId int32) {
	timer, ok := a.tickTimers[timerId]
	if ok {
		timer.done <- true
		delete(a.tickTimers, timerId)
	}
}

func (a *Actor) SetLifeTickTimer(du time.Duration, ticks int32, cb tickTimerCb) int32 {
	timerId := a.assignTimerId()
	ticker := time.NewTicker(du)
	timer := &tickTimer{
		done:   make(chan bool, 1),
		tickCb: cb,
		count:  0,
	}
	a.tickTimers[timerId] = timer
	go func() {
		for {
			select {
			case <-ticker.C:
				timer.count += 1
				if timer.count > ticks {
					a.UnsetTickTimer(timerId)
				} else {
					a.tickCbQueue <- timer
				}

			case <-timer.done:
				ticker.Stop()
				return
			}
		}
	}()
	return timerId
}

func (a *Actor) UnsetLifeTickTimer(timerId int32) {
	timer, ok := a.tickTimers[timerId]
	if ok {
		timer.done <- true
		delete(a.tickTimers, timerId)
	}
}

func (a *Actor) putCall(callRetChan chan *actorCallRet, function interface{}, params ...interface{}) error {
	method, err := a.getHostMethod(false, function, params)
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

func (a *Actor) putAsynCall(callRetChan chan *actorCallRet, function interface{}, params ...interface{}) error {
	method, err := a.getHostMethod(true, function, params)
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

func (a *Actor) getHostMethod(isAsyn bool, function interface{}, params []interface{}) (reflect.Method, error) {
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
	if isAsyn == false {
		if funcT.NumOut() != 1 {
			return method, errors.New("call function must return only 1 result")
		}
		if funcT.Out(0).Name() != "error" {
			return method, errors.New("call function's second return type must be error")
		}
	} else {
		if funcT.NumOut() != 0 {
			return method, errors.New("call function has no return result")
		}
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

func (a *Actor) startLoop() {
	go func() {
		for {
			select {
			case tickTimer, _ := <-a.tickCbQueue:
				tickTimer.tickCb(tickTimer.count)

			case timeoutTimer, _ := <-a.timeoutCbQueue:
				timeoutTimer.timeoutCb()

			case call, _ := <-a.callQueue:
				if nil != call {
					ret := call.method.Func.Call(call.params)
					if call.callRetChan != nil {
						call.callRetChan <- &actorCallRet{results: ret}
					}
				}

			case <-a.done:
				return

			default:
			}
		}
	}()
}
