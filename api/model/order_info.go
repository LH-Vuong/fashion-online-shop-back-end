package model

import "online_fashion_shop/api/model/payment"

type OrderInfo struct {
	Id             string          `bson:"_id" json:"id"`
	CustomerId     string          `bson:"customer_id" json:"customer_id"`
	Address        string          `bson:"address" json:"address"`
	CouponCode     string          `bson:"coupon_code" json:"coupon_code"`
	CouponDiscount int64           `bson:"coupon_discount"json:"coupon_discount"`
	PaymentInfo    *payment.Detail `json:"payment_info" bson:"payment_info"`
	TotalPrice     int64           `bson:"total" json:"total"`
	Items          []*CartItem     `json:"items"bson:"items"`
}

type OrderItem struct {
	orderId           string `bson:"order_id"`
	productQuantityId string `bson:"product_quantity_id"`
	productDetailId   string `bson:"product_detail_id"`
	quantity          int64  `bson:"quantity"`
	id                string `bson:"_id" json:"id"`
}
