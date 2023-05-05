package service

import (
	"fmt"
	"log"
	"online_fashion_shop/api/model/cart"
	"online_fashion_shop/api/model/coupon"
	"online_fashion_shop/api/model/order"
	"online_fashion_shop/api/model/payment"
	"online_fashion_shop/api/repository"
	"online_fashion_shop/api/worker"
	"online_fashion_shop/initializers/zalopay"
	"time"
)

type OrderService interface {
	// Create order with cart items of customer
	// If order is invalid, it will throw an error (invalid coupon,Invalid Cart_Item)
	// After creating order, the user's cart will be emptied
	Create(customerID string, paymentMethod payment.Method, addressInfo string, couponCode *string) (*order.OrderInfo, error)

	// ListByCustomerID Get list of orders (order history of customer)
	ListByCustomerID(customerID string, limit int, offset int) ([]*order.OrderInfo, int64, error)

	UpdateWithCallbackData(paymentId string, callbackData map[string]any, handle CallbackHandle) error
}

type OrderServiceImpl struct {
	CouponService CouponService
	CartService   CartService
	OrderRepo     repository.OrderRepository
	Processor     zalopay.Processor
}

func NewOrderServiceImpl(couponService CouponService,
	cartService CartService,
	orderRepo repository.OrderRepository,
	processor zalopay.Processor) OrderService {

	worker.AddTask(15*60, UpdateOrderTask, orderRepo, processor)
	return &OrderServiceImpl{
		CouponService: couponService,
		CartService:   cartService,
		OrderRepo:     orderRepo,
		Processor:     processor,
	}
}

type CallbackHandle func(info *order.OrderInfo, data map[string]any) error

func (svc *OrderServiceImpl) UpdateWithCallbackData(paymentId string, data map[string]any, handle CallbackHandle) error {

	orderInfo, err := svc.OrderRepo.GetOneByPaymentId(paymentId)
	if err != nil {
		return err
	}

	handle(orderInfo, data)

	err = svc.OrderRepo.Update(*orderInfo)
	if err != nil {
		return err
	}
	return nil
}

func (svc *OrderServiceImpl) Create(customerID string, paymentMethod payment.Method, addressInfo string, couponCode *string) (*order.OrderInfo, error) {
	// Check cart has any invalid Item
	invalidItems, err := svc.CartService.ListInvalidCartItem(customerID)
	if err != nil {
		return nil, err
	}
	if len(invalidItems) > 0 {
		return nil, fmt.Errorf("Invalid Cart Item")
	}

	// Get the customer's cart items.
	cartItems, err := svc.CartService.Get(customerID)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Check the coupon code, if provided.
	coupon, err := svc.getCoupon(couponCode)
	if err != nil {
		return nil, err
	}

	// Prepare the order information.
	total, err := calculateTotal(cartItems, coupon)
	if err != nil {
		return nil, err
	}
	paymentInfo := payment.PaymentDetail{
		PaymentId:     "",
		Status:        payment.StatusInit,
		OrderUrl:      nil,
		CreatedAt:     time.Now().UnixMilli(),
		PaymentMethod: paymentMethod,
	}

	orderInfo := order.OrderInfo{
		CustomerId:     customerID,
		Address:        addressInfo,
		CouponCode:     couponCode,
		CouponDiscount: coupon.DiscountAmount,
		TotalPrice:     total,
		Items:          cartItems,
		PaymentInfo:    &paymentInfo,
	}

	if paymentMethod == payment.ZaloPayMethod {
		err = svc.Processor.InitPayment(&orderInfo)
	}
	if paymentMethod == payment.CODMethod {
		orderInfo.PaymentInfo.PaymentMethod = paymentMethod
		orderInfo.PaymentInfo.Status = payment.StatusApproved
	}

	if err != nil {
		orderInfo.PaymentInfo.PaymentAt = time.Now().UnixMilli()
		orderInfo.PaymentInfo.Status = payment.StatusError
	} else {
		err := svc.CartService.DeleteAll(customerID)
		if err != nil {
			return nil, err
		}
	}

	return svc.OrderRepo.Create(orderInfo)

}

func (svc *OrderServiceImpl) getCoupon(couponCode *string) (*coupon.CouponInfo, error) {

	if couponCode == nil {
		return nil, nil
	}

	valid, err := svc.CouponService.Check(*couponCode)

	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, fmt.Errorf("invalid coupon code %s", couponCode)
	}
	return svc.CouponService.Get(*couponCode)
}

func calculateTotal(items []*cart.CartItem, coupon *coupon.CouponInfo) (int64, error) {
	var total int64

	for _, item := range items {
		total += int64(item.Quantity) * item.ProductDetail.Price
	}

	if coupon != nil {
		if coupon.DiscountAmount > 0 {
			total -= coupon.DiscountAmount
		} else if coupon.DiscountPercent > 0 {
			total *= int64(1.0 - coupon.DiscountPercent)
		}
	}

	if total < 0 {
		return 0, fmt.Errorf("invalid total price")
	}

	return total, nil
}

func (svc *OrderServiceImpl) ListByCustomerID(customerID string, limit int, offset int) ([]*order.OrderInfo, int64, error) {
	return svc.OrderRepo.ListByCustomerId(customerID, limit, offset)
}

func UpdateOrderTask(orderRepo repository.OrderRepository, processor zalopay.Processor) {

	orders, err := orderRepo.ListPendingOrder()
	if err != nil {
		log.Println("Error When svc.OrderRepo.ListPendingOrder")
	}
	for _, orderInfo := range orders {
		if orderInfo.PaymentInfo.PaymentMethod == payment.ZaloPayMethod {
			status, err := processor.GetPaymentStatus(orderInfo.PaymentInfo.PaymentId)
			if err != nil {
				log.Println("Error When svc.Processor.GetPaymentStatus")
			}
			orderInfo.PaymentInfo.Status = status
			err = orderRepo.Update(*orderInfo)
			if err != nil {
				log.Println("Error When err = orderRepo.Update(*orderInfo)")
			}
		}

	}

}
