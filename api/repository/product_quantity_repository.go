package repository

import "online_fashion_shop/api/model"

type ProductQuantityRepository interface {
	Get(productId string) (model.ProductQuantity, error)
}
