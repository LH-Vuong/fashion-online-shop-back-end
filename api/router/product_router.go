package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	controller2 "online_fashion_shop/api/controller"
	"online_fashion_shop/api/service"
)

func InitProductRouter(s *gin.Engine, c *dig.Container) {
	c.Invoke(func(ProductService service.ProductService) {
		controller := controller2.ProductController{Service: ProductService}
		s.GET("/product/:id", controller.Get)
	})

}
