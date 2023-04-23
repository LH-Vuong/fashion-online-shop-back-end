package router

import (
	"online_fashion_shop/api/controller"

	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controller.AuthController
}

func NewAuthRouteController(authController controller.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	router.POST("/sign-up", rc.authController.SignUp)
	router.POST("/sign-in", rc.authController.SignIn)
	router.GET("/verify", rc.authController.VerifyAccount)
}
