package router

import (
	"online_fashion_shop/api/controller"
	middleware "online_fashion_shop/api/middlewares"
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func InitUserRouter(s *gin.Engine, c *dig.Container) {
	err := c.Invoke(func(userService service.UserService) {
		controller := controller.UserController{Service: userService}
		router := s.Group("api/users")
		{
			router.Use(middleware.DeserializeUser())
			router.GET("/address", controller.GetUserAddressList)
			router.POST("/address", controller.CreateUserAddress)
			router.PUT("/address", controller.UpdateUserAddress)
			router.DELETE("/address", controller.DeleteUserAddress)
			router.GET("/wishlist", controller.GetUserWishlist)
			router.POST("/wishlist", controller.AddUserWishlist)
			router.DELETE("/wishlist", controller.DeleteUserWishlist)
			router.PATCH("", controller.UpdateUserInfo)
		}

		s.GET("api/users/me", middleware.DeserializeAdmin(), controller.GetMe)
		router1 := s.Group("api")
		{
			router1.GET("/provinces", controller.GetProvinces)
			router1.GET("/districts/:provinceId", controller.GetDistricts)
			router1.GET("/wards/:districtId", controller.GetWards)
		}

		adminRouter := s.Group("api/admin")
		{
			adminRouter.Use(middleware.DeserializeAdmin())
			adminRouter.GET("/users", controller.GetUsers)
		}
	})
	if err != nil {
		panic(err)
	}
}
