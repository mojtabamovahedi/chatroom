package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/mojtabamovahedi/chatroom/client/api"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	baseUrl = "localhost:8080"
	wsUrl   = fmt.Sprintf("ws://%s/chatroom", baseUrl)
)

func main() {
	var (
		uFlag = flag.String("u", "", "Your display name in chatroom.")
		cFlag = flag.String("c", "", "Name of chatroom you want to create.")
		iFlag = flag.String("i", "", "Chatroom ID you want to join.")
	)
	flag.Parse()

	cID, uID, err := validateFlags(*uFlag, *iFlag, *cFlag)
	if err != nil {
		log.Fatal(err)
	}

	connUrl := fmt.Sprintf("%s/%s?userId=%s", wsUrl, cID, uID)

	conn, _, err := websocket.DefaultDialer.Dial(connUrl, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	fmt.Println("Connected to WebSocket server.")
	fmt.Println("Type messages to send, or type '#exit' to quit.")
	fmt.Println("Share your code to join your friends:")
	fmt.Printf("###        CODE: %s        ###\n", cID)

	ch := make(chan struct{})

	// read data from server
	go func() {
		var (
			msg  message
			data []byte
			rErr error
		)
		for {
			_, data, rErr = conn.ReadMessage()
			if rErr != nil {
				if rErr == io.EOF {
				}
				fmt.Printf("\nError reading message: %v\n", rErr)
				break
			}

			if rErr = json.Unmarshal(data, &msg); rErr != nil {
				continue
			}

			fmt.Printf("#> %s: %s\n", msg.Name, msg.Msg)
		}
		ch <- struct{}{}
	}()

	// write data on server
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		var (
			wErr error
			text string
		)
		for {
			if !scanner.Scan() {
				fmt.Println("\nInput closed.")
				break
			}
			text = scanner.Text()
			if text == "#exit" {
				break
			}
			wErr = conn.WriteMessage(websocket.TextMessage, []byte(text))
			if wErr != nil {
				fmt.Printf("Error sending message: %v\n", wErr)
				break
			}
		}
		ch <- struct{}{}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		ch <- struct{}{}
	}()

	<-ch
	fmt.Println("Shutting down...")
	_ = conn.Close()
	close(ch)
}

func validateFlags(u, i, c string) (cID, uID string, err error) {
	if len(u) == 0 {
		return "", "", errors.New("you forgot about your display name (use -u)")
	}

	if len(c) == 0 && len(i) == 0 {
		return "", "", errors.New("you forgot about your chat room ID (use -i)")
	}

	if len(c) != 0 && len(i) != 0 {
		return "", "", errors.New("you can not create and join chatroom at the same time")
	}

	chat := api.NewChat(baseUrl)

	if len(c) != 0 {
		cBody, err := chat.CreateChatroom(u, c)
		if err != nil {
			return "", "", err
		}
		cID = cBody.ChatroomID
		uID = cBody.UserId
		return cID, uID, nil
	}

	if len(i) != 0 {
		jBody, err := chat.JoinChatroom(u, i)
		if err != nil {
			return "", "", err
		}
		cID = i
		uID = jBody.UserId
	}

	return cID, uID, nil

}

type message struct {
	Msg  string `json:"message"`
	Name string `json:"name"`
}
