package repository

import (
	"online_fashion_shop/api/model"
	"online_fashion_shop/initializers"

	"go.mongodb.org/mongo-driver/bson"
)

type ProductPhotoRepository interface {
	ListByProductId(productId string) ([]*model.ProductPhoto, error)
	ListByMultiProductId(productIds []string) ([]*model.ProductPhoto, error)
}

type ProductPhotoRepositoryImpl struct {
	PhotoCollection initializers.Collection
}

func (repository *ProductPhotoRepositoryImpl) ListByProductId(productId string) (photos []*model.ProductPhoto, err error) {
	ctx, cancel := initializers.InitContext()
	defer cancel()
	rs, err := repository.PhotoCollection.Find(ctx, bson.M{"product_id": productId})
	if err != nil {
		return nil, err
	}
	err = rs.All(ctx, photos)
	if err != nil {
		return nil, err
	}
	return
}

func (repository *ProductPhotoRepositoryImpl) ListByMultiProductId(productIds []string) (photos []*model.ProductPhoto, err error) {
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
