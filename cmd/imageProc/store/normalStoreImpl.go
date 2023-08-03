package store

import (
	"chat/pkg/util"
	"os"
)

type NormalFileStoreImpl struct {
}

func (impl NormalFileStoreImpl) SaveBytesToRemoteServer(fileBytes []byte, fileName string) (string, error) {
	err := util.EnsureFolderExist(FileSavePathPrefix)
	if err != nil {
		return "", err
	}

	savedFileName := SavedPath(FileSavePathPrefix, fileName)
	fd, err := os.Create(savedFileName)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	_, err = fd.Write(fileBytes)
	if err != nil {
		return "", err
	}

	return savedFileName, nil
}

func (impl NormalFileStoreImpl) DeleteFileFromRemoteServer(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
