package repository

import (
	"fmt"
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/model/product"
	"online_fashion_shop/initializers"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductDetailRepository interface {
	Get(id string) (*product.Product, error)
	List(ids []string,
		keyWord string,
		tags []string,
		brands []string,
		productType []string,
		gender []string,
		priceRange model.RangeValue[int64], startAt int, length int) ([]*product.Product, int64, error)

	ListBySearchOption(searchOption product.ProductSearchOption) ([]*product.Product, int64, error)
	ListByMultiId(ids []string) ([]*product.Product, error)
	Update(product *product.Product) error
	Create(product *product.Product) error
}

type ProductDetailRepositoryImpl struct {
	collection initializers.Collection
}

func (repo *ProductDetailRepositoryImpl) Update(product *product.Product) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	id, err := primitive.ObjectIDFromHex(product.Id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	product.Id = ""
	product.UpdatedAt = time.Now().UnixMilli()
	update := bson.M{
		"$set": *product,
	}

	_, err = repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
func (repo *ProductDetailRepositoryImpl) Create(product *product.Product) error {

	ctx, cancel := initializers.InitContext()
	defer cancel()
	now := time.Now().UnixMilli()
	product.CreatedAt = now
	product.UpdatedAt = now
	product.Id = ""
	// Insert the product document into the collection
	res, err := repo.collection.InsertOne(ctx, *product)
	if err != nil {
		return err
	}

	product.Id = res.(primitive.ObjectID).Hex()

	return nil
}

func (repository *ProductDetailRepositoryImpl) Get(id string) (product *product.Product, err error) {
	ctx, cancel := initializers.InitContext()
	objectId, _ := primitive.ObjectIDFromHex(id)
	defer cancel()
	rs := repository.collection.FindOne(ctx, bson.M{"_id": objectId})
	err = rs.Decode(&product)
	return
}

func (repository *ProductDetailRepositoryImpl) List(productIds []string, keyWord string,
	tags []string,
	brands []string,
	productTypes []string,
	genders []string,
	priceRange model.RangeValue[int64],
	beginAt int,
	length int,
) (products []*product.Product, totalDocs int64, err error) {
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
		genderFilter := bson.M{"gender": bson.M{"$in": genders}}
		filters = append(filters, genderFilter)
	}

	priceFilter := bson.M{"price": bson.M{"$gte": priceRange.From, "$lte": priceRange.To}}
	filters = append(filters, priceFilter)
	ctx, cancel := initializers.InitContext()
	defer cancel()

	var queryFilter primitive.M
	if len(filters) > 0 {
		queryFilter = bson.M{"$and": filters}
	} else {
		queryFilter = bson.M{}
	}

	opts := options.Find()
	opts.SetSkip(int64(beginAt))
	opts.SetLimit(int64(length))
	totalDocs, err = repository.collection.CountDocuments(ctx, queryFilter)

	rs, err := repository.collection.Find(ctx, queryFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	err = rs.All(ctx, &products)
	if err != nil {
		return nil, 0, err
	}
	return
}

func (repository *ProductDetailRepositoryImpl) ListBySearchOption(searchOption product.ProductSearchOption) ([]*product.Product, int64, error) {
	return repository.List(searchOption.Ids,
		searchOption.KeyWord,
		searchOption.Tags,
		searchOption.Brands,
		searchOption.ProductType,
		searchOption.Gender,
		searchOption.PriceRange,
		searchOption.StartAt,
		searchOption.Length,
	)

}

func (repository *ProductDetailRepositoryImpl) ListByMultiId(ids []string) (products []*product.Product, err error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	if len(ids) < 1 {
		return nil, fmt.Errorf("list with empty id array")
	}
	var objectIds []primitive.ObjectID
	for _, id := range ids {
		objectID, _ := primitive.ObjectIDFromHex(id)
		objectIds = append(objectIds, objectID)
	}
	rs, err := repository.collection.Find(ctx, bson.M{"_id": bson.M{"$in": objectIds}})
	if err != nil {
		return nil, err
	}
	rs.All(ctx, &products)
	return
}

func NewProductRepositoryImpl(productCollection initializers.Collection) ProductDetailRepository {
	return &ProductDetailRepositoryImpl{collection: productCollection}
}
