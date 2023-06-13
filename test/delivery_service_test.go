package flow

import (
	container "online_fashion_shop/api"
	"online_fashion_shop/api/service"
	"testing"
)

func TestDelivery(t *testing.T) {
	var deliveryService service.DeliveryService
	container.Inject(&deliveryService)
	t.Run("cal delivery fee", func(t *testing.T) {
		fee, err := deliveryService.CalculateFeeByCustomerAddress("64868cde86301cfb9920daa3")
		if err != nil {
			panic(err)
		}
		println(fee)
	})

}
