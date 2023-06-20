package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"online_fashion_shop/api/model/product"
	"online_fashion_shop/initializers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductQuantityRepository interface {
	GetId(size string, color string, productId string) (string, error)
	Get(productId string) (*product.ProductQuantity, error)
	GetBySearchOption(searchOption product.QuantitySearchOption) (*product.ProductQuantity, error)
	GetByDetailId(productDetailId string) ([]*product.ProductQuantity, error)
	MultiGet(productIds []string) ([]*product.ProductQuantity, error)
	Create(product product.ProductQuantity) error
	Update(quantity *product.ProductQuantity) error
	Delete(id string) error
	DeleteManyByDetailId(detailId string) error
}

func NewProductQuantityRepositoryImpl(quantityCollection initializers.Collection) ProductQuantityRepository {
	return &ProductQuantityRepositoryImpl{
		quantityCollection,
	}
}

type ProductQuantityRepositoryImpl struct {
	quantityCollection initializers.Collection
}

func (p *ProductQuantityRepositoryImpl) GetId(size string, color string, detailId string) (string, error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	query := bson.M{"detail_id": detailId, "size": size, "color": color}
	quantity := product.ProductQuantity{}
	err := p.quantityCollection.FindOne(ctx, query).Decode(&quantity)
	if err != nil {
		return "", err
	}
	return quantity.Id, nil
}

func (p *ProductQuantityRepositoryImpl) DeleteManyByDetailId(detailId string) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	_, err := p.quantityCollection.DeleteMany(ctx, bson.M{"detail_id": detailId})
	return err
}

func (p *ProductQuantityRepositoryImpl) GetBySearchOption(searchOption product.QuantitySearchOption) (*product.ProductQuantity, error) {
	var productQuantity product.ProductQuantity
	ctx, cancel := initializers.InitContext()
	defer cancel()
	// Perform search on quantityCollection based on provided parameters
	var filters []primitive.M

	if searchOption.DetailId != "" {
		detailFilter := bson.M{"detail_id": searchOption.DetailId}
		filters = append(filters, detailFilter)
	}

	if searchOption.Color != "" {
		colorFilter := bson.M{"color": searchOption.Color}
		filters = append(filters, colorFilter)
	}

	if searchOption.Size != "" {
		sizeFilter := bson.M{"size": searchOption.Size}
		filters = append(filters, sizeFilter)
	}
	var queryFilter primitive.M
	if len(filters) > 0 {
		queryFilter = bson.M{"$and": filters}
	} else {
		queryFilter = bson.M{}
	}

	err := p.quantityCollection.FindOne(ctx, queryFilter).Decode(&productQuantity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // return nil, nil if no matching document found
		}
		return nil, err
	}

	return &productQuantity, nil
}

func (p *ProductQuantityRepositoryImpl) Get(productId string) (*product.ProductQuantity, error) {
	var quantity product.ProductQuantity
	ctx, cancel := initializers.InitContext()
	defer cancel()
	id, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return nil, err
	}
	query := bson.M{"_id": id}
	err = p.quantityCollection.FindOne(ctx, query).Decode(&quantity)
	if err != nil {
		return nil, err
	}
	return &quantity, nil
}

func (p *ProductQuantityRepositoryImpl) GetByDetailId(productDetailId string) ([]*product.ProductQuantity, error) {
	var quantities []*product.ProductQuantity
	ctx, cancel := initializers.InitContext()
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

func (p *ProductQuantityRepositoryImpl) MultiGet(productIds []string) ([]*product.ProductQuantity, error) {
	var quantities []*product.ProductQuantity
	ctx, cancel := initializers.InitContext()
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

func (p *ProductQuantityRepositoryImpl) Create(quantity product.ProductQuantity) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	//ignore id
	quantity.Id = ""
	_, err := p.quantityCollection.InsertOne(ctx, &quantity)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductQuantityRepositoryImpl) Update(quantity *product.ProductQuantity) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	objectId, err := primitive.ObjectIDFromHex(quantity.Id)
	if err != nil {
		return err
	}
	query := bson.M{"_id": objectId}
	quantity.Id = ""
	update := bson.M{"$set": quantity}
	_, err = p.quantityCollection.UpdateOne(ctx, query, update)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductQuantityRepositoryImpl) Delete(id string) error {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	query := bson.M{"_id": objectId}
	_, err = p.quantityCollection.DeleteOne(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
