package routing

import (
	"chat/cmd/account/routing/handler"
	"chat/cmd/account/routing/path"
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

	group := eg.Group(path.Account)
	{
		group.POST(path.Register, handler.RegisterHandler())
		group.GET(path.Unregister, handler.UnRegisterHandler())
		group.POST(path.Login, handler.LoginHandler())
		group.GET(path.Logout, handler.LogoutHandler())
	}

}
