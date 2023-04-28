package service

import (
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/repository"
	"time"
)

type CouponService interface {
	Get(code string) (*model.CouponInfo, error)
	Check(code string) (bool, error)
}

type CouponServiceImpl struct {
	couponRepo repository.CouponRepository
}

func (c CouponServiceImpl) Get(couponCode string) (*model.CouponInfo, error) {
	return c.couponRepo.Get(couponCode)
}

func (svc CouponServiceImpl) Check(couponCode string) (bool, error) {
	coupon, err := svc.couponRepo.Get(couponCode)
	if err != nil {
		return false, err
	}
	if coupon.EndAt > time.Now().UnixMilli() && coupon.StartAt < time.Now().UnixMilli() {
		return true, nil
	}
	return false, nil
}
