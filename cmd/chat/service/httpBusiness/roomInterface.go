// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// room interface module: room's CRUD, the storage method of room's data can be implemented by different ways,
// for example: redis or mysql or other modules

package httpBusiness

type CreateRoomParam struct {
	CreatorID string
}

type CreateRoomResponse struct {
	RoomID string
}

type ChatRoom interface {
	CreateRoom(param CreateRoomParam) (CreateRoomResponse, error)
	DestroyRoom(roomID string) error
}
