package global

import (
	"chat/cmd/stock/config"
	"chat/pkg/http"
	"errors"
	"sync"
)

var once sync.Once
var gFatory *Factory

type Factory struct {
	httpServer *http.Server
}

func InitFactoryInstance(cfg config.Config) {
	once.Do(func() {
		httpServer := http.NewHttpServer(cfg.HttpServerConfig)

		gFatory = &Factory{
			httpServer: httpServer,
		}
	})
}

func StartHttpServer(manager http.RouterManagerInterface) error {
	if gFatory == nil {
		return errors.New("factory need init")
	}
	if gFatory.httpServer == nil {
		return errors.New("http server need init")
	}

	gFatory.httpServer.BindRouterManager(manager)
	return gFatory.httpServer.Start()
}
