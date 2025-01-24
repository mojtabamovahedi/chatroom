package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mojtabamovahedi/chatroom/server/app"
	"github.com/mojtabamovahedi/chatroom/server/config"
)


// Run initializes and starts the Fiber application
func Run(appContainer *app.App, cfg config.ServerConfig) error {
    router := fiber.New()

    // Use middleware for recovering from panics and logging requests
    router.Use(recover2.New())
    router.Use(logger.New())
    router.Use(rateLimiter())

    // Register API routes
    registerAPI(appContainer, router)

    // Start the server
    return router.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
}


// registerAPI registers the API routes
func registerAPI(appContainer *app.App, router fiber.Router) {
    chatGroup := router.Group("/api/v1/chat")

    // REST APIs for creating and joining chatrooms
    chatGroup.Post("/create", CreateChatRoom(appContainer.MapUser(), appContainer.MapChatroom()))
    chatGroup.Post("/join", JoinChatRoom(appContainer.MapUser(), appContainer.MapChatroom()))

    // WebSocket endpoint
    router.Use(upgradedWebSocket())
    router.Get("chatroom/:chatId", chatroomWebsocket(appContainer.Nats(), appContainer.MapUser(), appContainer.MapChatroom()))
}
