package service

import (
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/repository"
)

type ProductService interface {
	Get(id string) model.Product
}

type ProductServiceImpl struct {
	Repository repository.ProductRepository
}

func (service ProductServiceImpl) Get(id string) model.Product {
	return service.Repository.Get(id)
}

func NewProductServiceImpl(repository repository.ProductRepository) ProductService {
	return ProductServiceImpl{repository}
}
