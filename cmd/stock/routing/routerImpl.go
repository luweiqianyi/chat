package routing

import (
	"chat/cmd/stock/routing/handler"
	"chat/cmd/stock/routing/path"
	"github.com/gin-gonic/gin"
)

type RouterImpl struct {
}

func (impl RouterImpl) RegisterRouters(eg *gin.Engine) {
	g := eg.Group(path.StockPath)
	g.POST(path.CalTransaction, handler.CalTransactionHandler())
}
