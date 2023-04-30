package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"online_fashion_shop/api/controller"
	middleware "online_fashion_shop/api/middlewares"
)

func InitOrderRouter(s *gin.Engine, c *dig.Container) {

	c.Invoke(func(controller controller.OrderController) {

		//init an order
		s.PUT("/api/order", controller.Create)
		//list customer's order
		s.GET("/api/orders", controller.List)
		// listen zalo callback
		s.POST("/api/order/callback/zalo_pay", middleware.ValidateZaloPayCallback, controller.Callback)

	})

}
