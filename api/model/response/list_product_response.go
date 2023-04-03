package response

import "online_fashion_shop/api/model"

type ListProductResponse struct {
	Products  []model.Product `json:"products"`
	TotalPage int             `json:"total_page"`
	Page      int             `json:"page"`
}
