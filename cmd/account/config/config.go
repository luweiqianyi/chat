package config

import (
	"chat/pkg/http"
	"chat/pkg/mysql"
	"chat/pkg/redis"
	"chat/pkg/rpc"
	"sync"
)

type Config struct {
	RedisCfg            redis.Config
	MySqlCfg            mysql.Config
	HttpCfg             http.ServerConfig
	AccountServerRpcCfg rpc.ServerConfig
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
		RedisCfg: redis.Config{
			RedisAddr:     "localhost",
			RedisPort:     "6379",
			RedisPassword: "",
			DBIndex:       0,
		},
		MySqlCfg: mysql.Config{
			DSN: "root:123456@tcp(localhost:3306)/chat?charset=utf8",
		},
		HttpCfg: http.ServerConfig{
			HttpPort:       ":9011",
			UseHttps:       false,
			HttpsPort:      "",
			TlsCert:        "",
			TLSKey:         "",
			TrustedProxies: []string{"127.0.0.1"},
		},
		AccountServerRpcCfg: rpc.ServerConfig{
			Port: ":8088",
		},
	}
	return cfg
}

// LoadFromConfigFile TODO
func LoadFromConfigFile() Config {
	return Config{}
}
