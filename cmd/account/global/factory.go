package global

import (
	"chat/cmd/account/config"
	"chat/pkg/http"
	"chat/pkg/log"
	"chat/pkg/mysql"
	"chat/pkg/redis"
	"chat/pkg/rpc"
	"gorm.io/gorm"
	"sync"
)

var once sync.Once
var gClientFactory *ClientFactory

type ClientFactory struct {
	httpServer *http.Server
	redisCli   *redis.Client
	mysqlCli   *mysql.Client
	rpcServer  *rpc.GRpcServer
}

func InitClientFactoryInstance(cfg config.Config) {
	once.Do(func() {
		server := http.NewHttpServer(cfg.HttpCfg)
		redisCli := redis.NewClient(cfg.RedisCfg)
		mysqlCli := mysql.NewClient(cfg.MySqlCfg)
		rpcServer := rpc.NewRpcServer(cfg.AccountServerRpcCfg)

		gClientFactory = &ClientFactory{
			httpServer: server,
			redisCli:   redisCli,
			mysqlCli:   mysqlCli,
			rpcServer:  rpcServer,
		}
	})
}

func MySqlDB() *gorm.DB {
	return gClientFactory.mysqlCli.GetDB()
}

func RedisClient() *redis.Client {
	return gClientFactory.redisCli
}

func StartHttpServer(manager http.RouterManagerInterface) {
	if gClientFactory == nil {
		log.Panicf("ClientFactory required")
		return
	}

	if gClientFactory.httpServer == nil {
		log.Panicf("http server required")
		return
	}

	gClientFactory.httpServer.BindRouterManager(manager)

	if err := gClientFactory.httpServer.Start(); err != nil {
		log.Panicf("%v", err)
	}
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
