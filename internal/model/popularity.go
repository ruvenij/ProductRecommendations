package model

type Popularity struct {
	Id         string
	Category   string
	Popularity int
}

type PopularityParams struct {
	PurchaseCount int
	WishlistCount int
	ViewCount     int
}
