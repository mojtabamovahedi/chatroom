package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mojtabamovahedi/chatroom/server/app"
	"github.com/mojtabamovahedi/chatroom/server/config"
)

func Run(appContainer *app.App, cfg config.ServerConfig) error {
	router := fiber.New()

	router.Use(recover2.New())
	router.Use(logger.New())
	router.Use(rateLimiter())

	registerAPI(appContainer, router)

	return router.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
}

func registerAPI(appContainer *app.App, router fiber.Router) {
	chatGroup := router.Group("/api/v1/chat")

	chatGroup.Post("/create", CreateChatRoom(appContainer.MapUser(), appContainer.MapChatroom()))
	chatGroup.Post("/join", JoinChatRoom(appContainer.MapUser(), appContainer.MapChatroom()))

	// web socket
	router.Use(upgradedWebSocket())
	router.Get("chatroom/:chatId", chatroomWebsocket(appContainer.Nats(), appContainer.MapUser(), appContainer.MapChatroom()))
}
