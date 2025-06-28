package api

import (
	"ProductRecommendations/internal/app"
	"ProductRecommendations/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Api struct {
	app  *app.App
	echo *echo.Echo
}

func NewApi(app *app.App, e *echo.Echo) *Api {
	return &Api{
		app:  app,
		echo: e,
	}
}

func (a *Api) RegisterApiFunctions() {
	a.echo.GET("/api/recommendations", a.GetRecommendationsForUser)
	a.echo.POST("/api/activity", a.AddActivityForUser)

	a.echo.Logger.Fatal(a.echo.Start(":8080"))
}

func (a *Api) GetRecommendationsForUser(c echo.Context) error {
	userId := c.QueryParam("user_id")
	result, err := a.app.Analytics.GetRecommendationsForUser(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (a *Api) AddActivityForUser(c echo.Context) error {
	req := new(Activity)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON body"})
	}

	if req.UserId == "" || req.ProductId == "" || req.Type == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
	}

	v := &model.Activity{
		ProductId: req.ProductId,
		UserId:    req.UserId,
	}

	var err error
	switch req.Type {
	case "view":
		err = a.app.Store.AddView(v)
	case "purchase":
		v.Price = *req.Price
		v.Quantity = *req.Quantity
		err = a.app.Store.AddPurchase(v)
	case "wishlist":
		err = a.app.Store.AddWishlist(v)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid type"})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, req)
}
