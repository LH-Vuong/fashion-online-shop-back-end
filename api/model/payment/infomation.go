package payment

type PaymentDetail struct {
	PaymentId      string  `bson:"payment_id,omitempty"json:"payment_id"`
	Status         Status  `bson:"status,omitempty"json:"status"`
	OrderUrl       *string `bson:"order_url,omitempty"json:"order_url"`
	PaymentMethod  Method  `bson:"payment_method,omitempty"json:"payment_method"`
	CreatedAt      int64   `bson:"created_at,omitempty"json:"created_at"`
	UpdatedAt      int64   `bson:"updated_at,omitempty"json:"updated_at"`
	PaymentAt      int64   `bson:"payment_at,omitempty"json:"payment_at"`
	ReceivedAmount int64   `bson:"received_amount,omitempty"json:"received_amount"`
}
