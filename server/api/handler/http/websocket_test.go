package http

import (
	"github.com/gofiber/fiber/v2"
	chatMap "github.com/mojtabamovahedi/chatroom/server/pkg/map"
	"github.com/mojtabamovahedi/chatroom/server/pkg/map/types"
	"github.com/mojtabamovahedi/chatroom/server/pkg/nats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestChatroomWebsocket(t *testing.T) {
	app := fiber.New()

	userMap := chatMap.NewMap[string, *types.User]()
	chatroomMap := chatMap.NewMap[string, *types.ChatRoom]()

	natsClient := &nats.Nats{}

	chatroomId := "existingChatroomId"
	chatroom := types.NewChatRoom(chatroomId, "General", &types.User{Id: "adminId", Name: "Admin", Role: types.ADMIN})
	chatroomMap.Set(chatroomId, chatroom)

	userId := "userId"
	user := &types.User{Id: userId, Name: "John", Role: types.USER}
	userMap.Set(userId, user)

	app.Get("/ws/:chatId", chatroomWebsocket(natsClient, userMap, chatroomMap))

	t.Run("Valid WebSocket connection", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ws/"+chatroomId+"?userId="+userId, nil)
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Sec-WebSocket-Version", "13")
		req.Header.Set("Sec-WebSocket-Key", "x3JJHMbDL1EzLkh9GBhXDw==")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 101, resp.StatusCode)
	})

	t.Run("Invalid WebSocket connection - missing chatId", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ws/?userId="+userId, nil)
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Sec-WebSocket-Version", "13")
		req.Header.Set("Sec-WebSocket-Key", "x3JJHMbDL1EzLkh9GBhXDw==")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})

	t.Run("Invalid WebSocket connection - missing userId", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ws/"+chatroomId, nil)
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Sec-WebSocket-Version", "13")
		req.Header.Set("Sec-WebSocket-Key", "x3JJHMbDL1EzLkh9GBhXDw==")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 101, resp.StatusCode)
	})

	t.Run("Invalid WebSocket connection - non-existing chatroom", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ws/nonExistingChatroomId?userId="+userId, nil)
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Sec-WebSocket-Version", "13")
		req.Header.Set("Sec-WebSocket-Key", "x3JJHMbDL1EzLkh9GBhXDw==")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 101, resp.StatusCode)
	})

	t.Run("Invalid WebSocket connection - non-existing user", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ws/"+chatroomId+"?userId=nonExistingUserId", nil)
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Sec-WebSocket-Version", "13")
		req.Header.Set("Sec-WebSocket-Key", "x3JJHMbDL1EzLkh9GBhXDw==")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 101, resp.StatusCode)
	})
}
