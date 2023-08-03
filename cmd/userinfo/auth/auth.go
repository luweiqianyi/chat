package auth

import (
	"chat/cmd/userinfo/auth/pb"
	"chat/cmd/userinfo/global"
	"chat/pkg/log"
	"context"
)

type AccountParam struct {
	AccountName string `form:"accountName"`
	AccessToken string `form:"token"`
}

func IsTokenValid(param AccountParam) bool {
	rpcReply, err := pb.NewAccountValidateServiceClient(global.AccountServerRpcCliConn()).AccountValidate(
		context.Background(),
		&pb.AccountValidateReq{
			AccountName: param.AccountName,
			AccessToken: param.AccessToken,
		})
	if err != nil {
		log.Errorf("account[%v] validate failed,err: %v\n", param.AccountName, err)
		return false
	}

	return rpcReply.Success
}
