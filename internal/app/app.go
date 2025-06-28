package app

import (
	"ProductRecommendations/internal/model"
	"ProductRecommendations/internal/store"
	"time"
)

type App struct {
	Store     *store.Store
	Analytics *store.Analytics
}

func NewApp(ttlDuration time.Duration) *App {
	prodStore := store.NewStore()
	analyticsStore := store.NewAnalytics(prodStore, ttlDuration)
	return &App{
		Store:     prodStore,
		Analytics: analyticsStore,
	}
}

func (app *App) RegisterProducts(products []*model.Product) {
	for _, product := range products {
		_ = app.Store.AddProduct(product)
	}
}

func (app *App) RegisterUsers(users []*model.User) {
	for _, user := range users {
		_ = app.Store.AddUser(user)
	}
}
