package ws

type BusinessInterface interface {
	OnConnected(client *WebsocketClient)
	OnDisconnected(client *WebsocketClient)
	IsKeepAliveMessage(message []byte) bool
	OnMessageCallback(message []byte, client *WebsocketClient) []byte
}
