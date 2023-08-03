package ws

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
)

var (
	WebsocketClientsCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "WebsocketClientsCount",
		Help: "The number of websocket clients which connected successfully",
	}, []string{"WebsocketClientsCount"})

	WebsocketClientDetail = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "WebsocketClientDetail",
		Help: "The detail of one specific websocket client",
	}, []string{"ClientID", "ClientInfo"})
)

func UpdateWebsocketClientsCount(s *WebsocketServer) {
	WebsocketClientsCount.Reset()
	WebsocketClientsCount.WithLabelValues(strconv.Itoa(s.ClientsCount()))
}

func UpdateWebsocketClientDetail(s *WebsocketServer) {
	WebsocketClientDetail.Reset()
	s.clients.Range(func(key, value any) bool {
		client, ok := value.(*WebsocketClient)
		if !ok {
			return true
		}
		WebsocketClientDetail.WithLabelValues(client.ID(), "")
		return true
	})
}
