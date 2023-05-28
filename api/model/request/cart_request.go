package request

import (
	"online_fashion_shop/api/model/cart"
)

type CartRequest struct {
	CustomerId string          `json:"customer_id" binding:"required"`
	Items      []cart.CartItem `json:"items" binding:"required"`
}

type ModifyCartRequest struct {
	CustomerId string            `json:"customer_id" binding:"required"`
	Items      []CartItemUpdater `json:"items" binding:"required"`
}

type CartItemUpdater struct {
	ProductId string `json:"product_id"`
	Quantity  int    `json:"quantity" binding:"required"`
	Color     string `json:"color" binding:"required"`
	Size      string `json:"size" binding:"required"`
}

type CartCheckOutRequest struct {
	IsReady        bool            `json:"is_ready"`
	BackOrderItems []cart.CartItem `json:"able_to_order"`
	FunFilledItems []cart.CartItem `json:"fun_filled_items"`
}
