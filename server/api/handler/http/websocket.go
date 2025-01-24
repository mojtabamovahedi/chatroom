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
		// Retrieve chatId from URL parameters
        chatId := c.Params("chatId")
        // Retrieve userId from query parameters
        userId := c.Query("userId", "")

		// Check if chatId is provided
		if len(chatId) == 0 {
			_ = c.WriteJSON(fiber.Map{
				"Message": "please enter you chatroom UserID",
			})
			c.Close()
			return
		}

		// Retrieve the chat room from the map
		room, ok := chatroomMap.Get(chatId)
		if !ok {
			_ = c.WriteJSON(fiber.Map{
				"message": "this chatroom doesn't exist",
			})
			c.Close()
			return
		}

		// Check if userId is provided
		if len(userId) == 0 {
			_ = c.WriteJSON(fiber.Map{
				"message": "please enter your user UserID",
			})
			c.Close()
			return
		}

		currentUser, ok := userMap.Get(userId)
		if !ok {
			_ = c.WriteJSON(fiber.Map{
				"message": "this user doesn't exist",
			})
			c.Close()
			return
		}

		if currentUser.Role != types.ADMIN {
			room.AddChatter(currentUser)
		}
		var (
			uName   = currentUser.Name
			uID     = currentUser.Id
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
			c.Close()
			return
		}

		// function after user disconnected
		cleanup := func() {
			userMap.Remove(userId)
			isEmpty := room.RemoveChatter(currentUser)
			if isEmpty {
				chatroomMap.Remove(room.Id)
			}
			_ = subscribe.Unsubscribe()

			msg := Message{
				Msg:    fmt.Sprintf("'%s' left the chatroom", uName),
				Name:   Owner,
				UserID: OwnerID,
			}
			data, _ := json.Marshal(msg)
			_ = natsClient.Publish(subject, data)
			close(ch)
			cancel()
		}

		// welcome message
		joinMsg := Message{
			Msg:    fmt.Sprintf("'%s' join the chatroom", uName),
			Name:   Owner,
			UserID: OwnerID,
		}
		joinData, _ := json.Marshal(joinMsg)
		_ = natsClient.Publish(subject, joinData)

		wg.Add(2)

		// read message from client
		go func() {
			defer func() {
				wg.Done()
				cleanup()
			}()
			var (
				msg  Message
				err  error
				read []byte
				data []byte
			)
			for {
				_, read, err = c.ReadMessage()
				if err != nil {
					if !websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
						log.Println("read:", err)
					}
					return
				}

				if isCommand(string(read)) {
					var sb strings.Builder
					sb.WriteString("users in room:")
					for _, chatter := range room.GetChatters() {
						sb.WriteString(fmt.Sprintf("\n# '%s'", chatter.Name))
					}

					_ = c.WriteJSON(Message{
						Msg:    sb.String(),
						Name:   Owner,
						UserID: OwnerID,
					})
					continue
				}

				msg = Message{
					Msg:    strings.TrimSpace(string(read)),
					Name:   uName,
					UserID: uID,
				}
				data, err = json.Marshal(msg)
				if err != nil {
					continue
				}
				err = natsClient.Publish(subject, data)
				if err != nil {
					log.Println("error in publish Message:", err)
				}
			}
		}()

		// write messages to client
		go func() {
			defer func() {
				wg.Done()
				cleanup()
			}()
			var (
				msg  Message
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
					if msg.UserID == uID {
						continue
					}
					err = c.WriteJSON(fiber.Map{
						"message": msg.Msg,
						"Name":    msg.Name,
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

type Message struct {
	Msg    string `json:"message"`
	Name   string `json:"name"`
	UserID string `json:"id"`
}

func isCommand(command string) bool {
	command = strings.TrimSpace(command)
	if command == "#users" {
		return true
	}
	return false
}
