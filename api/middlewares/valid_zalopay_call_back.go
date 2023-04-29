package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"online_fashion_shop/initializers/payment"
)

func ValidateZaloPayCallback(ctx *gin.Context) {
	var cbData map[string]interface{}
	err := ctx.BindJSON(&cbData)
	if err != nil {
		return
	}
	mac := cbData["mac"].(string)
	dataStr := cbData["data"].(string)
	err = payment.IsValidCallback(mac, dataStr, "key2")
	if err != nil {
		return
	}
	var dataJSON map[string]interface{}
	json.Unmarshal([]byte(dataStr), &dataJSON)
	ctx.Set("payment_id", dataJSON["app_trans_id"])
	ctx.Next()
}
