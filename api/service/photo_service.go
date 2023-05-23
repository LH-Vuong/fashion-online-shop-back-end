package service

import (
	"io"
	"online_fashion_shop/api/model/product"
	"online_fashion_shop/api/repository"
	"online_fashion_shop/initializers/storage"
)

type PhotoService interface {
	ListByMultiProductId(productIds []string) (map[string][]*product.ProductPhoto, error)
	ListByProductId(productId string) ([]*product.ProductPhoto, error)
	UploadPhoto(file []io.Reader, productId string) ([]string, error)
	DeletePhotoByProductId(productId string) error
	DeleteOne(photoURL string, productId string) error
}

type PhotoServiceImpl struct {
	PhotoRepo    repository.ProductPhotoRepository
	PhotoStorage storage.PhotoStorage
}

func (p *PhotoServiceImpl) DeleteOne(photoURL string, productId string) error {
	if photo, err := p.PhotoRepo.GetOne(productId); err == nil && photo != nil {
		productPhotos := ConvertPhotosToUrls([]*product.ProductPhoto{photo})
		var updateData []string
		for _, productPhoto := range productPhotos {
			if productPhoto != photoURL {
				updateData = append(updateData, productPhoto)
			}
		}
		p.PhotoRepo.DeleteByProductId(productId)
		p.PhotoRepo.InsertOne(product.ProductPhoto{
			MainPhoto: "",
			SubPhotos: updateData,
			Color:     "",
			ProductId: productId,
		})
		return p.PhotoStorage.Delete(photoURL)
	} else {
		return err
	}
}

func (p *PhotoServiceImpl) DeletePhotoByProductId(productId string) error {
	photos, err := p.PhotoRepo.ListByProductId(productId)

	if err != nil {
		return err
	}
	photoUrls := ConvertPhotosToUrls(photos)

	p.PhotoStorage.DeleteMany(photoUrls)
	err = p.PhotoRepo.DeleteByProductId(productId)
	if err != nil {
		return err
	}
	return nil
}

func (p *PhotoServiceImpl) UploadPhoto(files []io.Reader, productId string) ([]string, error) {

	paths, err := p.PhotoStorage.MultiUpload(files)
	if err != nil {
		return nil, err
	}

	curPhoto, err := p.PhotoRepo.GetOne(productId)
	if curPhoto != nil {
		paths = append(paths, curPhoto.SubPhotos...)
	}
	err = p.PhotoRepo.DeleteByProductId(productId)
	if err != nil {
		return nil, err
	}
	photo := &product.ProductPhoto{
		MainPhoto: "",
		SubPhotos: paths,
		Color:     "",
		ProductId: productId,
	}
	err = p.PhotoRepo.InsertOne(*photo)
	if err != nil {
		return nil, err
	}
	return paths, nil

}

func NewPhotoServiceImpl(photoRepo repository.ProductPhotoRepository, photoStorage storage.PhotoStorage) PhotoService {
	return &PhotoServiceImpl{
		PhotoRepo: photoRepo, PhotoStorage: photoStorage,
	}
}

func (p *PhotoServiceImpl) ListByMultiProductId(productIds []string) (map[string][]*product.ProductPhoto, error) {

	photoMap := make(map[string][]*product.ProductPhoto, len(productIds))
	photos, err := p.PhotoRepo.ListByMultiProductId(productIds)
	if err != nil {
		return nil, err
	}
	for _, photo := range photos {
		photoMap[photo.ProductId] = append(photoMap[photo.ProductId], photo)
	}

	return photoMap, nil
}

func (p *PhotoServiceImpl) ListByProductId(productId string) ([]*product.ProductPhoto, error) {
	return p.PhotoRepo.ListByProductId(productId)
}
