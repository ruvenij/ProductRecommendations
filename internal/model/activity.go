package model

import "github.com/shopspring/decimal"

type Activity struct {
	Id        int
	UserId    string
	ProductId string
	Quantity  int
	Price     decimal.Decimal
}
