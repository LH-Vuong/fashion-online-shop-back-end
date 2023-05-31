package service

import (
	"fmt"
	"online_fashion_shop/api/model/cart"
	"online_fashion_shop/api/model/product"
	"online_fashion_shop/api/model/request"
	"online_fashion_shop/api/repository"
)

type CartService interface {
	Get(customerId string) ([]*cart.CartItem, error)

	Update(customerId string, items []request.CartItemUpdater) error

	AddMany(customerId string, items []request.CartItemUpdater) ([]*cart.CartItem, error)

	Add(customerId string, cartItem request.CartItemUpdater) (*cart.CartItem, error)

	DeleteOne(customerId string, productId string, size string, color string) error

	DeleteOneById(id string, customerId string) (string, error)

	DeleteAll(customerId string) ([]string, error)

	CheckOutAndDelete(customerId string) ([]string, error)

	ListInvalidCartItem(customerId string) ([]string, error)
}

func NewCartServiceImpl(cartRepo repository.CartRepository,
	quantityRepo repository.ProductQuantityRepository,
	productService ProductService) CartService {
	return &CartServiceImpl{
		cartRepo:      cartRepo,
		quantityRepo:  quantityRepo,
		detailService: productService,
	}
}

type CartServiceImpl struct {
	cartRepo      repository.CartRepository
	quantityRepo  repository.ProductQuantityRepository
	detailService ProductService
	//	detailRepo   repository.ProductDetailRepository
}

