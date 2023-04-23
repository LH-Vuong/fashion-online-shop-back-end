package router

import (
	"online_fashion_shop/api/controller"
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func InitCartRouter(s *gin.Engine, c *dig.Container) {
	err := c.Invoke(func(cartService service.CartService) {
		controller := controller.CartController{Service: cartService}
		s.GET("api/cart/:customer_id", controller.Get)
		s.PUT("api/cart", controller.Add)
		s.POST("api/cart", controller.Update)
		s.DELETE("api/cart/:customer_id/:product_id", controller.Delete)
		s.GET("api/cart/checkout/:customer_id", controller.CheckOut)
	})
	if err != nil {
		panic(err)
	}
}
