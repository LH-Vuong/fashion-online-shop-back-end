package service

import (
	"fmt"
	"online_fashion_shop/api/model/product"
	"online_fashion_shop/api/repository"
)

type ProductQuantityService interface {
	// Create if detail of product is not existed will return error, if product_quantity is existed will throw return error
	Create(newQuantity product.ProductQuantity) (*product.ProductQuantity, error)
	// if product_quantity is existed will throw return error
	Update(updateQuantity product.ProductQuantity) (*product.ProductQuantity, error)
	//id is product_quantity_id
	DeleteOne(id string) error
	//delete all product_quantity has same detail_id
	DeleteManyByDetailId(detailId string) error

	ListByDetailId(detailId string) ([]*product.ProductQuantity, error)

	Get(id string) (*product.ProductQuantity, error)
}

type ProductQuantityServiceImpl struct {
	QuantityRepo repository.ProductQuantityRepository
}

func NewProductQuantityServiceImpl(quantityRepo repository.ProductQuantityRepository) ProductQuantityService {
	return &ProductQuantityServiceImpl{
		QuantityRepo: quantityRepo,
	}
}

func (s *ProductQuantityServiceImpl) Get(id string) (*product.ProductQuantity, error) {
	return s.QuantityRepo.Get(id)
}

func (s *ProductQuantityServiceImpl) ListByDetailId(detailId string) ([]*product.ProductQuantity, error) {
	return s.QuantityRepo.GetByDetailId(detailId)
}

func (s *ProductQuantityServiceImpl) Create(newQuantity product.ProductQuantity) (*product.ProductQuantity, error) {

	searchOption := product.QuantitySearchOption{Color: newQuantity.Color, Size: newQuantity.Size, DetailId: newQuantity.DetailId}
	// Check if product quantity already exists
	if existing, err := s.QuantityRepo.GetBySearchOption(searchOption); err == nil {
		if existing != nil {
			return nil, fmt.Errorf("could not create quantity: quantity already exists")
		}
	}

	if err := s.QuantityRepo.Create(newQuantity); err != nil {
		return nil, fmt.Errorf("could not create quantity: %v", err)
	}
	return &newQuantity, nil
}

func (s *ProductQuantityServiceImpl) Update(updateQuantity product.ProductQuantity) (*product.ProductQuantity, error) {
	// Get the existing product quantity
	existing, err := s.QuantityRepo.GetBySearchOption(product.QuantitySearchOption{
		Id:       "",
		DetailId: updateQuantity.DetailId,
		Color:    updateQuantity.Color,
		Size:     updateQuantity.Size,
	})
	if existing != nil {
		return nil, fmt.Errorf("could not update quantity , because has a same one")
	}

	if err = s.QuantityRepo.Update(updateQuantity); err != nil {
		return nil, fmt.Errorf("could not update quantity: %v", err)
	}

	return &updateQuantity, nil
}

func (s *ProductQuantityServiceImpl) DeleteManyByDetailId(detailId string) error {
	return s.QuantityRepo.DeleteManyByDetailId(detailId)
}

func (s *ProductQuantityServiceImpl) DeleteOne(id string) error {
	return s.QuantityRepo.Delete(id)
}
