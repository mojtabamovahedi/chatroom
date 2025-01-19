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

	registerAPI(appContainer, router)

	return router.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
}

func registerAPI(appContainer *app.App, router fiber.Router) {
	router.Use(UpgradedWebSocket())
	router.Get("ws/:chatId", chatroomWebsocket(appContainer.Nats()))
}
