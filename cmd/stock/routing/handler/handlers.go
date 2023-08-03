package handler

import (
	"chat/cmd/stock/entity"
	"chat/cmd/stock/routing/api"
	"chat/pkg/http/common"
	"github.com/gin-gonic/gin"
)

func CalTransactionHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var param api.TransactionParam
		err := context.ShouldBind(&param)
		if err != nil {
			context.JSON(200, common.NewParameterErrorResponse())
			return
		}

		transaction := entity.NewTransaction(
			param.StockID,
			param.MarketType,
			param.BuyPrice,
			param.SellPrice,
			param.Number)

		context.JSON(200, api.TransactionResp{
			ResponseHeader: common.NewSuccessResponse(),
			InvestIn:       transaction.InvestIn(),
			BuyFee:         transaction.BuyFee(),
			SellFee:        transaction.SellFee(),
			TotalFee:       transaction.TotalFee(),
			Profit:         transaction.FinalProfit(),
		})

	}
}
