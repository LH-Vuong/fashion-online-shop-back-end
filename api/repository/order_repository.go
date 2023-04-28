package repository

import "online_fashion_shop/api/model"

type OrderRepository interface {
	ListByCustomerId(customerId string, limit int, offset int) ([]*model.OrderInfo, error)
	Create(info model.OrderInfo) (*model.OrderInfo, error)
}
