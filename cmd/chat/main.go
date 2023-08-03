package main

import (
	"chat/cmd/chat/config"
	"chat/cmd/chat/global"
	"chat/cmd/chat/groups"
	"chat/cmd/chat/service"
	"chat/cmd/chat/service/wsBusiness"
	"chat/cmd/chat/users"
	"chat/pkg/log"
	"chat/pkg/util"
)

func main() {

	globalCfg := config.GlobalCfg()
	global.InitClientFactoryInstance(globalCfg)

	groupManager := groups.NewGroupManager()

	userManager := users.NewUserManager()
	userManager.RegisterGroupManager(groupManager)
	wsBusiness.RegisterUserManager(userManager)

	err := global.StartHttpServer(new(service.RouterManagerInterfaceImpl))
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	util.MainProcessShutdownGracefully()
}
