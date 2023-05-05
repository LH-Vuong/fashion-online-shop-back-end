package product

import "online_fashion_shop/api/model"

type ProductSearchOption struct {
	Ids         []string
	KeyWord     string
	Tags        []string
	Brands      []string
	ProductType []string
	Gender      []string
	PriceRange  model.RangeValue[int64]
	StartAt     int
	Length      int
}
