package config

import (
	"chat/pkg/rpc"
	"sync"
)

// SftpServerConfig https://hub.docker.com/r/atmoz/sftp
type SftpServerConfig struct {
	ServerIP string
	SSHPort  string
	UserName string
	Password string
}

type Config struct {
	FileServerRpcCfg rpc.ServerConfig
	SftpServerConfig SftpServerConfig
}

var once sync.Once
var globalCfg Config

func GlobalCfg() Config {
	once.Do(func() {
		globalCfg = LoadFromMemory()
		// TODO parse from config file
	})
	return globalCfg
}

func LoadFromMemory() Config {
	cfg := Config{
		FileServerRpcCfg: rpc.ServerConfig{
			Port: ":8089",
		},
		SftpServerConfig: SftpServerConfig{
			ServerIP: "localhost",
			SSHPort:  "443",
			UserName: "LeeBai",
			Password: "123456",
		},
	}
	return cfg
}
