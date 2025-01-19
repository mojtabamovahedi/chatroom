package http

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/mojtabamovahedi/chatroom/server/pkg/nats"
	nats2 "github.com/nats-io/nats.go"
	"log"
	"sync"
)

func chatroomWebsocket(natsClient *nats.Nats) fiber.Handler {
	var (
		subject string
		wg      sync.WaitGroup
		ch      = make(chan string)
	)
	return websocket.New(func(c *websocket.Conn) {
		chatId := c.Params("chatId")
		_ = chatId
		subscribe, subErr := natsClient.Subscribe(subject, func(msg *nats2.Msg) {
			ch <- string(msg.Data)
		})

		if subErr != nil {
			_ = c.WriteJSON(fiber.Map{
				"error": subErr.Error(),
			})
			return
		}

		defer func() {
			fmt.Println("connection closed")
			_ = subscribe.Unsubscribe()
		}()

		wg.Add(2)

		go func() {
			var (
				mt  int
				msg []byte
				err error
			)
			_ = mt
			defer wg.Done()
			for {
				if mt, msg, err = c.ReadMessage(); err != nil {
					log.Println("read:", err)
					err = natsClient.Publish(subject, msg)
					if err != nil {
						log.Println("error in read message:", err)
					}
				}
			}
		}()

		go func() {
			defer wg.Done()
			for {
				v, ok := <-ch
				if !ok {
					break
				}
				err := c.WriteJSON(fiber.Map{
					"message": v,
				})
				if err != nil {
					log.Println("write:", err)
					break
				}
			}
		}()

		wg.Wait()

	})
}
