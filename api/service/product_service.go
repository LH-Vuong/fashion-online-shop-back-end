package service

import (
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/model/request"
	"online_fashion_shop/api/model/response"
	"online_fashion_shop/api/repository"
)

type ProductService interface {
	Get(id string) model.Product
	List(request request.ListProductsRequest) response.ListProductResponse
}

type ProductServiceImpl struct {
	ProductDetailRepository repository.ProductDetailRepository
	ProductPhotoRepository  repository.ProductPhotoRepository
	ProductRatingRepository repository.ProductRatingRepository
}

func (service *ProductServiceImpl) Get(id string) model.Product {

	product, err := service.ProductDetailRepository.Get(id)
	if err != nil {
		panic(err)
	}
	avr, err := service.getAvrRate(product.Id)
	if err != nil {
		panic(err)
	}
	product.AvrRate = avr
	return product
}

func (service *ProductServiceImpl) getAvrRate(productId string) (avr float64, err error) {

	rs, err := service.ProductRatingRepository.GetAvr(productId)
	if err != nil {
		return 5, err
	}
	avr = rs.AvrRate
	return
}

func getTotalPage(pageSize int, quantity int) (total int) {
	total = quantity/pageSize + 1 // equal round up
	return
}

func (service *ProductServiceImpl) listProductPhoto(products []model.Product) (photos []model.ProductPhoto, err error) {
	var productIds []string
	for _, product := range products {
		productIds = append(productIds, product.Id)
	}

	photos, err = service.ProductPhotoRepository.List(productIds)
	return
}

func addPhotosToProduct(photos []model.ProductPhoto, products *[]model.Product) {
	photoMap := make(map[string]model.ProductPhoto)
	for _, photo := range photos {
		photoMap[photo.ProductId] = photo
	}

	for i, product := range *products {
		if photo, ok := photoMap[product.Id]; ok {
			(*products)[i].Photos = append((*products)[i].Photos, photo)
		}
	}
}

func (service *ProductServiceImpl) List(productsRequest request.ListProductsRequest) response.ListProductResponse {

	priceRange := model.RangeValue[int64]{
		From: productsRequest.MinPrice,
		To:   productsRequest.MaxPrice,
	}

	products, err := service.ProductDetailRepository.List(
		[]string{},
		productsRequest.KeyWord,
		productsRequest.Tags,
		productsRequest.Brands,
		productsRequest.Types,
		productsRequest.Genders,
		priceRange, (productsRequest.Page)*productsRequest.PageSize, productsRequest.PageSize,
	)
	if err != nil {
		panic(err)
	}

	productPhotos, err := service.listProductPhoto(products)

	if err != nil {
		panic(err)
	}

	addPhotosToProduct(productPhotos, &products)
	return response.ListProductResponse{
		Products:  products,
		TotalPage: getTotalPage(productsRequest.PageSize, len(products)),
		Page:      productsRequest.Page,
	}
}
func NewProductServiceImpl(productDetailRepo repository.ProductDetailRepository,
	productRatingRepo repository.ProductRatingRepository,
	productPhotoRepo repository.ProductPhotoRepository) ProductService {
	return &ProductServiceImpl{
		productDetailRepo,
		productPhotoRepo,
		productRatingRepo,
	}
}
