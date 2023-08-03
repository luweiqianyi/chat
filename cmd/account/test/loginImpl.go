package test

import (
	"chat/cmd/account/pb"
	"context"
	"fmt"
)

type AccountValidateServiceImpl struct {
}

func (impl AccountValidateServiceImpl) AccountValidate(ctx context.Context, request *pb.AccountValidateReq) (*pb.AccountValidateResp, error) {
	fmt.Printf("received request: %v\n", request)
	username := request.AccountName
	password := request.AccessToken
	if username == "leebai" && password == "leebai-token" {
		return &pb.AccountValidateResp{
			Success: true,
		}, nil
	} else {
		return &pb.AccountValidateResp{
			Success: false,
		}, nil
	}
}
