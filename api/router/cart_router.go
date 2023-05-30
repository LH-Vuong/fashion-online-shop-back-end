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
		cartRouter := s.Group("api/cart")

		{
			cartRouter.Use(middleware.DeserializeUser())
			cartRouter.PUT("", controller.AddMany)
			cartRouter.POST("", controller.Update)
			cartRouter.DELETE("", controller.DeleteMany)
			cartRouter.DELETE("/:cart_id", controller.Delete)
			cartRouter.GET("", controller.Get)
			cartRouter.GET("/checkout", controller.CheckOut)
		}
	})
	if err != nil {
		panic(err)
	}
}
