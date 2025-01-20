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
			reqBody CreateChatRoomReq
			err     error
		)
		if err = c.BodyParser(&reqBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errorBodyResponse{Message: "bad request"})
		}

		if !reqBody.Validate() {
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

		chatroom := types.NewChatRoom(chatroomId, reqBody.ChatRoomName, *user)
		chatroomMap.Set(chatroomId, chatroom)

		// will add jwt token
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"chatroomId": chatroomId,
		})
	}
}

type CreateChatRoomReq struct {
	Creator      string `json:"creator"`
	ChatRoomName string `json:"chatroomName"`
}

func (r CreateChatRoomReq) Validate() bool {
	if len(r.ChatRoomName) == 0 || len(r.Creator) == 0 {
		return false
	}
	return true
}

type errorBodyResponse struct {
	Message string `json:"message"`
}

func JoinChatRoom(userMap *chatMap.Map[string, *types.User], chatroomMap *chatMap.Map[string, *types.ChatRoom]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			reqBody JoinChatRoomReq
			err     error
		)
		if err = c.BodyParser(&reqBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errorBodyResponse{Message: "bad request"})
		}

		if !reqBody.Validate() {
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

type JoinChatRoomReq struct {
	Name       string `json:"name"`
	ChatRoomId string `json:"chatroomId"`
}

func (r JoinChatRoomReq) Validate() bool {
	if len(r.Name) == 0 || len(r.ChatRoomId) == 0 {
		return false
	}
	return true
}
