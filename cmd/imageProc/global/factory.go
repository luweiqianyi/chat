package global

import (
	"chat/cmd/imageProc/config"
	"chat/pkg/log"
	"chat/pkg/rpc"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"sync"
)

var once sync.Once
var gClientFactory *ClientFactory

type ClientFactory struct {
	rpcServer *rpc.GRpcServer

	sftpClient *sftp.Client
}

func InitClientFactoryInstance(cfg config.Config) {
	once.Do(func() {
		rpcServer := rpc.NewRpcServer(cfg.FileServerRpcCfg)
		sftpClient, _ := newSftpClient(cfg.SftpServerConfig)

		gClientFactory = &ClientFactory{
			rpcServer:  rpcServer,
			sftpClient: sftpClient,
		}
	})
}

func StartRpcSever(registerRpcServiceCallback rpc.RegisterRpcServiceCallback) {
	defer func() {
		if p := recover(); p != nil {
			log.Panicf("%v", p)
		}
	}()

	if gClientFactory.rpcServer != nil {
		err := gClientFactory.rpcServer.Start(registerRpcServiceCallback)
		if err != nil {
			log.Panicf("%v", err)
		}
	}
}

func newSftpClient(serverConfig config.SftpServerConfig) (*sftp.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: serverConfig.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(serverConfig.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", serverConfig.ServerIP, serverConfig.SSHPort), sshConfig)
	if err != nil {
		return nil, err
	}

	client, err := sftp.NewClient(conn)
	return client, err
}

func SftpClient() *sftp.Client {
	return gClientFactory.sftpClient
}

func RecycleResources() {
	defer func(sftpClient *sftp.Client) {
		err := sftpClient.Close()
		if err != nil {
			log.Errorf("%v", err)
		}
	}(gClientFactory.sftpClient)
}
