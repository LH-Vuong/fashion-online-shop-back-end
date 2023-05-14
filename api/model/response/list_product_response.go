package response

import (
	"online_fashion_shop/api/model/product"
)

type ListProductResponse struct {
	Products  []product.Product `json:"products"`
	TotalPage int               `json:"total_page"`
	Page      int               `json:"page"`
}
