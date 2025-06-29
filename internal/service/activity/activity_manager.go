package activity

import (
	"ProductRecommendations/internal/model"
	"ProductRecommendations/internal/util"
)

type UserActivityManager struct {
	ActivityHandlers map[util.ActivityType]Handler
}

func NewUserActivityManager() *UserActivityManager {
	handlers := make(map[util.ActivityType]Handler)
	handlers[util.PurchaseActivity] = NewPurchaseActivityHandler()
	handlers[util.ViewActivity] = NewViewActivityHandler()
	handlers[util.WishlistActivity] = NewWishlistActivityHandler()

	return &UserActivityManager{
		ActivityHandlers: handlers,
	}
}

func (ua *UserActivityManager) AddViewActivity(activity *model.ViewActivity) error {
	return ua.ActivityHandlers[util.ViewActivity].ProcessActivity(activity)
}

func (ua *UserActivityManager) AddPurchaseActivity(activity *model.PurchaseActivity) error {
	return ua.ActivityHandlers[util.PurchaseActivity].ProcessActivity(activity)
}

func (ua *UserActivityManager) AddWishlistActivity(activity *model.WishlistActivity) error {
	return ua.ActivityHandlers[util.WishlistActivity].ProcessActivity(activity)
}

func (ua *UserActivityManager) GetActivityCountForUserAndCategory(activityType util.ActivityType, userId string, category string) int {
	return ua.ActivityHandlers[activityType].GetActivityCountForUserAndCategory(userId, category)
}

func (ua *UserActivityManager) GetActivityCountForProduct(activityType util.ActivityType, productId string) int {
	return ua.ActivityHandlers[activityType].GetActivityCountForProduct(productId)
}

func (ua *UserActivityManager) IsUserAlreadyPurchasedProduct(userId string, productId string) bool {
	activity := ua.ActivityHandlers[util.PurchaseActivity].(*PurchaseActivityHandler)
	return activity.IsUserAlreadyPurchasedProduct(userId, productId)
}
