package util

import (
	"chat/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func MainProcessShutdownGracefully() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	exitCode := <-ch
	log.Infof("main process exit Code:%v", exitCode)
}
