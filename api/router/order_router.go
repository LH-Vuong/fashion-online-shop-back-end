package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"online_fashion_shop/api/controller"
	middleware "online_fashion_shop/api/middlewares"
	"online_fashion_shop/api/service"
)

func InitOrderRouter(s *gin.Engine, c *dig.Container) {

	c.Invoke(func(orderService service.OrderService) {

		c := controller.OrderController{
			Service: orderService,
		}
		//init an order
		s.PUT("/api/order", middleware.DeserializeUser(), c.Create)
		//list customer's order
		s.GET("/api/orders/", middleware.DeserializeUser(), c.List)
		//checkout order info
		s.POST("/api/order/checkout", middleware.DeserializeUser(), c.Checkout)
		// listen zalo callback
		s.POST("/api/order/callback/zalo_pay", middleware.ValidateZaloPayCallback, c.Callback)
	})

}
