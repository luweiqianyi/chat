// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// redis-group CRUD module: crud operations to manipulate data with redis.
// group is a concept in specific business scene

package groups

import (
	"chat/cmd/chat/api"
	"chat/cmd/chat/global"
	"chat/pkg/redis"
	"encoding/json"
	"errors"
)

type GroupManager struct {
	redisCli *redis.Client
}

func NewGroupManager() *GroupManager {
	return &GroupManager{
		redisCli: global.ClientFactoryInstance().RedisClient(),
	}
}

func (groupManager *GroupManager) CreateGroup(param CreateGroupParam) error {
	if groupManager.redisCli == nil {
		return errors.New("redis cli required")
	}

	key := api.GenRedisKey(api.FormatGroup, api.PrefixGroupKey, param.GroupID)
	data, _ := json.Marshal(param)
	return groupManager.redisCli.HSet(key, api.FieldGroupBaseInfo, string(data))
}

func (groupManager *GroupManager) DeleteGroup(param DeleteGroupParam) error {
	if groupManager.redisCli == nil {
		return errors.New("redis cli required")
	}

	key := api.GenRedisKey(api.FormatGroup, api.PrefixGroupKey, param.GroupID)
	return groupManager.redisCli.HDelAllFields(key)
}

func (groupManager *GroupManager) Add2Group(param Add2GroupParam) error {
	if groupManager.redisCli == nil {
		return errors.New("redis cli required")
	}

	key := api.GenRedisKey(api.FormatGroup, api.PrefixGroupKey, param.GroupID)
	return groupManager.redisCli.HSet(key, param.AddedUserID, "")
}

func (groupManager *GroupManager) RemoveFromGroup(param RemoveFromGroupParam) error {
	if groupManager.redisCli == nil {
		return errors.New("redis cli required")
	}

	key := api.GenRedisKey(api.FormatGroup, api.PrefixGroupKey, param.GroupID)
	return groupManager.redisCli.HDel(key, param.RemovedUserID)
}
