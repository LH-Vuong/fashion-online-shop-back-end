package router

import (
	"online_fashion_shop/api/controller"
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func InitAuthRouter(s *gin.Engine, c *dig.Container) {
	err := c.Invoke(func(userService service.UserService) {
		controller := controller.AuthController{Service: userService}

		s.POST("api/auth/sign-up", controller.SignUp)
		s.POST("api/auth/sign-in", controller.SignIn)
		s.GET("api/auth/verify", controller.VerifyAccount)
	})
	if err != nil {
		panic(err)
	}
}
