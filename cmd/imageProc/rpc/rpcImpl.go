package rpc

import (
	"chat/cmd/imageProc/pb"
	"chat/cmd/imageProc/store"
	"google.golang.org/grpc"

	"context"
)

func RegisterFileService(s *grpc.Server) {
	pb.RegisterFileServiceServer(s, new(FileServiceImpl))
}

type FileServiceImpl struct {
}

func (impl FileServiceImpl) UploadFile(ctx context.Context, request *pb.FileUploadRequest) (*pb.FileUploadResponse, error) {
	savedPath, err := new(store.NormalFileStoreImpl).SaveBytesToRemoteServer(request.FileData, request.FileName)
	if err != nil {
		return &pb.FileUploadResponse{
			Success:       false,
			FileSavedPath: "",
		}, err
	}

	return &pb.FileUploadResponse{
		Success:       true,
		FileSavedPath: savedPath,
	}, nil
}

func (impl FileServiceImpl) DeleteFile(ctx context.Context, request *pb.FileDeleteRequest) (*pb.FileDeleteResponse, error) {
	err := new(store.NormalFileStoreImpl).DeleteFileFromRemoteServer(request.FilePath)
	if err != nil {
		return &pb.FileDeleteResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &pb.FileDeleteResponse{
		Success: true,
		Error:   "",
	}, nil
}
