package service

import (
	"online_fashion_shop/api/model/rating"
)

type RateService interface {
	Get(id string) rating.Rating
	List(productIds []string) rating.Rating
}
