package request

import "online_fashion_shop/api/model"

type CartRequest struct {
	CustomerId string           `json:"customer_id" binding:"required"`
	Items      []model.CartItem `json:"items" binding:"required"`
}

type UpdateCartRequest struct {
	CustomerId string       `json:"customer_id" binding:"required"`
	Items      []UpdateItem `json:"items" binding:"required"`
}

type AddItemRequest struct {
	CustomerId string `json:"customer_id" binding:"required"`
	ProductId  string `bson:"product_id" binding:"required"`
	Quantity   int    `bson:"quantity" binding:"required"`
}

type UpdateItem struct {
	ProductId string `bson:"product_id" binding:"required"`
	Quantity  int    `bson:"quantity" binding:"required"`
}

type CartCheckOutRequest struct {
	IsReady        bool             `json:"is_ready"`
	BackOrderItems []model.CartItem `json:"able_to_order"`
	FunFilledItems []model.CartItem `json:"fun_filled_items"`
}
