package global

import (
	"chat/cmd/chat/config"
	"chat/cmd/chat/metrics"
	"chat/pkg/http"
	"chat/pkg/redis"
	"chat/pkg/rpc"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var once sync.Once
var gClientFactory *ClientFactory

type ClientFactory struct {
	redisCli *redis.Client

	registry *prometheus.Registry

	rpcClientConn *rpc.GRpcClientConnWrapper

	httpServer *http.Server
}

func InitClientFactoryInstance(cfg config.Config) {
	once.Do(func() {
		redisCli := redis.NewClient(cfg.RedisCfg)
		rpcClientConn := rpc.NewGRpcClientConnWrapper(cfg.RpcConnCfg)
		httpServer := http.NewHttpServer(cfg.HttpCfg)

		gClientFactory = &ClientFactory{
			redisCli:      redisCli,
			registry:      metrics.GetRegistry(),
			rpcClientConn: rpcClientConn,
			httpServer:    httpServer,
		}
	})
}

func ClientFactoryInstance() *ClientFactory {
	return gClientFactory
}

func (factory *ClientFactory) PrometheusRegistry() *prometheus.Registry {
	return factory.registry
}

func (factory *ClientFactory) RedisClient() *redis.Client {
	return factory.redisCli
}

func (factory *ClientFactory) RpcClient() *rpc.GRpcClientConnWrapper {
	return factory.rpcClientConn
}

func StartHttpServer(routerManagerInterface http.RouterManagerInterface) error {
	server := gClientFactory.httpServer
	if server == nil {
		return errors.New("http server required")
	}
	server.BindRouterManager(routerManagerInterface)
	return server.Start()
}
