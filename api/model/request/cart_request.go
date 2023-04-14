package request

import "online_fashion_shop/api/model"

type CartRequest struct {
	CustomerId string           `json:"customer_id"`
	Items      []model.CartItem `json:"items"`
}

type CartCheckOutRequest struct {
	IsReady        bool             `json:"is_ready"`
	BackOrderItems []model.CartItem `json:"able_to_order"`
	FunFilledItems []model.CartItem `json:"fun_filled_items"`
}
