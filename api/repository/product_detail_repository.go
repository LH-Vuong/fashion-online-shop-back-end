package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"online_fashion_shop/api/dbs"
	"online_fashion_shop/api/model"
)

type ProductDetailRepository interface {
	Get(id string) (model.Product, error)
	List(ids []string,
		keyWord string,
		tags []string,
		brands []string,
		productType []string,
		gender []string,
		priceRange model.RangeValue[int64], startAt int, length int) ([]model.Product, error)
}

type ProductDetailRepositoryImpl struct {
	collection dbs.Collection
}

func (repository ProductDetailRepositoryImpl) Get(id string) (product model.Product, err error) {
	ctx, cancel := dbs.InitContext()
	defer cancel()
	rs := repository.collection.FindOne(ctx, bson.D{})
	err = rs.Decode(&product)
	return
}

func (repository ProductDetailRepositoryImpl) List(productIds []string, keyWord string,
	tags []string,
	brands []string,
	productTypes []string,
	genders []string,
	priceRange model.RangeValue[int64],
	beginAt int,
	length int,
) (products []model.Product, err error) {
	var filters []primitive.M

	var objectIds []primitive.ObjectID
	for _, id := range productIds {
		objectID, _ := primitive.ObjectIDFromHex(id)
		objectIds = append(objectIds, objectID)
	}

	if len(objectIds) > 0 {
		idFilter := bson.M{"_id": bson.M{"$in": objectIds}}
		filters = append(filters, idFilter)
	}

	if keyWord != "" {
		keyWordFilter := bson.M{"name": primitive.Regex{Pattern: keyWord, Options: ""}}
		filters = append(filters, keyWordFilter)
	}

	if len(tags) > 0 {
		tagFilter := bson.M{"tags": bson.M{"$in": tags}}
		filters = append(filters, tagFilter)
	}
	if len(brands) > 0 {
		brandFilter := bson.M{"brand": bson.M{"$in": brands}}
		filters = append(filters, brandFilter)
	}
	if len(productTypes) > 0 {
		typeFilter := bson.M{"types": bson.M{"$in": productTypes}}
		filters = append(filters, typeFilter)

	}
	if len(genders) > 0 {
		genderFilter := bson.M{"genders": bson.M{"$in": genders}}
		filters = append(filters, genderFilter)
	}
	println("price")
	println(priceRange.From)
	println(priceRange.To)
	priceFilter := bson.M{"price": bson.M{"$gte": priceRange.From, "$lte": priceRange.To}}
	filters = append(filters, priceFilter)
	ctx, cancel := dbs.InitContext()
	defer cancel()

	var queryFilter primitive.M
	if len(filters) > 0 {
		println("filter")
		println(len(filters))
		queryFilter = bson.M{"$and": filters}
	} else {
		queryFilter = bson.M{}
	}

	opts := options.Find()
	opts.SetSkip(int64(beginAt))
	opts.SetLimit(int64(length))

	rs, err := repository.collection.Find(ctx, queryFilter)
	if err != nil {
		return nil, err
	}
	err = rs.All(ctx, &products)
	if err != nil {
		return nil, err
	}
	return
}

func NewProductRepositoryImpl(productCollection dbs.Collection) ProductDetailRepository {
	return ProductDetailRepositoryImpl{collection: productCollection}
}
