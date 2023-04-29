package order

type PaymentDetail struct {
	Id             string
	OrderId        string
	Status         Status
	OrderUrl       *string
	PaymentMethod  Method
	CreatedAt      int64
	UpdatedAt      int64
	PaymentAt      int64
	ReceivedAmount int64
}
