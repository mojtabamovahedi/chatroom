package api

type createReqBody struct {
	Creator      string `json:"creator"`
	ChatRoomName string `json:"chatroomName"`
}

type joinReqBody struct {
	Name       string `json:"name"`
	ChatroomID string `json:"chatroomId"`
}

type Message struct {
	Msg  string `json:"message"`
	Name string `json:"name"`
}

type createRespBody struct {
	Message    string `json:"message"`
	ChatroomID string `json:"chatroomId"`
	UserId     string `json:"userId"`
}

type joinRespBody struct {
	Message string `json:"message"`
	UserId  string `json:"userId"`
	Name    string `json:"name"`
}
