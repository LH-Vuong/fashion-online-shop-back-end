package response

import (
	"online_fashion_shop/api/model/cart"
)

type UpdateCartResponse struct {
	CartItems []cart.CartItem `json:"cart_items"`
	isSuccess bool            `json:"is_success"`
}
