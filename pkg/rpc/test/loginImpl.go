package test

import (
	"chat/pkg/rpc/test/pb"
	"context"
	"fmt"
)

type LoginServiceImpl struct {
}

func (impl LoginServiceImpl) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	fmt.Printf("received request: %v\n", request)
	username := request.Username
	password := request.Password
	if username == "leebai" && password == "123456" {
		return &pb.LoginResponse{
			Code: 200,
			Msg:  "success",
		}, nil
	} else {
		return &pb.LoginResponse{
			Code: 201,
			Msg:  "fail",
		}, nil
	}
}
