package model

type CategoryScore struct {
	Category   string
	Popularity int
}

type ProductScore struct {
	Id         string
	Popularity int
	Category   string
}

type PopularityParams struct {
	PurchaseCount int
	WishlistCount int
	ViewCount     int
}
