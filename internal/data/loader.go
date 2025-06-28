package data

import (
	"ProductRecommendations/internal/model"
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
)

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) LoadProducts(productsFilePath string) ([]*model.Product, error) {
	file, err := os.Open(productsFilePath)
	if err != nil {
		return []*model.Product{}, err
	}

	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	_, err = reader.Read()
	if err != nil {
		return []*model.Product{}, err
	}

	result := make([]*model.Product, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error occurred while reading csv, error : ", err)
			continue
		}

		product := parseProductRecord(record)
		result = append(result, product)
	}

	return result, nil
}

func (l *Loader) LoadUsers(usersFilePath string) ([]*model.User, error) {
	file, err := os.Open(usersFilePath)
	if err != nil {
		return []*model.User{}, err
	}

	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	_, err = reader.Read()
	if err != nil {
		return []*model.User{}, err
	}

	result := make([]*model.User, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error occurred while reading csv, error : ", err)
			continue
		}

		user := parseUserRecord(record)
		result = append(result, user)
	}

	return result, nil
}

func parseProductRecord(record []string) *model.Product {
	return &model.Product{
		Id:       record[0],
		Name:     record[1],
		Category: record[2],
	}
}

func parseUserRecord(record []string) *model.User {
	return &model.User{
		Id:   record[0],
		Name: record[1],
	}
}
