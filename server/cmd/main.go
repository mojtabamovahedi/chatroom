package main

import (
	"github.com/mojtabamovahedi/chatroom/server/api/handler/http"
	"github.com/mojtabamovahedi/chatroom/server/app"
	"github.com/mojtabamovahedi/chatroom/server/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg := config.MustReadConfig("config.json")
	appContainer := app.MustNewApp(cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	server := make(chan error, 1)

	go func() {
		server <- http.Run(appContainer, cfg.Server)
	}()

	select {
	case <-quit:
		log.Println("shutting down...")
	case err := <-server:
		if err != nil {
			log.Fatalf("Server error: %s", err)
		}
	}

	appContainer.Shutdown()
	log.Println("Server stopped")
}
