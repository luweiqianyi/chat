// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// room db operation: implementation of roomInterface.go, use redis as a store medium, save room data to redis

package httpBusiness

import (
	"chat/cmd/chat/api"
	"chat/cmd/chat/global"
	"chat/cmd/chat/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

func GenerateChatRoomID() string {
	return uuid.NewString()
}

type ChatRoomRedis struct {
}

// CreateRoom one person can only create one chat room, field=creatorID, value=string(json(api.RedisCreateRoomData))
func (room *ChatRoomRedis) CreateRoom(param CreateRoomParam) (CreateRoomResponse, error) {
	if global.ClientFactoryInstance().RedisClient() == nil {
		return CreateRoomResponse{}, errors.New("redis cli required")
	}

	value, _ := global.ClientFactoryInstance().RedisClient().HGet(api.RoomsRedisKey(), param.CreatorID)
	if value == nil {
		roomID := GenerateChatRoomID()
		data := api.RedisCreateRoomData{
			RoomID:          roomID,
			CreatorID:       param.CreatorID,
			CreateTimestamp: util.GenerateTimestamp(),
		}
		jsonBytes, _ := json.Marshal(data)
		err := global.ClientFactoryInstance().RedisClient().HSet(api.RoomsRedisKey(), param.CreatorID, string(jsonBytes))
		if err != nil {
			return CreateRoomResponse{}, err
		} else {
			return CreateRoomResponse{RoomID: roomID}, nil
		}
	} else {
		data, ok := value.(string)
		if !ok {
			return CreateRoomResponse{}, errors.New("redis room data type error")
		}
		roomData := &api.RedisCreateRoomData{}
		err := json.Unmarshal([]byte(data), roomData)
		if err != nil {
			return CreateRoomResponse{}, err
		}
		return CreateRoomResponse{RoomID: roomData.RoomID}, nil
	}
}

func (room *ChatRoomRedis) DestroyRoom(roomID string) error {
	if global.ClientFactoryInstance().RedisClient() == nil {
		return errors.New("redis cli required")
	}

	roomMap, err := global.ClientFactoryInstance().RedisClient().HGetAll(api.RoomsRedisKey())
	if err != nil {
		return err
	}

	found := false
	creatorID := ""
	for _, value := range roomMap {
		roomData := &api.RedisCreateRoomData{}
		err = json.Unmarshal([]byte(value), roomData)
		if err != nil {
			continue
		}
		if roomData.RoomID == roomID {
			found = true
			creatorID = roomData.CreatorID
			break
		}
	}

	if !found {
		return fmt.Errorf("room[%v] not exist", roomID)
	}

	err = global.ClientFactoryInstance().RedisClient().HDel(api.RoomsRedisKey(), creatorID)
	if err != nil {
		return err
	}

	return nil
}
