package http

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

// upgradedWebSocket checks if the request is a WebSocket upgrade request
func upgradedWebSocket() fiber.Handler {
    return func(c *fiber.Ctx) error {
        if websocket.IsWebSocketUpgrade(c) {
            return c.Next()
        }
        return fiber.ErrUpgradeRequired
    }
}

// rateLimiter sets up a rate limiter middleware
func rateLimiter() fiber.Handler {
    return limiter.New(limiter.Config{
        Next: func(c *fiber.Ctx) bool {
            return c.IP() == "127.0.0.1"
        },
        Max:        10, // Maximum number of requests
        Expiration: 1 * time.Minute, // Time window for rate limiting
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(fiber.StatusTooManyRequests).JSON(errorBodyResponse{Message: "Take a break!"})
        },
    })
}