package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"online_fashion_shop/api/controller"
	"online_fashion_shop/api/service"
)

func InitCouponRouter(s *gin.Engine, c *dig.Container) {
	c.Invoke(func(couponService service.CouponService) {
		controller := controller.CouponController{Service: couponService}
		s.GET("api/coupon/:code", controller.Get)
		s.DELETE("api/coupon/:code", controller.Delete)
		s.POST("api/coupon", controller.Update)
		s.PUT("api/coupon", controller.Create)
	})
}
