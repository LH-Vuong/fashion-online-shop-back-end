package controller

import (
	"net/http"
	"online_fashion_shop/api/common/errs"
	"online_fashion_shop/api/model/cart"
	"online_fashion_shop/api/model/request"
	"online_fashion_shop/api/model/response"
	model "online_fashion_shop/api/model/user"
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	Service service.CartService
}

func ToCartResponse(cartItems []*cart.CartItem) []*response.CartItem {
	responseItem := make([]*response.CartItem, len(cartItems))

	for index, item := range cartItems {
		var image string
		if len(item.ProductDetail.Photos) > 0 {
			image = item.ProductDetail.Photos[0]
		}

		responseItem[index] = &response.CartItem{
			Id:              item.InventoryId,
			Name:            item.ProductDetail.Name,
			Image:           image,
			Price:           float64(item.ProductDetail.Price),
			DiscountPercent: item.ProductDetail.DiscountPercent,
			Size:            item.Size,
			Color:           item.Color,
			Quantity:        item.Quantity,
		}
	}
	return responseItem
}

// Update Cart Items of User
//
//	@Summary		Update the cart of the current customer with the items received in the request body(replay)
//	@Description	Delete all the previous cart items of the customer by using their access token then add the items received in the request body to their cart.
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Param          CartRequest   body       []request.CartItemUpdater    true    "Array of cart items to be added to the customer's cart"
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart [post]
func (controller CartController) Update(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(model.User)
	var updaters []request.CartItemUpdater
	c.BindJSON(&updaters)

	err := controller.Service.Update(currentUser.Id, updaters)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"status": "success"})

}

// Get Cart Items of Customer
//
//	@Summary		Get the cart items of the current user
//	@Description	Retrieve a list of cart items for the current Customer by using their access token.
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Success		200				{object}	response.BaseResponse[[]response.CartItem]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart [get]
func (controller CartController) Get(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(model.User)
	cartItems, err := controller.Service.Get(currentUser.Id)
	responseItems := ToCartResponse(cartItems)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(200, response.BaseResponse[[]*response.CartItem]{
		Data:    responseItems,
		Message: "",
		Status:  "success",
	})

}

// AddMany items to card
//
//	@Summary		Add multiple items to the cart
//	@Description	Adds multiple items to the cart and returns the information of the added items ,by using their access token.If an item already exists in the cart,its quantity will be updated.
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Param          CartRequest   body      []request.CartItemUpdater   true    "Array of cart items to be added to the cart"
//	@Success		200				{object}	response.BaseResponse[[]response.CartItem]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart [put]
func (controller CartController) AddMany(c *gin.Context) {
	var items []request.CartItemUpdater
	currentUser := c.MustGet("currentUser").(model.User)
	err := c.BindJSON(&items)

	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}

	addedItems, err := controller.Service.AddMany(currentUser.Id, items)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(200, response.BaseResponse[[]*response.CartItem]{
		Data:    ToCartResponse(addedItems),
		Message: "",
		Status:  "success",
	})

}

// Delete one item in cart
//
//	@Summary		Delete a single item from the cart
//	@Description	Deletes a cart item specified by the product ID and the customer authentication code.
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Param          cart_id  path       string    true    "cart id"
//	@Success		200				{object}	response.BaseResponse[string]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart/{cart_id} [delete]
func (controller CartController) Delete(c *gin.Context) {
	id := c.Param("cart_id")
	currentUser := c.MustGet("currentUser").(model.User)
	deletedCartId, err := controller.Service.DeleteOneById(id, currentUser.Id)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, response.BaseResponse[string]{
		Data:    deletedCartId,
		Message: "",
		Status:  "success",
	})
}

// DeleteMany  items in customer cart
//
//	@Summary		Delete all items from the cart
//	@Description	Deletes all items from the customer's cart by using their access token
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Success		200				{object}	response.BaseResponse[[]string] "Array of items' id was deleted"
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart [delete]
func (controller CartController) DeleteMany(c *gin.Context) {

	currentUser := c.MustGet("currentUser").(model.User)

	ids, err := controller.Service.DeleteAll(currentUser.Id)

	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, response.BaseResponse[[]string]{
		Data:    ids,
		Message: "",
		Status:  "success",
	})
}

// CheckOut items in cart
//
//	@Summary		Validate and modify the items in the cart
//	@Description	Validates the items in the customer's cart and modifies them if any items are invalid, such as sold-out items.
//	Returns the list of sold-out items' IDs that have been removed from the cart.
//	Use this method before placing an order to ensure that the order is valid.
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Success		200				{object}	 response.BaseResponse[[]string]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart/checkout [get]
func (controller CartController) CheckOut(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(model.User)
	soldItemIds, err := controller.Service.CheckOutAndDelete(currentUser.Id)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, response.BaseResponse[[]string]{
		Data:    soldItemIds,
		Message: "list of deleted item",
		Status:  "success",
	})
}
