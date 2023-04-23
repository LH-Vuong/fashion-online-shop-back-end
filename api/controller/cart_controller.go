package controller

import (
	"net/http"
	"online_fashion_shop/api/common/errs"
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/model/request"
	"online_fashion_shop/api/service"
	"os/user"

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

	customer, isExit := c.MustGet("currentUser").(user.User)
	if isExit {
		errs.HandleFailStatus(c, "can not detect customer id", http.StatusUnauthorized)
		return
	}
	cartItem := make([]model.CartItem, len(updateCartRequest.Items))

	for index := range updateCartRequest.Items {
		updateItem := updateCartRequest.Items[index]
		cartItem[index] = model.CartItem{
			ProductId: updateItem.ProductId,
			Quantity:  updateItem.Quantity,
		}
	}

	err := controller.Service.Update(customer.Uid, cartItem)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	} else {
		c.JSON(200, gin.H{"status": "success", "message": "Update Cart successfully"})
	}

}

// Get Cart Items of User
//
//	@Summary		get cart item
//	@Description	get cart item by customer's id
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Success		200				{object}	[]model.CartItem
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/product [get]
func (controller CartController) Get(c *gin.Context) {

	customer, isExit := c.MustGet("currentUser").(user.User)
	if isExit {
		errs.HandleFailStatus(c, "can not detect customer id", http.StatusUnauthorized)
		return
	}
	cartItems, err := controller.Service.Get(customer.Uid)

	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	} else {
		c.JSON(200, gin.H{"status": "success", "data": cartItems})
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

	customer, isExit := c.MustGet("currentUser").(user.User)
	if isExit {
		errs.HandleFailStatus(c, "can not detect customer id", http.StatusUnauthorized)
		return
	}
	var rq request.AddItemRequest
	err := c.BindJSON(&rq)

	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	addedItem, err := controller.Service.Add(customer.Uid, model.CartItem{ProductId: rq.ProductId, Quantity: rq.Quantity})
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
//	@Param          product_id    path       string    true    "product's id"
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart/:product_id [delete]
func (controller CartController) Delete(c *gin.Context) {

	customer, isExit := c.MustGet("currentUser").(user.User)
	if isExit {
		errs.HandleFailStatus(c, "can not detect customer id", http.StatusUnauthorized)
		return
	}
	productId := c.Param("product_id")
	err := controller.Service.DeleteOne(customer.Uid, productId)

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
//	@Success		200				{object}	[]string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart/checkout [get]
func (controller CartController) CheckOut(c *gin.Context) {
	customer, isExit := c.MustGet("currentUser").(user.User)
	if isExit {
		errs.HandleFailStatus(c, "can not detect customer id", http.StatusUnauthorized)
		return
	}
	soldItemIds, err := controller.Service.CheckOut(customer.Uid)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": soldItemIds})
}
