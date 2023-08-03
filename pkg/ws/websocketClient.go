package ws

import (
	"chat/pkg/log"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"sync/atomic"
)

type ClientNetworkInfo struct {
	RemoteIP   string
	RemotePort string
	// TODO other properties
}

type WebsocketClient struct {
	id                string
	clientNetworkInfo ClientNetworkInfo

	owner *WebsocketServer

	ctx    context.Context
	cancel context.CancelFunc

	conn *websocket.Conn

	receiveChan                chan []byte
	normalMessageHandleChan    chan []byte
	keepaliveMessageHandleChan chan []byte
	responseChan               chan []byte
	messageNotifyChan          chan []byte
	keepAliveAckChan           chan []byte

	callbacks BusinessInterface
	isLogin   atomic.Bool
}

func (c *WebsocketClient) IsLogin() bool {
	return c.isLogin.Load()
}

func (c *WebsocketClient) SetLogin(login bool) {
	c.isLogin.Store(login)
}

func (c *WebsocketClient) ID() string {
	return c.id
}

func (c *WebsocketClient) getOwner() *WebsocketServer {
	return c.owner
}

func NewWebsocketClient(ID string, owner *WebsocketServer, networkInfo ClientNetworkInfo, conn *websocket.Conn) *WebsocketClient {
	client := &WebsocketClient{
		id:                         ID,
		owner:                      owner,
		clientNetworkInfo:          networkInfo,
		conn:                       conn,
		receiveChan:                make(chan []byte, 1024),
		normalMessageHandleChan:    make(chan []byte, 1024),
		keepaliveMessageHandleChan: make(chan []byte, 1024),
		responseChan:               make(chan []byte, 1024),
		keepAliveAckChan:           make(chan []byte, 1024),
		messageNotifyChan:          make(chan []byte, 1024),
	}

	client.ctx, client.cancel = context.WithCancel(context.Background())
	return client
}

func (c *WebsocketClient) Cancel() {
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *WebsocketClient) ReadWebsocketMessageHandler() {
	defer func() {
		if p := recover(); p != nil {
			panic(p)
		}
	}()

	defer func(c *WebsocketClient) {
		err := c.ActiveShutdown()
		if err != nil {
			log.Errorf("client[%v] active shutdown failed, err: %v", c.ID(), err)
		}
		c.Cancel()
		log.Infof("client[%s] end to read message...", c.ID())
	}(c)

	log.Infof("client[%s] start to read message...", c.ID())
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Errorf("client[%v] read failed, err: %v", c.ID(), err)
			return
		}
		err = c.PutToReceiveChan(message)
		if err != nil {
			log.Errorf("client[%v] put: %v to receiveChan, err: %v", c.ID(), string(message), err)
			return
		} else {
			log.Debugf("client[%v] put: %v to receiveChan success", c.ID(), string(message))
		}
	}
}

func (c *WebsocketClient) ProcessWebsocketMessageHandler() {
	defer func() {
		if p := recover(); p != nil {
			panic(p)
		}

		log.Infof("client[%s] end to process message...", c.ID())
	}()

	log.Infof("client[%s] start to process message...", c.ID())
	for {
		select {
		case <-c.ctx.Done():
			return
		case message := <-c.receiveChan:
			if c.callbacks != nil && c.callbacks.IsKeepAliveMessage(message) {
				err := c.PutToKeepaliveMessageHandleChan(message)
				if err != nil {
					log.Errorf("%v", err)
				}
			} else {
				err := c.PutToNormalMessageHandleChan(message)
				if err != nil {
					log.Errorf("%v", err)
				}
			}
		case message := <-c.normalMessageHandleChan:
			c.onNormalMessage(message)
		case keepaliveMessage := <-c.keepaliveMessageHandleChan:
			c.OnKeepAliveMessage(keepaliveMessage)
		case response := <-c.responseChan:
			err := c.sendResponseToRemote(response)
			if err != nil {
				log.Errorf("client[%v] send response: %v to remote failed, err: %v", c.ID(), response, err)
				break
			}
		case keepAliveAck := <-c.keepAliveAckChan:
			err := c.sendKeepAliveAckToRemote(keepAliveAck)
			if err != nil {
				log.Errorf("client[%v] send keepAliveAck: %v to remote failed, err: %v", c.ID(), keepAliveAck, err)
				break
			}
		case notifyMessage := <-c.messageNotifyChan:
			err := c.sendNotifyMessageToRemote(notifyMessage)
			if err != nil {
				log.Errorf("client[%v] send notifyMessage: %v to remote failed, err: %v", c.ID(), notifyMessage, err)
				break
			}
		}
	}
}

func (c *WebsocketClient) StartHandleRequests() {
	go c.ReadWebsocketMessageHandler()
	go c.ProcessWebsocketMessageHandler()

	c.onConnected()
}

func (c *WebsocketClient) ActiveShutdown() error {
	c.onDisconnected()

	return c.conn.Close()
}

