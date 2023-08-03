package handler

import (
	"chat/cmd/account/routing/api"
	"chat/cmd/account/service"
	"chat/pkg/http/common"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RegisterHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param api.RegisterParam
		if err := context.ShouldBind(&param); err != nil {
			context.JSON(200, common.NewParameterErrorResponse())
			return
		}
		if param.AccountName == "" {
			context.JSON(200, common.NewCustomizeFailedResponse(fmt.Sprintf("accountName empty")))
			return
		}

		retCode := service.Register(param)
		switch retCode {
		case service.Success:
			context.JSON(200, common.NewSuccessResponse())
		case service.AccountAlreadyExist:
			context.JSON(200, common.NewCustomizeFailedResponse(fmt.Sprintf("account[%v] already exist", param.AccountName)))
		default:
			context.JSON(200, common.NewFailedResponse())
		}
	}
}

func UnRegisterHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param api.UnRegisterParam
		if err := context.ShouldBind(&param); err != nil {
			context.JSON(200, common.NewParameterErrorResponse())
			return
		}

		retCode := service.UnRegister(param)
		switch retCode {
		case service.Success:
			context.JSON(200, common.NewSuccessResponse())
		case service.AccountNotExist:
			context.JSON(200, common.NewCustomizeFailedResponse(fmt.Sprintf("account[%v] already unregister", param.AccountName)))
		default:
			context.JSON(200, common.NewFailedResponse())
		}
	}
}

func LoginHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param api.LoginParam
		if err := context.ShouldBind(&param); err != nil {
			context.JSON(200, common.NewParameterErrorResponse())
			return
		}

		ret, token := service.Login(param)
		switch ret {
		case service.Success:
			context.JSON(200, api.LoginResponse{
				ResponseHeader: common.NewSuccessResponse(),
				Token:          token,
			})
		case service.AccountNotExist:
			context.JSON(200, common.NewCustomizeFailedResponse(fmt.Sprintf("account[%v] not registered", param.AccountName)))
		default:
			context.JSON(200, common.NewFailedResponse())
		}
	}
}

func LogoutHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param api.LogoutParam
		if err := context.ShouldBind(&param); err != nil {
			context.JSON(200, common.NewParameterErrorResponse())
			return
		}

		retCode := service.Logout(param)
		switch retCode {
		case service.Success:
			context.JSON(200, common.NewSuccessResponse())
		case service.AccountNotExist:
			context.JSON(200, common.NewCustomizeFailedResponse(fmt.Sprintf("account[%v] not registered", param.AccountName)))
		default:
			context.JSON(200, common.NewFailedResponse())
		}
	}
}
