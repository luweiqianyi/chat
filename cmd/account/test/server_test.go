package test

import (
	"chat/cmd/account/pb"
	"chat/pkg/rpc"
	"google.golang.org/grpc"
	"testing"
)

// TestStartRpcServer 模拟外部业务模块调用pkg/rpc模块
func TestStartRpcServer(t *testing.T) {
	cfg := rpc.ServerConfig{
		Port: ":8089",
	}
	server := rpc.NewRpcServer(cfg)
	//if server != nil {
	//	server.RegisterRpcServiceImpl(RegisterAccountValidateService)
	//}
	server.Start(RegisterAccountValidateService)
}

// RegisterAccountValidateService 外部业务模块代码,通过脚本生成rpc请求的客户端和服务端, 然后由业务代码实现具体rpc服务端处理逻辑,比如LoginServiceImpl
func RegisterAccountValidateService(s *grpc.Server) {
	pb.RegisterAccountValidateServiceServer(s, new(AccountValidateServiceImpl))
}
