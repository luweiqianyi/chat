package api

import "chat/pkg/http/common"

type RegisterParam struct {
	AccountName string `form:"accountName"`
	Password    string `form:"password"`
}

type UnRegisterParam struct {
	AccountName string `form:"accountName"`
}

type LoginParam struct {
	AccountName string `form:"accountName"`
	Password    string `form:"password"` // ok
}

type LogoutParam struct {
	AccountName string `form:"accountName"`
}

type LoginResponse struct {
	common.ResponseHeader
	Token string `json:"token"`
}
