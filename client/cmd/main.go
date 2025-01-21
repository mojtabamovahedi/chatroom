package main

import (
	"bufio"
	"fmt"
	"github.com/fasthttp/websocket"
	"log"
	"os"
	"os/signal"
	"sync"
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

	done := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("\nError reading message: %v\n", err)
				done <- struct{}{}
			}
			fmt.Printf("Server: %s\n", message)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				fmt.Println("\nInput closed.")
				close(done)
				return
			}
			text := scanner.Text()
			if text == "exit" {
				close(done)
				return
			}
			err := conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				fmt.Printf("Error sending message: %v\n", err)
				close(done)
				return
			}
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		done <- struct{}{}
	}()

	<-done
	close(done)
	fmt.Println("Shutting down...")
	conn.Close()
	wg.Wait()
}
