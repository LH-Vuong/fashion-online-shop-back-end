package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"online_fashion_shop/api/model/order"
	"online_fashion_shop/api/model/request"
	"online_fashion_shop/api/model/response"
	"online_fashion_shop/api/service"
	"online_fashion_shop/initializers/zalopay"
	"strconv"
)

type OrderController struct {
	Service service.OrderService
}

func (controller *OrderController) Create(ctx *gin.Context) {
	var createRequest request.CreateOrderRequest
	ctx.BindJSON(&createRequest)

	info, err := controller.Service.Create(createRequest.CustomerId, createRequest.PaymentMethod, createRequest.AddressInfo, createRequest.CouponCode)
	if err != nil {
		return
	}
	ctx.JSON(200, gin.H{"status": "success", "data": info})
}

func (controller *OrderController) List(ctx *gin.Context) {
	offset, _ := strconv.Atoi(ctx.Param("off_set"))
	limit, _ := strconv.Atoi(ctx.Param("limit"))
	customerId := ctx.Param("customer_id")
	infos, total, err := controller.Service.ListByCustomerID(customerId, limit, offset)

	if err != nil {
		return
	}

	ctx.JSON(200, response.PagingResponse[*order.OrderInfo]{
		Length: len(infos),
		Total:  total,
		Status: "success",
		Data:   infos,
	})
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
