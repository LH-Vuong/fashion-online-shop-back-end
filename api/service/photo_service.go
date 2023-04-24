package service

import (
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/repository"
)

type PhotoService interface {
	ListByMultiProductId(productIds []string) (map[string][]*model.ProductPhoto, error)
	ListByProductId(productId string) ([]*model.ProductPhoto, error)
}

type PhotoServiceImpl struct {
	PhotoRepo repository.ProductPhotoRepository
}

func NewPhotoServiceImpl(photoRepo repository.ProductPhotoRepository) PhotoService {
	return PhotoServiceImpl{
		PhotoRepo: photoRepo,
	}
}

func (p PhotoServiceImpl) ListByMultiProductId(productIds []string) (map[string][]*model.ProductPhoto, error) {

	photoMap := make(map[string][]*model.ProductPhoto, len(productIds))
	photos, err := p.PhotoRepo.ListByMultiProductId(productIds)
	if err != nil {
		return nil, err
	}
	for _, photo := range photos {
		photoMap[photo.ProductId] = append(photoMap[photo.ProductId], photo)
	}

	return photoMap, nil
}

func (p PhotoServiceImpl) ListByProductId(productId string) ([]*model.ProductPhoto, error) {
	return p.PhotoRepo.ListByProductId(productId)
}
