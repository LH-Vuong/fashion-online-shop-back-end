package service

import "online_fashion_shop/api/model"

type PhotoService interface {
	Get(id string) model.ProductPhoto
	List(ids string) []model.ProductPhoto
}
