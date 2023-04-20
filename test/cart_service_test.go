package test_test

import (
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/repository"
	"online_fashion_shop/test"
	"testing"
)

func TestCartRepository(t *testing.T) {

	var cartRepo repository.CartRepository
	test.Inject(&cartRepo)
	// create some test data
	customerId := "tester_01"
	cartItem := model.CartItem{
		CustomerId: customerId,
		ProductId:  "mock_product_id",
		Quantity:   2,
	}
	cartSearchOption := repository.CartSearchOption{
		CustomerId: "tester_01",
		ProductId:  "mock_product_id",
		Id:         "",
	}

	t.Run("Create", func(t *testing.T) {
		id, err := cartRepo.Create("tester", cartItem)
		if err != nil {
			t.Errorf("Unexpected error during Create(): %v", err)
		}
		cartSearchOption.Id = id
	})

	t.Run("MultiCreate", func(t *testing.T) {
		items := []model.CartItem{cartItem}
		ids, err := cartRepo.MultiCreate(customerId, items)
		if err != nil {
			t.Errorf("Unexpected error during MultiCreate(): %v", err)
		}
		if len(ids) < 0 {
			t.Failed()
		}
	})

	t.Run("ListByCustomerId", func(t *testing.T) {
		items, err := cartRepo.ListByCustomerId(customerId)
		if err != nil {
			t.Errorf("Unexpected error during ListByCustomerId(): %v", err)
		}
		if len(items) < 0 {
			t.Errorf("Unexpected result from ListByCustomerId()")
		}
	})

	t.Run("GetBySearchOption", func(t *testing.T) {
		item, err := cartRepo.GetBySearchOption(cartSearchOption)
		if err != nil {
			t.Errorf("Unexpected error during GetBySearchOption(): %v", err)
		}
		if item == nil || item.ProductId != cartItem.ProductId {
			t.Errorf("Unexpected result from GetBySearchOption()")
		}
	})

	t.Run("DeleteByCustomerId", func(t *testing.T) {
		err := cartRepo.DeleteByCustomerId(customerId)
		if err != nil {
			t.Errorf("Unexpected error during DeleteByCustomerId(): %v", err)
		}
	})

	t.Run("DeleteOne", func(t *testing.T) {
		err := cartRepo.DeleteOne(cartItem.CustomerId, cartItem.ProductId)
		if err != nil {
			t.Errorf("Unexpected error during DeleteOne(): %v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		err := cartRepo.Update(customerId, cartItem)
		if err != nil {
			t.Errorf("Unexpected error during Update(): %v", err)
		}
	})
}
