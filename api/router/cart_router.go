package router

import (
	"online_fashion_shop/api/controller"
	middleware "online_fashion_shop/api/middlewares"
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func InitCartRouter(s *gin.Engine, c *dig.Container) {
	err := c.Invoke(func(cartService service.CartService) {
		controller := controller.CartController{Service: cartService}
		s.GET("api/cart", middleware.DeserializeUser(), controller.Get)
		s.PUT("api/cart", middleware.DeserializeUser(), controller.AddMany)
		s.POST("api/cart", middleware.DeserializeUser(), controller.Update)
		s.DELETE("api/cart", middleware.DeserializeUser(), controller.DeleteMany)
		s.DELETE("api/cart/product_id", middleware.DeserializeUser(), controller.Delete)
		s.GET("api/cart/checkout", middleware.DeserializeUser(), controller.CheckOut)
	})
	if err != nil {
		panic(err)
	}
}
