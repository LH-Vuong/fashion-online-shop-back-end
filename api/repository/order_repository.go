package repository

import (
	"online_fashion_shop/api/model/order"
)

type OrderRepository interface {
	ListByCustomerId(customerId string, limit int, offset int) ([]*order.OrderInfo, int, error)
	Create(info order.OrderInfo) (*order.OrderInfo, error)
	GetOneByPaymentId(paymentId string) (*order.OrderInfo, error)
	Update(info order.OrderInfo) error
}
