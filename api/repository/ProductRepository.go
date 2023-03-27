package repository

import "online_fashion_shop/api/model"

type ProductRepository interface {
	Get(id string) model.Product
}

type MockProductRepository struct {
}

func (repository MockProductRepository) Get(id string) model.Product {
	return model.Product{
		Name:  "mock_name",
		Color: "mock_color",
		Size:  "mock_size",
	}
}

func NewMockProductRepository() ProductRepository {
	return MockProductRepository{}
}
