package service

import (
	"math"
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/repository"
)

type CartService interface {
	Get(customerId string) ([]*model.CartItem, error)

	Update(customerId string, items []model.CartItem) error

	Add(customerId string, cartItem model.CartItem) (*model.CartItem, error)

	DeleteOne(customerId string, productId string) error

	DeleteAll(customerId string) error

	CheckOut(customerId string) ([]string, error)
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

func toCartItemMap(items []*model.CartItem) map[string]*model.CartItem {
	cartItemMap := make(map[string]*model.CartItem, len(items))

	for _, item := range items {
		cartItemMap[item.ProductId] = item
	}
	return cartItemMap
}

func toProductQuantityMap(quantities []*model.ProductQuantity) map[string]*model.ProductQuantity {

	quantityMap := make(map[string]*model.ProductQuantity, len(quantities))

	for _, quantity := range quantities {
		quantityMap[quantity.Id] = quantity
	}
	return quantityMap
}

func (service *CartServiceImpl) CheckOut(customerId string) ([]string, error) {

	items, err := service.cartRepo.ListByCustomerId(customerId)
	if err != nil {
		return nil, err
	}
	cartItemMap := toCartItemMap(items)
	productIds := make([]string, len(items))
	quantities, err := service.quantityRepo.MultiGet(productIds)
	if err != nil {
		return nil, err
	}
	quantityMap := toProductQuantityMap(quantities)

	var soldOutItemIds []string
	for productId, cartItem := range cartItemMap {
		if quantityMap[productId] == nil {

		}
		if cartItem.Quantity > quantityMap[productId].Quantity {
			soldOutItemIds = append(soldOutItemIds, cartItem.ProductId)
		}
	}

	service.cartRepo.DeleteAll(customerId, soldOutItemIds)

	return soldOutItemIds, nil
}

func (service *CartServiceImpl) Get(customerId string) (cartItems []*model.CartItem, err error) {
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
	quantityMap := make(map[string]*model.ProductQuantity, len(productQuantities))
	for _, item := range productQuantities {

		quantityMap[item.Id] = item
	}

	productDetailMap, err := service.getProductDetailMap(productQuantities)
	if err != nil {
		return nil, err
	}

	for index := range cartItems {
		item := cartItems[index]
		productQuantity := quantityMap[item.ProductId]
		if productQuantity == nil {
			//with item not found product quantity will be fill by empty values
			continue
		}
		item.ProductDetail = *productDetailMap[productQuantity.DetailId]
		item.Color = productQuantity.Color
		item.Size = productQuantity.Size
	}
	return cartItems, err
}

func (service *CartServiceImpl) getProductDetailMap(productQuantities []*model.ProductQuantity) (map[string]*model.Product, error) {

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
	productDetailMap := make(map[string]*model.Product, len(productDetails))

	for _, productDetail := range productDetails {

		productDetailMap[productDetail.Id] = &productDetail
	}

	return productDetailMap, nil

}

func (service *CartServiceImpl) Update(customerId string, items []model.CartItem) error {

	err := service.cartRepo.DeleteByCustomerId(customerId)
	if err != nil {
		return err
	}
	_, err = service.cartRepo.MultiCreate(customerId, items)
	if err != nil {
		return err
	}

	return nil
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

func (service *CartServiceImpl) DeleteOne(customerId string, productId string) error {

	err := service.cartRepo.DeleteOne(customerId, productId)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return err
}

func (service *CartServiceImpl) DeleteAll(customerId string) error {
	return service.cartRepo.DeleteByCustomerId(customerId)
}
