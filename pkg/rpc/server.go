package rpc

import (
	"chat/pkg/log"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"sync"
)

var once sync.Once
var gRpcServer *GRpcServer

type ServerConfig struct {
	Port string
}

type GRpcServer struct {
	srv      *grpc.Server
	listener net.Listener
}

type RegisterRpcServiceCallback func(server *grpc.Server)

func NewRpcServer(cfg ServerConfig) *GRpcServer {
	once.Do(func() {
		listener, s := newRpcServer(cfg)
		gRpcServer = &GRpcServer{
			srv:      s,
			listener: listener,
		}
	})
	return gRpcServer
}

func newRpcServer(cfg ServerConfig) (net.Listener, *grpc.Server) {
	s := grpc.NewServer()

	listener, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Errorf("rpc listen port %s failed, err: %v\n", cfg.Port, err)
		return nil, nil
	}

	return listener, s
}

func (srv *GRpcServer) Start(registerRpcServiceCallback RegisterRpcServiceCallback) error {
	if srv.srv == nil {
		return errors.New("grpc.Server required")
	}

	registerRpcServiceCallback(srv.srv)

	err := srv.srv.Serve(srv.listener)
	if err != nil {
		log.Errorf("rpc server start failed, err: %v", err)
		return fmt.Errorf("rpc server start failed, err: %v", err)
	}
	return nil
}
