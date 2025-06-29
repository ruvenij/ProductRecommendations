package analytics

import (
	"ProductRecommendations/internal/model"
	"ProductRecommendations/internal/service/activity"
	"ProductRecommendations/internal/store"
	"ProductRecommendations/internal/util"
	"github.com/labstack/gommon/log"
	"sort"
	"sync"
	"time"
)

type Recommendation struct {
	ExpiryAt        time.Time
	Recommendations []*model.Product
}

type Analytics struct {
	mu                  sync.RWMutex
	UserRecommendations map[string]*Recommendation
	UserCacheTtl        time.Duration
	Store               *store.Store
	UserActivityManager *activity.UserActivityManager
}

func NewAnalyticsHandler(store *store.Store, activityManager *activity.UserActivityManager, ttlDuration time.Duration) *Analytics {
	return &Analytics{
		UserRecommendations: make(map[string]*Recommendation),
		UserCacheTtl:        ttlDuration,
		Store:               store,
		UserActivityManager: activityManager,
	}
}

func (a *Analytics) GetRecommendationsForUser(userId string, productLimit int) ([]*model.Product, error) {
	recommendationsFromCache, err := a.readRecommendationsFromCache(userId, productLimit)
	if err != nil {
		return nil, err
	}

	if recommendationsFromCache != nil && len(recommendationsFromCache) > 0 {
		return recommendationsFromCache, nil
	}

	// else compute again and update the cache
	return a.computeRecommendations(userId, productLimit)
}

func (a *Analytics) readRecommendationsFromCache(userId string, productLimit int) ([]*model.Product, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// check whether ttl has expired
	if rec, ok := a.UserRecommendations[userId]; ok {
		if time.Now().Before(rec.ExpiryAt) {
			log.Debug("Current time is within the ttl active time. No need to compute recommendations")
			return rec.Recommendations, nil
		}
	}

	return []*model.Product{}, nil
}

func (a *Analytics) computeRecommendations(userId string, productLimit int) ([]*model.Product, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// find the most liked and most purchased categories for user id
	popularityIndexForCategories := a.getMostPopularCategoriesBasedOnUserInteraction(userId, util.CategoryLimitForPopularity)
	log.Debug("Popularity Index ", popularityIndexForCategories)

	// find the products for chosen categories
	popularityIndexForProducts := a.findProductsForCategories(userId, productLimit, popularityIndexForCategories)

	// get the product details based on the chosen product ids
	result := make([]*model.Product, 0)
	for _, pop := range popularityIndexForProducts {
		result = append(result, a.Store.GetProduct(pop.Id))
	}

	// add the recommendations to the cache, for future requests
	r := &Recommendation{
		ExpiryAt:        time.Now().Add(a.UserCacheTtl), // computed time changes to current time
		Recommendations: result,
	}
	a.UserRecommendations[userId] = r

	return result, nil
}

func (a *Analytics) getMostPopularCategoriesBasedOnUserInteraction(userId string, categoryCount int) []*model.CategoryScore {
	categoryIndex := make(map[string]*model.CategoryScore)
	popularityIndex := make([]*model.CategoryScore, 0)
	productCategories := a.Store.GetProductCategories()

	for _, category := range productCategories {
		// calculate popularity score
		popularityScore := getPopularityScore(&model.PopularityParams{
			PurchaseCount: a.UserActivityManager.GetActivityCountForUserAndCategory(util.PurchaseActivity, userId, category),
			WishlistCount: a.UserActivityManager.GetActivityCountForUserAndCategory(util.WishlistActivity, userId, category),
			ViewCount:     a.UserActivityManager.GetActivityCountForUserAndCategory(util.ViewActivity, userId, category),
		})

		if _, ok := categoryIndex[category]; !ok {
			p := &model.CategoryScore{
				Category: category,
			}
			categoryIndex[category] = p
			popularityIndex = append(popularityIndex, p)
		}

		categoryIndex[category].Popularity = popularityScore
	}

	// sort the categories based on the popularity
	sort.Slice(popularityIndex, func(i, j int) bool {
		if popularityIndex[i].Popularity > 0 || popularityIndex[j].Popularity > 0 {
			return popularityIndex[i].Popularity > popularityIndex[j].Popularity
		}

		// else sort by category
		return popularityIndex[i].Category < popularityIndex[j].Category
	})

	// get the most popular categories out of existing
	if len(popularityIndex) > categoryCount {
		popularityIndex = popularityIndex[:categoryCount]
	}

	return popularityIndex
}

// Calculate the popularity score, we can change this later if we have a better approach
func getPopularityScore(params *model.PopularityParams) int {
	return params.PurchaseCount*util.PurchaseCountScore +
		params.WishlistCount*util.WishlistCountScore +
		params.ViewCount*util.ViewCountScore
}

func (a *Analytics) findProductsForCategories(userId string, productCount int, popularCategories []*model.CategoryScore) []*model.ProductScore {
	popularityIndexForProducts := make([]*model.ProductScore, 0)
	for _, category := range popularCategories {
		productsForCategory := a.Store.GetProductsForCategory(category.Category)
		for _, prodId := range productsForCategory {

			// is this already purchased by user before, if so, ignore the product
			if a.UserActivityManager.IsUserAlreadyPurchasedProduct(userId, prodId) {
				log.Info("Ignored the product for category as user has purchased the product before ", category.Category, prodId)
				continue
			}

			// calculate popularity score
			popularityScore := getPopularityScore(&model.PopularityParams{
				PurchaseCount: a.UserActivityManager.GetActivityCountForProduct(util.PurchaseActivity, prodId),
				WishlistCount: a.UserActivityManager.GetActivityCountForProduct(util.WishlistActivity, prodId),
				ViewCount:     a.UserActivityManager.GetActivityCountForProduct(util.ViewActivity, prodId),
			})

			popForProduct := &model.ProductScore{
				Id:         prodId,
				Category:   category.Category,
				Popularity: popularityScore,
			}

			popularityIndexForProducts = append(popularityIndexForProducts, popForProduct)
		}
	}

	sort.Slice(popularityIndexForProducts, func(i, j int) bool {
		if popularityIndexForProducts[i].Popularity != 0 || popularityIndexForProducts[j].Popularity != 0 {
			return popularityIndexForProducts[i].Popularity > popularityIndexForProducts[j].Popularity
		}

		// sort based on the name since both values are zero
		return popularityIndexForProducts[i].Id < popularityIndexForProducts[j].Id

	})

	if len(popularityIndexForProducts) > productCount {
		popularityIndexForProducts = popularityIndexForProducts[:productCount]
	}

	return popularityIndexForProducts
}
