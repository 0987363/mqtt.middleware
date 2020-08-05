package mqtt_middleware

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var engine *Engine

func init() {
	engine = New()
}

type HandlerFunc func(*Context, mqtt.Client, mqtt.Message)
type HandlersChain []HandlerFunc

type HandlersMiddware struct {
	Handlers HandlersChain
	Index    int
}

type Engine struct {
	handlers HandlersChain
	m        map[string]*subscribe
	stop     chan int
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

func (s *Engine) Use(middleware ...HandlerFunc) {
	s.handlers = append(s.handlers, middleware...)
}

func (s *Engine) Subscribe(topic string, qos byte, callback HandlerFunc) {
	s.m[topic] = &subscribe{topic, qos, callback}
}

func (s *Engine) Stop() {
	s.stop <- 1
}

func (s *Engine) Run(client mqtt.Client) error {
	for _, v := range s.m {
		if token := client.Subscribe(v.topic, v.qos, func(client mqtt.Client, msg mqtt.Message) {
			c := &Context{
				handlers: make(HandlersChain, len(engine.handlers)),
				index:    -1,
			}
			copy(c.handlers, engine.handlers)

			c.Next(c, client, msg)
			v.f(c, client, msg)
		}); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	<-s.stop
	return nil
}
