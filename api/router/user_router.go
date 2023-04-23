package router

import (
	"online_fashion_shop/api/controller"
	middleware "online_fashion_shop/api/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controller.UserController
}

func NewUserRouteController(userController controller.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (rc *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("/users")
	router.Use(middleware.DeserializeUser())
	router.GET("/me", rc.userController.GetMe)
}
