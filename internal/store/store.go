package store

import (
	"ProductRecommendations/internal/model"
	"errors"
	"fmt"
)

type Store struct {
	products            map[string]*model.Product
	users               map[string]*model.User
	productsForCategory map[string][]string // key - category, value - product id slice

	// by user
	viewsByUserId     map[string]map[string]int               // key1 - userid, key2 - product id, value - view count
	purchasesByUserId map[string]map[string][]*model.Purchase // key1 - userid, key2 - product id, value - slice of purchases
	wishlistByUserId  map[string]map[string]struct{}          // key1 - userid, key2 - product id

	categoryViewsByUserId    map[string]map[string]int // key1 - userid, key2 - category, value - view count
	categoryTxnCountByUserId map[string]map[string]int // key1 - userid, key2 - category, value - txn count
	categoryWishListByUserId map[string]map[string]int // key1 - userid, key2 - category, value - wishlist count

	// by product
	purchasesTxnCountByProductId map[string]int
	viewsByProductId             map[string]int
	wishlistQtyByProductId       map[string]int
}

func NewStore() *Store {
	return &Store{
		products:                     make(map[string]*model.Product),
		users:                        make(map[string]*model.User),
		viewsByUserId:                make(map[string]map[string]int),
		purchasesByUserId:            make(map[string]map[string][]*model.Purchase),
		viewsByProductId:             make(map[string]int),
		wishlistByUserId:             make(map[string]map[string]struct{}),
		purchasesTxnCountByProductId: make(map[string]int),
		categoryViewsByUserId:        make(map[string]map[string]int),
		categoryTxnCountByUserId:     make(map[string]map[string]int),
		categoryWishListByUserId:     make(map[string]map[string]int),
		productsForCategory:          make(map[string][]string),
		wishlistQtyByProductId:       make(map[string]int),
	}
}

func (s *Store) AddProduct(prod *model.Product) error {
	s.products[prod.Id] = prod
	s.productsForCategory[prod.Category] = append(s.productsForCategory[prod.Category], prod.Id)

	return nil
}

func (s *Store) AddUser(user *model.User) error {
	s.users[user.Id] = user
	return nil
}

func (s *Store) AddView(view *model.Activity) error {
	fmt.Println("Adding view")
	if !s.IsValidUser(view.UserId) || !s.IsValidProduct(view.ProductId) {
		return errors.New("invalid user or/or purchases")
	}

	product := s.products[view.ProductId]

	if _, exists := s.viewsByUserId[view.UserId]; !exists {
		s.viewsByUserId[view.UserId] = make(map[string]int)
	}

	if _, exists := s.categoryTxnCountByUserId[view.UserId]; !exists {
		s.categoryTxnCountByUserId[view.UserId] = make(map[string]int)
	}

	s.viewsByProductId[view.ProductId] += 1
	s.viewsByUserId[view.UserId][view.ProductId] += 1
	s.categoryTxnCountByUserId[view.UserId][product.Category] += 1
	return nil
}

func (s *Store) AddWishlist(view *model.Activity) error {
	if !s.IsValidUser(view.UserId) || !s.IsValidProduct(view.ProductId) {
		return errors.New("invalid user or/or purchases")
	}

	if _, exists := s.wishlistByUserId[view.UserId]; !exists {
		s.wishlistByUserId[view.UserId] = make(map[string]struct{})
	}

	if _, exists := s.categoryWishListByUserId[view.UserId]; !exists {
		s.categoryWishListByUserId[view.UserId] = make(map[string]int)
	}

	product := s.products[view.ProductId]
	s.categoryWishListByUserId[view.UserId][product.Category] += 1
	s.wishlistByUserId[view.UserId][view.ProductId] = struct{}{}
	s.wishlistQtyByProductId[view.ProductId] += 1
	return nil
}

func (s *Store) AddPurchase(view *model.Activity) error {
	if !s.IsValidUser(view.UserId) || !s.IsValidProduct(view.ProductId) {
		return errors.New("invalid user or/or purchases")
	}

	if view.Price.IsZero() {
		return errors.New("invalid price")
	}

	if view.Quantity <= 0 {
		return errors.New("invalid quantity")
	}

	product := s.products[view.ProductId]

	if _, exists := s.purchasesByUserId[view.UserId]; !exists {
		s.purchasesByUserId[view.UserId] = make(map[string][]*model.Purchase)
	}

	if _, exists := s.purchasesByUserId[view.UserId][view.ProductId]; !exists {
		s.purchasesByUserId[view.UserId][view.ProductId] = make([]*model.Purchase, 0)
	}

	s.purchasesByUserId[view.UserId][view.ProductId] = append(s.purchasesByUserId[view.UserId][view.ProductId], &model.Purchase{
		Product:  product,
		Quantity: view.Quantity,
		Price:    view.Price,
	})

	s.purchasesTxnCountByProductId[view.ProductId] += 1

	return nil
}

func (s *Store) IsValidProduct(id string) bool {
	if product, exists := s.products[id]; !exists || product == nil || product.Id == "" {
		return false
	}

	return true
}

func (s *Store) IsValidUser(id string) bool {
	if user, exists := s.users[id]; !exists || user == nil || user.Id == "" {
		return false
	}

	return true
}
