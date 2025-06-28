package store

import (
	"ProductRecommendations/internal/model"
	"fmt"
	"sort"
	"time"
)

type Recommendation struct {
	ComputedTime    time.Time
	Recommendations []*model.Product
}

type Analytics struct {
	Store               *Store
	UserRecommendations map[string]*Recommendation
	UserCacheTtl        time.Duration
}

type Popularity struct {
	Id         string
	Category   string
	Popularity int
}

func NewAnalytics(store *Store, ttlDuration time.Duration) *Analytics {
	return &Analytics{
		Store:               store,
		UserRecommendations: make(map[string]*Recommendation),
		UserCacheTtl:        ttlDuration,
	}
}

func (a *Analytics) GetRecommendationsForUser(userId string) ([]*model.Product, error) {
	// check whether ttl has expired
	if rec, ok := a.UserRecommendations[userId]; ok {
		if time.Since(rec.ComputedTime).Seconds() < a.UserCacheTtl.Seconds() {
			return rec.Recommendations, nil
		}
	}

	// else compute again and update the cache
	return a.ComputeRecommendations(userId)
}

func (a *Analytics) ComputeRecommendations(userId string) ([]*model.Product, error) {
	// find the most liked and most purchased categories for user id
	categoryIndex := make(map[string]*Popularity)
	popularityIndex := make([]*Popularity, 0)
	categoryPurchasesForUser := a.Store.categoryTxnCountByUserId[userId]
	categoryWishlistForUser := a.Store.categoryWishListByUserId[userId]
	categoryViewsForUser := a.Store.categoryTxnCountByUserId[userId]

	for category, _ := range a.Store.productsForCategory {
		popularity := 0
		purchaseCount := categoryPurchasesForUser[category]
		wishlistCount := categoryWishlistForUser[category]
		viewCount := categoryViewsForUser[category]

		popularity = purchaseCount*3 + wishlistCount*2 + viewCount*1
		if popularity > 0 {
			if _, ok := categoryIndex[category]; !ok {
				p := &Popularity{
					Category: category,
				}
				categoryIndex[category] = p
				popularityIndex = append(popularityIndex, p)
			}

			categoryIndex[category].Popularity = popularity
		}
	}

	sort.Slice(popularityIndex, func(i, j int) bool {
		return popularityIndex[i].Popularity > popularityIndex[j].Popularity
	})

	// top four
	if len(popularityIndex) > 4 {
		popularityIndex = popularityIndex[:4]
	}

	fmt.Println("Popularity Index ", popularityIndex)

	// find the products for chosen categories
	userPurchases := a.Store.purchasesByUserId[userId]
	popularityIndexForProducts := make([]*Popularity, 0)
	for _, pop := range popularityIndex {
		for _, prodId := range a.Store.productsForCategory[pop.Category] {

			// is this already purchased by user before, if so, ignore the product
			if userPurchases != nil {
				if _, ok := userPurchases[prodId]; ok {
					fmt.Println("Ignored the product for category ", pop.Category, prodId)
					continue
				}
			}

			popularity := 0
			purchaseCount := a.Store.purchasesTxnCountByProductId[prodId]
			viewCount := a.Store.viewsByProductId[prodId]
			wishlistCount := a.Store.wishlistQtyByProductId[prodId]

			popularity = purchaseCount*3 + wishlistCount*2 + viewCount*1
			popForProduct := &Popularity{
				Id:         prodId,
				Category:   pop.Category,
				Popularity: popularity,
			}

			popularityIndexForProducts = append(popularityIndexForProducts, popForProduct)
		}
	}

	sort.Slice(popularityIndexForProducts, func(i, j int) bool {
		if popularityIndexForProducts[i].Popularity != 0 || popularityIndexForProducts[j].Popularity != 0 {
			return popularityIndexForProducts[i].Popularity > popularityIndexForProducts[j].Popularity
		}

		// sort based on the name since both values are zero
		return popularityIndexForProducts[0].Id > popularityIndexForProducts[j].Id

	})

	if len(popularityIndexForProducts) > 10 {
		popularityIndexForProducts = popularityIndexForProducts[:10]
	}

	result := make([]*model.Product, 0)
	for _, pop := range popularityIndexForProducts {
		result = append(result, a.Store.products[pop.Id])
	}

	// add the recommendations to the cache
	r := &Recommendation{
		ComputedTime:    time.Now(),
		Recommendations: result,
	}
	a.UserRecommendations[userId] = r

	return result, nil
}
