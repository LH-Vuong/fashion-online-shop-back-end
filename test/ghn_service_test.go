package flow

import (
	container "online_fashion_shop/api"
	"online_fashion_shop/external_services"
	"testing"
)

func TestProduct(t *testing.T) {
	var ghn external_services.GHNService
	container.Inject(&ghn)
	t.Run("Get Product", func(t *testing.T) {

	})

}
