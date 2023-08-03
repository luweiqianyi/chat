package ws

import (
	"chat/pkg/log"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type WebsocketServer struct {
	id      string
	clients *sync.Map

	registry *prometheus.Registry
}

type WebSocketServerConfig struct {
	ServerIP   string
	ServerPort string

	Registry *prometheus.Registry
}

func NewWebsocketServer(cfg WebSocketServerConfig) *WebsocketServer {
	server := &WebsocketServer{
		clients:  &sync.Map{},
		id:       fmt.Sprintf("%v:%v", cfg.ServerIP, cfg.ServerPort),
		registry: cfg.Registry,
	}

	server.registerMetricsIndicators()
	return server
}

func (s *WebsocketServer) registerMetricsIndicators() {
	if s.registry != nil {
		s.registry.MustRegister(WebsocketClientsCount)
		s.registry.MustRegister(WebsocketClientDetail)
	}
}

func (s *WebsocketServer) ID() string {
	return s.id
}

func (s *WebsocketServer) Register(ID string, client *WebsocketClient) {
	s.clients.Store(ID, client)

	// TODO optimize:performance promote
	UpdateWebsocketClientsCount(s)
	UpdateWebsocketClientDetail(s)

	log.Infof(">>> server[%v] add client[%v]", s.ID(), ID)
}

func (s *WebsocketServer) UnRegister(ID string) {
	s.clients.Delete(ID)

	// TODO optimize:performance promote
	UpdateWebsocketClientsCount(s)
	UpdateWebsocketClientDetail(s)

	log.Infof("<<< server[%v] delete client[%v]", s.ID(), ID)
}

func (s *WebsocketServer) UniCast(ID string, message []byte) {
	value, found := s.clients.Load(ID)
	if found {
		client, ok := value.(*WebsocketClient)
		if ok {
			err := client.PutToResponseChan(message)
			if err != nil {
				log.Errorf("client[%v] put response:%v to response channel failed, err: %v", client.ID(), message, err)
			}
		}
	}
}

func (s *WebsocketServer) Broadcast(message []byte) {
	s.clients.Range(func(key, value any) bool {
		client, ok := value.(*WebsocketClient)
		if !ok {
			return true
		}

		err := client.PutToResponseChan(message)
		if err != nil {
			log.Errorf("client[%v] put response:%v to response channel failed, err: %v", client.ID(), string(message), err)
		}
		return true
	})
}

func (s *WebsocketServer) ClientsCount() int {
	count := 0
	s.clients.Range(func(key, value any) bool {
		count++
		return true
	})

	return count
}
