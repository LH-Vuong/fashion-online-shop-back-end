package coupon

type CouponInfo struct {
	Code            string  `json:"code" bson:"code,omitempty"`
	StartAt         int64   `json:"start_at"bson:"start_at,omitempty"`
	EndAt           int64   `json:"end_at" bson:"end_at,omitempty"`
	DiscountPercent float32 `json:"discount_percent" bson:"discount_percent,omitempty"`
}
