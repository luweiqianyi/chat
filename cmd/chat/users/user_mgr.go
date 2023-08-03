// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// manager module

package users

import (
	"chat/cmd/chat/groups"
	"chat/cmd/chat/util"
	"errors"
	"fmt"
	"sync"
)

type UserManager struct {
	users  sync.Map
	groups sync.Map

	groupManager *groups.GroupManager
}

func NewUserManager() *UserManager {
	return &UserManager{}
}

func (userManager *UserManager) RegisterGroupManager(manager *groups.GroupManager) {
	userManager.groupManager = manager
}

func (userManager *UserManager) AddUser(user *User) {
	userManager.users.Store(user.AccountName, user)
}

func (userManager *UserManager) DeleteUser(accountName string) {
	userManager.users.Delete(accountName)
}

func (userManager *UserManager) GetUser(accountName string) (*User, bool) {
	value, found := userManager.users.Load(accountName)
	if !found {
		return nil, false
	}

	user, ok := value.(*User)
	if !ok {
		return nil, false
	}
	return user, true
}

func (userManager *UserManager) UserCount() int {
	count := 0
	userManager.users.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

func (userManager *UserManager) GetUsers() []*User {
	var users []*User
	userManager.users.Range(func(key, value any) bool {
		user, ok := value.(*User)
		if ok {
			users = append(users, user)
		}

		return true
	})
	return users
}

func (userManager *UserManager) CreateGroup(groupID string, createUserID string) error {
	// 1. add to remote
	err := userManager.groupManager.CreateGroup(groups.CreateGroupParam{
		GroupID:         groupID,
		CreateUserID:    createUserID,
		CreateTimestamp: util.GenerateTimestamp(),
	})
	if err != nil {
		return err
	}

	// 2. add to local
	userManager.groups.Store(groupID, NewGroup(groupID))
	return nil
}

func (userManager *UserManager) DeleteGroup(groupID string) error {
	// 1. delete from remote
	err := userManager.groupManager.DeleteGroup(groups.DeleteGroupParam{
		GroupID: groupID,
	})
	if err != nil {
		return err
	}

	// 2. delete from local
	userManager.groups.Delete(groupID)
	return nil
}

func (userManager *UserManager) AddUser2Group(groupID string, userID string) error {
	// 1. add to remote
	err := userManager.groupManager.Add2Group(groups.Add2GroupParam{
		GroupID:     groupID,
		AddedUserID: userID,
	})
	if err != nil {
		return err
	}

	// 2. add to local, TODO: think about rollback
	value, found := userManager.groups.Load(groupID)
	if !found {
		return fmt.Errorf("group[%v] not exist", groupID)
	}

	group, ok := value.(*Group)
	if !ok {
		return errors.New("group data type error")
	}

	user, found := userManager.GetUser(userID)
	if !found {
		return fmt.Errorf("user[%v] not exist", userID)
	}

	if group != nil {
		group.AddUser(user)
	}
	return nil
}

func (userManager *UserManager) RemoveUserFromGroup(groupID string, accountName string) error {
	// 1. remove from remote
	err := userManager.groupManager.RemoveFromGroup(groups.RemoveFromGroupParam{
		GroupID:       groupID,
		RemovedUserID: accountName,
	})
	if err != nil {
		return err
	}

	// 2. remove from local
	value, found := userManager.groups.Load(groupID)
	if !found {
		return fmt.Errorf("group[%v] not exist", groupID)
	}

	group, ok := value.(*Group)
	if !ok {
		return errors.New("group data type error")
	}

	if group != nil {
		group.DelUser(accountName)
	}

	return nil
}

func (userManager *UserManager) SendGroupMessage(groupID string, message []byte) error {
	value, found := userManager.groups.Load(groupID)
	if !found {
		return fmt.Errorf("group[%v] not exist", groupID)
	}

	group, ok := value.(*Group)
	if !ok {
		return errors.New("group data type error")
	}

	if group != nil {
		group.Broadcast(message)
	}

	return nil
}

func (userManager *UserManager) SendMessage2Receiver(accountName string, message []byte) error {
	value, found := userManager.users.Load(accountName)
	if !found {
		return fmt.Errorf("user[%v] not exist", accountName)
	}

	user, ok := value.(*User)
	if !ok {
		return errors.New("user data type error")
	}

	if user != nil {
		user.SendMessage(message)
	}
	return nil
}

func (userManager *UserManager) SendBroadcastMessage(message []byte) {
	userManager.users.Range(func(key, value any) bool {
		user, ok := value.(*User)
		if !ok {
			return true
		}

		if user != nil {
			user.SendMessage(message)
		}
		return true
	})
}

func (userManager *UserManager) DeleteUserByClientID(clientID string) {
	accountName := ""
	found := false
	userManager.users.Range(func(key, value any) bool {
		user, ok := value.(*User)
		if !ok {
			return true
		}
		if user.Client.ID() == clientID {
			accountName = key.(string)
			found = true
			return false
		}
		return true
	})

	if found {
		userManager.DeleteUser(accountName)
	}
}
