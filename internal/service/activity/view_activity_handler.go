package activity

import (
	"ProductRecommendations/internal/model"
	"errors"
)

type ViewActivityHandler struct {
	transactionsByUserId                map[string]map[string]int // key1 - userid, key2 - product id, value - view count
	transactionCountByUserIdAndCategory map[string]map[string]int // key1 - userid, key2 - category, value - view count

	txnCountByProductId map[string]int // key - product id, value - total view count by all the users
}

func NewViewActivityHandler() *ViewActivityHandler {
	return &ViewActivityHandler{
		transactionsByUserId:                make(map[string]map[string]int),
		txnCountByProductId:                 make(map[string]int),
		transactionCountByUserIdAndCategory: make(map[string]map[string]int),
	}
}

func (v *ViewActivityHandler) ProcessActivity(activity model.Activity) error {
	vAct, ok := activity.(*model.ViewActivity)
	if !ok {
		return errors.New("Incompatible activity type received by view activity handler ")
	}

	if _, exists := v.transactionsByUserId[vAct.UserId]; !exists {
		v.transactionsByUserId[vAct.UserId] = make(map[string]int)
	}
	v.transactionsByUserId[vAct.UserId][vAct.ProductId] += 1

	v.txnCountByProductId[vAct.ProductId] += 1

	if _, exists := v.transactionCountByUserIdAndCategory[vAct.UserId]; !exists {
		v.transactionCountByUserIdAndCategory[vAct.UserId] = make(map[string]int)
	}
	v.transactionCountByUserIdAndCategory[vAct.UserId][vAct.Category] += 1

	return nil
}

func (v *ViewActivityHandler) GetActivityCountForUserAndCategory(userId string, category string) int {
	if _, exists := v.transactionCountByUserIdAndCategory[userId]; !exists {
		return 0
	}

	return v.transactionCountByUserIdAndCategory[userId][category]
}

func (v *ViewActivityHandler) GetActivityCountForProduct(productId string) int {
	return v.txnCountByProductId[productId]
}
