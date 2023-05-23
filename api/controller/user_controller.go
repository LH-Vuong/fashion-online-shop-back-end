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
//
//	@Summary		get user info
//	@Description	get the current user info
//	@Tags			users
//	@Param		Authorization	header		string	true	"Access Token"
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
//
//	@Summary		create user's address
//	@Description	create new user's address
//	@Tags			users
//	@Param			Authorization	header		string	true	"Access Token"
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
//
//	@Summary		update user's address
//	@Description	update user's address
//	@Tags			users
//	@Param			Authorization	header		string	true	"Access Token"
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
//
//	@Summary		Detele user's address
//	@Description	Detele user's address
//	@Tags			users
//	@Param			Authorization	header		string	true	"Access Token"
//	@Accept			json
//	@Produce		json
//	@Param			DeletedId		body		string	true	"User's address id"
//	@Success		200				{object}	model.UserAddress
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/users/address [delete]
func (uc *UserController) DeleteUserAddress(ctx *gin.Context) {
	deleteId := ctx.PostForm("DeletedId")

	uc.Service.DeleteUserAddress(ctx, deleteId)
}

// Get User's Address List
//
//	@Summary		Get user's address list
//	@Description	Get user's address list
//	@Tags			users
//	@Param			Authorization	header		string	true	"Access Token"
//	@Accept			json
//	@Produce		json
//	@Param			GetUserAddressListModel		query		model.GetUserAddressListModel	true	"User's address filter"
//	@Success		200				{object}	[]model.UserAddress
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/users/address [get]
func (uc *UserController) GetUserAddressList(ctx *gin.Context) {
	var payload model.GetUserAddressListModel

	if err := ctx.ShouldBindQuery(&payload); err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
		return
	}

	uc.Service.GetUserAddressList(ctx, payload)
}

// Add User's Wishlist Item
//
//	@Summary		add user's wishlist item
//	@Description	add user's wishlist item
//	@Tags			users
//	@Param			Authorization	header		string	true	"Access Token"
//	@Accept			json
//	@Produce		json
//	@Param			AddWishListModel		body		model.AddWishListModel	true	"User's wishlist item"
//	@Success		200				{object}	model.UserWishlist
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/users/wishlist [post]
func (uc *UserController) AddUserWishlist(ctx *gin.Context) {
	var payload model.AddWishListModel

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
		return
	}

	uc.Service.AddUserWishlistItem(ctx, payload)
}

// Delete User's wishlist item
//
//	@Summary		Detele user's wishlist item
//	@Description	Detele user's wishlist item
//	@Tags			users
//	@Param			Authorization	header		string	true	"Access Token"
//	@Accept			json
//	@Produce		json
//	@Param			DeleteIds		body		model.DeleteWishListModel	true	"User's wishlist item id"
//	@Success		200				{object}	[]string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/users/wishlist [delete]
func (uc *UserController) DeleteUserWishlist(ctx *gin.Context) {
	var payload model.DeleteWishListModel
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
		return
	}

	uc.Service.DeleteUserWishlistItems(ctx, payload)
}

// Get User's Address List
//
//	@Summary		Get user's address list
//	@Description	Get user's address list
//	@Tags			users
//	@Param			Authorization	header		string	true	"Access Token"
//	@Accept			json
//	@Produce		json
//	@Param			GetUserAddressListModel		query		model.GetUserAddressListModel	true	"User's address filter"
//	@Success		200				{object}	[]model.UserWishlist
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/users/wishlist [get]
func (uc *UserController) GetUserWishlist(ctx *gin.Context) {
	var payload model.GetUserWishlistModel

	if err := ctx.ShouldBindQuery(&payload); err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
		return
	}

	uc.Service.GetUserWishlist(ctx, payload)
}
