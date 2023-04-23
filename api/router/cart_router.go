package router

import (
	"online_fashion_shop/api/controller"
	middleware "online_fashion_shop/api/middlewares"
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func InitCartRouter(s *gin.Engine, c *dig.Container) {
	router := s.Group("/cart")
	router.Use(middleware.DeserializeUser())
	err := c.Invoke(func(cartService service.CartService) {
		controller := controller.CartController{Service: cartService}
		s.GET("api/cart", controller.Get)
		s.PUT("api/cart", controller.Add)
		s.POST("api/cart", controller.Update)
		s.DELETE("api/cart/:product_id", controller.Delete)
		s.GET("api/cart/checkout", controller.CheckOut)
	})
	if err != nil {
		panic(err)
	}
}