func (c *WebsocketClient) onConnected() {
	log.Infof("client[%v] connected", c.ID())
	if owner := c.getOwner(); owner != nil {
		owner.Register(c.ID(), c)
	}
	if c.callbacks != nil && c.callbacks.OnConnected != nil {
		c.callbacks.OnConnected(c)
	}
}

func (c *WebsocketClient) onDisconnected() {
	log.Infof("client[%v] disconnected", c.ID())
	if owner := c.getOwner(); owner != nil {
		owner.UnRegister(c.ID())
	}
	if c.callbacks != nil && c.callbacks.OnDisconnected != nil {
		c.callbacks.OnDisconnected(c)
	}
}

func (c *WebsocketClient) onNormalMessage(messageData []byte) {
	if c.callbacks != nil && c.callbacks.OnMessageCallback != nil {
		response := c.callbacks.OnMessageCallback(messageData, c)
		err := c.PutToResponseChan(response)
		if err != nil {
			log.Errorf("client[%v] put response: %v to response channel failed, err: %v", c.ID(), string(response), err)
		}
	}
}

func (c *WebsocketClient) OnKeepAliveMessage(keepalive []byte) {
	if c.callbacks != nil && c.callbacks.OnMessageCallback != nil {
		keepAliveAck := c.callbacks.OnMessageCallback(keepalive, c)
		err := c.PutToKeepAliveAckChan(keepAliveAck)
		if err != nil {
			log.Errorf("client[%v] put keepAliveAck:%v to response channel failed, err: %v", c.ID(), string(keepAliveAck), err)
		}
	}
}

func (c *WebsocketClient) PutToReceiveChan(message []byte) error {
	if len(c.receiveChan) == cap(c.receiveChan) {
		return fmt.Errorf("client[%v] receiveChan full", c.ID())
	}
	c.receiveChan <- message

	return nil
}

func (c *WebsocketClient) PutToNormalMessageHandleChan(message []byte) error {
	if len(c.normalMessageHandleChan) == cap(c.normalMessageHandleChan) {
		return fmt.Errorf("client[%v] normalMessageHandleChan full", c.ID())
	}
	c.normalMessageHandleChan <- message

	return nil
}

func (c *WebsocketClient) PutToKeepaliveMessageHandleChan(message []byte) error {
	if len(c.keepaliveMessageHandleChan) == cap(c.keepaliveMessageHandleChan) {
		return fmt.Errorf("client[%v] keepaliveMessageHandleChan full", c.ID())
	}
	c.keepaliveMessageHandleChan <- message

	return nil
}

func (c *WebsocketClient) PutToResponseChan(response []byte) error {
	if len(c.responseChan) == cap(c.responseChan) {
		return fmt.Errorf("client[%v] responseChan full", c.ID())
	}
	c.responseChan <- response
	return nil
}

func (c *WebsocketClient) PutToKeepAliveAckChan(response []byte) error {
	if len(c.keepAliveAckChan) == cap(c.keepAliveAckChan) {
		return fmt.Errorf("client[%v] keepAliveAckChan full", c.ID())
	}
	c.keepAliveAckChan <- response
	return nil
}

func (c *WebsocketClient) PutToMessageNotifyChan(response []byte) error {
	if len(c.messageNotifyChan) == cap(c.messageNotifyChan) {
		return fmt.Errorf("client[%v] messageNotifyChan full", c.ID())
	}
	c.messageNotifyChan <- response
	return nil
}

func (c *WebsocketClient) RegisterBusinessCallbacks(clientInterface BusinessInterface) {
	c.callbacks = clientInterface
}

func (c *WebsocketClient) sendResponseToRemote(response []byte) error {
	err := c.conn.WriteMessage(websocket.TextMessage, response)
	if err != nil {
		log.Errorf("client[%v] send response: %v failed, err: %v", c.ID(), string(response), err)
	} else {
		log.Infof("client[%v] send response: %v success", c.ID(), string(response))
	}
	return err
}

func (c *WebsocketClient) sendKeepAliveAckToRemote(keepaliveAck []byte) error {
	err := c.conn.WriteMessage(websocket.TextMessage, keepaliveAck)
	if err != nil {
		log.Errorf("client[%v] send keepaliveAck: %v failed, err: %v", c.ID(), string(keepaliveAck), err)
	} else {
		log.Infof("client[%v] send keepaliveAck: %v success", c.ID(), string(keepaliveAck))
	}
	return err
}

func (c *WebsocketClient) sendNotifyMessageToRemote(notifyMessage []byte) error {
	err := c.conn.WriteMessage(websocket.TextMessage, notifyMessage)
	if err != nil {
		log.Errorf("client[%v] send notifyMessage: %v failed, err: %v", c.ID(), string(notifyMessage), err)
	} else {
		log.Infof("client[%v] send notifyMessage: %v success", c.ID(), string(notifyMessage))
	}
	return err
}
