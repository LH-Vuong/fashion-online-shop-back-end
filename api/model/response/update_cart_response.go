package response

import "online_fashion_shop/api/model"

type UpdateCartResponse struct {
	CartItems []model.CartItem `json:"cart_items"`
	isSuccess bool             `json:"is_success"`
}
