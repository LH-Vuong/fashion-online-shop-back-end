package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online_fashion_shop/api/common/errs"
	"online_fashion_shop/api/model/response"
	"online_fashion_shop/api/service"
)

type DeliveryController struct {
	DeliveryService service.DeliveryService
}

// CalculateFee of delivery service
//
//	@Summary		calculate fee
//	@Description	calculate and return fee base on addressId of customer
//	@Tags			Delivery
//	@Accept			json
//	@Param			address_id	path	string	true	"addressID"
//	@Produce		json
//	@Success		200				{object}	response.BaseResponse[int]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/delivery/cal_fee/{address_id} [get]
func (deliveryCtr DeliveryController) CalculateFee(ctx *gin.Context) {
	addressId := ctx.Param("address_id")
	fee, err := deliveryCtr.DeliveryService.CalculateFeeByCustomerAddress(addressId)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, response.BaseResponse[int]{
		Data:    fee,
		Message: "delivery fee",
		Status:  "success",
	})

}
