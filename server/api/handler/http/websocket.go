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
	"strings"
	"sync"
)

const (
	Owner   = "SNAPP CHAT"
	OwnerID = "-1"
)

func chatroomWebsocket(
	natsClient *nats.Nats,
	userMap *chatMap.Map[string, *types.User],
	chatroomMap *chatMap.Map[string, *types.ChatRoom],
) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		chatId := c.Params("chatId")
		userId := c.Query("userId", "")
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

		room, ok := chatroomMap.Get(chatId)
		if !ok {
			_ = c.WriteJSON(fiber.Map{
				"message": "this chatroom doesn't exist",
			})
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
			userMap.Remove(userId)
			isEmpty := room.RemoveChatter(currentUser)
			if isEmpty {
				chatroomMap.Remove(room.Id)
			}
			_ = subscribe.Unsubscribe()

			msg := message{
				msg:  fmt.Sprintf("'%s' left the chatroom", name),
				name: Owner,
				id:   OwnerID,
			}
			data, _ := json.Marshal(msg)
			_ = natsClient.Publish(subject, data)
			close(ch)
		}()

		joinMsg := message{
			msg:  fmt.Sprintf("'%s' join the chatroom", name),
			name: Owner,
			id:   OwnerID,
		}
		joinData, _ := json.Marshal(joinMsg)
		_ = natsClient.Publish(subject, joinData)

		wg.Add(2)

		go func() {
			defer func() {
				wg.Done()
				cancel()
			}()
			var (
				msg  message
				err  error
				read []byte
				data []byte
			)
			for {
				_, read, err = c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					break
				}
				msg = message{
					msg:  strings.TrimSpace(string(read)),
					name: name,
					id:   id,
				}
				data, err = json.Marshal(msg)
				if err != nil {
					continue
				}
				err = natsClient.Publish(subject, data)
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
			var (
				msg  message
				err  error
				data []byte
			)
			for {
				select {
				case data, ok = <-ch:
					if !ok {
						break
					}

					err = json.Unmarshal(data, &msg)
					if msg.id == id {
						continue
					}
					err = c.WriteJSON(fiber.Map{
						"message": msg.msg,
						"name":    msg.name,
					})
					if err != nil {
						log.Println("write:", err)
						break
					}
				case <-ctx.Done():
					break
				}
			}
		}()

		wg.Wait()
	})
}

type message struct {
	msg  string
	name string
	id   string
}
