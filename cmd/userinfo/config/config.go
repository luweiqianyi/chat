package config

import (
	"chat/pkg/http"
	"chat/pkg/mysql"
	"chat/pkg/rpc"
	"sync"
)

type Config struct {
	MySqlCfg                mysql.Config
	HttpCfg                 http.ServerConfig
	AccountServerRpcConnCfg rpc.ConnConfig
	FileServerRpcConnCfg    rpc.ConnConfig
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
		MySqlCfg: mysql.Config{
			DSN: "root:123456@tcp(localhost:3306)/chat?charset=utf8",
		},
		HttpCfg: http.ServerConfig{
			HttpPort:       ":9013",
			UseHttps:       false,
			HttpsPort:      "",
			TlsCert:        "",
			TLSKey:         "",
			TrustedProxies: []string{"127.0.0.1"},
		},
		AccountServerRpcConnCfg: rpc.ConnConfig{
			Host: "localhost",
			Port: "8088",
		},
		FileServerRpcConnCfg: rpc.ConnConfig{
			Host: "localhost",
			Port: "8089",
		},
	}
	return cfg
}

// LoadFromConfigFile TODO
func LoadFromConfigFile() Config {
	return Config{}
}
