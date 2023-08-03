// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// define some data structure which outside module want to communicate with this module

package groups

type CreateGroupParam struct {
	GroupID         string
	CreateUserID    string
	CreateTimestamp int64
}

type DeleteGroupParam struct {
	GroupID string
}

type Add2GroupParam struct {
	GroupID     string
	AddedUserID string
}

type RemoveFromGroupParam struct {
	GroupID       string
	RemovedUserID string
}
