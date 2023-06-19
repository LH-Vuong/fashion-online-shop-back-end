package repository

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"online_fashion_shop/api/model/coupon"
	"online_fashion_shop/initializers"
)

type CouponRepository interface {
	Get(couponCode string) (*coupon.CouponInfo, error)
	List(codes []string) ([]*coupon.CouponInfo, error)
	ListBySearchOptions(searchOptions coupon.SearchOption) ([]*coupon.CouponInfo, error)
	Delete(couponCode string) error
	Update(create *coupon.CouponInfo) error
	Create(create *coupon.CouponInfo) error
}

type CouponRepositoryImpl struct {
	CouponCollection initializers.Collection
}

func (repo *CouponRepositoryImpl) ListBySearchOptions(searchOptions coupon.SearchOption) (coupons []*coupon.CouponInfo, err error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	var filters []primitive.M

	if searchOptions.Code != "" {
		keyWordFilter := bson.M{"code": primitive.Regex{Pattern: searchOptions.Code, Options: ""}}
		filters = append(filters, keyWordFilter)
	}
	if searchOptions.From > 0 {
		startTimeFilter := bson.M{"start_at": bson.M{"$gte": searchOptions.From}}
		filters = append(filters, startTimeFilter)
	}

	if searchOptions.To > 0 {
		endTimeFilter := bson.M{"end_at": bson.M{"$lte": searchOptions.To}}
		filters = append(filters, endTimeFilter)
	}

	var queryFilter primitive.M
	if len(filters) > 0 {
		queryFilter = bson.M{"$and": filters}
	} else {
		queryFilter = bson.M{}
	}

	rs, err := repo.CouponCollection.Find(ctx, queryFilter)
	if err != nil {
		return nil, err
	}
	err = rs.All(ctx, &coupons)
	if err != nil {
		return nil, err
	}
	return
}

func (repo *CouponRepositoryImpl) List(codes []string) (coupons []*coupon.CouponInfo, err error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	rs, err := repo.CouponCollection.Find(ctx, bson.M{"code": bson.M{"$in": codes}})
	if err != nil {
		return nil, err
	}
	err = rs.All(ctx, &coupons)
	if err != nil {
		return nil, err
	}
	return
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
	objId, err := primitive.ObjectIDFromHex(couponInfo.Id)
	couponInfo.Id = ""
	_, err = repo.CouponCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": *couponInfo})
	if err != nil {
		return errors.New("failed to update coupon: " + err.Error())
	}
	return nil
}

func (repo *CouponRepositoryImpl) Create(couponInfo *coupon.CouponInfo) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	couponInfo.Id = ""
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
