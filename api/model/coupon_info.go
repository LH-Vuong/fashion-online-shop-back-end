package model

type CouponInfo struct {
	Code            string
	StartAt         int64
	EndAt           int64
	DiscountAmount  int64
	DiscountPercent float32
}
