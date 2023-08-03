// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Kernel module: implementation of various business functions

package wsBusiness

import (
	"chat/cmd/chat/api"
	"chat/cmd/chat/users"
	"chat/pkg/log"
	"chat/pkg/ws"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/prometheus/client_golang/prometheus"
)

var gUserManager *users.UserManager

func RegisterUserManager(manager *users.UserManager) {
	gUserManager = manager
}

type WebsocketBusinessImpl struct {
}

func NewWebsocketBusinessImpl() WebsocketBusinessImpl {
	return WebsocketBusinessImpl{}
}

func (impl WebsocketBusinessImpl) RegisterMetricsToRegistry(registry *prometheus.Registry) {
	RegisterMonitorIndicators(registry)
}

func (impl WebsocketBusinessImpl) OnConnected(client *ws.WebsocketClient) {
	log.Infof("remoteClient[%v] connected", client.ID())
}

func (impl WebsocketBusinessImpl) OnDisconnected(client *ws.WebsocketClient) {
	log.Infof("remoteClient[%v] disconnected", client.ID())

	if gUserManager != nil {
		gUserManager.DeleteUserByClientID(client.ID())
		RefreshUsersInfoMetrics()
	}
}

func (impl WebsocketBusinessImpl) IsKeepAliveMessage(message []byte) bool {
	requestParam := &api.AppKeepAliveParam{}
	err := json.Unmarshal(message, requestParam)
	if err != nil {
		return false
	}

	_, err = govalidator.ValidateStruct(requestParam)
	if err != nil {
		return false
	}

	if requestParam.Method == api.Keepalive {
		return true
	}
	return false
}

func (impl WebsocketBusinessImpl) OnMessageCallback(message []byte, client *ws.WebsocketClient) []byte {
	requestParam := &api.RequestHeader{}
	err := json.Unmarshal(message, requestParam)
	if err != nil {
		response := api.NewErrorResponseWithTip(
			"",
			"",
			api.ParameterFormatError,
			string(message)+": "+api.ResponseString(api.ParameterFormatError))
		responseJsonBytes, _ := json.Marshal(response)
		return responseJsonBytes
	}

	_, err = govalidator.ValidateStruct(requestParam)
	if err != nil {
		response := api.NewErrorResponse(
			requestParam.MsgID,
			requestParam.Method,
			api.ParameterFormatError)
		responseJsonBytes, _ := json.Marshal(response)
		return responseJsonBytes
	}

	if !client.IsLogin() && requestParam.Method != api.Login {
		response := api.NewErrorResponse(
			requestParam.MsgID,
			requestParam.Method,
			api.UserNotLogin)
		responseJsonBytes, _ := json.Marshal(response)
		return responseJsonBytes
	}

	var response interface{}
	switch requestParam.Method {
	case api.Login:
		response = LoginHandler(message, client, func() {
			client.SetLogin(true)

			RefreshUsersInfoMetrics()
		})
	case api.Logout:
		response = LogoutHandler(message, func() {
			client.SetLogin(false)

			RefreshUsersInfoMetrics()
		})
	case api.Keepalive:
		response = KeepAliveHandler(message)
	case api.CreateGroup:
		response = CreateGroupHandler(message)
	case api.DeleteGroup:
		response = DeleteGroupHandler(message)
	case api.AddToGroup:
		response = Add2GroupHandler(message)
	case api.RemoveFromGroup:
		response = RemoveFromGroupHandler(message)
	case api.One2One:
		response = Send2OneHandler(message)
	case api.One2Group:
		response = Send2GroupHandler(message)
	}

	responseJsonBytes, err := json.Marshal(response)
	if err != nil {
		response = api.NewErrorResponse(
			requestParam.MsgID,
			requestParam.Method,
			api.ServerInternalError)
		responseJsonBytes, _ = json.Marshal(response)
		return responseJsonBytes
	}

	return responseJsonBytes
}

