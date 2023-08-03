// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// image process service: handle upload、delete、modification of images

package main

import (
	"chat/cmd/imageProc/config"
	"chat/cmd/imageProc/global"
	"chat/cmd/imageProc/rpc"
	"chat/pkg/util"
)

func main() {
	cfg := config.GlobalCfg()
	global.InitClientFactoryInstance(cfg)

	go global.StartRpcSever(rpc.RegisterFileService)

	util.MainProcessShutdownGracefully()
	defer global.RecycleResources()
}
