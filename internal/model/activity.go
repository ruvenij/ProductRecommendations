package model

import "github.com/shopspring/decimal"

type Activity interface {
	GetActivityType() string
}

type BaseActivity struct {
	Id           int
	ActivityType string
	UserId       string
	ProductId    string
	Category     string
	Quantity     int
	Price        decimal.Decimal
}
type PurchaseActivity struct {
	Id        int
	UserId    string
	ProductId string
	Category  string
	Quantity  int
	Price     decimal.Decimal
}

func (p *PurchaseActivity) GetActivityType() string {
	return "purchase"
}

type ViewActivity struct {
	Id        int
	UserId    string
	ProductId string
	Category  string
}

func (p *ViewActivity) GetActivityType() string {
	return "view"
}

type WishlistActivity struct {
	Id        int
	UserId    string
	ProductId string
	Category  string
}

func (p *WishlistActivity) GetActivityType() string {
	return "wishlist"
}
