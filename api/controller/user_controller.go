package controller

import (
	"net/http"
	model "online_fashion_shop/api/model/user"
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return UserController{service}
}

// Get Me
//	@Summary		get user info
//	@Description	get the current user info
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200				{object}	model.User
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/users/me [get]
func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(model.User)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": currentUser}})
}
