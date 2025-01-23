package main

import (
	"bufio"
	"fmt"
	"github.com/fasthttp/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	serverUrl := "ws://localhost:8080/chatroom/12"

	conn, _, err := websocket.DefaultDialer.Dial(serverUrl, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	fmt.Println("Connected to WebSocket server.")
	fmt.Println("Type messages to send, or type 'exit' to quit.")

	ch := make(chan struct{})

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("\nError reading message: %v\n", err)
				break
			}
			fmt.Printf("Server: %s", message)
		}
		ch <- struct{}{}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				fmt.Println("\nInput closed.")
				break
			}
			text := scanner.Text()
			if text == "exit" {
				break
			}
			err := conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				fmt.Printf("Error sending message: %v\n", err)
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
