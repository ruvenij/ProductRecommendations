package activity

import (
	"ProductRecommendations/internal/model"
	"errors"
	"sync"
)

type WishlistActivityHandler struct {
	mu                                  sync.RWMutex
	transactionsByUserId                map[string]map[string]int // key1 - userid, key2 - product id, value - wishlist count
	transactionCountByUserIdAndCategory map[string]map[string]int // key1 - userid, key2 - category, value - wishlist count

	txnCountByProductId map[string]int // key - product id, value - total view count by all the users
}

func NewWishlistActivityHandler() *ViewActivityHandler {
	return &ViewActivityHandler{
		transactionsByUserId:                make(map[string]map[string]int),
		txnCountByProductId:                 make(map[string]int),
		transactionCountByUserIdAndCategory: make(map[string]map[string]int),
	}
}

func (w *WishlistActivityHandler) ProcessActivity(activity model.Activity) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	wAct, ok := activity.(*model.WishlistActivity)
	if !ok {
		return errors.New("Incompatible activity type received by wishlist activity handler ")
	}

	if _, exists := w.transactionsByUserId[wAct.UserId]; !exists {
		w.transactionsByUserId[wAct.UserId] = make(map[string]int)
	}
	w.transactionsByUserId[wAct.UserId][wAct.ProductId] += 1

	w.txnCountByProductId[wAct.ProductId] += 1

	if _, exists := w.transactionCountByUserIdAndCategory[wAct.UserId]; !exists {
		w.transactionCountByUserIdAndCategory[wAct.UserId] = make(map[string]int)
	}
	w.transactionCountByUserIdAndCategory[wAct.UserId][wAct.Category] += 1

	return nil
}

func (w *WishlistActivityHandler) GetActivityCountForUserAndCategory(userId string, category string) int {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if _, exists := w.transactionCountByUserIdAndCategory[userId]; !exists {
		return 0
	}

	return w.transactionCountByUserIdAndCategory[userId][category]
}

func (w *WishlistActivityHandler) GetActivityCountForProduct(productId string) int {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.txnCountByProductId[productId]
}
