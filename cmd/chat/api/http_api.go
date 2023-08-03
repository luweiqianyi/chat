package api

import "chat/pkg/http/common"

const (
	CreateChatRoomPath  = "/chatRoom/create"
	DestroyChatRoomPath = "/chatRoom/destroy"
)

type CreateRoomParam struct {
	CreatorID string `form:"creatorID"`
}

type CreateRoomResponse struct {
	common.ResponseHeader
	RoomID string `json:"roomID"`
}

type DestroyRoomParam struct {
	RoomID string `form:"roomID"`
}
