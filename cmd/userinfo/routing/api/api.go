package api

import (
	"chat/cmd/userinfo/auth"
	"chat/pkg/http/common"
	"mime/multipart"
)

const (
	Male   = 1
	Female = 0
)

type UserInfoParam struct {
	NickName string `form:"nickname"`
	Birthday string `form:"birthday"`
	Gender   int    `form:"gender"`
}

type AddUserInfoParam struct {
	auth.AccountParam
	UserInfoParam
}

type DelUserInfoParam struct {
	auth.AccountParam
}

type QueryUserInfoParam struct {
	auth.AccountParam
}

type QueryUserInfoResponse struct {
	common.ResponseHeader
	UserInfoResp `json:"userInfo"`
}

type UserInfoResp struct {
	NickName string `json:"nickname"`
	Birthday string `json:"birthday"`
	Gender   int    `json:"gender"`
}

func NewQueryUserInfoResponse(userInfo UserInfoResp) QueryUserInfoResponse {
	return QueryUserInfoResponse{
		ResponseHeader: common.NewSuccessResponse(),
		UserInfoResp:   userInfo,
	}
}

// UpdateUserInfoParam update all information of user
type UpdateUserInfoParam struct {
	auth.AccountParam
	UserInfoParam
}

// below: update one information of user

type UpdateNickNameParam struct {
	auth.AccountParam
	NickName string `form:"nickname"`
}

type UpdateBirthdayParam struct {
	auth.AccountParam
	Birthday string `form:"birthday"`
}

type UpdateAvatarParam struct {
	auth.AccountParam
	AvatarFile *multipart.FileHeader `form:"avatarName"`
}

type UpdateAvatarResp struct {
	common.ResponseHeader
	AvatarRemotePath string `json:"avatarRemotePath"`
}

func NewSuccessUpdateAvatarResp(path string) UpdateAvatarResp {
	return UpdateAvatarResp{
		ResponseHeader:   common.NewSuccessResponse(),
		AvatarRemotePath: path,
	}
}
