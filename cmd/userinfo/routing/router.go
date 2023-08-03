package routing

import (
	"chat/cmd/userinfo/routing/handler"
	"chat/cmd/userinfo/routing/path"
	"chat/pkg/log"
	"github.com/gin-gonic/gin"
)

type RouterHandlerImpl struct {
}

func (impl RouterHandlerImpl) RegisterRouters(eg *gin.Engine) {
	if eg == nil {
		log.Errorf("gin.Engine required")
		return
	}

	eg.Use(handler.AccountValidateInterceptor())

	group := eg.Group(path.UserInfoPath)
	{
		group.POST(path.AddUserInfo, handler.AddUserInfoHandler())
		group.GET(path.DelUserInfo, handler.DelUserInfoHandler())
		group.GET(path.QueryUserInfo, handler.QueryUserInfoHandler())
		group.POST(path.UpdateUserInfo, handler.UpdateUserInfoHandler())
		group.POST(path.UpdateNickName, handler.UpdateNickNameHandler())
		group.POST(path.UpdateBirthday, handler.UpdateBirthdayHandler())
		group.POST(path.UpdateAvatar, handler.UpdateAvatarHandler())
	}
}
