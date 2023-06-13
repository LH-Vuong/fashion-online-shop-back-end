package flow

import (
	container "online_fashion_shop/api"
	"online_fashion_shop/external_services"
	"testing"
)

func TestGHN(t *testing.T) {
	var ghn external_services.GHNService
	container.Inject(&ghn)
	t.Run("cal delivery fee", func(t *testing.T) {
		fee, err := ghn.CalculateFee(1566, "510101")
		if err != nil {
			panic(err)
		}
		println(fee)
	})

}
