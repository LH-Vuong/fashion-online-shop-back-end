package request

import (
	"online_fashion_shop/api/model/cart"
)

type CartRequest struct {
	CustomerId string          `json:"customer_id" binding:"required"`
	Items      []cart.CartItem `json:"items" binding:"required"`
}

type ModifyCartRequest struct {
	CustomerId string     `json:"customer_id" binding:"required"`
	Items      []CartItem `json:"items" binding:"required"`
}

type CartItem struct {
	ProductId string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type CartCheckOutRequest struct {
	IsReady        bool            `json:"is_ready"`
	BackOrderItems []cart.CartItem `json:"able_to_order"`
	FunFilledItems []cart.CartItem `json:"fun_filled_items"`
}
