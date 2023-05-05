package flow_test

import (
	"fmt"
	container "online_fashion_shop/api"
	"online_fashion_shop/api/service"
	"testing"
)

func TestProduct(t *testing.T) {
	productId := "64292d1b01f4fe96f272b208"
	var productService service.ProductService
	container.Inject(&productService)
	t.Run("Get Product", func(t *testing.T) {

		product, err := productService.Get(productId)
		if err != nil {
			t.Errorf(err.Error())
		}
		fmt.Printf("%+v\n", *product)
	})

}
