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

		s.GET("/api/admin/orders", middleware.DeserializeUser(), c.List)

		//list customer's order
		s.GET("/api/orders", middleware.DeserializeUser(), c.CustomerList)
		//checkout order info
		s.POST("/api/order/checkout", middleware.DeserializeUser(), c.Checkout)
		// listen zalo callback
		s.POST("/api/order/callback/zalo_pay", middleware.ValidateZaloPayCallback, c.Callback)

		s.GET("/api/order/reject/:order_id", middleware.DeserializeUser(), c.Reject)
		s.GET("/api/order/approve/:order_id", middleware.DeserializeUser(), c.Approve)
	})

}
