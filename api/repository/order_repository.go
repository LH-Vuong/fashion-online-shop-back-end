package repository

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"online_fashion_shop/api/model/order"
	"online_fashion_shop/initializers"
)

type OrderRepository interface {
	ListByCustomerId(customerId string, limit int, offset int) (orders []*order.OrderInfo, total int64, err error)
	Create(info order.OrderInfo) (*order.OrderInfo, error)
	GetOneByPaymentId(paymentId string) (*order.OrderInfo, error)
	Update(info order.OrderInfo) error
	ListPendingOrder() (orders []*order.OrderInfo, err error)
	GetOneByOrderId(orderId string) (*order.OrderInfo, error)
	ListBySearchOptions(searchOptions order.SearchOptions) (orders []*order.OrderInfo, total int64, err error)
}

type OrderRepositoryImpl struct {
	collection initializers.Collection
}

func (repo *OrderRepositoryImpl) ListBySearchOptions(searchOptions order.SearchOptions) (orders []*order.OrderInfo, total int64, err error) {
	var filters []primitive.M
	if searchOptions.Status != "" {
		tagFilter := bson.M{"status": searchOptions.Status}
		filters = append(filters, tagFilter)
	}
	ctx, cancel := initializers.InitContext()
	defer cancel()
	var queryFilter primitive.M
	if len(filters) > 0 {
		queryFilter = bson.M{"$and": filters}
	} else {
		queryFilter = bson.M{}
	}

	total, err = repo.collection.CountDocuments(ctx, queryFilter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := repo.collection.Find(ctx, queryFilter, options.Find().
		SetLimit(int64(searchOptions.Limit)).
		SetSkip(int64(searchOptions.Offset)))
	cursor.All(ctx, &orders)
	return
}

func (repo *OrderRepositoryImpl) GetOneByOrderId(orderId string) (*order.OrderInfo, error) {
	objectId, _ := primitive.ObjectIDFromHex(orderId)
	filter := bson.M{"_id": objectId}
	ctx, cancel := initializers.InitContext()
	defer cancel()
	var orderInfo order.OrderInfo
	rs := repo.collection.FindOne(ctx, filter)
	err := rs.Decode(&orderInfo)
	if err != nil {
		return nil, err
	}
	return &orderInfo, nil
}

func NewOrderRepositoryImpl(collection initializers.Collection) OrderRepository {
	return &OrderRepositoryImpl{collection: collection}
}

func (repo *OrderRepositoryImpl) ListPendingOrder() ([]*order.OrderInfo, error) {
	filter := bson.M{"payment_info.status": "PENDING"}
	ctx, cancel := initializers.InitContext()
	defer cancel()
	var orderInfos []*order.OrderInfo
	rs, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	rs.All(ctx, &orderInfos)
	return orderInfos, nil
}

func (repo *OrderRepositoryImpl) Create(info order.OrderInfo) (*order.OrderInfo, error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	if info.Items == nil || len(info.Items) == 0 {
		return nil, fmt.Errorf("no items to order")
	}
	if info.PaymentInfo == nil {
		return nil, fmt.Errorf("no zalopay info provided")
	}

	// Insert the order into the database.
	_, err := repo.collection.InsertOne(ctx, info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (repo *OrderRepositoryImpl) ListByCustomerId(customerID string, limit int, offset int) ([]*order.OrderInfo, int64, error) {
	// Construct the filter to find orders for the customer.

	ctx, cancel := initializers.InitContext()
	defer cancel()
	filter := bson.M{"customer_id": customerID}

	// Count the total number of orders for the customer.
	total, err := repo.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find the orders for the customer with limit and offset.
	cursor, err := repo.collection.Find(ctx, filter, options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode the orders into order.OrderInfo objects.
	orders := make([]*order.OrderInfo, limit)

	err = cursor.All(ctx, &orders)

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (repo *OrderRepositoryImpl) GetOneByPaymentId(paymentID string) (*order.OrderInfo, error) {
	filter := bson.M{"payment_info.payment_id": paymentID}
	ctx, cancel := initializers.InitContext()
	defer cancel()
	var orderInfo order.OrderInfo
	err := repo.collection.FindOne(ctx, filter).Decode(&orderInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("order not found for zalopay id %s", paymentID)
		}
		return nil, err
	}
	return &orderInfo, nil
}

func (repo *OrderRepositoryImpl) Update(info order.OrderInfo) error {

	ctx, cancel := initializers.InitContext()
	defer cancel()
	// Create the filter and update document for the order.
	objectId, err := primitive.ObjectIDFromHex(info.Id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectId}
	//ignore _id field
	info.Id = ""
	update := bson.M{"$set": info}

	// Update the order in the database.
	result, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("order not found for id %s", info.Id)
	}
	return nil
}
