package _test

import (
	"online_fashion_shop/api"
	"testing"
)

var ct = container.BuildContainer()

/*
	{
	  "_id": { "$oid": "6429383901f4fe96f272b209" },
	  "name": "mock_test",
	  " tags": ["HOT", "NEW", "SALE"],
	  "type": ["dress", "coach"],
	  "brand": "H&M",
	  "discount_amount": "1000",
	  "discount_percent": "0.3",
	  "gender": "MEN",
	  "price": { "$numberLong": "1000000" },
	  "description": "this product use for test purpose",
	  "created_at": {
	    "$date": { "$numberLong": "1001869200000" }
	  },
	  "created_by": "tester",
	  "updated_at": "",
	  "updated_by": "no_one"
	}
*/

//func TestGetProduct(t *testing.T) {
//
//	ct.Invoke(func(repo repository.ProductDetailRepository) {
//		repo.Get("64292d1b01f4fe96f272b208")
//		t.Fatalf("hello")
//	})
//}

/*
	{
	  "_id": { "$oid": "642941a3d2aea624550b98cc" },
	  "name": "mock_test",
	  "brand": "H&M",
	  "discount_amount": "1000",
	  "discount_percent": "0.3",
	  "gender": "MEN",
	  "price": { "$numberLong": "1000000" },
	  "description": "this product use for test purpose",
	  "created_at": {
	    "$date": { "$numberLong": "1001869200000" }
	  },
	  "created_by": "tester",
	  "updated_at": "",
	  "updated_by": "no_one",
	  "tags": ["HOT", "NEW", "SALE"],
	  "types": ["dress", "coach"]
	}
*/
func TestListProduct(t *testing.T) {
	//
	//ct.Invoke(func(repo repository.ProductDetailRepository) {
	//	products := repo.List(
	//		"mock",
	//		[]string{"HOT"},
	//		[]string{"H&M"},
	//		[]string{"dress"},
	//		[]string{"MEN"},
	//		model.RangeValue[int64]{To: 10000000000000, From: 0})
	//	println(products[0].Name)
	//	println(products[0].Brand)
	//	t.Fatalf("%s,%s", products[0].Name, products[0].Brand)
	//})
}
