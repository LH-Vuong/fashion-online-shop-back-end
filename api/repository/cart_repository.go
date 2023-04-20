package repository

import (
	"errors"
	"online_fashion_shop/api/model"
	"online_fashion_shop/initializers"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartRepository interface {
	Create(customerId string, item model.CartItem) (string, error)
	MultiCreate(customerId string, items []model.CartItem) ([]string, error)
	ListByCustomerId(customerId string) ([]model.CartItem, error)
	GetBySearchOption(searchOption CartSearchOption) (*model.CartItem, error)
	DeleteByCustomerId(customerId string) error
	Delete(item model.CartItem) error
	Update(customerId string, cartItem model.CartItem) error
}

type CartSearchOption struct {
	CustomerId string `bson:"customer_id"`
	ProductId  string `bson:"product_id"`
	Id         string `bson:"id"`
}

func (searchOption CartSearchOption) ToQuery() primitive.M {

	filters := make([]primitive.M, 0, 3)
	if searchOption.Id != "" {
		filters = append(filters, bson.M{"_id": searchOption.Id})
	}
	if searchOption.CustomerId != "" {
		filters = append(filters, bson.M{"customer_id": searchOption.CustomerId})
	}
	if searchOption.ProductId != "" {
		filters = append(filters, bson.M{"product_id": searchOption.ProductId})
	}

	if len(filters) > 0 {
		return bson.M{"$and": filters}
	}
	return bson.M{}
}

// CartRepositoryImpl represents an implementation of the CartRepository interface
type CartRepositoryImpl struct {
	cartCollection initializers.Collection
}

// NewCartRepositoryImpl creates a new instance of the CartRepositoryImpl
func NewCartRepositoryImpl(cartCollection initializers.Collection) CartRepository {
	return &CartRepositoryImpl{
		cartCollection: cartCollection,
	}
}

// Create inserts a new CartItem into the cartCollection
func (cri *CartRepositoryImpl) Create(customerID string, item model.CartItem) (string, error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	item.CustomerId = customerID
	item.CreatedAt = time.Now().UnixMilli()
	item.CreateBy = customerID
	res, err := cri.cartCollection.InsertOne(ctx, item)
	if err != nil {
		return "", errors.New("failed to insert CartItem: " + err.Error())
	}
	insertedID := res.(primitive.ObjectID).Hex()
	return insertedID, nil
}

// MultiCreate inserts multiple CartItems into the cartCollection
func (cri *CartRepositoryImpl) MultiCreate(customerID string, items []model.CartItem) ([]string, error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	var documents []interface{}
	for _, item := range items {
		item.CustomerId = customerID
		item.CreatedAt = time.Now().UnixMilli()
		item.CreateBy = customerID
		documents = append(documents, item)
	}
	res, err := cri.cartCollection.InsertMany(ctx, documents)
	if err != nil {
		return nil, errors.New("failed to insert multiple CartItems: " + err.Error())
	}
	insertedIDs := make([]string, len(res))
	for i, id := range res {
		insertedIDs[i] = id.(primitive.ObjectID).Hex()
	}
	return insertedIDs, nil
}

// ListByCustomerId fetches all CartItems associated with a customerID
func (cri *CartRepositoryImpl) ListByCustomerId(customerID string) ([]model.CartItem, error) {

	ctx, cancel := initializers.InitContext()
	defer cancel()
	var cartItems []model.CartItem
	query := bson.M{"customer_id": customerID}
	cursor, err := cri.cartCollection.Find(ctx, query)
	if err != nil {
		return nil, errors.New("failed to fetch CartItems: " + err.Error())
	}

	if err = cursor.All(ctx, &cartItems); err != nil {
		return nil, errors.New("failed to parse CartItems: " + err.Error())
	}

	return cartItems, nil
}

// GetBySearchOption fetches a CartItem based on a CartSearchOption object
func (cri *CartRepositoryImpl) GetBySearchOption(searchOption CartSearchOption) (*model.CartItem, error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	var cartItem *model.CartItem
	query := searchOption.ToQuery()
	err := cri.cartCollection.FindOne(ctx, query).Decode(&cartItem)
	if err != nil {
		return nil, errors.New("failed to fetch CartItem: " + err.Error())
	}

	return cartItem, nil
}

// DeleteByCustomerId removes all CartItems associated with a customerID
func (cri *CartRepositoryImpl) DeleteByCustomerId(customerID string) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	query := bson.M{"customer_id": customerID}
	_, err := cri.cartCollection.DeleteMany(ctx, query)
	if err != nil {
		return errors.New("failed to remove CartItems: " + err.Error())
	}
	return nil
}

// Update updates a CartItem in the cartCollection
func (cri *CartRepositoryImpl) Update(customerID string, cartItem model.CartItem) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	cartItem.UpdatedAt = time.Now().UnixMilli()
	cartItem.UpdatedBy = customerID
	cartItem.CustomerId = customerID
	query := bson.M{"customer_id": customerID, "product_id": cartItem.ProductId}
	update := bson.M{"$set": cartItem}
	result, err := cri.cartCollection.UpdateOne(ctx, query, update)
	if err != nil {
		return errors.New("failed to update CartItem: " + err.Error())
	}
	if result.ModifiedCount == 0 {
		return errors.New("no CartItem found to update")
	}

	return nil
}

func (cri *CartRepositoryImpl) Delete(item model.CartItem) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	query := bson.M{"customer_id": item.CustomerId, "product_id": item.ProductId}
	_, err := cri.cartCollection.DeleteMany(ctx, query)
	if err != nil {
		return errors.New("failed to remove CartItems: " + err.Error())
	}
	return nil

}
