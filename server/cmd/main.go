package main

import (
	"github.com/mojtabamovahedi/chatroom/server/config"
	"github.com/mojtabamovahedi/chatroom/server/pkg/nats"
	nats2 "github.com/nats-io/nats.go"
	"log"
)

func main() {

	cfg := config.MustReadConfig("config.json")
	ch := make(chan string)
	n, err := nats.New(cfg.Nats.Host, cfg.Nats.Port)
	if err != nil {
		log.Fatal(err)
	}
	_, err = n.Subscribe(cfg.Nats.Subject, func(msg *nats2.Msg) {
		ch <- string(msg.Data)
	})

}
