package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"online_fashion_shop/api/model"
	"online_fashion_shop/initializers"
)

type CouponRepository interface {
	Get(couponCode string) (*model.CouponInfo, error)
}

type CouponRepositoryImpl struct {
	CouponCollection initializers.Collection
}

func (repo *CouponRepositoryImpl) Get(couponCode string) (*model.CouponInfo, error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	filter := bson.M{"code": couponCode}
	var coupon model.CouponInfo
	err := repo.CouponCollection.FindOne(ctx, filter).Decode(&coupon)
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func NewCouponRepositoryImpl(collection initializers.Collection) CouponRepository {
	return &CouponRepositoryImpl{CouponCollection: collection}
}
