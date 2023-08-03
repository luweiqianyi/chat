package rpc

import (
	"chat/cmd/account/pb"
	"chat/cmd/account/service"
	"context"
	"google.golang.org/grpc"
)

type AccountValidateServiceImpl struct {
}

func (impl AccountValidateServiceImpl) AccountValidate(ctx context.Context, request *pb.AccountValidateReq) (*pb.AccountValidateResp, error) {
	success := service.ValidateAccount(request.AccountName, request.AccessToken)
	return &pb.AccountValidateResp{
		Success: success,
	}, nil
}

func RegisterAccountValidateService(s *grpc.Server) {
	pb.RegisterAccountValidateServiceServer(s, new(AccountValidateServiceImpl))
}
