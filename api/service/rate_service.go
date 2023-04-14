package service

import "online_fashion_shop/api/model"

type RateService interface {
	Get(id string) model.Rating
	List(productIds []string) model.Rating
}
