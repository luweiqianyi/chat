package http

import (
	"chat/pkg/log"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

type ServerConfig struct {
	HttpPort       string
	UseHttps       bool
	HttpsPort      string
	TlsCert        string
	TLSKey         string
	TrustedProxies []string
}

type RouterManagerInterface interface {
	RegisterRouters(eg *gin.Engine)
}

type Server struct {
	eg            *gin.Engine
	cfg           ServerConfig
	routerManager RouterManagerInterface
}

func NewHttpServer(cfg ServerConfig) *Server {
	return &Server{
		eg:  gin.Default(),
		cfg: cfg,
	}
}

func (s *Server) BindRouterManager(manager RouterManagerInterface) {
	s.routerManager = manager
}

func (s *Server) Start() error {
	if s.eg == nil {
		return errors.New("http engine required")
	}

	s.routerManager.RegisterRouters(s.eg)

	if !strings.Contains(s.cfg.HttpPort, ":") {
		s.cfg.HttpPort = ":" + s.cfg.HttpPort
	}
	if !strings.Contains(s.cfg.HttpsPort, ":") {
		s.cfg.HttpsPort = ":" + s.cfg.HttpsPort
	}

	err := s.eg.SetTrustedProxies(s.cfg.TrustedProxies)
	if err != nil {
		log.Errorf("server set trusted proxies[%v] failed, err: %v", s.cfg.TrustedProxies, err)
	}

	err = s.eg.Run(s.cfg.HttpPort)
	if err != nil {
		log.Panicf("server start http failed, listen port: %v", s.cfg.HttpPort)
		return err
	} else {
		log.Infof("server start http success, listen port: %v", s.cfg.HttpPort)
	}

	if s.cfg.UseHttps {
		go func(eg *gin.Engine, cfg ServerConfig) {
			defer func() {
				if p := recover(); p != nil {
					log.Panicf("%v", p)
				}
			}()

			err = eg.RunTLS(cfg.HttpsPort, cfg.TlsCert, cfg.TLSKey)
			if err != nil {
				log.Errorf("server start https failed, listen port: %v", cfg.HttpsPort)
				return
			} else {
				log.Infof("server start https success, listen port: %v", cfg.HttpsPort)
			}
		}(s.eg, s.cfg)
	}

	return err
}
