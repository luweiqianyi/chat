package service

import (
	"chat/cmd/account/routing/api"
	"chat/cmd/account/service/dao"
	"chat/cmd/account/service/redis"
	"chat/pkg/log"
	"chat/pkg/util"
	"time"
)

const (
	Success             = 0
	Failed              = 1
	AccountNotExist     = 2
	AccountAlreadyExist = 3

	TokenSecretKey = "chat"
)

func Register(param api.RegisterParam) int {
	exist := dao.FindAccount(param.AccountName)
	if exist {
		return AccountAlreadyExist
	}

	success := dao.AddAccount(param.AccountName, param.Password)
	if success {
		return Success
	} else {
		return Failed
	}
}

func UnRegister(param api.UnRegisterParam) int {
	exist := dao.FindAccount(param.AccountName)
	if !exist {
		return AccountNotExist
	}

	dao.DelAccount(param.AccountName)
	err := redis.DelToken(param.AccountName) // TODO 确保正确删除
	if err != nil {
		log.Errorf("account[%v]'s token del failed, err: %v", param.AccountName, err)
	}
	return Success
}

func Login(param api.LoginParam) (int, string) {
	exist := dao.FindAccount(param.AccountName)
	if !exist {
		return AccountNotExist, ""
	}

	exist, token := redis.QueryToken(param.AccountName)
	if exist {
		return Success, token
	}

	token, err := util.GenerateToken(param.AccountName, TokenSecretKey)
	if err != nil {
		return Failed, ""
	}
	err = redis.AddToken(param.AccountName, token, time.Hour*24*15)
	if err != nil {
		log.Errorf("login, account[%v] add token failed, err: %v", param.AccountName, err)
		return Failed, ""
	}

	return Success, token
}

func Logout(param api.LogoutParam) int {
	exist := dao.FindAccount(param.AccountName)
	if !exist {
		return AccountNotExist
	}

	err := redis.DelToken(param.AccountName)
	if err != nil {
		log.Errorf("account[%v] logout failed, err: %v", param.AccountName, err)
		return Failed
	}
	return Success
}

func ValidateAccount(accountName, accessToken string) bool {
	exist, token := redis.QueryToken(accountName)
	if !exist {
		return false
	}

	return token == accessToken
}
