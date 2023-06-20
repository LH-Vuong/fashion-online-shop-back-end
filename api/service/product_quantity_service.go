package service

import (
	"fmt"
	"online_fashion_shop/api/model/product"
	"online_fashion_shop/api/repository"
)

type ProductQuantityService interface {
	Create(newQuantity product.ProductQuantity) (*product.ProductQuantity, error)
	Update(updateQuantity *product.ProductQuantity) error
	DeleteOne(id string) error
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

func (s *ProductQuantityServiceImpl) Update(updateQuantity *product.ProductQuantity) (err error) {
	if err = s.QuantityRepo.Update(updateQuantity); err != nil {
		return fmt.Errorf("could not update quantity: %v", err)
	}
	return nil
}

func (s *ProductQuantityServiceImpl) DeleteManyByDetailId(detailId string) error {
	return s.QuantityRepo.DeleteManyByDetailId(detailId)
}

func (s *ProductQuantityServiceImpl) DeleteOne(id string) error {
	return s.QuantityRepo.Delete(id)
}
