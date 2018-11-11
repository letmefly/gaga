package main

import (
	"utils"
)

type InHandlers func(*clientMsg) error
type OutHandlers func(*clientMsg) error

type group struct {
	id          string
	in          chan *clientMsg
	out         chan *clientMsg
	inHandlers  InHandlers
	outHandlers OutHandlers
	outList     map[string]chan *clientMsg
}

func newGroup() *group {
	g := &group{
		id:  utils.CreateUUID(),
		in:  make(chan *clientMsg, 1000),
		out: make(chan *clientMsg, 1000),
	}
	return g
}

func distoryGroup(g *group) {
	close(g.in)
	close(g.out)
}

func (g *group) start() {
	go func() {
		for {
			select {
			case inMsg, ok := <-g.in:
				if ok {
					g.inHandlers(inMsg)
				}
			}
		}
	}()
}

func (g *group) stop() {

}

func (g *group) handleMsg(msg *clientMsg) {
	select {
	case g.in <- msg:
		g.outList[msg.userId] = msg.out
	}
}

func (g *group) sendMsg(userId string, msg *clientMsg) {
	out := g.outList[userId]
	out <- msg
}
