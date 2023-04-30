package model

type CouponInfo struct {
	Code            string  `json:"code" bson:"code"`
	StartAt         int64   `json:"start_at"bson:"start_at"`
	EndAt           int64   `json:"end_at" bson:"end_at"`
	DiscountAmount  int64   `json:"discount_amount" bson:"discount_amount"`
	DiscountPercent float32 `json:"discount_percent" bson:"discount_percent"`
}
