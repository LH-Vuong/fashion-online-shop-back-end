package request

import (
	"online_fashion_shop/api/model/payment"
)

type CreateOrderRequest struct {
	CustomerId    string
	CouponCode    *string
	PaymentMethod payment.Method
	AddressInfo   string
}
