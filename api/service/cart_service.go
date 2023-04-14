package service

import (
	"math"
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/repository"
)

type CartService interface {
	Get(customerId string) ([]model.CartItem, error)

	Update(customerId string, items []model.CartItem) ([]model.CartItem, error)

	Add(customerId string, cartItem model.CartItem) (*model.CartItem, error)

	DeleteOne(customerId string, productId string) ([]model.CartItem, error)

	DeleteAll(customerId string) error
}

func NewCartServiceImpl(cartRepo repository.CartRepository,
	quantityRepo repository.ProductQuantityRepository,
	detailRepo repository.ProductDetailRepository) CartService {
	return &CartServiceImpl{
		cartRepo:     cartRepo,
		quantityRepo: quantityRepo,
		detailRepo:   detailRepo,
	}
}

type CartServiceImpl struct {
	cartRepo     repository.CartRepository
	quantityRepo repository.ProductQuantityRepository
	detailRepo   repository.ProductDetailRepository
}

func (service *CartServiceImpl) Get(customerId string) (cartItems []model.CartItem, err error) {
	cartItems, err = service.cartRepo.ListByCustomerId(customerId)
	if err != nil {
		return nil, err
	}
	productIds := make([]string, len(cartItems))
	for index := range cartItems {
		productIds[index] = cartItems[index].ProductId
	}
	productQuantities, err := service.quantityRepo.MultiGet(productIds)
	if err != nil {
		return nil, err
	}
	quantityMap := make(map[string]model.ProductQuantity, len(productQuantities))
	for _, item := range productQuantities {
		quantityMap[item.Id] = item
	}

	productDetailMap, err := service.getProductDetailMap(productQuantities)
	if err != nil {
		return nil, err
	}

	for index := range cartItems {
		item := &cartItems[index]
		productQuantity := quantityMap[item.ProductId]
		item.ProductDetail = productDetailMap[productQuantity.DetailId]
		item.Color = productQuantity.Color
		item.Size = productQuantity.Size
	}
	return cartItems, err
}

func (service *CartServiceImpl) getProductDetailMap(productQuantities []model.ProductQuantity) (map[string]model.Product, error) {

	detailIds := make([]string, len(productQuantities))
	for index, item := range productQuantities {
		detailIds[index] = item.DetailId
	}

	productDetails, err := service.detailRepo.ListBySearchOption(model.ProductSearchOption{Ids: detailIds, PriceRange: model.RangeValue[int64]{
		From: 0,
		To:   math.MaxInt64,
	}, StartAt: -1, Length: 10})
	if err != nil {
		panic(err)
		return nil, err
	}
	productDetailMap := make(map[string]model.Product, len(productDetails))

	for _, productDetail := range productDetails {

		productDetailMap[productDetail.Id] = productDetail
	}

	return productDetailMap, nil

}

func (service *CartServiceImpl) Update(customerId string, items []model.CartItem) ([]model.CartItem, error) {

	err := service.cartRepo.DeleteByCustomerId(customerId)
	if err != nil {
		return items, err
	}
	_, err = service.cartRepo.MultiCreate(customerId, items)
	if err != nil {
		return items, err
	}

	return items, nil
}

func (service *CartServiceImpl) Add(customerId string, addItem model.CartItem) (*model.CartItem, error) {

	item, err := service.cartRepo.GetBySearchOption(repository.CartSearchOption{CustomerId: customerId, ProductId: addItem.ProductId})
	if item == nil {
		_, err = service.cartRepo.Create(customerId, addItem)
	} else {
		addItem.Quantity = item.Quantity + addItem.Quantity
		err = service.cartRepo.Update(customerId, addItem)
	}

	if err != nil {
		return nil, err
	}

	return &addItem, nil
}

func (service *CartServiceImpl) DeleteOne(customerId string, productId string) ([]model.CartItem, error) {
	deleteItem := model.CartItem{
		CustomerId: customerId,
		ProductId:  productId,
	}
	err := service.cartRepo.Delete(deleteItem)
	if err != nil {
		return nil, err
	}

	afterDeletedItems, err := service.cartRepo.ListByCustomerId(customerId)
	if err != nil {
		return nil, err
	}
	return afterDeletedItems, err
}

func (service *CartServiceImpl) DeleteAll(customerId string) error {
	return service.cartRepo.DeleteByCustomerId(customerId)
}
