package controller

import (
	"net/http"
	"online_fashion_shop/api/common/errs"
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/model/request"
	"online_fashion_shop/api/model/response"
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	Service service.CartService
}

// Update Cart Items of User
//
//	@Summary		update cart of customer by received cart item
//	@Description	Update cart item by delete all old cart items, then add received item to customer card
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Param          CartRequest   body       request.UpdateCartRequest    true    "access token received after login"
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart [post]
func (controller CartController) Update(c *gin.Context) {
	var updateCartRequest request.UpdateCartRequest
	c.BindJSON(updateCartRequest)

	cartItem := make([]model.CartItem, len(updateCartRequest.Items))

	for index := range updateCartRequest.Items {
		updateItem := updateCartRequest.Items[index]
		cartItem[index] = model.CartItem{
			CustomerId: updateCartRequest.CustomerId,
			ProductId:  updateItem.ProductId,
			Quantity:   updateItem.Quantity,
		}
	}

	err := controller.Service.Update(updateCartRequest.CustomerId, cartItem)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	} else {
		c.JSON(200, gin.H{"status": "success", "message": "Update Cart successfully"})
	}

}

// Get Cart Items of User
//
//	@Summary		List customer's cart item
//	@Description	get List cart item by access_token of user
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Success		200				{object}	response.BaseResponse[[]model.CartItem]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart [get]
func (controller CartController) Get(c *gin.Context) {
	cartItems, err := controller.Service.Get(c.Param("customer_id"))

	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	} else {
		c.JSON(200, response.BaseResponse[[]*model.CartItem]{
			Data:    cartItems,
			Message: "",
			Status:  "success",
		})
	}
}

// Add new item to card
//
//	@Summary		add item to cart
//	@Description	add item to cart, return info of added item
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Param          CartRequest   body       request.AddItemRequest    true    "Add item request"
//	@Success		200				{object}	model.CartItem
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart [put]
func (controller CartController) Add(c *gin.Context) {
	var rq request.AddItemRequest
	err := c.BindJSON(&rq)

	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	addedItem, err := controller.Service.Add(rq.CustomerId, model.CartItem{ProductId: rq.ProductId, Quantity: rq.Quantity})
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": addedItem})
	}

}

// Delete new item to card
//
//	@Summary		delete card item of customer cart
//	@Description	Delete cart item by product id of customer
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Param          customer_id   path       string    true    "customer's id"
//	@Param          product_id    path       string    true    "product's id"
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart [delete]
func (controller CartController) Delete(c *gin.Context) {

	customerId := c.Param("customer_id")
	productId := c.Param("product_id")
	err := controller.Service.DeleteOne(customerId, productId)

	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "DeleteOne successfully"})
}

// CheckOut items in cart
//
//	@Summary		Use to validate items of cart then modify it if invalid
//	@Description	Check out customer cart then delete invalid item(sold out item) return list of sold-out Items' ID
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Param          customer_id   path       string    true    "customer's id"
//	@Success		200				{object}	[]string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart/checkout [get]
func (controller CartController) CheckOut(c *gin.Context) {
	customerId := c.Param("customer_id")
	soldItemIds, err := controller.Service.CheckOutAndDelete(customerId)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": soldItemIds})
}
