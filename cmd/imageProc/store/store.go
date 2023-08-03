// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// image store module: interface, other module should implement the interface below,
// for example: sftpStoreImpl.go is an implementation of file storage using sftp
// developer can also implement below interface using other method, e.g. alibaba oss

package store

type FileStoreInterface interface {
	SaveBytesToRemoteServer(fileBytes []byte, fileName string) (string, error)
	DeleteFileFromRemoteServer(path string) error
}