func (service *CartServiceImpl) DeleteOneById(id string, customerId string) (string, error) {
	err := service.cartRepo.DeleteOne(customerId, id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (service *CartServiceImpl) AddMany(customerId string, updateInfos []request.CartItemUpdater) ([]*cart.CartItem, error) {
	curItems, err := service.cartRepo.ListByCustomerId(customerId)
	updatedItems := make([]*cart.CartItem, len(updateInfos))
	for index := range updatedItems {
		updatedItems[index], err = service.toCartItem(updateInfos[index], customerId)
		if err != nil {
			return nil, err
		}
	}

	curItemAsMap := toCartItemMap(curItems)
	if err != nil {
		return nil, err
	}

	//check item has exited one before add,
	for _, addItem := range updatedItems {
		if exitedItem, ok := curItemAsMap[addItem.InventoryId]; ok {
			exitedItem.Quantity = exitedItem.Quantity + addItem.Quantity
		} else {
			curItems = append(curItems, addItem)
		}
	}

	_, err = service.DeleteAll(customerId)
	if err != nil {
		return nil, err
	}
	_, err = service.cartRepo.MultiCreate(customerId, curItems)
	if err != nil {
		return nil, err
	}
	err = service.addDetail(updatedItems)
	if err != nil {
		return nil, err
	}
	return curItems, nil
}

func toCartItemMap(items []*cart.CartItem) map[string]*cart.CartItem {
	cartItemMap := make(map[string]*cart.CartItem, len(items))

	for _, item := range items {
		cartItemMap[item.InventoryId] = item
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

	cartItems, err := service.cartRepo.ListByCustomerId(customerId)
	if err != nil {
		return nil, err
	}
	cartItemAsMap := toCartItemMap(cartItems)
	inventoryIds := make([]string, len(cartItems))
	inventoryInfos, err := service.quantityRepo.MultiGet(inventoryIds)
	if err != nil {
		return nil, err
	}
	inventoryAsMap := toProductQuantityMap(inventoryInfos)

	var invalidProduct []string
	for cartItemId, cartItem := range cartItemAsMap {
		if inventoryAsMap[cartItemId] != nil {
			if cartItem.Quantity > inventoryAsMap[cartItemId].Quantity {
				invalidProduct = append(invalidProduct, cartItem.InventoryId)
			}

			if cartItem.Quantity < 1 {
				invalidProduct = append(invalidProduct, cartItem.InventoryId)
			}
		}
	}
	return invalidProduct, err
}

func (service *CartServiceImpl) addDetail(cartItems []*cart.CartItem) error {
	if len(cartItems) == 0 {
		return nil
	}
	productIds := make([]string, len(cartItems))
	for index := range cartItems {
		productIds[index] = cartItems[index].InventoryId
	}
	productQuantities, err := service.quantityRepo.MultiGet(productIds)

	if err != nil {
		return err
	}
	quantityMap := make(map[string]*product.ProductQuantity, len(productQuantities))
	for _, item := range productQuantities {
		quantityMap[item.Id] = item
	}

	productDetailMap, err := service.getProductDetailMap(productQuantities)
	if err != nil {
		return err
	}

	for index := range cartItems {
		item := cartItems[index]
		productQuantity := quantityMap[item.InventoryId]
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
	return nil
}

func (service *CartServiceImpl) Get(customerId string) (cartItems []*cart.CartItem, err error) {
	cartItems, err = service.cartRepo.ListByCustomerId(customerId)
	service.addDetail(cartItems)
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

	productDetails, err := service.detailService.ListMultiId(detailIds)
	if err != nil {
		return nil, err
	}

	for _, productDetail := range productDetails {
		productDetailMap[productDetail.Id] = productDetail
	}

	return productDetailMap, nil

}

// to get inventory info of product
func (service *CartServiceImpl) getQuantityId(size string, color string, productId string) (string, error) {
	return service.quantityRepo.GetId(size, color, productId)
}

func (service *CartServiceImpl) Update(customerId string, updaters []request.CartItemUpdater) error {
	newItem := make([]*cart.CartItem, len(updaters))
	for index, updater := range updaters {
		inventoryId, err := service.getQuantityId(updater.Size, updater.Color, updater.ProductId)

		if err != nil {
			fmt.Errorf("encountered an error while trying to retrieve the inventory information for this item. It's possible that the item does not exist")
		}

		newItem[index] = &cart.CartItem{
			CustomerId:  customerId,
			InventoryId: inventoryId,
			Quantity:    updater.Quantity,
			Color:       updater.Color,
			Size:        updater.Size,
		}
	}
	_, err := service.DeleteAll(customerId)
	if err != nil {
		return err
	}
	_, err = service.AddMany(customerId, updaters)
	if err != nil {
		return err
	}

	return nil
}

func (service *CartServiceImpl) toCartItem(updateInfo request.CartItemUpdater, customerId string) (*cart.CartItem, error) {
	inventoryId, err := service.getQuantityId(updateInfo.Size, updateInfo.Color, updateInfo.ProductId)
	if err != nil {
		return nil, fmt.Errorf("encountered an error(%s) while defining your cart item. It's possible that the information for your cart item is incorrect", err.Error())
	}
	cartItem := cart.CartItem{
		CustomerId:  customerId,
		InventoryId: inventoryId,
		Quantity:    updateInfo.Quantity,
		Color:       updateInfo.Color,
		Size:        updateInfo.Size,
	}
	return &cartItem, nil
}

func (service *CartServiceImpl) Add(customerId string, addItem request.CartItemUpdater) (*cart.CartItem, error) {
	cartItem, err := service.toCartItem(addItem, customerId)
	item, err := service.cartRepo.GetBySearchOption(repository.CartSearchOption{CustomerId: customerId, InventoryId: cartItem.InventoryId})
	if item == nil {
		_, err = service.cartRepo.Create(customerId, *cartItem)
	} else {
		addItem.Quantity = item.Quantity + addItem.Quantity
		err = service.cartRepo.Update(customerId, *cartItem)
	}

	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (service *CartServiceImpl) DeleteOne(customerId string, productId string, size string, color string) error {
	inventoryId, err := service.getQuantityId(size, color, productId)
	if err != nil {
		fmt.Errorf("encountered an error while trying to define your delete item")
	}
	err = service.cartRepo.DeleteOne(customerId, inventoryId)
	if err != nil {
		return err
	}
	return err
}

func (service *CartServiceImpl) DeleteAll(customerId string) ([]string, error) {
	deleteItems, err := service.cartRepo.ListByCustomerId(customerId)
	if err != nil {
		return nil, err
	}
	Ids := make([]string, len(deleteItems))
	for index := range deleteItems {
		Ids[index] = deleteItems[index].InventoryId
	}
	return Ids, service.cartRepo.DeleteByCustomerId(customerId)
}
