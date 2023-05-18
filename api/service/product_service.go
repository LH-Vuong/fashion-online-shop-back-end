package service

import (
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/model/product"
	"online_fashion_shop/api/model/request"
	"online_fashion_shop/api/repository"
)

type ProductService interface {
	Get(id string) (*product.Product, error)
	List(request request.ListProductsRequest) ([]*product.Product, int64, error)
	Update(updateInfo *product.Product) error
	Create(productInfo *product.Product) error
}

type ProductServiceImpl struct {
	ProductDetailRepository repository.ProductDetailRepository
	PhotoService            PhotoService
	ProductRatingRepository repository.ProductRatingRepository
	QuantityService         ProductQuantityService
}

func ConvertPhotosToUrls(photos []*product.ProductPhoto) []string {

	var urls []string

	for _, photo := range photos {
		if photo.MainPhoto != "" {
			urls = append(urls, photo.MainPhoto)
		}
		for _, subPhoto := range photo.SubPhotos {
			urls = append(urls, subPhoto)
		}
	}
	return urls

}

func (service *ProductServiceImpl) Update(updateInfo *product.Product) error {
	return service.ProductDetailRepository.Update(updateInfo)
}

func (service *ProductServiceImpl) Create(productInfo *product.Product) error {
	return service.ProductDetailRepository.Create(productInfo)
}

func (service *ProductServiceImpl) Get(id string) (*product.Product, error) {

	product, err := service.ProductDetailRepository.Get(id)
	if err != nil {
		return nil, err
	}
	avr, err := service.getAvrRate(product.Id)
	if err != nil {
		return nil, err
	}
	product.AvrRate = avr
	photos, err := service.PhotoService.ListByProductId(id)
	if err != nil {
		return nil, err
	}
	quantities, err := service.QuantityService.ListByDetailId(id)
	if err != nil {
		return nil, err
	}
	product.ProductQuantities = quantities
	product.Photos = ConvertPhotosToUrls(photos)
	return product, nil
}

func (service *ProductServiceImpl) getAvrRate(productId string) (avr float64, err error) {

	rs, err := service.ProductRatingRepository.GetAvr(productId)
	if err != nil {
		return 5, err
	}
	avr = rs.AvrRate
	return
}

func (s *ProductServiceImpl) addPhotosToProduct(products []*product.Product) error {

	productIds := make([]string, len(products))
	for index := range products {
		productIds[index] = products[index].Id
	}

	photoMap, err := s.PhotoService.ListByMultiProductId(productIds)
	if err != nil {
		return err
	}

	for _, product := range products {
		if photos, ok := photoMap[product.Id]; ok {
			product.Photos = ConvertPhotosToUrls(photos)
		}
	}
	return nil
}

func (service *ProductServiceImpl) List(productsRequest request.ListProductsRequest) ([]*product.Product, int64, error) {

	priceRange := model.RangeValue[int64]{
		From: productsRequest.MinPrice,
		To:   productsRequest.MaxPrice,
	}

	products, totalDocs, err := service.ProductDetailRepository.List(
		[]string{},
		productsRequest.KeyWord,
		productsRequest.Tags,
		productsRequest.Brands,
		productsRequest.Types,
		productsRequest.Genders,
		priceRange, (productsRequest.Page)*productsRequest.PageSize, productsRequest.PageSize,
	)
	if err != nil {
		return nil, 0, err
	}
	service.addPhotosToProduct(products)

	return products, totalDocs, nil
}
func NewProductServiceImpl(productDetailRepo repository.ProductDetailRepository,
	productRatingRepo repository.ProductRatingRepository,
	photoService PhotoService, quantityService ProductQuantityService) ProductService {
	return &ProductServiceImpl{
		productDetailRepo,
		photoService,
		productRatingRepo,
		quantityService,
	}
}
