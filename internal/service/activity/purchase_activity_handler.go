package activity

import (
	"ProductRecommendations/internal/model"
	"errors"
	"sync"
)

type PurchaseActivityHandler struct {
	mu                                  sync.RWMutex
	transactionsByUserId                map[string]map[string][]*model.Purchase // key1 - userid, key2 - product id, value - slice of purchases
	transactionCountByUserIdAndCategory map[string]map[string]int               // key1 - userid, key2 - category, value - txn count

	txnCountByProductId map[string]int // key - product id, value - total purchase count by all the users
}

func NewPurchaseActivityHandler() *PurchaseActivityHandler {
	return &PurchaseActivityHandler{
		transactionsByUserId:                make(map[string]map[string][]*model.Purchase),
		txnCountByProductId:                 make(map[string]int),
		transactionCountByUserIdAndCategory: make(map[string]map[string]int),
	}
}

func (p *PurchaseActivityHandler) ProcessActivity(activity model.Activity) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pAct, ok := activity.(*model.PurchaseActivity)
	if !ok {
		return errors.New("Incompatible activity type received by purchase activity handler ")
	}

	if _, exists := p.transactionsByUserId[pAct.UserId]; !exists {
		p.transactionsByUserId[pAct.UserId] = make(map[string][]*model.Purchase)
	}

	if _, exists := p.transactionsByUserId[pAct.UserId][pAct.ProductId]; !exists {
		p.transactionsByUserId[pAct.UserId][pAct.ProductId] = make([]*model.Purchase, 0)
	}

	p.transactionsByUserId[pAct.UserId][pAct.ProductId] =
		append(p.transactionsByUserId[pAct.UserId][pAct.ProductId], &model.Purchase{
			Quantity:  pAct.Quantity,
			Price:     pAct.Price,
			ProductId: pAct.ProductId,
		})

	p.txnCountByProductId[pAct.ProductId] += 1

	if _, exists := p.transactionCountByUserIdAndCategory[pAct.UserId]; !exists {
		p.transactionCountByUserIdAndCategory[pAct.UserId] = make(map[string]int)
	}

	p.transactionCountByUserIdAndCategory[pAct.UserId][pAct.Category] += 1
	return nil
}

func (p *PurchaseActivityHandler) GetActivityCountForUserAndCategory(userId string, category string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if _, exists := p.transactionCountByUserIdAndCategory[userId]; !exists {
		return 0
	}

	return p.transactionCountByUserIdAndCategory[userId][category]
}

func (p *PurchaseActivityHandler) GetActivityCountForProduct(productId string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.txnCountByProductId[productId]
}

func (p *PurchaseActivityHandler) IsUserAlreadyPurchasedProduct(userId string, productId string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if _, ok := p.transactionsByUserId[userId]; !ok {
		return false
	}

	if _, ok := p.transactionsByUserId[userId][productId]; !ok {
		return false
	}

	return true
}
