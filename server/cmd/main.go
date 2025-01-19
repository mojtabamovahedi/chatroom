package main

import (
	"github.com/mojtabamovahedi/chatroom/server/app"
	"github.com/mojtabamovahedi/chatroom/server/config"
)

func main() {

	cfg := config.MustReadConfig("config.json")
	chatroomApp := app.MustNewApp(cfg)

	_ = chatroomApp
}
