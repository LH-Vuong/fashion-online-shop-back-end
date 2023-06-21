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

		productRouter := s.Group("api/product")
		{
			productRouter.GET("/:id", controller.Get)
			productRouter.POST("", controller.Update)
			productRouter.PUT("", controller.Create)
		}

		multiProductRouter := s.Group("api/products")
		{
			multiProductRouter.GET("", controller.List)
			multiProductRouter.GET("/brands", controller.ListBrands)
			multiProductRouter.GET("/types", controller.ListType)

		}

	})

	if err != nil {
		panic(err)
	}

}
