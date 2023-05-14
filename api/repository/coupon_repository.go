package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"online_fashion_shop/api/model/coupon"
	"online_fashion_shop/initializers"
)

type CouponRepository interface {
	Get(couponCode string) (*coupon.CouponInfo, error)
}

type CouponRepositoryImpl struct {
	CouponCollection initializers.Collection
}

func (repo *CouponRepositoryImpl) Get(couponCode string) (*coupon.CouponInfo, error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	filter := bson.M{"code": couponCode}
	var coupon coupon.CouponInfo
	err := repo.CouponCollection.FindOne(ctx, filter).Decode(&coupon)
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func NewCouponRepositoryImpl(collection initializers.Collection) CouponRepository {
	return &CouponRepositoryImpl{CouponCollection: collection}
}
