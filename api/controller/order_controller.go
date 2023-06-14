package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"online_fashion_shop/api/common/errs"
	"online_fashion_shop/api/model/order"
	"online_fashion_shop/api/model/request"
	"online_fashion_shop/api/model/response"
	model "online_fashion_shop/api/model/user"
	"online_fashion_shop/api/service"
	"online_fashion_shop/initializers/zalopay"
	"strconv"
)

type OrderController struct {
	Service service.OrderService
}

// Create Order
//
//	@Summary		Creat order
//	@Description	Create order by customer's. Will delete all items of customer's cart
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Access Token"
//	@Param          OrderRequest		body       request.CreateOrderRequest	  true 	"access token received after login"
//	@Success		200				{object}	order.OrderInfo
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/order/ [put]
func (controller *OrderController) Create(ctx *gin.Context) {
	var createRequest request.CreateOrderRequest
	err := ctx.BindJSON(&createRequest)
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
	}
	currentUser := ctx.MustGet("currentUser").(model.User)
	createRequest.CustomerId = currentUser.Id
	info, err := controller.Service.Create(createRequest.CustomerId, createRequest.PaymentMethod, createRequest.AddressId, createRequest.CouponCodes)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.JSON(200, gin.H{"status": "success", "data": info})
}

// List Order
//
//	@Summary		list of customer's order
//	@Description	List order by customer id
//	@Tags			order
//	@Param		Authorization	header		string	true	"Access Token"
//	@Accept			json
//	@Produce		json
//	@Param          off_set	query       int		  false 	"index of first item, default is 0"
//	@Param          limit	query       int		  false		"max length of response, default is 10"
//	@Success		200				{object}	response.PagingResponse[order.OrderInfo]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/orders/ [get]
func (controller *OrderController) List(ctx *gin.Context) {
	offset, err := strconv.Atoi(ctx.Param("off_set"))
	if err != nil || offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(ctx.Param("limit"))
	if err != nil || limit > request.PageMaximum {
		limit = 10
	}
	currentUser := ctx.MustGet("currentUser").(model.User)
	infos, total, err := controller.Service.ListByCustomerID(currentUser.Id, limit, offset)

	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
	}

	ctx.JSON(200, response.PagingResponse[*order.OrderInfo]{
		Length: len(infos),
		Total:  total,
		Status: "success",
		Data:   infos,
	})
}

// Checkout  is able to create order
//
//	@Summary		checkout order request is valid
//	@Description	Validates order info if any invalid info, such as sold-out cart items, invalid coupon. Use this method before placing an order to ensure that the order is valid.If the order status is "failed," the reason for the failure will be displayed in the "message" field, and any issues will be indicated in the "data" field.
//	@Tags			order
//	@Param			Authorization	header		string	true	"Access Token"
//	@Param          coupon_code	path       string		  false		"code want to apply for order"
//	@Accept			json
//	@Produce		json
//	@Success		200				{object}	 response.BaseResponse[[]string]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/order/checkout/:coupon_code [get]
func (controller *OrderController) Checkout(ctx *gin.Context) {
	couponCode := ctx.Param("coupon_code")
	currentUser := ctx.MustGet("currentUser").(model.User)
	invalidData, err := controller.Service.IsAbleCreateOrder(currentUser.Id, couponCode)
	if err != nil {
		ctx.JSON(200, response.BaseResponse[[]string]{
			Status:  "failed",
			Data:    invalidData,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{"status": "success"})

}

func (controller *OrderController) Callback(ctx *gin.Context) {
	var cbData map[string]interface{}
	err := ctx.BindJSON(&cbData)
	if err != nil {
		return
	}
	dataStr := cbData["data"].(string)
	var dataJSON map[string]interface{}
	json.Unmarshal([]byte(dataStr), &dataJSON)

	err = controller.Service.UpdateWithCallbackData(ctx.Param("payment_id"), dataJSON, zalopay.HandleCallback)
	if err != nil {
		return
	}

}
