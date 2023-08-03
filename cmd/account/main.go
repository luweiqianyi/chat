package main

import (
	"chat/cmd/account/config"
	"chat/cmd/account/global"
	"chat/cmd/account/routing"
	"chat/cmd/account/service/dao"
	"chat/cmd/account/service/rpc"
	"chat/pkg/log"
	"chat/pkg/util"
)

func main() {
	cfg := config.GlobalCfg()
	global.InitClientFactoryInstance(cfg)

	err := dao.CreateTable()
	if err != nil {
		log.Infof("%v", err)
		return
	}

	go global.StartRpcSever(rpc.RegisterAccountValidateService)
	global.StartHttpServer(new(routing.RouterHandlerImpl))

	util.MainProcessShutdownGracefully()
}
