package flow_test

import (
	"encoding/json"
	"fmt"
	container "online_fashion_shop/api"
	"online_fashion_shop/api/model/cart"
	"online_fashion_shop/api/model/payment"
	"online_fashion_shop/api/repository"
	"online_fashion_shop/api/service"
	"online_fashion_shop/api/worker"
	"online_fashion_shop/initializers/zalopay"
	"testing"
	"time"
)

func TestFlow(t *testing.T) {
	var cartService service.CartService
	container.Inject(&cartService)
	customerId := "test_customer"
	var orderService service.OrderService
	container.Inject(&orderService)
	couponCode := "TEST"
	var processor zalopay.Processor
	container.Inject(&processor)
	var orderRepo repository.OrderRepository
	container.Inject(&orderRepo)
	var paymentID string
	item1 := cart.CartItem{
		CustomerId: customerId,
		ProductId:  "642965d8d2aea624550b98ce",
		Quantity:   10,
		Color:      "RED",
		Size:       "XXL",
	}
	item2 := cart.CartItem{
		CustomerId: customerId,
		ProductId:  "643909b67f100dc4309bcd71",
		Quantity:   10,
		Color:      "RED",
		Size:       "XXL",
	}

	t.Run("Add Item", func(t *testing.T) {
		cartService.Add(customerId, item1)
		cartService.Add(customerId, item2)
	})

	t.Run("List", func(t *testing.T) {
		_, _, err := orderRepo.ListByCustomerId(customerId, 10, 0)
		if err != nil {
			t.Errorf("false")
		}
	})

	t.Run("Create Order", func(t *testing.T) {
		info, err := orderService.Create(customerId, payment.ZaloPayMethod, "KHU_A", &couponCode)
		if err != nil {
			panic(err)
		}
		paymentID = info.PaymentInfo.PaymentId
		s, _ := json.MarshalIndent(info, "", "\t")
		fmt.Print(string(s))
	})

	t.Run("GET Orders", func(t *testing.T) {
		orderService.ListByCustomerID(customerId, 10, 0)
	})
	println(paymentID)
	t.Run("get Status Orders", func(t *testing.T) {
		status, err := processor.GetPaymentStatus("230430_823299")
		if err != nil {
			panic(err)
		}
		println(status)
	})

	t.Run("Test update status ", func(t *testing.T) {
		service.UpdateOrderTask(orderRepo, processor)
	})

	t.Run("worker_flow", func(t *testing.T) {
		worker.Run()
		var i = 0
		worker.AddTask(3, func(msg string) {
			println(msg)
			i++
		}, "my_task")
		time.Sleep(7 * time.Second)
		if i != 3 {
			t.Errorf("worker")
		}
	})
}
