package events

import (
	"log"
	"net"

	"github.com/nsqio/go-nsq"
)

type EventCB func(string)

type consumerHandler struct {
	eventCB EventCB
}

func (h *consumerHandler) HandleMessage(message *nsq.Message) error {
	h.eventCB("json-format")
	return nil
}

type EventManager struct {
	producer  *nsq.Producer
	consumers map[string]*nsq.Consumer
	addr      string
}

var defaultEventManager EventManager

func Init(addr string) {
	config := nsq.NewConfig()
	config.LocalAddr, _ = net.ResolveTCPAddr("tcp", addr+":0")
	producer, _ := nsq.NewProducer(addr, config)
	defaultEventManager.consumers = make(map[string]*nsq.Consumer, 0)
	defaultEventManager.producer = producer
	defaultEventManager.addr = addr
}

func Register(topic string, channel string, eventCB EventCB) {
	key := topic + channel
	config := nsq.NewConfig()
	consumer, _ := nsq.NewConsumer(topic, channel, config)
	handler := &consumerHandler{
		eventCB: eventCB,
	}
	consumer.AddHandler(handler)
	err := consumer.ConnectToNSQD(defaultEventManager.addr)
	if err != nil {
		log.Fatal(err)
	}
	preConsumer, ok := defaultEventManager.consumers[key]
	if ok {
		preConsumer.Stop()
	}
	defaultEventManager.consumers[key] = consumer
}

func Unregister(topic string, channel string) {
	key := topic + channel
	consumer, ok := defaultEventManager.consumers[key]
	if ok {
		consumer.Stop()
	}
	delete(defaultEventManager.consumers, key)
}

func Publish(topic string, eventData []byte) {
	err := defaultEventManager.producer.Publish(topic, eventData)
	if err != nil {
		log.Fatal(err)
	}
}
