// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// http route module: define some http request path and handlers to handle request

package service

import (
	"chat/cmd/chat/api"
	"chat/cmd/chat/config"
	"chat/cmd/chat/global"
	"chat/cmd/chat/service/httpBusiness"
	myWebsocket "chat/cmd/chat/service/wsBusiness"
	"chat/pkg/http/common"
	"chat/pkg/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
	"sync"
	"time"
)

var chatRoomOnce sync.Once
var gChatRoomStore *httpBusiness.ChatRoomRedis

func ChatRoomStoreInstance() *httpBusiness.ChatRoomRedis {
	chatRoomOnce.Do(func() {
		gChatRoomStore = &httpBusiness.ChatRoomRedis{}
	})
	return gChatRoomStore
}

type RouterManagerInterfaceImpl struct {
}

func (impl *RouterManagerInterfaceImpl) RegisterRouters(eg *gin.Engine) {
	eg.POST(api.CreateChatRoomPath, createChatRoomHandler())
	eg.POST(api.DestroyChatRoomPath, destroyChatRoomHandler())

	websocketBusinessImpl := myWebsocket.NewWebsocketBusinessImpl()
	websocketBusinessImpl.RegisterMetricsToRegistry(global.ClientFactoryInstance().PrometheusRegistry())
	eg.GET("/ws", WsHandler(websocketBusinessImpl))

	if config.GlobalCfg().PrometheusEnabled() {
		eg.GET("/metrics", gin.WrapH(promhttp.HandlerFor(
			global.ClientFactoryInstance().PrometheusRegistry(),
			promhttp.HandlerOpts{})))
	}
}

func createChatRoomHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param api.CreateRoomParam
		err := context.ShouldBind(&param)
		if err != nil {
			context.JSON(200, common.NewParameterErrorResponse())
			return
		}

		response, err := ChatRoomStoreInstance().CreateRoom(httpBusiness.CreateRoomParam{CreatorID: param.CreatorID})
		if err != nil {
			context.JSON(200, common.NewCustomizeFailedResponse(err.Error()))
		} else {
			context.JSON(200, api.CreateRoomResponse{
				ResponseHeader: common.NewSuccessResponse(),
				RoomID:         response.RoomID,
			})
		}
	}
}

func destroyChatRoomHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param api.DestroyRoomParam
		err := context.ShouldBind(&param)
		if err != nil {
			context.JSON(200, common.NewParameterErrorResponse())
			return
		}

		err = ChatRoomStoreInstance().DestroyRoom(param.RoomID)
		if err != nil {
			context.JSON(200, common.NewCustomizeFailedResponse(err.Error()))
		} else {
			context.JSON(200, common.NewSuccessResponse())
		}
	}
}

// websocket

var once sync.Once
var gWsServer *ws.WebsocketServer

var gWsUpgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //allow Cross-Origin Resource Sharing(CORS)
	},
}

var gWsResponseHeader = http.Header{}

const (
	timeFormat = "2006-01-02 15:04:05.000000"
)

func WsHandler(businessImpl ws.BusinessInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := gWsUpgrade.Upgrade(c.Writer, c.Request, gWsResponseHeader)
		if err != nil {
			//c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("websocket upgrade failed, err: %v", err)})
			return
		}

		// 1. create websocket server
		once.Do(func() {
			ip, _ := ws.GetServerIP()
			gWsServer = ws.NewWebsocketServer(ws.WebSocketServerConfig{
				ServerIP:   ip,
				ServerPort: config.GlobalCfg().HttpCfg.HttpPort,
				Registry:   global.ClientFactoryInstance().PrometheusRegistry(),
			})
		})

		// 2. create websocket client, bind client to server
		remoteIP, remotePort, _ := net.SplitHostPort(c.Request.RemoteAddr)
		connTime := time.Now().Format(timeFormat)
		websocketClient := ws.NewWebsocketClient(
			ws.GenerateWebsocketClientID(remoteIP, remotePort, connTime),
			gWsServer,
			ws.ClientNetworkInfo{
				RemoteIP:   remoteIP,
				RemotePort: remotePort,
			},
			conn,
		)

		// 3. register business implementation to websocket client
		websocketClient.RegisterBusinessCallbacks(businessImpl)

		// 4. start websocket client's logic handler
		websocketClient.StartHandleRequests()
	}
}
