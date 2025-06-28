package api

import "github.com/shopspring/decimal"

type Activity struct {
	UserId    string           `json:"user_id"`
	Type      string           `json:"type"`
	ProductId string           `json:"product_id"`
	Quantity  *int             `json:"quantity"`
	Price     *decimal.Decimal `json:"price"`
}
