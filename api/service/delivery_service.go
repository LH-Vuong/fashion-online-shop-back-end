package service

type DeliveryService interface {
	CalculateFeeByCustomerAddress(addressId string) (float64, error)
}
