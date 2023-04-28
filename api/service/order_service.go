package service

import (
	"fmt"
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/model/payment"
	"online_fashion_shop/api/repository"
	"time"
)

type OrderService interface {
	// Create order with cart items of customer
	// If order is invalid, it will throw an error (invalid coupon,Invalid Cart_Item)
	// After creating order, the user's cart will be emptied
	Create(customerID string, paymentMethod payment.Method, addressInfo string, couponCode *string) (*model.OrderInfo, error)

	// ListByCustomerID Get list of orders (order history of customer)
	ListByCustomerID(customerID string, limit int, offset int) ([]*model.OrderInfo, error)
}

type OrderServiceImpl struct {
	CouponService  CouponService
	CartService    CartService
	OrderRepo      repository.OrderRepository
	paymentService PaymentService
}

// init payment and modify order_info
func (svc *OrderServiceImpl) PaymentProcessing(orderInfo *model.OrderInfo) error {

	return nil
}

func (svc *OrderServiceImpl) Create(customerID string, paymentMethod payment.Method, addressInfo string, couponCode *string) (*model.OrderInfo, error) {
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
	paymentInfo := payment.Detail{
		Id:            "",
		OrderId:       "",
		Status:        payment.StatusInit,
		OrderUrl:      nil,
		CreatedAt:     time.Now().UnixMilli(),
		PaymentMethod: paymentMethod,
	}

	order := model.OrderInfo{
		CustomerId:     customerID,
		Address:        addressInfo,
		CouponCode:     *couponCode,
		CouponDiscount: coupon.DiscountAmount,
		TotalPrice:     total,
		Items:          cartItems,
		PaymentInfo:    &paymentInfo,
	}

	if err = svc.paymentService.InitPayment(&order); err != nil {
		order.PaymentInfo.LastUpdateAt = time.Now().UnixMilli()
		order.PaymentInfo.Status = payment.StatusError
	}

	return svc.OrderRepo.Create(order)

}

func (svc *OrderServiceImpl) getCoupon(couponCode *string) (*model.CouponInfo, error) {

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

func calculateTotal(items []*model.CartItem, coupon *model.CouponInfo) (int64, error) {
	var total int64

	for _, item := range items {
		total += int64(item.Quantity) * item.ProductDetail.Price
	}

	if coupon != nil {
		if coupon.DiscountAmount > 0 {
			total -= coupon.DiscountAmount
		} else if coupon.DiscountPercent > 0 {
			total *= int64(1.0 - coupon.DiscountPercent/100.0)
		}
	}

	if total < 0 {
		return 0, fmt.Errorf("invalid total price")
	}

	return total, nil
}

func (svc *OrderServiceImpl) ListByCustomerID(customerID string, limit int, offset int) ([]*model.OrderInfo, error) {
	return svc.OrderRepo.ListByCustomerId(customerID, limit, offset)
}
