package store

import (
	"chat/cmd/userinfo/global"
	"chat/cmd/userinfo/store/pb"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
)

// SaveAvatarFile initiate(start) a remote rpc call to upload avatar to remote image-proc server
func SaveAvatarFile(fileName string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", errors.New("file open failed")
	}
	defer src.Close()

	buffer := make([]byte, file.Size)
	_, err = src.Read(buffer)
	if err != nil {
		return "", errors.New("file store to buffer failed")
	}

	reply, err := pb.NewFileServiceClient(global.FileServerRpcCliConn()).UploadFile(
		context.Background(),
		&pb.FileUploadRequest{
			FileData: buffer,
			FileName: fileName,
		})
	if err != nil {
		return "", err
	}

	if !reply.Success {
		return "", fmt.Errorf("%v", reply.Error)
	}
	return reply.FileSavedPath, nil
}

func DeleteAvatarFile(path string) error {
	_, err := pb.NewFileServiceClient(global.FileServerRpcCliConn()).DeleteFile(
		context.Background(),
		&pb.FileDeleteRequest{
			FilePath: path,
		})
	if err != nil {
		return err
	}
	return nil
}
