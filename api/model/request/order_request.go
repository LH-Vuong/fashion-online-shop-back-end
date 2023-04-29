package request

import "online_fashion_shop/api/model/order"

type CreateOrderRequest struct {
	CustomerId    string
	CouponCode    *string
	PaymentMethod order.Method
	AddressInfo   string
}
