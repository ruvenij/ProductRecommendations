package store

import (
	"ProductRecommendations/internal/model"
	"github.com/labstack/gommon/log"
	"sync"
)

type Store struct {
	mu                  sync.RWMutex
	products            map[string]*model.Product
	users               map[string]*model.User
	productsForCategory map[string][]string // key - category, value - product ids belong to the category
}

func NewStore() *Store {
	return &Store{
		products:            make(map[string]*model.Product),
		users:               make(map[string]*model.User),
		productsForCategory: make(map[string][]string),
	}
}

func (s *Store) AddProduct(prod *model.Product) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.products[prod.Id] = prod
	s.productsForCategory[prod.Category] = append(s.productsForCategory[prod.Category], prod.Id)
	log.Debugf("Added product %s to store", prod.Id)
}

func (s *Store) AddUser(user *model.User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[user.Id] = user
	log.Debugf("Added user %s to store", user.Id)
}

func (s *Store) IsValidProduct(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if product, exists := s.products[id]; !exists || product == nil || product.Id == "" {
		return false
	}

	return true
}

func (s *Store) IsValidUser(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if user, exists := s.users[id]; !exists || user == nil || user.Id == "" {
		return false
	}

	return true
}

func (s *Store) GetProduct(id string) *model.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.products[id]
}

func (s *Store) GetProductCategories() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	categories := make([]string, 0)
	for category, _ := range s.productsForCategory {
		categories = append(categories, category)
	}

	return categories
}

func (s *Store) GetProductsForCategory(category string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.productsForCategory[category]
}
