package test

import (
	"chat/pkg/rpc/test/pb"
	"google.golang.org/grpc"
	"testing"
)

// TestStartRpcLoginServer 模拟外部业务模块调用pkg/rpc模块
func TestStartRpcLoginServer(t *testing.T) {
	//cfg := rpc.ServerConfig{
	//	Port: ":8089",
	//}
	//rpc.NewRpcServer(cfg, RegisterLoginService)
}

// RegisterLoginService 外部业务模块代码,通过脚本生成rpc请求的客户端和服务端, 然后由业务代码实现具体rpc服务端处理逻辑,比如LoginServiceImpl
func RegisterLoginService(s *grpc.Server) {
	pb.RegisterLoginServiceServer(s, new(LoginServiceImpl))
}
