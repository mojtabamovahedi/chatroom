package http

import (
	"github.com/gofiber/fiber/v2"
	chatMap "github.com/mojtabamovahedi/chatroom/server/pkg/map"
	"github.com/mojtabamovahedi/chatroom/server/pkg/map/types"
	"github.com/mojtabamovahedi/chatroom/server/pkg/nanoId"
)

func CreateChatRoom(userMap *chatMap.Map[string, *types.User], chatroomMap *chatMap.Map[string, *types.ChatRoom]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			reqBody createChatroomReq
			err     error
		)
		if err = c.BodyParser(&reqBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errorBodyResponse{Message: "bad request"})
		}

		if !reqBody.validate() {
			return c.Status(fiber.StatusBadRequest).JSON(errorBodyResponse{Message: "wrong request body"})
		}

		userId, err := nanoId.GenerateId()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errorBodyResponse{Message: "can not create id for user"})
		}

		chatroomId, err := nanoId.GenerateId()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errorBodyResponse{Message: "can not create id for chatroom"})
		}

		user := &types.User{
			Id:   userId,
			Name: reqBody.Creator,
			Role: types.ADMIN,
		}
		userMap.Set(userId, user)

		chatroom := types.NewChatRoom(chatroomId, reqBody.ChatroomName, user)
		chatroomMap.Set(chatroomId, chatroom)

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"chatroomId": chatroomId,
			"userID":     userId,
		})
	}
}

func JoinChatRoom(userMap *chatMap.Map[string, *types.User], chatroomMap *chatMap.Map[string, *types.ChatRoom]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			reqBody joinChatroomReq
			err     error
		)
		if err = c.BodyParser(&reqBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errorBodyResponse{Message: "bad request"})
		}

		if !reqBody.validate() {
			return c.Status(fiber.StatusBadRequest).JSON(errorBodyResponse{Message: "wrong request body"})
		}

		room, ok := chatroomMap.Get(reqBody.ChatRoomId)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(errorBodyResponse{Message: "chat room not found"})
		}

		_ = room

		userId, err := nanoId.GenerateId()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errorBodyResponse{Message: "can not create id for user"})
		}

		user := &types.User{
			Id:   userId,
			Name: reqBody.Name,
			Role: types.USER,
		}

		userMap.Set(userId, user)

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"userId": userId,
			"name":   user.Name,
		})

	}
}

type createChatroomReq struct {
	Creator      string `json:"creator"`
	ChatroomName string `json:"chatroomName"`
}

func (c *createChatroomReq) validate() bool {
	if len(c.ChatroomName) == 0 || len(c.Creator) == 0 {
		return false
	}
	return true
}

type joinChatroomReq struct {
	Name       string `json:"name"`
	ChatRoomId string `json:"chatroomId"`
}

func (j joinChatroomReq) validate() bool {
	if len(j.Name) == 0 || len(j.ChatRoomId) == 0 {
		return false
	}
	return true
}

type errorBodyResponse struct {
	Message string `json:"message"`
}
