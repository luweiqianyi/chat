package entity

import (
	"fmt"
	"testing"
)

func TestNewTransaction(t *testing.T) {

	history := NewTransactionHistory()

	// 预估
	expectedSellPrice := 13.80
	t1 := NewTransaction("002436", Shenzhen, 14.70, expectedSellPrice, 1400)

	// 已完成 20230802
	fmt.Printf("=============\n\n")
	DealSellPrice := 13.70
	t2 := NewTransaction("002436", Shenzhen, 13.39, DealSellPrice, 300)

	// 已完成  20230802
	fmt.Printf("=============\n\n")
	expectBuyPrice := 13.80
	t3 := NewTransaction("002436", Shenzhen, expectBuyPrice, 13.71, 1400)

	expectedSellPrice = 13.80
	t4 := NewTransaction("002436", Shenzhen, 13.57, expectedSellPrice, 300)

	// ==============================================================================================
	// add to history
	history.AddTransaction(t1)
	history.AddTransaction(t2)
	history.AddTransaction(t3)
	history.AddTransaction(t4)

	history.PrintHistory()
}

func TestSellClass(t *testing.T) {
	profit := 199 * 2000
	fmt.Println(profit)

	fmt.Println(-1415.1374999999978 + 86.85824999999961 + 573.3595000000021) //- 754.9197499999962
}

func Test3200(t *testing.T) {
	fmt.Printf("3200: %v\n", 3200*1.05) // 3360
}
