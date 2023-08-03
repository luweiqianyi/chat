package entity

import "math"

const (
	CommissionRate  = 0.00025 // both
	StampDutyRate   = 0.001   // both
	TransferFeeRate = 0.00001 // ShangHai

	MinCommission = 5
	Eps           = 0.0000001

	ShangHai = 1
	Shenzhen = 2
)

type Transaction struct {
	StockID    string  `json:"stockID"`
	MarketType int     `json:"marketType"`
	BuyPrice   float64 `json:"buyPrice"`
	SellPrice  float64 `json:"sellPrice"`
	Number     int     `json:"number"`
}

func NewTransaction(stockID string, marketType int, buyPrice float64, sellPrice float64, number int) Transaction {
	return Transaction{
		StockID:    stockID,
		MarketType: marketType,
		BuyPrice:   buyPrice,
		SellPrice:  sellPrice,
		Number:     number,
	}
}

// InvestIn 总投入
func (t Transaction) InvestIn() float64 {
	return t.BuyPrice * float64(t.Number)
}

// CostPrice 成本价
func (t Transaction) CostPrice() float64 {
	return (t.BuyPrice*float64(t.Number) + t.BuyFee()) / float64(t.Number)
}

// FinalProfit 最终收益
func (t Transaction) FinalProfit() float64 {
	return (t.SellPrice-t.BuyPrice)*float64(t.Number) - t.TotalFee()
}

// PositionGainAndLoss 持仓盈亏(股票还未卖出时的预计盈亏)
func (t Transaction) PositionGainAndLoss(currentPrice float64) float64 {
	return (currentPrice - t.CostPrice()) * float64(t.Number)
}

// TotalFee 总费用
func (t Transaction) TotalFee() float64 {
	return t.BuyFee() + t.SellFee()
}

// BuyFee 买入费用
func (t Transaction) BuyFee() float64 {
	return t.BuyCommission() + t.BuyTransferFee()
}

// SellFee 卖出费用
func (t Transaction) SellFee() float64 {
	return t.SellCommission() + t.SellTransferFee() + t.StampDuty()
}

// BuyCommission 买入佣金
func (t Transaction) BuyCommission() float64 {
	fee := t.BuyPrice * float64(t.Number) * CommissionRate
	if MinCommission-fee > 0 && math.Abs(MinCommission-fee) < Eps {
		return MinCommission
	}
	return fee
}

// SellCommission 卖出佣金
func (t Transaction) SellCommission() float64 {
	fee := t.SellPrice * float64(t.Number) * CommissionRate
	if MinCommission-fee > 0 && math.Abs(MinCommission-fee) < Eps {
		return MinCommission
	}
	return fee
}

// StampDuty 印花税
func (t Transaction) StampDuty() float64 {
	fee := t.SellPrice * float64(t.Number) * StampDutyRate

	return fee
}

// BuyTransferFee 买入过户费
func (t Transaction) BuyTransferFee() float64 {
	if t.MarketType == Shenzhen {
		return 0.0
	}

	fee := t.BuyPrice * float64(t.Number) * TransferFeeRate
	return fee
}

// SellTransferFee 售出过户费
func (t Transaction) SellTransferFee() float64 {
	if t.MarketType == Shenzhen {
		return 0.0
	}

	fee := t.SellPrice * float64(t.Number) * TransferFeeRate
	return fee
}
