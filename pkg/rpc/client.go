package rpc

import (
	"chat/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnConfig struct {
	Host string
	Port string
}

type GRpcClientConnWrapper struct {
	clientConn *grpc.ClientConn
}

func NewGRpcClientConnWrapper(cfg ConnConfig) *GRpcClientConnWrapper {
	targetHost := cfg.Host + ":" + cfg.Port
	conn, err := grpc.Dial(targetHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("Dial %s failed,err: %v\n", targetHost, err)
	}
	return &GRpcClientConnWrapper{
		clientConn: conn,
	}
}

func (w *GRpcClientConnWrapper) UnderlyingConnection() *grpc.ClientConn {
	return w.clientConn
}

func (w *GRpcClientConnWrapper) Close() {
	err := w.clientConn.Close()
	if err != nil {
		log.Errorf("rpc connection close failed, err: %v", err)
		return
	}
}

// TODO create connection pool, get connection obj from pool
