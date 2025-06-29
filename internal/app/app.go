package app

import (
	"ProductRecommendations/internal/model"
	"ProductRecommendations/internal/service/activity"
	"ProductRecommendations/internal/service/analytics"
	"ProductRecommendations/internal/store"
	"time"
)

type App struct {
	Store               *store.Store
	Analytics           *analytics.Analytics
	UserActivityManager *activity.UserActivityManager
}

func NewApp(ttlDuration time.Duration) *App {
	prodStore := store.NewStore()
	activityManager := activity.NewUserActivityManager()
	analyticsStore := analytics.NewAnalyticsHandler(prodStore, activityManager, ttlDuration)
	return &App{
		Store:               prodStore,
		Analytics:           analyticsStore,
		UserActivityManager: activityManager,
	}
}

func (app *App) RegisterProducts(products []*model.Product) {
	for _, product := range products {
		app.Store.AddProduct(product)
	}
}

func (app *App) RegisterUsers(users []*model.User) {
	for _, user := range users {
		app.Store.AddUser(user)
	}
}

func (app *App) IsValidUser(userId string) bool {
	return app.Store.IsValidUser(userId)
}

func (app *App) IsValidProduct(productId string) bool {
	return app.Store.IsValidProduct(productId)
}

func (app *App) GetRecommendationsForUser(userId string, productLimit int) ([]*model.Product, error) {
	return app.Analytics.GetRecommendationsForUser(userId, productLimit)
}
