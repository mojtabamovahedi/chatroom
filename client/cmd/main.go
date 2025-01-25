package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/mojtabamovahedi/chatroom/client/api"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	baseUrl = "localhost:8080"
	// flags
	uFlag   = flag.String("u", "", "Your display name in chatroom.")
	cFlag   = flag.String("c", "", "Name of chatroom you want to create.")
	iFlag   = flag.String("i", "", "Chatroom ID you want to join.")
	urlFlag = flag.String("url", "", fmt.Sprintf("URL of your chatroom. Default is %s", baseUrl))
)

func main() {
	flag.Parse()

	if len(*urlFlag) != 0 {
		baseUrl = *urlFlag
	}

	wsUrl := fmt.Sprintf("ws://%s/chatroom", baseUrl)

	// validate flag to get what user need it
	cID, uID, err := validateFlags(*uFlag, *iFlag, *cFlag)
	if err != nil {
		log.Fatal(err)
	}

	connUrl := fmt.Sprintf("%s/%s?userId=%s", wsUrl, cID, uID)

	// connect to ws server
	conn, _, err := websocket.DefaultDialer.Dial(connUrl, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	fmt.Println("Connected to WebSocket server.")
	fmt.Println("Type messages to send, or type '#exit' to quit.")
	fmt.Println("Share your code to join your friends:")
	fmt.Printf("###        CODE: %s        ###\n", cID)

	ch := make(chan struct{})

	// create goroutine for read data from server
	go func() {
		var (
			msg  message
			data []byte
			rErr error
		)
		for {
			_, data, rErr = conn.ReadMessage()
			if rErr != nil {
				fmt.Printf("\nError reading message: %v\n", rErr)
				break
			}

			rErr = json.Unmarshal(data, &msg)
			if rErr != nil {
				continue
			}

			fmt.Printf("#> %s: %s\n", msg.Name, msg.Msg)
		}
		ch <- struct{}{}
	}()

	// create goroutine for write data on server
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
			// check for exit command
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

	// Create a channel to listen for OS signals
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

	chatAPI := api.NewChat(baseUrl)

	if len(c) != 0 {
		cBody, err := chatAPI.CreateChatroom(u, c)
		if err != nil {
			return "", "", err
		}
		cID = cBody.ChatroomID
		uID = cBody.UserId
		return cID, uID, nil
	}

	if len(i) != 0 {
		jBody, err := chatAPI.JoinChatroom(u, i)
		if err != nil {
			return "", "", err
		}
		cID = i
		uID = jBody.UserId
	}

	return cID, uID, nil

}

// message type for communication with server
type message struct {
	Msg  string `json:"message"`
	Name string `json:"name"`
}
