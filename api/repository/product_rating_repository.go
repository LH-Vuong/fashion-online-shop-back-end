package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"online_fashion_shop/api/dbs"
	"online_fashion_shop/api/model"
)

type ProductRatingRepository interface {
	GetAvr(productId string) (model.AvrRate, error)
	List(productIds []string, value model.RangeValue[int]) ([]model.Rating, error)
	ListWithAvrRate(productIds []string) ([]model.AvrRate, error)
}

type ProductRatingRepositoryImpl struct {
	ratingCollection dbs.Collection
}

func (repository *ProductRatingRepositoryImpl) GetAvr(productId string) (avr model.AvrRate, err error) {

	ctx, cancel := dbs.InitContext()
	defer cancel()
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"rate_for", productId}}}},
		bson.D{{"$group", bson.D{{"_id", "$rate_for"}, {"avr_rate", bson.D{{"$avg", "$rate"}}}}}},
	}
	cursor, err := repository.ratingCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return avr, err
	}

	if cursor.Next(ctx) {
		err = cursor.Decode(&avr)
	}

	if err != nil {
		return avr, err
	}
	return
}

func (repository *ProductRatingRepositoryImpl) ListWithAvrRate(productIds []string) (ratings []model.AvrRate, err error) {

	ctx, cancel := dbs.InitContext()
	defer cancel()

	var filters []primitive.M
	if len(productIds) > 0 {
		filters = append(filters, bson.M{"product_id": bson.M{"$in": productIds}})
	}
	pipeline := bson.D{{"$group", bson.D{{"_id", "$rate_for"}, {"avr_rate", bson.D{{"$avg", "$rate"}}}}}}
	cursor, err := repository.ratingCollection.Aggregate(ctx, mongo.Pipeline{pipeline})
	if err != nil {
		panic(err)
	}

	cursor.All(ctx, &ratings)
	return
}

func (repository *ProductRatingRepositoryImpl) List(productIds []string, value model.RangeValue[int]) (ratings []model.Rating, err error) {
	ctx, cancel := dbs.InitContext()
	defer cancel()

	var filters []primitive.M

	rateFilter := bson.M{"rate": bson.M{"$gte": value.From, "$lte": value.To}}
	filters = append(filters, rateFilter)

	if len(productIds) > 0 {
		productIdsFilter := bson.M{"rate_for": bson.M{"$in": productIds}}
		filters = append(filters, productIdsFilter)
	}

	rs, err := repository.ratingCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = rs.All(ctx, &ratings)
	if err != nil {
		return nil, err
	}
	return
}

func NewProductRatingRepositoryImpl(collection dbs.Collection) ProductRatingRepository {
	return &ProductRatingRepositoryImpl{collection}
}
