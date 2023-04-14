package model

type ProductSearchOption struct {
	Ids         []string
	KeyWord     string
	Tags        []string
	Brands      []string
	ProductType []string
	Gender      []string
	PriceRange  RangeValue[int64]
	StartAt     int
	Length      int
}
