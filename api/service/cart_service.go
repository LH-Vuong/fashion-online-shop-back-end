package service

import (
	"online_fashion_shop/api/model/cart"
	"online_fashion_shop/api/model/product"
	"online_fashion_shop/api/repository"
)

type CartService interface {
	Get(customerId string) ([]*cart.CartItem, error)

	Update(customerId string, items []*cart.CartItem) error

	AddMany(customerId string, items []*cart.CartItem) ([]*cart.CartItem, error)

	Add(customerId string, cartItem cart.CartItem) (*cart.CartItem, error)

	DeleteOne(customerId string, productId string) error

	DeleteAll(customerId string) error

	CheckOutAndDelete(customerId string) ([]string, error)

	ListInvalidCartItem(customerId string) ([]string, error)
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

func (service *CartServiceImpl) AddMany(customerId string, items []*cart.CartItem) ([]*cart.CartItem, error) {
	curItems, err := service.cartRepo.ListByCustomerId(customerId)
	itemMap := toCartItemMap(curItems)
	if err != nil {
		return nil, err
	}

	//check item has exited one before add,
	for _, addItem := range items {
		if exitedItem, ok := itemMap[addItem.ProductId]; ok {
			exitedItem.Quantity = exitedItem.Quantity + addItem.Quantity
		} else {
			curItems = append(curItems, addItem)
		}
	}

	err = service.DeleteAll(customerId)
	if err != nil {
		return nil, err
	}
	_, err = service.cartRepo.MultiCreate(customerId, curItems)
	if err != nil {
		return nil, err
	}
	curItems, err = service.Get(customerId)
	if err != nil {
		return nil, err
	}
	return curItems, nil
}

func toCartItemMap(items []*cart.CartItem) map[string]*cart.CartItem {
	cartItemMap := make(map[string]*cart.CartItem, len(items))

	for _, item := range items {
		cartItemMap[item.ProductId] = item
	}
	return cartItemMap
}

func toProductQuantityMap(quantities []*product.ProductQuantity) map[string]*product.ProductQuantity {

	quantityMap := make(map[string]*product.ProductQuantity, len(quantities))

	for _, quantity := range quantities {
		quantityMap[quantity.Id] = quantity
	}
	return quantityMap
}

func (service *CartServiceImpl) CheckOutAndDelete(customerId string) ([]string, error) {
	invalidItems, err := service.ListInvalidCartItem(customerId)
	if err != nil {
		return nil, err
	}
	err = service.cartRepo.DeleteAll(customerId, invalidItems)
	if err != nil {
		return nil, err
	}
	return invalidItems, nil
}

func (service *CartServiceImpl) ListInvalidCartItem(customerId string) ([]string, error) {

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

	var invalidProduct []string
	for productId, cartItem := range cartItemMap {
		if quantityMap[productId] != nil {
			if cartItem.Quantity > quantityMap[productId].Quantity {
				invalidProduct = append(invalidProduct, cartItem.ProductId)
			}

			if cartItem.Quantity < 1 {
				invalidProduct = append(invalidProduct, cartItem.ProductId)
			}
		}
	}
	return invalidProduct, err
}
func (service *CartServiceImpl) Get(customerId string) (cartItems []*cart.CartItem, err error) {
	cartItems, err = service.cartRepo.ListByCustomerId(customerId)

	if len(cartItems) == 0 {
		return nil, nil
	}
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
	quantityMap := make(map[string]*product.ProductQuantity, len(productQuantities))
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

		if detailId, ok := productDetailMap[productQuantity.DetailId]; ok {
			item.ProductDetail = *detailId
		}

		item.Color = productQuantity.Color
		item.Size = productQuantity.Size
	}

	return cartItems, err
}

func (service *CartServiceImpl) getProductDetailMap(productQuantities []*product.ProductQuantity) (map[string]*product.Product, error) {

	detailIds := make([]string, len(productQuantities))
	for index, item := range productQuantities {
		detailIds[index] = item.DetailId
	}
	productDetailMap := make(map[string]*product.Product, len(detailIds))
	if len(detailIds) < 1 {
		return productDetailMap, nil
	}

	productDetails, err := service.detailRepo.ListByMultiId(detailIds)
	if err != nil {
		return nil, err
	}

	for _, productDetail := range productDetails {
		productDetailMap[productDetail.Id] = productDetail
	}

	return productDetailMap, nil

}

func (service *CartServiceImpl) Update(customerId string, items []*cart.CartItem) error {

	err := service.DeleteAll(customerId)
	if err != nil {
		return err
	}
	println("get item len")
	println(len(items))
	_, err = service.cartRepo.MultiCreate(customerId, items)
	if err != nil {
		return err
	}

	return nil
}

func (service *CartServiceImpl) Add(customerId string, addItem cart.CartItem) (*cart.CartItem, error) {

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
