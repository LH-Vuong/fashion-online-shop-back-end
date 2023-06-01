package service

import (
	"fmt"
	"online_fashion_shop/api/model/coupon"
	"online_fashion_shop/api/repository"
	"time"
)

type CouponService interface {
	Get(code string) (*coupon.CouponInfo, error)
	Check(code string) (bool, error)
	Delete(code string) error
	Update(info *coupon.CouponInfo) error
	Create(info *coupon.CouponInfo) error
}

type CouponServiceImpl struct {
	couponRepo repository.CouponRepository
}

func (c CouponServiceImpl) Update(info *coupon.CouponInfo) error {
	return c.couponRepo.Update(info)
}

func (c CouponServiceImpl) Create(info *coupon.CouponInfo) error {
	return c.couponRepo.Create(info)
}

func (c CouponServiceImpl) Delete(code string) error {
	err := c.couponRepo.Delete(code)
	if err != nil {
		return fmt.Errorf("'%v' coupon do not exist", code)
	}
	return nil
}

func NewCouponServiceImpl(repo repository.CouponRepository) CouponService {
	return &CouponServiceImpl{couponRepo: repo}
}

func (c CouponServiceImpl) Get(couponCode string) (*coupon.CouponInfo, error) {
	couponInfo, err := c.couponRepo.Get(couponCode)
	if err != nil {
		return nil, fmt.Errorf("Your '%s'code is invalid ", couponCode)
	}

	if couponInfo.EndAt < time.Now().UnixMilli() {
		expiredAt := time.UnixMilli(couponInfo.EndAt).Format("02/01/2006 15:04:05")
		return nil, fmt.Errorf("Your '%s' coupon  was expired at %s", couponCode, expiredAt)
	}

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
