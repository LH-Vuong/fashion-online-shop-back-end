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
			router.GET("/me", controller.GetMe)
			router.GET("/address", controller.GetUserAddressList)
			router.POST("/address", controller.CreateUserAddress)
			router.PUT("/address", controller.UpdateUserAddress)
			router.DELETE("/address", controller.DeleteUserAddress)
			router.GET("/wishlist", controller.GetUserWishlist)
			router.POST("/wishlist", controller.AddUserWishlist)
			router.DELETE("/wishlist", controller.DeleteUserWishlist)
		}
	})
	if err != nil {
		panic(err)
	}
}
