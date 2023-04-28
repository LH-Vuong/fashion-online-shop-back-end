package controller

import (
	"net/http"
	"online_fashion_shop/api/common/errs"
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
//	@Success		200				{object}	model.User
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/users/me [get]
func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(model.User)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": currentUser}})
}

// Create User's Address
//	@Summary		create user's address
//	@Description	create new user's address
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			CreateUserAddressModel		body		model.CreateUserAddressModel	true	"User's address"
//	@Success		200				{object}	model.UserAddress
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/users/address [post]
func (uc *UserController) CreateUserAddress(ctx *gin.Context) {
	var payload model.CreateUserAddressModel

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
		return
	}

	uc.Service.CreateUserAddress(ctx, payload)
}

// Update User's Address
//	@Summary		update user's address
//	@Description	update user's address
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			UpdateUserAddressModel		body		model.UpdateUserAddressModel	true	"User's address"
//	@Success		200				{object}	model.UserAddress
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/users/address [put]
func (uc *UserController) UpdateUserAddress(ctx *gin.Context) {
	var payload model.UpdateUserAddressModel

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
		return
	}

	uc.Service.UpdateUserAddress(ctx, payload)
}

// Delete User's Address
//	@Summary		Detele user's address
//	@Description	Detele user's address
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			DeletedId		body		string	true	"User's address id"
//	@Success		200				{object}	model.UserAddress
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/users/address [delete]
func (uc *UserController) DeleteUserAddress(ctx *gin.Context) {
	var deleteId string
	if err := ctx.ShouldBindJSON(&deleteId); err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
		return
	}

	uc.Service.DeleteUserAddress(ctx, deleteId)
}

// Get User's Address List
//	@Summary		Get user's address list
//	@Description	Get user's address list
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			GetUserAddressListModel		query		GetUserAddressListModel	true	"User's address filter"
//	@Success		200				{object}	[]*model.UserAddress
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/users/address [delete]
func (uc *UserController) GetUserAddressList(ctx *gin.Context) {
	var payload model.GetUserAddressListModel

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
		return
	}

	uc.Service.GetUserAddressList(ctx, payload)
}
