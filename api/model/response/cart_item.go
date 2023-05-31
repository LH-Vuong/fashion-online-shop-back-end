package response

type CartItem struct {
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	Image           string  `json:"image"`
	Price           float64 `json:"price"`
	DiscountPercent float64 `bson:"discount_percent"`
	Size            string  `json:"size"`
	Color           string  `json:"color"`
	Quantity        int     `json:"quantity"`
}
