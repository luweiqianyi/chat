package main

import (
	"chat/cmd/userinfo/config"
	"chat/cmd/userinfo/global"
	"chat/cmd/userinfo/routing"
	"chat/cmd/userinfo/service/dao"
	"chat/pkg/log"
	"chat/pkg/util"
)

func main() {
	cfg := config.GlobalCfg()
	global.InitClientFactoryInstance(cfg)

	err := dao.CreateTable()
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	global.StartHttpServer(new(routing.RouterHandlerImpl))

	util.MainProcessShutdownGracefully()
}
