package request

import (
	"online_fashion_shop/api/model/payment"
)

type CreateOrderRequest struct {
	CustomerId    string         `json:"-"`
	CouponCodes   []string       `json:"coupon_codes"`
	PaymentMethod payment.Method `json:"payment_method"`
	AddressId     string         `json:"address_id"`
}
