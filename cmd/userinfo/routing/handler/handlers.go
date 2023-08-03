package handler

import (
	"chat/cmd/userinfo/auth"
	"chat/cmd/userinfo/routing/api"
	"chat/cmd/userinfo/service"
	"chat/cmd/userinfo/store"
	"chat/pkg/http/common"
	"chat/pkg/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
)

func AccountValidateInterceptor() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param auth.AccountParam
		if err := context.ShouldBind(&param); err != nil {
			context.JSON(http.StatusOK, common.NewParameterErrorResponse())
			context.Abort()
			return
		}

		pass := isTokenValid(param)
		if !pass {
			context.JSON(http.StatusOK, common.NewCustomizeFailedResponse("account invalid"))
			context.Abort()
			return
		}

		context.Next()
	}
}

// isTokenValid remoteRpcCall to validate token is a valid token
func isTokenValid(param auth.AccountParam) bool {
	return auth.IsTokenValid(param)
}

func AddUserInfoHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param api.AddUserInfoParam
		if err := context.ShouldBind(&param); err != nil {
			context.JSON(http.StatusOK, common.NewParameterErrorResponse())
			return
		}

		success := service.AdduserInfo(param)
		if success {
			context.JSON(http.StatusOK, common.NewSuccessResponse())
		} else {
			context.JSON(http.StatusOK, common.NewFailedResponse())
		}
	}
}

func DelUserInfoHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param api.DelUserInfoParam
		if err := context.ShouldBind(&param); err != nil {
			context.JSON(http.StatusOK, common.NewParameterErrorResponse())
			return
		}

		success := service.DelUserInfo(param.AccountParam.AccountName)
		if success {
			context.JSON(http.StatusOK, common.NewSuccessResponse())
		} else {
			context.JSON(http.StatusOK, common.NewFailedResponse())
		}
	}
}

func QueryUserInfoHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param api.QueryUserInfoParam
		if err := context.ShouldBind(&param); err != nil {
			context.JSON(http.StatusOK, common.NewParameterErrorResponse())
			return
		}

		userInfo, err := service.QueryUserInfo(param.AccountName)
		if err != nil {
			context.JSON(http.StatusOK, common.NewFailedResponse())
		} else {
			context.JSON(http.StatusOK, api.NewQueryUserInfoResponse(userInfo))
		}
	}
}

func UpdateUserInfoHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param api.UpdateUserInfoParam
		if err := context.ShouldBind(&param); err != nil {
			context.JSON(http.StatusOK, common.NewParameterErrorResponse())
			return
		}

		success := service.UpdateUserInfo(param)
		if success {
			context.JSON(http.StatusOK, common.NewSuccessResponse())
		} else {
			context.JSON(http.StatusOK, common.NewFailedResponse())
		}
	}
}

func UpdateNickNameHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param api.UpdateNickNameParam
		if err := context.ShouldBind(&param); err != nil {
			context.JSON(http.StatusOK, common.NewParameterErrorResponse())
			return
		}

		success := service.UpdateNickName(param.AccountName, param.NickName)
		if success {
			context.JSON(http.StatusOK, common.NewSuccessResponse())
		} else {
			context.JSON(http.StatusOK, common.NewFailedResponse())
		}
	}
}

func UpdateBirthdayHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, common.NewCustomizeFailedResponse("future support"))
	}
}

func UpdateAvatarHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param api.UpdateAvatarParam
		err := context.ShouldBind(&param)
		if err != nil {
			context.JSON(http.StatusOK, common.NewParameterErrorResponse())
			return
		}

		// 随机生成文件名, 避免业务失败时的删除操作晚于重试时业务成功的上传操作,导致图片丢失的现象
		// TODO 优化 随机文件名存在的问题，如果进行重复测试，那么同一份数据在服务端上会保存多份，需要一定的机制保证删除垃圾图片
		fileNameInImageServer := fmt.Sprintf(
			"%v-%v%v",
			param.AccountName,
			uuid.NewString(),
			filepath.Ext(param.AvatarFile.Filename))
		avatarSavedPath, err := store.SaveAvatarFile(fileNameInImageServer, param.AvatarFile)
		if err != nil {
			context.JSON(http.StatusOK, common.NewCustomizeFailedResponse(fmt.Sprintf("%v", err)))
			return
		}

		success := service.UpdateAvatarPath(param.AccountName, avatarSavedPath)
		if success {
			context.JSON(http.StatusOK, api.NewSuccessUpdateAvatarResp(avatarSavedPath))
		} else {
			// 任务失败, 删除文件服务器上的图片
			go func() {
				err := store.DeleteAvatarFile(avatarSavedPath)
				if err != nil {
					log.Errorf("%v", err)
					return
				}
			}()
			context.JSON(http.StatusOK, common.NewCustomizeFailedResponse("error: avatar path update failed"))
		}
	}
}
