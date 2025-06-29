package model

import "github.com/shopspring/decimal"

type Purchase struct {
	Id        int             `json:"id"`
	Quantity  int             `json:"quantity"`
	Price     decimal.Decimal `json:"price"`
	ProductId string          `json:"product_id"`
}
