package mqtt_middleware

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var engine *Engine

func init() {
	engine = New()
}

type SubscriberFunc func(*Context)
type SubscribersChain []SubscriberFunc
type HandlerFunc func(*Context, mqtt.Client, mqtt.Message)

type Engine struct {
	subscribers SubscribersChain
	m           map[string]*subscribe
	stop        chan int
}

type subscribe struct {
	topic string
	qos   byte
	f     HandlerFunc
}

func New() *Engine {
	return &Engine{
		m:    make(map[string]*subscribe),
		stop: make(chan int, 1),
	}
}

func (s *Engine) Use(middleware ...SubscriberFunc) {
	s.subscribers = append(s.subscribers, middleware...)
}

func (s *Engine) Subscribe(topic string, qos byte, callback HandlerFunc) {
	s.m[topic] = &subscribe{topic, qos, callback}
}

func (s *Engine) Stop() {
	s.stop <- 1
}

func (s *Engine) Run(client mqtt.Client) error {
	for _, v := range s.m {
		f := v.f
		if token := client.Subscribe(v.topic, v.qos, func(client mqtt.Client, msg mqtt.Message) {
			c := &Context{
				subscribers: make(SubscribersChain, len(s.subscribers)),
				index:       -1,
			}
			copy(c.subscribers, s.subscribers)
			c.subscribers = append(c.subscribers, func(c *Context) {
				f(c, client, msg)
			})

			c.Next()
		}); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	<-s.stop
	return nil
}
