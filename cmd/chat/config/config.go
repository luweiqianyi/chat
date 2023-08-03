package config

import (
	"chat/pkg/http"
	"chat/pkg/redis"
	"chat/pkg/rpc"
	"sync"
)

// TODO define every config structures which this service uses, eg. mysql redis rpc etc

type PrometheusConfig struct {
	Enable bool
}

type Config struct {
	HttpCfg          http.ServerConfig
	RedisCfg         redis.Config
	PrometheusConfig PrometheusConfig
	RpcConnCfg       rpc.ConnConfig
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
	return Config{
		HttpCfg: http.ServerConfig{
			HttpPort:       ":9012",
			UseHttps:       false,
			HttpsPort:      "",
			TlsCert:        "",
			TLSKey:         "",
			TrustedProxies: []string{"localhost"},
		},
		RedisCfg: redis.Config{
			RedisAddr:     "localhost",
			RedisPort:     "6379",
			RedisPassword: "",
			DBIndex:       0,
		},
		PrometheusConfig: PrometheusConfig{
			Enable: true,
		},
	}
}

func LoadFromConfigFile() Config {
	return Config{}
}

func (cfg Config) PrometheusEnabled() bool {
	return cfg.PrometheusConfig.Enable
}
