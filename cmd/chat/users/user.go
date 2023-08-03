// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// entity User and its operations

package users

import (
	"chat/pkg/log"
	"chat/pkg/ws"
)

type User struct {
	AccountName string //user's unique identify

	UserInfo
}

type UserInfo struct {
	Client *ws.WebsocketClient

	// TODO expand in the future
}

func (user User) SendMessage(message []byte) {
	err := user.Client.PutToMessageNotifyChan(message)
	if err != nil {
		log.Errorf("%v", err)
	}
}
