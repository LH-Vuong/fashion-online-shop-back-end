package repository

import "online_fashion_shop/api/model"

type CouponRepository interface {
	Get(couponCode string) (*model.CouponInfo, error)
}
