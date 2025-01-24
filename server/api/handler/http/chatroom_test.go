package http

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	chatMap "github.com/mojtabamovahedi/chatroom/server/pkg/map"
	"github.com/mojtabamovahedi/chatroom/server/pkg/map/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateChatRoom(t *testing.T) {
	app := fiber.New()

	userMap := chatMap.NewMap[string, *types.User]()
	chatroomMap := chatMap.NewMap[string, *types.ChatRoom]()

	app.Post("/create", CreateChatRoom(userMap, chatroomMap))

	tests := []struct {
		description  string
		requestBody  createChatroomReq
		expectedCode int
	}{
		{
			description: "Valid request",
			requestBody: createChatroomReq{
				Creator:      "John",
				ChatroomName: "General",
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			description: "Invalid request - missing fields",
			requestBody: createChatroomReq{
				Creator:      "",
				ChatroomName: "",
			},
			expectedCode: fiber.StatusBadRequest,
		},
	}

	for _, test := range tests {
		body, _ := json.Marshal(test.requestBody)
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		utils.AssertEqual(t, test.expectedCode, resp.StatusCode, test.description)
	}
}

func TestJoinChatRoom(t *testing.T) {
	app := fiber.New()

	userMap := chatMap.NewMap[string, *types.User]()
	chatroomMap := chatMap.NewMap[string, *types.ChatRoom]()

	chatroomId := "existingChatroomId"
	chatroom := types.NewChatRoom(chatroomId, "General", &types.User{Id: "adminId", Name: "Admin", Role: types.ADMIN})
	chatroomMap.Set(chatroomId, chatroom)

	app.Post("/join", JoinChatRoom(userMap, chatroomMap))

	tests := []struct {
		description  string
		requestBody  joinChatroomReq
		expectedCode int
	}{
		{
			description: "Valid request",
			requestBody: joinChatroomReq{
				Name:       "Jane",
				ChatRoomId: chatroomId,
			},
			expectedCode: fiber.StatusAccepted,
		},
		{
			description: "Invalid request - missing fields",
			requestBody: joinChatroomReq{
				Name:       "",
				ChatRoomId: "",
			},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			description: "Invalid request - chatroom not found",
			requestBody: joinChatroomReq{
				Name:       "Jane",
				ChatRoomId: "nonExistingChatroomId",
			},
			expectedCode: fiber.StatusBadRequest,
		},
	}

	for _, test := range tests {
		body, _ := json.Marshal(test.requestBody)
		req := httptest.NewRequest(http.MethodPost, "/join", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		utils.AssertEqual(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
