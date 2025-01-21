package http

import (
	"context"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/mojtabamovahedi/chatroom/server/pkg/nats"
	nats2 "github.com/nats-io/nats.go"
	"log"
	"sync"
)

func chatroomWebsocket(natsClient *nats.Nats) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		chatId := c.Params("chatId")
		var (
			subject = fmt.Sprintf("%s.%s", nats.BaseSubject, chatId)
			wg      sync.WaitGroup
			ch      = make(chan []byte)
		)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		subscribe, subErr := natsClient.Subscribe(subject, func(msg *nats2.Msg) {
			select {
			case ch <- msg.Data:
			case <-ctx.Done():
			}
		})

		if subErr != nil {
			_ = c.WriteJSON(fiber.Map{
				"error": subErr.Error(),
			})
			return
		}

		defer func() {
			log.Println("connection closed")
			_ = subscribe.Unsubscribe()
			close(ch)
		}()

		wg.Add(2)

		go func() {
			defer func() {
				wg.Done()
				cancel()
			}()
			for {
				_, msg, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					break
				}
				err = natsClient.Publish(subject, msg)
				if err != nil {
					log.Println("error in publish message:", err)
				}
			}
		}()

		go func() {
			defer func() {
				wg.Done()
				cancel()
			}()
			for {
				select {
				case v, ok := <-ch:
					if !ok {
						return
					}
					err := c.WriteJSON(fiber.Map{
						"message": string(v),
					})
					if err != nil {
						log.Println("write:", err)
						return
					}
				case <-ctx.Done(): // Stop if the context is canceled
					return
				}
			}
		}()

		wg.Wait()
		fmt.Println("we are closed")
	})
}
