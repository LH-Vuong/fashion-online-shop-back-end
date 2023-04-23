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
		s.GET("/cart/:customer_id", controller.Get)
		s.PUT("/cart", controller.Add)
		s.POST("/cart", controller.Update)
		s.DELETE("/cart/:customer_id/:product_id", controller.Delete)
	})
	if err != nil {
		panic(err)
	}
}
