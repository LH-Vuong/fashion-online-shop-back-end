package service

import (
	"math"
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/model/request"
	"online_fashion_shop/api/repository"
)

type ProductService interface {
	Get(id string) (*model.Product, error)
	List(request request.ListProductsRequest) ([]*model.Product, int, error)
}

type ProductServiceImpl struct {
	ProductDetailRepository repository.ProductDetailRepository
	PhotoService            PhotoService
	ProductRatingRepository repository.ProductRatingRepository
}

func (service *ProductServiceImpl) Get(id string) (*model.Product, error) {

	product, err := service.ProductDetailRepository.Get(id)
	if err != nil {
		return nil, err
	}
	avr, err := service.getAvrRate(product.Id)
	if err != nil {
		return nil, err
	}
	product.AvrRate = avr
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

func getTotalPage(pageSize int, quantity int) (total int) {
	total = quantity/pageSize + 1 // equal round up
	return
}

func (s *ProductServiceImpl) addPhotosToProduct(products []*model.Product) error {

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
			product.Photos = photos
		}
	}
	return nil
}

func (service *ProductServiceImpl) List(productsRequest request.ListProductsRequest) ([]*model.Product, int, error) {

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
	numPages := int(math.Ceil(float64(totalDocs) / float64(productsRequest.PageSize)))
	return products, numPages, nil
}
func NewProductServiceImpl(productDetailRepo repository.ProductDetailRepository,
	productRatingRepo repository.ProductRatingRepository,
	photoService PhotoService) ProductService {
	return &ProductServiceImpl{
		productDetailRepo,
		photoService,
		productRatingRepo,
	}
}
