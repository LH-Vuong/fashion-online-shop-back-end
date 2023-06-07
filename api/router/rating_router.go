package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"online_fashion_shop/api/controller"
	"online_fashion_shop/api/service"
)

func InitRatingRouter(s *gin.Engine, c *dig.Container) {
	err := c.Invoke(func(ratingService service.RatingService) {
		ratingController := controller.RatingController{Service: ratingService}
		s.GET("api/ratings/:product_id", ratingController.List)
		s.GET("api/rating/:id", ratingController.Get)
		s.PUT("api/rating", ratingController.Create)
		s.POST("api/rating", ratingController.UpdateOne)
		s.DELETE("api/rating/:id", ratingController.Delete)
	})
	if err != nil {
		panic(err)
	}
}