func KeepAliveHandler(message []byte) interface{} {
	requestParam := &api.AppKeepAliveParam{}
	err := json.Unmarshal(message, requestParam)
	if err != nil {
		return api.NewErrorResponseWithTip(
			"",
			"",
			api.ParameterFormatError,
			string(message)+": "+api.ResponseString(api.ParameterFormatError))
	}

	_, err = govalidator.ValidateStruct(requestParam)
	if err != nil {
		return api.NewErrorResponse(
			requestParam.MsgID,
			requestParam.Method,
			api.ParameterFormatError)
	}

	return api.NewSuccessResponse(requestParam.MsgID, requestParam.Method)
}

func LoginValidate(accountName string, accessToken string) bool {
	return true
}

func LoginHandler(message []byte, client *ws.WebsocketClient, onSuccessCallback func()) interface{} {
	requestParam := &api.LoginParam{}
	err := json.Unmarshal(message, requestParam)
	if err != nil {
		return api.NewErrorResponseWithTip(
			"",
			"",
			api.ParameterFormatError,
			string(message)+": "+api.ResponseString(api.ParameterFormatError))
	}

	_, err = govalidator.ValidateStruct(requestParam)
	if err != nil {
		return api.NewErrorResponse(
			requestParam.MsgID,
			requestParam.Method,
			api.ParameterFormatError)
	}

	pass := LoginValidate(requestParam.AccountName, requestParam.AccessToken)
	if !pass {
		return api.NewErrorResponse(requestParam.MsgID, requestParam.Method, api.LoginFailed)
	}

	if gUserManager != nil {
		gUserManager.AddUser(&users.User{
			AccountName: requestParam.AccountName,
			UserInfo: users.UserInfo{
				Client: client,
			},
		})
	}

	onSuccessCallback()
	return api.NewSuccessResponse(requestParam.MsgID, requestParam.Method)
}

func LogoutHandler(message []byte, onSuccessCallback func()) interface{} {
	requestParam := &api.LogoutParam{}
	err := json.Unmarshal(message, requestParam)
	if err != nil {
		return api.NewErrorResponseWithTip(
			"",
			"",
			api.ParameterFormatError,
			string(message)+": "+api.ResponseString(api.ParameterFormatError))
	}

	_, err = govalidator.ValidateStruct(requestParam)
	if err != nil {
		return api.NewErrorResponse(
			requestParam.MsgID,
			requestParam.Method,
			api.ParameterFormatError)
	}

	if gUserManager != nil {
		gUserManager.DeleteUser(requestParam.AccountName)
	}

	onSuccessCallback()
	return api.NewSuccessResponse(requestParam.MsgID, requestParam.Method)
}

func CreateGroupHandler(message []byte) interface{} {
	requestParam := &api.CreateGroupParam{}
	err := json.Unmarshal(message, requestParam)
	if err != nil {
		return api.NewErrorResponseWithTip(
			"",
			"",
			api.ParameterFormatError,
			string(message)+": "+api.ResponseString(api.ParameterFormatError))
	}

	_, err = govalidator.ValidateStruct(requestParam)
	if err != nil {
		return api.NewErrorResponse(
			requestParam.MsgID,
			requestParam.Method,
			api.ParameterFormatError)
	}

	err = gUserManager.CreateGroup(requestParam.GroupID, requestParam.SenderID)
	if err != nil {
		log.Errorf("%v: %v", api.ResponseString(api.GroupCreateFailed), err)
		return api.NewErrorResponse(requestParam.MsgID, requestParam.Method, api.GroupCreateFailed)
	}

	return api.NewSuccessResponse(requestParam.MsgID, requestParam.Method)
}

func DeleteGroupHandler(message []byte) interface{} {
	requestParam := &api.DeleteGroupParam{}
	err := json.Unmarshal(message, requestParam)
	if err != nil {
		return api.NewErrorResponseWithTip(
			"",
			"",
			api.ParameterFormatError,
			string(message)+": "+api.ResponseString(api.ParameterFormatError))
	}

	_, err = govalidator.ValidateStruct(requestParam)
	if err != nil {
		return api.NewErrorResponse(
			requestParam.MsgID,
			requestParam.Method,
			api.ParameterFormatError)
	}

	err = gUserManager.DeleteGroup(requestParam.GroupID)
	if err != nil {
		log.Errorf("%v: %v", api.ResponseString(api.GroupDestroyFailed), err)
		return api.NewErrorResponse(requestParam.MsgID, requestParam.Method, api.GroupDestroyFailed)
	}

	return api.NewSuccessResponse(requestParam.MsgID, requestParam.Method)
}

