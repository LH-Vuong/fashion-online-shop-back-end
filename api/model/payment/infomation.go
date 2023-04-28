package payment

type Detail struct {
	Id            string
	OrderId       string
	Status        Status
	OrderUrl      *string
	PaymentMethod Method
	CreatedAt     int64
	LastUpdateAt  int64
}
