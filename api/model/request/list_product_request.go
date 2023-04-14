package request

type ListProductsRequest struct {
	Brands   []string `json:"brands"`
	Colors   []string `json:"colors"`
	Rate     int      `json:"rate"`
	Tags     []string `json:"tags"`
	Genders  []string `json:"genders"`
	Types    []string `json:"types"`
	MinPrice int64    `json:"min_price"`
	MaxPrice int64    `json:"max_price"`
	KeyWord  string   `json:"key_word"`
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
}

const PageMinimum = 1

const PageMaximum = 10000
