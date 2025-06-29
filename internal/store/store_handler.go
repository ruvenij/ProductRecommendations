package store

import (
	"ProductRecommendations/internal/model"
)

type Store struct {
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

func (s *Store) AddProduct(prod *model.Product) error {
	s.products[prod.Id] = prod
	s.productsForCategory[prod.Category] = append(s.productsForCategory[prod.Category], prod.Id)

	return nil
}

func (s *Store) AddUser(user *model.User) error {
	s.users[user.Id] = user
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

func (s *Store) GetProduct(id string) *model.Product {
	return s.products[id]
}

func (s *Store) GetProductCategories() []string {
	categories := make([]string, 0)
	for category, _ := range s.productsForCategory {
		categories = append(categories, category)
	}

	return categories
}

func (s *Store) GetProductsForCategory(category string) []string {
	return s.productsForCategory[category]
}
