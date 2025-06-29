package app

import (
	"ProductRecommendations/internal/model"
	"ProductRecommendations/internal/util"
	"errors"
)

func (app *App) ProcessActivity(activity *model.BaseActivity) error {
	product := app.Store.GetProduct(activity.ProductId)
	if product == nil {
		return errors.New("product not found")
	}

	var err error
	switch activity.ActivityType {
	case util.PurchaseActivityType:
		err = app.UserActivityManager.AddPurchaseActivity(&model.PurchaseActivity{
			Id:        activity.Id,
			UserId:    activity.UserId,
			ProductId: activity.ProductId,
			Category:  product.Category,
			Quantity:  activity.Quantity,
			Price:     activity.Price,
		})
	case util.ViewActivityType:
		err = app.UserActivityManager.AddViewActivity(&model.ViewActivity{
			Id:        activity.Id,
			UserId:    activity.UserId,
			ProductId: activity.ProductId,
			Category:  product.Category,
		})
	case util.WishlistActivityType:
		err = app.UserActivityManager.AddWishlistActivity(&model.WishlistActivity{
			Id:        activity.Id,
			UserId:    activity.UserId,
			ProductId: activity.ProductId,
			Category:  product.Category,
		})
	default:
		return errors.New("activity type not supported")
	}

	return err
}
