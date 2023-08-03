// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// api module: request parameters and response format definition, interact with front-end

package api

import "chat/cmd/chat/util"

const (
	ServerApiVersion = "1.0.0"

	Login           = "login"
	Logout          = "logout"
	Keepalive       = "keepalive"
	CreateGroup     = "createGroup"
	DeleteGroup     = "deleteGroup"
	AddToGroup      = "addToGroup"
	RemoveFromGroup = "removeFromGroup"
	One2One         = "One2One"
	One2Group       = "One2Group"
)

const (
	Success              = 1000
	ParameterFormatError = 1001
	ServerInternalError  = 1002

	UserNotLogin            = 1003
	LoginFailed             = 1004
	GroupCreateFailed       = 1005
	GroupDestroyFailed      = 1006
	Add2GroupFailed         = 1007
	RemoveFromGroupFailed   = 1008
	Send2OneMessageFailed   = 1009
	Send2GroupMessageFailed = 1010
)

func ResponseString(code int) string {
	switch code {
	case Success:
		return "success"
	case ParameterFormatError:
		return "parameter format error"
	case ServerInternalError:
		return "server internal error"
	case LoginFailed:
		return "login: validation failed"
	case GroupCreateFailed:
		return "group create failed"
	case GroupDestroyFailed:
		return "group destroy failed"
	case Add2GroupFailed:
		return "add to group failed"
	case RemoveFromGroupFailed:
		return "remove from group failed"
	case Send2OneMessageFailed:
		return "send message to one failed"
	case Send2GroupMessageFailed:
		return "send group message failed"
	case UserNotLogin:
		return "user not login"
	}
	return "unknown error"
}

type RequestHeader struct {
	Version   string `json:"version"`
	MsgID     string `json:"msgID"`
	Method    string `json:"method"`
	Timestamp int64  `json:"timestamp"`
}

type LoginParam struct {
	RequestHeader
	AccountName string `json:"accountName"`
	AccessToken string `json:"accessToken"`
}

type LogoutParam struct {
	RequestHeader
	AccountName string `json:"accountName"`
	AccessToken string `json:"accessToken"`
}

type One2OneMessageParam struct {
	RequestHeader
	SenderID   string `json:"senderID"`
	ReceiverID string `json:"receiverID"`
	Message    string `json:"message"`
}

type One2GroupMessageParam struct {
	RequestHeader
	SenderID string `json:"senderID"`
	GroupID  string `json:"groupID"`
	Message  string `json:"message"`
}

type CreateGroupParam struct {
	RequestHeader
	SenderID string `json:"senderID"`
	GroupID  string `json:"groupID"`
}

type DeleteGroupParam struct {
	RequestHeader
	SenderID string `json:"senderID"`
	GroupID  string `json:"groupID"`
}

type Add2GroupParam struct {
	RequestHeader
	SenderID   string `json:"senderID"`
	ReceiverID string `json:"receiverID"`
	GroupID    string `json:"groupID"`
}

type RemoveFromGroupParam struct {
	RequestHeader
	SenderID   string `json:"senderID"`
	ReceiverID string `json:"receiverID"`
	GroupID    string `json:"groupID"`
}

type AppKeepAliveParam struct {
	RequestHeader
}

type ResponseHeader struct {
	Version   string `json:"version"`
	MsgID     string `json:"msgID"`
	Method    string `json:"method"`
	Timestamp int64  `json:"timestamp"`
}

type CommonResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	ResponseHeader
	CommonResponseData
}

func NewSuccessResponse(msgID string, method string) SuccessResponse {
	return SuccessResponse{
		ResponseHeader: ResponseHeader{
			Version:   ServerApiVersion,
			MsgID:     msgID,
			Method:    method,
			Timestamp: util.GenerateTimestamp(),
		},
		CommonResponseData: CommonResponseData{
			Code:    Success,
			Message: ResponseString(Success),
		},
	}
}

type ErrorResponse struct {
	ResponseHeader
	CommonResponseData
}

func NewErrorResponse(msgID string, method string, code int) ErrorResponse {
	return ErrorResponse{
		ResponseHeader: ResponseHeader{
			Version:   ServerApiVersion,
			MsgID:     msgID,
			Method:    method,
			Timestamp: util.GenerateTimestamp(),
		},
		CommonResponseData: CommonResponseData{
			Code:    code,
			Message: ResponseString(code),
		},
	}
}

func NewErrorResponseWithTip(msgID string, method string, code int, message string) ErrorResponse {
	return ErrorResponse{
		ResponseHeader: ResponseHeader{
			Version:   ServerApiVersion,
			MsgID:     msgID,
			Method:    method,
			Timestamp: util.GenerateTimestamp(),
		},
		CommonResponseData: CommonResponseData{
			Code:    code,
			Message: message,
		},
	}
}
