package api

import (
	"ProductRecommendations/internal/app"
	"ProductRecommendations/internal/model"
	"ProductRecommendations/internal/util"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
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
	limit := c.QueryParam("limit")

	log.Info("Incoming recommendation request for user : ", userId, " , limit : ", limit)

	// validate the request first
	limitInt, err := a.validateRecommendationParams(userId, limit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// get the recommendations for the params
	result, err := a.app.GetRecommendationsForUser(userId, limitInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	log.Debug("Get Recommendations result : ", result)
	return c.JSON(http.StatusOK, result)
}

func (a *Api) validateRecommendationParams(userId string, limit string) (int, error) {
	if userId == "" {
		return -1, errors.New("User_id is required ")
	}

	if !a.app.IsValidUser(userId) {
		return -1, errors.New("Invalid user_id received ")
	}

	limitInt := util.RecommendationLimit
	var err error
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return -1, err
		}
	}

	return limitInt, nil
}

func (a *Api) AddActivityForUser(c echo.Context) error {
	req := new(Activity)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New("Invalid incoming json body for activity "))
	}

	log.Info("Incoming activity request for user : ", req.UserId, ", Content : ", req)

	// validate the request first
	err := a.validateActivityParams(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ba := &model.BaseActivity{
		ProductId:    req.ProductId,
		UserId:       req.UserId,
		ActivityType: req.Type,
	}

	if req.Price != nil {
		ba.Price, _ = decimal.NewFromString(*req.Price)
	}

	if req.Quantity != nil {
		ba.Quantity = *req.Quantity
	}

	// process the activity
	err = a.app.ProcessActivity(ba)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, req)
}

func (a *Api) validateActivityParams(req *Activity) error {
	if req.UserId == "" {
		return errors.New("UserId is required ")
	}

	if req.ProductId == "" {
		return errors.New("ProductId is required ")
	}

	if req.Type == "" {
		return errors.New("Type is required ")
	}

	if !a.app.IsValidProduct(req.ProductId) {
		return errors.New("Received invalid product id ")
	}

	if !a.app.IsValidUser(req.UserId) {
		return errors.New("Received invalid user id ")
	}

	if req.Type != util.PurchaseActivityType && req.Type != util.ViewActivityType &&
		req.Type != util.WishlistActivityType {
		return errors.New("Received invalid activity type ")
	}

	if req.Type == util.PurchaseActivityType {
		if req.Price == nil || req.Quantity == nil {
			return errors.New("Received invalid price or quantity for purchase activity type ")
		}

		price, err := decimal.NewFromString(*req.Price)
		if err != nil {
			return err
		}

		if lo.FromPtr(req.Quantity) == 0 {
			return errors.New("Received zero quantity for purchase activity type ")
		}

		if price == decimal.Zero {
			return errors.New("Received zero price for purchase activity type ")
		}
	}

	return nil
}