func Add2GroupHandler(message []byte) interface{} {
	requestParam := &api.Add2GroupParam{}
	err := json.Unmarshal(message, requestParam)
	if err != nil {
		return api.NewErrorResponseWithTip(
			"",
			"",
			api.ParameterFormatError,
			string(message)+": "+api.ResponseString(api.ParameterFormatError))
	}

	_, err = govalidator.ValidateStruct(requestParam)
	if err != nil {
		return api.NewErrorResponse(
			requestParam.MsgID,
			requestParam.Method,
			api.ParameterFormatError)
	}

	err = gUserManager.AddUser2Group(requestParam.GroupID, requestParam.ReceiverID)
	if err != nil {
		return api.NewErrorResponse(requestParam.MsgID, requestParam.Method, api.Add2GroupFailed)
	}

	return api.NewSuccessResponse(requestParam.MsgID, requestParam.Method)
}

func RemoveFromGroupHandler(message []byte) interface{} {
	requestParam := &api.RemoveFromGroupParam{}
	err := json.Unmarshal(message, requestParam)
	if err != nil {
		return api.NewErrorResponseWithTip(
			"",
			"",
			api.ParameterFormatError,
			string(message)+": "+api.ResponseString(api.ParameterFormatError))
	}

	_, err = govalidator.ValidateStruct(requestParam)
	if err != nil {
		return api.NewErrorResponse(
			requestParam.MsgID,
			requestParam.Method,
			api.ParameterFormatError)
	}

	err = gUserManager.RemoveUserFromGroup(requestParam.GroupID, requestParam.ReceiverID)
	if err != nil {
		return api.NewErrorResponse(requestParam.MsgID, requestParam.Method, api.RemoveFromGroupFailed)
	}

	return api.NewSuccessResponse(requestParam.MsgID, requestParam.Method)
}

func Send2OneHandler(message []byte) interface{} {
	requestParam := &api.One2OneMessageParam{}
	err := json.Unmarshal(message, requestParam)
	if err != nil {
		return api.NewErrorResponseWithTip(
			"",
			"",
			api.ParameterFormatError,
			string(message)+": "+api.ResponseString(api.ParameterFormatError))
	}

	_, err = govalidator.ValidateStruct(requestParam)
	if err != nil {
		return api.NewErrorResponse(
			requestParam.MsgID,
			requestParam.Method,
			api.ParameterFormatError)
	}

	err = gUserManager.SendMessage2Receiver(requestParam.ReceiverID, []byte(requestParam.Message))
	if err != nil {
		return api.NewErrorResponseWithTip(
			requestParam.MsgID,
			requestParam.Method,
			api.Send2OneMessageFailed,
			fmt.Sprintf("%v", err))
	}
	return api.NewSuccessResponse(requestParam.MsgID, requestParam.Method)
}

func Send2GroupHandler(message []byte) interface{} {
	requestParam := &api.One2GroupMessageParam{}
	err := json.Unmarshal(message, requestParam)
	if err != nil {
		return api.NewErrorResponseWithTip(
			"",
			"",
			api.ParameterFormatError,
			string(message)+": "+api.ResponseString(api.ParameterFormatError))
	}

	_, err = govalidator.ValidateStruct(requestParam)
	if err != nil {
		return api.NewErrorResponse(
			requestParam.MsgID,
			requestParam.Method,
			api.ParameterFormatError)
	}

	err = gUserManager.SendGroupMessage(requestParam.GroupID, []byte(requestParam.Message))
	if err != nil {
		return api.NewErrorResponseWithTip(
			requestParam.MsgID,
			requestParam.Method,
			api.Send2GroupMessageFailed,
			fmt.Sprintf("%v", err))
	}

	return api.NewSuccessResponse(requestParam.MsgID, requestParam.Method)
}
