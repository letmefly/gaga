package services

import (
	"utils"
)

type InMsgHandlers func(*clientMsg) error
type OutMsgHandlers func(*clientMsg) error

type clientMsg struct {
	name string
	data []byte
}

type Module struct {
	id          string
	in          chan *clientMsg
	out         chan *clientMsg
	inHandlers  InMsgHandlers
	outHandlers OutMsgHandlers
}

func NewModule(inHandlers InMsgHandlers, outHandlers OutMsgHandlers) *Module {
	m := &Module{
		id:          utils.CreateUUID(),
		in:          make(chan *clientMsg, 1000),
		out:         make(chan *clientMsg, 1000),
		inHandlers:  inHandlers,
		outHandlers: outHandlers,
	}
	return m
}

func DistroyModule(m *Module) {
	close(m.in)
	close(m.out)
}

func (m *Module) Start() {
	go func() {
		for {
			select {
			case inMsg, ok1 := <-m.in:
				if ok1 {
					m.inHandlers(inMsg)
				}
			case outMsg, ok2 := <-m.out:
				if ok2 {
					m.outHandlers(outMsg)
				}
			}
		}
	}()
}

func (m *Module) Stop() {

}
