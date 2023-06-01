package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online_fashion_shop/api/common/errs"
	"online_fashion_shop/api/model/coupon"
	"online_fashion_shop/api/model/response"
	"online_fashion_shop/api/service"
)

type CouponController struct {
	Service service.CouponService
}

// Get coupon code
//
//	@Summary		Get info of coupon
//	@Description	Try to retry coupon
//	@Tags			Coupon
//	@Accept			json
//	@Param			code	path	string	true	"coupon's code"
//	@Produce		json
//	@Success		200				{object}	response.BaseResponse[coupon.CouponInfo]
//	@Failure		400				{object}	string "code is invalid or expired"
//	@Failure		401				{object}	string
//	@Router			/coupon/{code} [get]
func (controller CouponController) Get(ctx *gin.Context) {
	couponCode := ctx.Param("code")
	couponInfo, err := controller.Service.Get(couponCode)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	ctx.JSON(200, response.BaseResponse[*coupon.CouponInfo]{
		Data:   couponInfo,
		Status: "success",
	})
}

// Delete Get coupon code
//
//	@Summary		Delete coupon
//	@Description	Delete coupon
//	@Tags			Coupon
//	@Accept			json
//	@Param			code	path	string	true	"coupon's code"
//	@Produce		json
//	@Success		200				{object}	string
//	@Failure		400				{object}	string "code is invalid or expired"
//	@Failure		401				{object}	string
//	@Router			/coupon/{code} [delete]
func (controller CouponController) Delete(ctx *gin.Context) {
	couponCode := ctx.Param("code")
	err := controller.Service.Delete(couponCode)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	ctx.JSON(200, gin.H{"status": "success"})
}

// Create coupon code
//
//	@Summary		Create coupon
//	@Description	Create coupon
//	@Tags			Coupon
//	@Accept			json
//	@Param			coupon	body	coupon.CouponInfo	true	"coupon's info"
//	@Produce		json
//	@Success		200				{object}	string
//	@Failure		400				{object}	string "Invalid code"
//	@Failure		401				{object}	string
//	@Router			/coupon [put]
func (controller CouponController) Create(ctx *gin.Context) {
	var couponInfo coupon.CouponInfo
	err := ctx.BindJSON(&couponInfo)
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
		return
	}
	err = controller.Service.Create(&couponInfo)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.JSON(200, response.BaseResponse[*coupon.CouponInfo]{
		Data:    &couponInfo,
		Message: "",
		Status:  "success",
	})
}

// Update coupon code
//
//	@Summary		update coupon
//	@Description	update coupon
//	@Tags			Coupon
//	@Accept			json
//	@Param			coupon	body	coupon.CouponInfo	true	"coupon's info"
//	@Produce		json
//	@Success		200				{object}	string
//	@Failure		400				{object}	string "Invalid Code"
//	@Failure		401				{object}	string
//	@Router			/coupon [post]
func (controller CouponController) Update(ctx *gin.Context) {
	var couponInfo coupon.CouponInfo
	err := ctx.BindJSON(&couponInfo)
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
	}
	err = controller.Service.Update(&couponInfo)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.JSON(200, response.BaseResponse[*coupon.CouponInfo]{
		Data:    &couponInfo,
		Message: "",
		Status:  "success",
	})
}
