package config

import (
	"chat/pkg/http"
	"sync"
)

type Config struct {
	HttpServerConfig http.ServerConfig
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

		HttpServerConfig: http.ServerConfig{
			HttpPort:       ":9014",
			UseHttps:       false,
			HttpsPort:      "",
			TlsCert:        "",
			TLSKey:         "",
			TrustedProxies: []string{"localhost"},
		},
	}
	return cfg
}
