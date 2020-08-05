package mqtt_middleware

import (
	"math"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Context struct {
	handlers HandlersChain
	index    int

	Keys map[string]interface{}
}

const abortIndex int = math.MaxInt8 / 2

func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Keys[key]
	return
}

func (c *Context) Next(ctx *Context, client mqtt.Client, msg mqtt.Message) {
	c.index++
	for s := len(c.handlers); c.index < s; c.index++ {
		c.handlers[c.index](c, client, msg)
	}
}
