package repository

import (
	"online_fashion_shop/api/model/product"
	"online_fashion_shop/initializers"

	"go.mongodb.org/mongo-driver/bson"
)

type ProductPhotoRepository interface {
	ListByProductId(productId string) ([]*product.ProductPhoto, error)
	GetOne(productId string) (*product.ProductPhoto, error)
	ListByMultiProductId(productIds []string) ([]*product.ProductPhoto, error)
	InsertOne(photo product.ProductPhoto) error
	DeleteByProductId(productId string) error
}

type ProductPhotoRepositoryImpl struct {
	PhotoCollection initializers.Collection
}

func (repository *ProductPhotoRepositoryImpl) GetOne(productId string) (*product.ProductPhoto, error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	rs := repository.PhotoCollection.FindOne(ctx, bson.M{"product_id": productId})
	var photo product.ProductPhoto
	err := rs.Decode(&photo)
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

func (repository *ProductPhotoRepositoryImpl) InsertOne(photo product.ProductPhoto) error {

	ctx, cancel := initializers.InitContext()
	defer cancel()
	_, err := repository.PhotoCollection.InsertOne(ctx, photo)
	return err
}

func (repository *ProductPhotoRepositoryImpl) DeleteByProductId(productId string) error {

	ctx, cancel := initializers.InitContext()
	defer cancel()

	_, err := repository.PhotoCollection.DeleteOne(ctx, bson.M{"product_id": productId})
	if err != nil {
		return err
	}

	return nil
}

func (repository *ProductPhotoRepositoryImpl) ListByProductId(productId string) (photos []*product.ProductPhoto, err error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	rs, err := repository.PhotoCollection.Find(ctx, bson.M{"product_id": productId})
	if err != nil {
		return nil, err
	}
	err = rs.All(ctx, &photos)
	if err != nil {
		return nil, err
	}
	return
}

func (repository *ProductPhotoRepositoryImpl) ListByMultiProductId(productIds []string) (photos []*product.ProductPhoto, err error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()

	if len(productIds) > 0 {
		rs, err := repository.PhotoCollection.Find(ctx, bson.M{"product_id": bson.M{"$in": productIds}})
		if err != nil {
			return nil, err
		}
		err = rs.All(ctx, &photos)
		if err != nil {
			return nil, err
		}
	}

	return
}

func NewProductPhotoRepository(photoCollection initializers.Collection) ProductPhotoRepository {
	return &ProductPhotoRepositoryImpl{
		photoCollection,
	}
}
