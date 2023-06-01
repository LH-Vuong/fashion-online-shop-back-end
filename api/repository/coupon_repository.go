package repository

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"online_fashion_shop/api/model/coupon"
	"online_fashion_shop/initializers"
)

type CouponRepository interface {
	Get(couponCode string) (*coupon.CouponInfo, error)
	Delete(couponCode string) error
	Update(create *coupon.CouponInfo) error
	Create(create *coupon.CouponInfo) error
}

type CouponRepositoryImpl struct {
	CouponCollection initializers.Collection
}

func (repo *CouponRepositoryImpl) Delete(couponCode string) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	_, err := repo.CouponCollection.DeleteOne(ctx, bson.M{"code": couponCode})
	if err != nil {
		return errors.New("failed to delete coupon: " + err.Error())
	}
	return nil
}

func (repo *CouponRepositoryImpl) Update(couponInfo *coupon.CouponInfo) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	_, err := repo.CouponCollection.UpdateOne(ctx, bson.M{"code": couponInfo.Code}, bson.M{"$set": *couponInfo})
	if err != nil {
		return errors.New("failed to update coupon: " + err.Error())
	}
	return nil
}

func (repo *CouponRepositoryImpl) Create(couponInfo *coupon.CouponInfo) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	_, err := repo.CouponCollection.InsertOne(ctx, couponInfo)
	if err != nil {
		return errors.New("failed to create coupon: " + err.Error())
	}
	return nil
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
