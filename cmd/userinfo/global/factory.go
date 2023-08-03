package global

import (
	"chat/cmd/userinfo/config"
	"chat/pkg/http"
	"chat/pkg/log"
	"chat/pkg/mysql"
	"chat/pkg/rpc"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"sync"
)

var once sync.Once
var gClientFactory *ClientFactory

type ClientFactory struct {
	server            *http.Server
	mysqlCli          *mysql.Client
	fileServerConn    *rpc.GRpcClientConnWrapper
	accountServerConn *rpc.GRpcClientConnWrapper
}

func InitClientFactoryInstance(cfg config.Config) {
	once.Do(func() {
		server := http.NewHttpServer(cfg.HttpCfg)
		mysqlCli := mysql.NewClient(cfg.MySqlCfg)
		fileServerConn := rpc.NewGRpcClientConnWrapper(cfg.FileServerRpcConnCfg)
		accountServerConn := rpc.NewGRpcClientConnWrapper(cfg.AccountServerRpcConnCfg)

		gClientFactory = &ClientFactory{
			server:            server,
			mysqlCli:          mysqlCli,
			fileServerConn:    fileServerConn,
			accountServerConn: accountServerConn,
		}
	})
}

func GetMySqlDB() *gorm.DB {
	return gClientFactory.mysqlCli.GetDB()
}

func AccountServerRpcCliConn() *grpc.ClientConn {
	return gClientFactory.accountServerConn.UnderlyingConnection()
}

func FileServerRpcCliConn() *grpc.ClientConn {
	return gClientFactory.fileServerConn.UnderlyingConnection()
}

func StartHttpServer(manager http.RouterManagerInterface) {
	if gClientFactory == nil {
		log.Panicf("ClientFactory required")
		return
	}

	if gClientFactory.server == nil {
		log.Panicf("http server required")
		return
	}

	gClientFactory.server.BindRouterManager(manager)

	if err := gClientFactory.server.Start(); err != nil {
		log.Panicf("%v", err)
	}
}
