package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"online_fashion_shop/api/dbs"
	"online_fashion_shop/api/model"
)

type ProductQuantityRepository interface {
	Get(productId string) (*model.ProductQuantity, error)
	GetByDetailId(productDetailId string) ([]*model.ProductQuantity, error)
	MultiGet(productIds []string) ([]*model.ProductQuantity, error)
	Create(product model.ProductQuantity) error
	Update(quantity model.ProductQuantity) error
	Delete(id string) error
}

func NewProductQuantityRepositoryImpl(quantityCollection dbs.Collection) ProductQuantityRepository {
	return &ProductQuantityRepositoryImpl{
		quantityCollection,
	}
}

type ProductQuantityRepositoryImpl struct {
	quantityCollection dbs.Collection
}

func (p *ProductQuantityRepositoryImpl) Get(productId string) (*model.ProductQuantity, error) {
	var quantity model.ProductQuantity
	ctx, cancel := dbs.InitContext()
	defer cancel()
	id, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		panic(err)
	}
	query := bson.M{"_id": id}
	err = p.quantityCollection.FindOne(ctx, query).Decode(&quantity)
	if err != nil {
		panic(err)
	}
	return &quantity, nil
}

func (p *ProductQuantityRepositoryImpl) GetByDetailId(productDetailId string) ([]*model.ProductQuantity, error) {
	var quantities []*model.ProductQuantity
	ctx, cancel := dbs.InitContext()
	defer cancel()

	query := bson.M{"detail_id": productDetailId}
	cursor, err := p.quantityCollection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &quantities); err != nil {
		return nil, err
	}

	return quantities, nil
}

func (p *ProductQuantityRepositoryImpl) MultiGet(productIds []string) ([]*model.ProductQuantity, error) {
	var quantities []*model.ProductQuantity
	ctx, cancel := dbs.InitContext()
	defer cancel()
	objectIds := make([]primitive.ObjectID, len(productIds))
	for index := range productIds {
		objectIds[index], _ = primitive.ObjectIDFromHex(productIds[index])
	}

	query := bson.M{"_id": bson.M{"$in": objectIds}}
	cursor, err := p.quantityCollection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &quantities); err != nil {
		return nil, err
	}
	return quantities, nil
}

func (p *ProductQuantityRepositoryImpl) Create(quantity model.ProductQuantity) error {
	ctx, cancel := dbs.InitContext()
	defer cancel()

	_, err := p.quantityCollection.InsertOne(ctx, &quantity)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductQuantityRepositoryImpl) Update(quantity model.ProductQuantity) error {
	ctx, cancel := dbs.InitContext()
	defer cancel()

	query := bson.M{"_id": quantity.Id}
	update := bson.M{"$set": quantity}
	_, err := p.quantityCollection.UpdateOne(ctx, query, update)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductQuantityRepositoryImpl) Delete(id string) error {
	ctx, cancel := dbs.InitContext()
	defer cancel()

	query := bson.M{"_id": id}
	_, err := p.quantityCollection.DeleteOne(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
