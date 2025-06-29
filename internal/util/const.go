package util

const (
	CategoryLimitForPopularity = 4

	PurchaseCountScore = 3
	WishlistCountScore = 2
	ViewCountScore     = 1
)

type ActivityType int

const (
	PurchaseActivity ActivityType = iota
	ViewActivity
	WishlistActivity
)

const (
	PurchaseActivityType = "purchase"
	ViewActivityType     = "view"
	WishlistActivityType = "wishlist"
)

const (
	RecommendationLimit = 10 // make this a config later
)
