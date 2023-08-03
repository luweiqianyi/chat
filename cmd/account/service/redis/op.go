package redis

import (
	"chat/cmd/account/global"
	"time"
)

const (
	// 命名规则 “域名:业务名称:服务名称”
	servicePrefix = "www.company.com:chat:account:"
)

func CommonKey(key string) string {
	return servicePrefix + key
}

func QueryToken(accountName string) (bool, string) {
	token, err := global.RedisClient().Get(CommonKey(accountName))
	if err != nil {
		return false, ""
	}
	return true, token
}

func AddToken(accountName string, token string, expire time.Duration) error {
	return global.RedisClient().Set(CommonKey(accountName), token, expire)
}

func DelToken(accountName string) error {
	return global.RedisClient().Del(CommonKey(accountName))
}
