package model

import "github.com/shopspring/decimal"

type Purchase struct {
	Id       int             `json:"id"`
	Product  *Product        `json:"product"`
	Quantity int             `json:"quantity"`
	Price    decimal.Decimal `json:"price"`
}
