// tstSubscr.go
package main

import (
	"fmt"
	"time"

	"./Models"
	nats "github.com/nats-io/go-nats"
)

var nc *nats.Conn

func init() {
	Models.Cfg.Kind = "postgres"
	Models.Cfg.Transport = "tcp"
	Models.Cfg.ServerIP = "127.0.0.1"
	Models.Cfg.ServerPort = "5432"
	Models.Cfg.DbName = "mytest"
	Models.Cfg.User = "postgres"
	Models.Cfg.Password = "bobics"
	Models.Start()
}

func subscriber1(m *nats.Msg) {
	fmt.Println("1", string(m.Data), m.Subject)
	//	Models.New(data, news)

}
func subscriber2(m *nats.Msg) {
	fmt.Println("2", string(m.Data), m.Subject)
	//	Models.GetId(id)
}
func subscriber3(m *nats.Msg) {
	fmt.Println("3", string(m.Data), m.Subject)
	//	Models.GetData(data)
	nc.Publish("help", []byte("I can help!"))
}
func main() {
	var err error
	fmt.Println("Test subscriber!")
	nc, err = nats.Connect("nats://192.168.99.100:4222", nats.PingInterval(20*time.Second))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Nats is here!")
		// Simple Async Subscriber
		nc.Subscribe("new.*", subscriber1)
		nc.Subscribe("id.*", subscriber2)
		nc.Subscribe("data.*", subscriber3)

		for {
			time.Sleep(1 * time.Second)
		}

	}
}
