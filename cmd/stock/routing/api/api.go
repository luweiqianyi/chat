package api

import (
	"chat/cmd/stock/auth"
	"chat/pkg/http/common"
)

type TransactionParam struct {
	auth.AccountParam

	StockID    string  `form:"stockID"`
	MarketType int     `form:"marketType"`
	BuyPrice   float64 `form:"buyPrice"`
	SellPrice  float64 `form:"sellPrice"`
	Number     int     `form:"number"`
}

type TransactionResp struct {
	common.ResponseHeader

	InvestIn float64 `json:"investIn"`
	BuyFee   float64 `json:"buyFee"`
	SellFee  float64 `json:"sellFee"`
	TotalFee float64 `json:"totalFee"`
	Profit   float64 `json:"profit"`
}
