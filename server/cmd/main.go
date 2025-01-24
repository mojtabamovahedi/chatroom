package main

import (
	"flag"
	"github.com/mojtabamovahedi/chatroom/server/api/handler/http"
	"github.com/mojtabamovahedi/chatroom/server/app"
	"github.com/mojtabamovahedi/chatroom/server/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Flag for config file
var configPath = flag.String("config", "config.json", "Path to config file")

func main() {

	flag.Parse()

	// Read the configuration file
	cfg := config.MustReadConfig(*configPath)

	// Initialize the application container
	appContainer := app.MustNewApp(cfg)

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Create a channel to listen for server errors
	server := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		server <- http.Run(appContainer, cfg.Server)
	}()

	// Wait for a signal or server error
	select {
	case <-quit:
		log.Println("shutting down...")
	case err := <-server:
		if err != nil {
			log.Fatalf("Server error: %s", err)
		}
	}

	// Shutdown the application container
	appContainer.Shutdown()
	log.Println("Server stopped")
}
