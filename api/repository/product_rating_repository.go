package repository

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"online_fashion_shop/api/model/rating"
	"online_fashion_shop/initializers"
	"time"
)

type ProductRatingRepository interface {
	GetAvr(productId string) (int, error)
	List(option rating.RateSearchOption) ([]*rating.Rating, error)
	ListAvrRate(productIds []string) (map[string]int, error)
	Get(id string) (*rating.Rating, error)
	InsertOne(rateInfo *rating.Rating) error
	DeleteOne(id string) error
	Update(updateInfo *rating.Rating) error
}

type ProductRatingRepositoryImpl struct {
	ratingCollection initializers.Collection
}

func (r *ProductRatingRepositoryImpl) DeleteOne(id string) error {
	objId, err := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": objId}
	ctx, cancel := initializers.InitContext()
	defer cancel()
	_, err = r.ratingCollection.DeleteOne(ctx, filter)
	return err
}

func (r *ProductRatingRepositoryImpl) Update(updateInfo *rating.Rating) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	updateInfo.UpdatedAt = time.Now().UnixMilli()
	objId, err := primitive.ObjectIDFromHex(updateInfo.Id)
	if err != nil {
		return err
	}
	query := bson.M{"_id": objId}
	updateInfo.Id = ""
	update := bson.M{"$set": &updateInfo}
	_, err = r.ratingCollection.UpdateOne(ctx, query, update)
	return err
}

func (r *ProductRatingRepositoryImpl) GetAvr(productId string) (int, error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()

	num, _ := r.ratingCollection.CountDocuments(ctx, bson.M{"rate_for": productId})
	if num < 1 {
		return 5, nil
	}
	matchStage := bson.D{{"$match", bson.D{{"rate_for", productId}}}}

	groupState := bson.D{{"$group", bson.D{
		{"_id", "$rate_for"},
		{"average", bson.D{{"$avg", "$rate"}}}}}}

	cursor, err := r.ratingCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupState})
	if err != nil {
		return -1, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			Average float64 `bson:"average"`
		}
		if err := cursor.Decode(&result); err != nil {
			return -1, err
		}
		return int(math.Ceil(result.Average)), nil
	}
	return -1, fmt.Errorf("retrieve empty value")
}

func (r *ProductRatingRepositoryImpl) List(option rating.RateSearchOption) ([]*rating.Rating, error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()

	filter := bson.M{}

	if len(option.RateFor) > 0 {
		filter["rate_for"] = bson.M{"$in": option.RateFor}
	}
	if len(option.RateBy) > 0 {
		filter["rate_by"] = bson.M{"$in": option.RateBy}
	}
	if len(option.Rates) > 0 {
		filter["rate"] = bson.M{"$in": option.Rates}
	}

	if option.RateTime.From != 0 || option.RateTime.To != 0 {
		filter["created_at"] = bson.M{}
		if option.RateTime.From != 0 {
			filter["created_at"].(bson.M)["$gt"] = option.RateTime.From
		}
		if option.RateTime.To != 0 {
			filter["created_at"].(bson.M)["$lt"] = option.RateTime.To
		}
	}

	sortBy := bson.M{}
	if option.Order {
		sortBy["rate"] = 1
	} else {
		sortBy["rate"] = -1
	}

	cursor, err := r.ratingCollection.Find(ctx, filter, options.Find().SetSort(sortBy))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*rating.Rating
	for cursor.Next(ctx) {
		var rating rating.Rating
		if err := cursor.Decode(&rating); err != nil {
			return nil, err
		}
		results = append(results, &rating)
	}

	return results, nil
}

func (r *ProductRatingRepositoryImpl) ListAvrRate(productIds []string) (map[string]int, error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	filterStatement := bson.D{{"rate_for", bson.M{"$in": productIds}}}
	groupStatement := bson.D{{"$group", bson.D{{"_id", "$rate_for"}, {"average", bson.D{{"$avg", "$rate"}}}}}}
	cursor, err := r.ratingCollection.Aggregate(ctx, mongo.Pipeline{filterStatement, groupStatement})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	resultMap := make(map[string]int)
	for cursor.Next(ctx) {
		var result struct {
			ID      string  `bson:"_id"`
			Average float64 `bson:"average"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		resultMap[result.ID] = int(result.Average)
	}

	return resultMap, nil
}

func (r *ProductRatingRepositoryImpl) Get(id string) (*rating.Rating, error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objId}
	var result rating.Rating
	if err := r.ratingCollection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("rating not found")
		} else {
			return nil, err
		}
	}
	return &result, nil
}

func (r *ProductRatingRepositoryImpl) InsertOne(rateInfo *rating.Rating) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	insertRs, err := r.ratingCollection.InsertOne(ctx, *rateInfo)
	rateInfo.Id = insertRs.(primitive.ObjectID).Hex()
	return err
}
func NewProductRatingRepositoryImpl(collection initializers.Collection) ProductRatingRepository {
	return &ProductRatingRepositoryImpl{collection}
}
