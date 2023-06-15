package order

import (
	"online_fashion_shop/api/model/cart"
	"online_fashion_shop/api/model/payment"
)

type OrderInfo struct {
	Id             string                 `bson:"_id,omitempty" json:"id"`
	CustomerId     string                 `bson:"customer_id,omitempty" json:"customer_id"`
	Address        string                 `bson:"address,omitempty" json:"address"`
	CouponCode     []string               `bson:"coupon_code,omitempty" json:"coupon_code"`
	CouponDiscount int64                  `bson:"coupon_discount,omitempty"json:"coupon_discount"`
	PaymentInfo    *payment.PaymentDetail `json:"payment_info" bson:"payment_info,omitempty"`
	TotalPrice     int64                  `bson:"total,omitempty" json:"total"`
	Items          []*cart.CartItem       `json:"items" bson:"items,omitempty"`
	Status         Status                 `json:"status" bson:"status"`
}
