package entity

import "fmt"

type TransactionHistory struct {
	curPos  int
	History map[int]Transaction
}

func NewTransactionHistory() *TransactionHistory {
	return &TransactionHistory{
		curPos:  0,
		History: make(map[int]Transaction),
	}
}

func (h *TransactionHistory) AddTransaction(transaction Transaction) {
	h.History[h.curPos] = transaction
	h.curPos++
}

func (h *TransactionHistory) ClearAll() {
	h.History = make(map[int]Transaction)
	h.curPos = 0
}

func (h *TransactionHistory) ToTalFee() float64 {
	fee := 0.0
	for _, transaction := range h.History {
		fee = fee + transaction.TotalFee()
	}
	return fee
}

func (h *TransactionHistory) FinalProfit() float64 {
	finalProfit := 0.0
	for _, transaction := range h.History {
		finalProfit = finalProfit + transaction.FinalProfit()
	}
	return finalProfit
}

func (h *TransactionHistory) PrintHistory() {
	for key, transaction := range h.History {
		fmt.Printf("transaction[%v] info: {buyPrice=%v sellPrice=%v number=%v}\n", key, transaction.BuyPrice, transaction.SellPrice, transaction.Number)
		fmt.Printf("	invest in: %f\n", transaction.InvestIn())
		fmt.Printf("	buy fee: %v\n", transaction.BuyFee())
		fmt.Printf("	sell fee: %v\n", transaction.SellFee())
		fmt.Printf("	total fee: %v\n", transaction.TotalFee())
		fmt.Printf("	costPrice: %v\n", transaction.CostPrice())
		fmt.Printf("	※※※ final profit: %v\n\n", transaction.FinalProfit())
	}

	fmt.Printf("all transactions total fee: %v\n", h.ToTalFee())
	fmt.Printf("all transactions final profit: %v\n\n", h.FinalProfit())
}
