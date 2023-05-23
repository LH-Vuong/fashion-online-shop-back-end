package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"online_fashion_shop/api/controller"
	"online_fashion_shop/api/service"
)

func InitPhotoRouter(s *gin.Engine, c *dig.Container) {
	err := c.Invoke(func(photoService service.PhotoService) {
		controller := controller.PhotoController{PhotoService: photoService}
		s.DELETE("api/product_photo", controller.DeleteOne)
		s.DELETE("api/product_photos/:product_id", controller.DeleteAll)
		s.POST("api/product_photos/upload/:product_id", controller.Upload)
	})

	if err != nil {
		panic(err)
	}

}
