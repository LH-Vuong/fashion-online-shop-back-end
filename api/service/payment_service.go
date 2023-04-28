package service

import (
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/model/payment"
)

type PaymentService interface {
	InitPayment(orderInfo *model.OrderInfo) error
	GetStatus(paymentId string, method payment.Method) (payment.Status, error)
}
