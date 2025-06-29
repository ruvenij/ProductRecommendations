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
	result, err := a.app.Analytics.GetRecommendationsForUser(userId, 10)
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

	ba := &model.BaseActivity{
		ProductId:    req.ProductId,
		UserId:       req.UserId,
		ActivityType: req.Type,
	}

	if req.Price != nil {
		ba.Price = *req.Price
	}

	if req.Quantity != nil {
		ba.Quantity = *req.Quantity
	}

	err := a.app.ProcessActivity(ba)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, req)
}
