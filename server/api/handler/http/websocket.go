package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	chatMap "github.com/mojtabamovahedi/chatroom/server/pkg/map"
	"github.com/mojtabamovahedi/chatroom/server/pkg/map/types"
	"github.com/mojtabamovahedi/chatroom/server/pkg/nats"
	nats2 "github.com/nats-io/nats.go"
	"log"
	"sync"
)

func chatroomWebsocket(
	natsClient *nats.Nats,
	userMap *chatMap.Map[string, *types.User],
	chatroomMap *chatMap.Map[string, *types.ChatRoom],
) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		chatId := c.Params("chatId")
		userId := c.Query("userId", "")
		fmt.Println("user id: ", userId)
		if len(chatId) == 0 {
			_ = c.WriteJSON(fiber.Map{
				"message": "please enter you chatroom id",
			})
			return
		}

		if _, ok := chatroomMap.Get(chatId); !ok {
			_ = c.WriteJSON(fiber.Map{
				"message": "this chatroom doesn't exist",
			})
			return
		}

		if len(userId) == 0 {
			_ = c.WriteJSON(fiber.Map{
				"message": "please enter your user id",
			})
			return
		}

		currentUser, ok := userMap.Get(userId)
		if !ok {
			_ = c.WriteJSON(fiber.Map{
				"message": "this user doesn't exist",
			})
			return
		}
		var (
			name    = currentUser.Name
			id      = currentUser.Id
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
			// userMap.Remove("")
			_ = subscribe.Unsubscribe()
			close(ch)
			log.Println("connection closed")
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
					if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
						log.Println("read:", err)
					}
					break
				}
				message := Message{
					Msg:  string(msg),
					Name: name,
					Id:   id,
				}
				arr, _ := json.Marshal(message)
				err = natsClient.Publish(subject, arr)
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
					var v2 Message
					_ = json.Unmarshal(v, &v2)
					if v2.Id == id {
						continue
					}
					err := c.WriteJSON(fiber.Map{
						"message": v2.Msg,
						"name":    v2.Name,
					})
					if err != nil {
						log.Println("write:", err)
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()

		wg.Wait()
	})
}

type Message struct {
	Msg  string
	Name string
	Id   string
}
