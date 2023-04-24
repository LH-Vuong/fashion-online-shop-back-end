package test

import (
	"online_fashion_shop/api/repository"
	"testing"
)

func TestProductRepo(t *testing.T) {

	var productDetailRepo repository.ProductDetailRepository
	err := Inject(&productDetailRepo)
	if err != nil {
		panic(err)
	}

	//searchOption := model.ProductSearchOption{
	//	Ids:         nil,
	//	KeyWord:     "",
	//	Tags:        nil,
	//	Brands:      nil,
	//	ProductType: nil,
	//	Gender:      nil,
	//	PriceRange: model.RangeValue[int64]{
	//		From: 0,
	//		To:   math.MaxInt64,
	//	},
	//	StartAt: 0,
	//	Length:  10,
	//}

	//t.Run("List Product By Search Option", func(t *testing.T) {
	//	products, err := productDetailRepo.ListBySearchOption(searchOption)
	//	if err != nil {
	//		t.Errorf(err.Error())
	//	}
	//
	//	for _, item := range products {
	//		t.Fatalf("product %+v\n", item)
	//	}
	//
	//	if len(products) < 1 {
	//		t.Errorf("Not found any product")
	//	}
	//})

}
