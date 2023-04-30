package request

import (
	"online_fashion_shop/api/model/payment"
)

type CreateOrderRequest struct {
	CustomerId    string         `json:"-"`
	CouponCode    *string        `json:"coupon_code"`
	PaymentMethod payment.Method `json:"payment_method"`
	AddressInfo   string         `json:"address_info"`
}
