package http

import (
	"github.com/gofiber/fiber/v2"
	chatMap "github.com/mojtabamovahedi/chatroom/server/pkg/map"
	"github.com/mojtabamovahedi/chatroom/server/pkg/map/types"
	"github.com/mojtabamovahedi/chatroom/server/pkg/nanoId"
)

// CreateChatRoom handles the creation of a new chat room
func CreateChatRoom(userMap *chatMap.Map[string, *types.User], chatroomMap *chatMap.Map[string, *types.ChatRoom]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			reqBody createChatroomReq
			err     error
		)

		// Parse the request body
		if err = c.BodyParser(&reqBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errorBodyResponse{Message: "bad request"})
		}

		// Validate the request body
		if !reqBody.validate() {
			return c.Status(fiber.StatusBadRequest).JSON(errorBodyResponse{Message: "wrong request body"})
		}

		// Generate a unique user ID
		userId, err := nanoId.GenerateId()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errorBodyResponse{Message: "can not create UserID for user"})
		}

		// Generate a unique chat room ID
		chatroomId, err := nanoId.GenerateId()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errorBodyResponse{Message: "can not create UserID for chatroom"})
		}

		// Create a new user and add to the user map
		user := &types.User{
			Id:   userId,
			Name: reqBody.Creator,
			Role: types.ADMIN,
		}
		userMap.Set(userId, user)

		// Create a new chat room and add to the chat room map
		chatroom := types.NewChatRoom(chatroomId, reqBody.ChatroomName, user)
		chatroomMap.Set(chatroomId, chatroom)

		// Return the chat room ID and user ID
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"chatroomId": chatroomId,
			"userID":     userId,
		})
	}
}

// JoinChatRoom handles the joining of an existing chat room
func JoinChatRoom(userMap *chatMap.Map[string, *types.User], chatroomMap *chatMap.Map[string, *types.ChatRoom]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			reqBody joinChatroomReq
			err     error
		)

		// Parse the request body
		if err = c.BodyParser(&reqBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errorBodyResponse{Message: "bad request"})
		}

		// Validate the request body
		if !reqBody.validate() {
			return c.Status(fiber.StatusBadRequest).JSON(errorBodyResponse{Message: "wrong request body"})
		}

		// Retrieve the chat room from the map
		_, ok := chatroomMap.Get(reqBody.ChatRoomId)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(errorBodyResponse{Message: "chat room not found"})
		}

		// Generate a unique user ID
		userId, err := nanoId.GenerateId()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errorBodyResponse{Message: "can not create UserID for user"})
		}

		// Create a new user and add to the user map
		user := &types.User{
			Id:   userId,
			Name: reqBody.Name,
			Role: types.USER,
		}
		userMap.Set(userId, user)

		// Return the user ID and name
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"userId": userId,
			"Name":   user.Name,
		})
	}
}

// createChatroomReq represents the request body for creating a chat room
type createChatroomReq struct {
	Creator      string `json:"creator"`
	ChatroomName string `json:"chatroomName"`
}

// validate represent if one of field in createChatroomReq is empty return false
func (c *createChatroomReq) validate() bool {
	if len(c.ChatroomName) == 0 || len(c.Creator) == 0 {
		return false
	}
	return true
}

// joinChatroomReq represents the request body for joining a chat room
type joinChatroomReq struct {
	ChatRoomId string `json:"chatRoomId"`
	Name       string `json:"name"`
}

// validate represent if one of field in joinChatroomReq is empty return false
func (j joinChatroomReq) validate() bool {
	if len(j.Name) == 0 || len(j.ChatRoomId) == 0 {
		return false
	}
	return true
}

// errorBodyResponse represents the error response body
type errorBodyResponse struct {
	Message string `json:"message"`
}
