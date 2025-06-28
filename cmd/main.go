package main

import (
	"ProductRecommendations/internal/api"
	"ProductRecommendations/internal/app"
	"ProductRecommendations/internal/data"
	"fmt"
	"github.com/labstack/echo/v4"
	"time"
)

func main() {
	newApp := app.NewApp(2 * time.Second)
	e := echo.New()
	newApi := api.NewApi(newApp, e)
	newLoader := data.NewLoader()

	products, err := newLoader.LoadProducts("./external/data/products.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	users, err := newLoader.LoadUsers("./external/data/users.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	newApp.RegisterProducts(products)
	newApp.RegisterUsers(users)
	newApi.RegisterApiFunctions()
}
