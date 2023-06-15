package controller

import (
	"context"
	"fmt"
	"net/http"
	"online_fashion_shop/api/common/errs"
	model "online_fashion_shop/api/model/user"
	"online_fashion_shop/api/service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
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

// Update User's Address@
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

// Get User's Address CustomerList
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

// Get User's Address CustomerList
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

// Get Provinces
//
//	@Summary		Get provinces
//	@Description	Get provinces
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/provinces [get]
func (uc *UserController) GetProvinces(ctx *gin.Context) {
	var resp struct {
		Districts []struct {
			ProvinceId    int      `json:"ProvinceID"`
			ProvinceName  string   `json:"ProvinceName"`
			CountryId     int      `json:"CountryID"`
			NameExtension []string `json:"NameExtension"`
		} `json:"data"`
	}

	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r := resty.New()
	res, err := r.R().
		SetContext(timeout).
		SetResult(&resp).
		SetHeader("Content-Type", "application/json").
		SetHeader("Token", "cc01798d-e4cf-11ed-bc91-ba0234fcde32").
		ForceContentType("application/json").
		Get("http://dev-online-gateway.ghn.vn/shiip/public-api/master-data/province")

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "GetDistrict")
	}

	if res.StatusCode() == http.StatusOK {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   resp,
		})
	}
}

// Get District
//
//	@Summary		Get districts
//	@Description	Get districts
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param 			provinceId 		path 			string 		true "Province's Id"
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/districts/{provinceId} [get]
func (uc *UserController) GetDistricts(ctx *gin.Context) {
	provinceId := ctx.Param("provinceId")
	var resp struct {
		Districts []struct {
			DistrictId    int      `json:"DistrictID"`
			ProvinceId    int      `json:"ProvinceID"`
			DistrictName  string   `json:"DistrictName"`
			NameExtension []string `json:"NameExtension,omitempty"`
		} `json:"data"`
	}

	r := resty.New()
	res, err := r.R().
		SetResult(&resp).
		SetHeader("Content-Type", "application/json").
		SetHeader("token", "cc01798d-e4cf-11ed-bc91-ba0234fcde32").
		SetBody(strings.NewReader(fmt.Sprintf(`{
			"province_id":%s 
		}`, provinceId))).
		Get("http://dev-online-gateway.ghn.vn/shiip/public-api/master-data/district")

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "GetDistricts")
		return
	}

	if res.StatusCode() == http.StatusOK {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   resp,
		})
	}

}

// Get Ward
//
//	@Summary		Get wards
//	@Description	Get wards
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param 			districtId 		path 			string 		true "District's Id"
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/wards/{districtId} [get]
func (uc *UserController) GetWards(ctx *gin.Context) {
	districtId := ctx.Param("districtId")
	var resp struct {
		Districts []struct {
			WardCode      string   `json:"WardCode"`
			DistrictId    int      `json:"DistrictID"`
			WardName      string   `json:"WardName"`
			NameExtension []string `json:"NameExtension,omitempty"`
		} `json:"data"`
	}

	r := resty.New()
	res, err := r.R().
		SetResult(&resp).
		SetHeader("Content-Type", "application/json").
		SetHeader("token", "cc01798d-e4cf-11ed-bc91-ba0234fcde32").
		SetQueryParam("district_id", districtId).
		Get("http://dev-online-gateway.ghn.vn/shiip/public-api/master-data/ward")

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "GetDistricts")
		return
	}

	if res.StatusCode() == http.StatusOK {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   resp,
		})
	}

}
