package mqtt_middleware

import (
	"fmt"
	"sync"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const TOPIC = "druid/topic/test"

func Test_accepting_new_client_callback(t *testing.T) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://192.168.88.16:1883")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
	}

	fmt.Println("init logger")
	engine.Use(func(c *Context) {
		c.Set("logger", "logger value")
		fmt.Println("Init logger success.")
	})

	var wg sync.WaitGroup
	wg.Add(1)

	fmt.Println("init subscribe")
	engine.Subscribe(TOPIC, 0, func(c *Context, client mqtt.Client, msg mqtt.Message) {
		logger, _ := c.Get("logger")
		fmt.Println("value: ", logger.(string), msg.Topic())
		if string(msg.Payload()) != "mymessage" {
			t.Fatalf("want mymessage, got %s", msg.Payload())
		}
		fmt.Println("recv msg:", string(msg.Payload()))
		wg.Done()
	})

	go func() {
		time.Sleep(time.Second)
		token := client.Publish(TOPIC, 1, false, "mymessage")
		fmt.Println("token wait:", token.Wait())
		if token.Error() != nil {
			t.Fatal(token.Error())
		}
		fmt.Println("publish msg finish.")
		engine.Stop()
	}()

	fmt.Println("run subscribe")
	engine.Run(client)

	fmt.Println("wait subscribe")
	wg.Wait()
}
