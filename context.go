package mqtt_middleware

import (
	"math"
)

type Context struct {
	subscribers SubscribersChain
	index       int
	subscribe   *subscribe

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

func (c *Context) Topic() string {
	return c.subscribe.topic
}

func (c *Context) Next() {
	c.index++
	for s := len(c.subscribers); c.index < s; c.index++ {
		c.subscribers[c.index](c)
	}
}
