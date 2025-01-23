package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Chat struct {
	baseUrl string
}

func NewChat(url string) Chat {
	return Chat{
		baseUrl: url,
	}
}

func (c *Chat) CreateChatroom(user, chat string) (*createRespBody, error) {
	reqBody := createReqBody{
		Creator:      user,
		ChatRoomName: chat,
	}
	reqJson, err := json.Marshal(reqBody)
	_ = err

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/%s", c.baseUrl, "api/v1/chat/create"), bytes.NewReader(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/reqJson")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var respBody *createRespBody
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	if len(respBody.Message) != 0 {
		return nil, errors.New(respBody.Message)
	}

	return respBody, nil
}

func (c *Chat) JoinChatroom(user, chatID string) (*joinRespBody, error) {
	reqBody := joinReqBody{
		Name:       user,
		ChatroomID: chatID,
	}
	reqJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/%s", c.baseUrl, "api/v1/chat/join"), bytes.NewReader(reqJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/reqJson")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respBody *joinRespBody

	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	if len(respBody.Message) != 0 {
		return nil, errors.New(respBody.Message)
	}

	return respBody, nil
}
