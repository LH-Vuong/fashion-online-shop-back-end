package model

type PaymentMethod string

const (
	ZALO_PAY PaymentMethod = "ZALO_PAY"
	COD      PaymentMethod = "COD"
)

type OrderStatus string

const (
	PENDING OrderStatus = "PENDING"
	CANCEL  OrderStatus = "CANCEL"
	ERROR   OrderStatus = "ERROR"
)

type OrderInfo struct {
	Id             string        `bson:"_id" json:"id"`
	PaymentMethod  PaymentMethod `bson:"payment_method" json:"payment_method"`
	CustomerId     string        `bson:"customer_id" json:"customer_id"`
	Address        string        `bson:"address" json:"address"`
	CouponCode     string        `bson:"coupon_code" json:"coupon_code"`
	CouponDiscount int64         `bson:"coupon_discount"json:"coupon_discount"`
	Status         OrderStatus   `bson:"status" json:"status"`
	TotalPrice     int64         `bson:"total" json:"total"`
	RedirectUrl    string        `bson:"redirect_url" json:"redirect_url"`
	Items          []*CartItem   `json:"items"bson:"items"`
}

type OrderItem struct {
	orderId           string `bson:"order_id"`
	productQuantityId string `bson:"product_quantity_id"`
	productDetailId   string `bson:"product_detail_id"`
	quantity          int64  `bson:"quantity"`
	id                string `bson:"_id" json:"id"`
}
