package service

import (
	"online_fashion_shop/api/model/coupon"
	"online_fashion_shop/api/repository"
	"time"
)

type CouponService interface {
	Get(code string) (*coupon.CouponInfo, error)
	Check(code string) (bool, error)
}

type CouponServiceImpl struct {
	couponRepo repository.CouponRepository
}

func NewCouponService(repo repository.CouponRepository) CouponService {
	return &CouponServiceImpl{couponRepo: repo}
}

func (c CouponServiceImpl) Get(couponCode string) (*coupon.CouponInfo, error) {
	return c.couponRepo.Get(couponCode)
}

func (svc CouponServiceImpl) Check(couponCode string) (bool, error) {
	coupon, err := svc.couponRepo.Get(couponCode)
	if err != nil {
		panic(err)
		return false, err
	}
	if coupon.EndAt > time.Now().UnixMilli() && coupon.StartAt < time.Now().UnixMilli() {
		return true, nil
	}
	return true, nil
}
