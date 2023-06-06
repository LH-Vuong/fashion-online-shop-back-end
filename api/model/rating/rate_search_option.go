package rating

import "online_fashion_shop/api/model"

type RateSearchOption struct {
	RateFor  []string
	RateBy   []string
	Rates    []string
	RateTime model.RangeValue[int64]
	Order    bool
}
