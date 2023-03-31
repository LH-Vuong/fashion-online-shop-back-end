package router

import (
	controller2 "online_fashion_shop/api/controller"
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func InitProductRouter(s *gin.Engine, c *dig.Container) {
	v1 := s.Group("/api/v1")
	{
		c.Invoke(func(ProductService service.ProductService) {
			controller := controller2.ProductController{Service: ProductService}
			v1.GET("/product/:id", controller.Get)
		})
	}

}
