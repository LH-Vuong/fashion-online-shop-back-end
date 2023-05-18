package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"online_fashion_shop/api/controller"
	"online_fashion_shop/api/service"
)

func InitProductDetailRouter(s *gin.Engine, c *dig.Container) {
	err := c.Invoke(func(productService service.ProductService) {
		controller := controller.ProductController{Service: productService}
		s.GET("api/product/:id", controller.Get)
		s.GET("api/products", controller.List)
		s.POST("api/product/:product_id", controller.Update)
		s.PUT("api/product", controller.Create)
	})

	if err != nil {
		panic(err)
	}

}
