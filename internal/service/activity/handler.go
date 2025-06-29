package activity

import (
	"ProductRecommendations/internal/model"
)

type Handler interface {
	ProcessActivity(activity model.Activity) error
	GetActivityCountForUserAndCategory(userId string, category string) int
	GetActivityCountForProduct(productId string) int
}
