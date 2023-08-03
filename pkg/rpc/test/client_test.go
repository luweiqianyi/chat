package test

import (
	"chat/pkg/rpc/test/pb"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
	"time"
)

// TestRpcLogin 模拟外部业务代码发起一次rpc请求
func TestRpcLogin(t *testing.T) {
	targetHost := "localhost:8089"
	conn, err := grpc.Dial(targetHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Dial %s failed,err=%v\n", targetHost, err)
	}
	defer conn.Close()

	// 构建一次用户名和密码匹配的rpc请求
	c := pb.NewLoginServiceClient(conn)
	t1 := time.Now()
	rpcReply, err := c.Login(context.Background(), &pb.LoginRequest{
		Username: "leebai",
		Password: "123456",
	})
	log.Printf("spend time:%v ms", time.Now().Sub(t1).Milliseconds())
	if err != nil {
		log.Printf("login failed,err=%v\n", err)
	} else {
		log.Printf("login response: %v\n", rpcReply)
	}
}
