// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// entity Group and its operations

package users

import (
	"sync"
)

type Group struct {
	groupID string
	users   sync.Map
}

func NewGroup(groupID string) *Group {
	return &Group{
		groupID: groupID,
	}
}

func (g Group) AddUser(user *User) {
	g.users.Store(user.AccountName, user)
}

func (g Group) DelUser(accountName string) {
	g.users.Delete(accountName)
}

func (g Group) Broadcast(message []byte) {
	g.users.Range(func(key, value any) bool {
		user, ok := value.(*User)
		if !ok {
			return true
		}

		user.SendMessage(message)
		return true
	})
}
