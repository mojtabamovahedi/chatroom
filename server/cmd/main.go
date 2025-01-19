package main

import (
	"github.com/mojtabamovahedi/chatroom/server/api/handler/http"
	"github.com/mojtabamovahedi/chatroom/server/app"
	"github.com/mojtabamovahedi/chatroom/server/config"
	"log"
)

func main() {

	cfg := config.MustReadConfig("config.json")
	appContainer := app.MustNewApp(cfg)

	log.Fatal(http.Run(appContainer, cfg.Server))
}
