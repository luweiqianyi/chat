package main

import (
	"chat/cmd/stock/config"
	"chat/cmd/stock/global"
	"chat/cmd/stock/routing"
	"chat/pkg/log"
	"chat/pkg/util"
)

func main() {
	global.InitFactoryInstance(config.GlobalCfg())

	err := global.StartHttpServer(new(routing.RouterImpl))
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	util.MainProcessShutdownGracefully()
}
