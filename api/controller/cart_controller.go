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
			Id:       item.ProductId,
			Name:     item.ProductDetail.Name,
			Image:    image,
			Price:    float64(item.ProductDetail.Price),
			Size:     item.Size,
			Color:    item.Color,
			Quantity: item.Quantity,
		}
	}
	return responseItem
}

// Update Cart Items of User
//
//	@Summary		update cart of customer by received cart item
//	@Description	Update cart item by delete all old cart items, then add received item to customer card
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Param          CartRequest   body       []request.CartItem    true    "access token received after login"
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart [post]
func (controller CartController) Update(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(model.User)
	var newCartItems []request.CartItem
	c.BindJSON(&newCartItems)
	cartItem := make([]*cart.CartItem, len(newCartItems))

	for index := range newCartItems {
		cartItem[index] = &cart.CartItem{
			CustomerId: currentUser.Id,
			ProductId:  newCartItems[index].ProductId,
			Quantity:   newCartItems[index].Quantity,
		}
	}

	err := controller.Service.Update(currentUser.Id, cartItem)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"status": "success", "message": "Update Cart successfully"})

}

// Get Cart Items of User
//
//	@Summary		List customer's cart item
//	@Description	get List cart item by access_token of user
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

// AddMany new items to card
//
//	@Summary		add items to cart
//	@Description	add items to cart, return info of added item
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Param          CartRequest   body      []request.CartItem    true    "AddMany item request"
//	@Success		200				{object}	response.BaseResponse[[]response.CartItem]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart [put]
func (controller CartController) AddMany(c *gin.Context) {
	var items []request.CartItem
	currentUser := c.MustGet("currentUser").(model.User)
	err := c.BindJSON(&items)

	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}

	cartItem := make([]*cart.CartItem, len(items))
	for index, item := range items {
		cartItem[index] = &cart.CartItem{
			CustomerId: currentUser.Id,
			ProductId:  item.ProductId,
			Quantity:   item.Quantity,
		}
	}

	addedItems, err := controller.Service.AddMany(currentUser.Id, cartItem)
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
//	@Summary		delete one item in cart
//	@Description	Delete cart item which specifies by product_id and customer_auth code
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Param          product_id    query       string    true    "product's id"
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart/{product_id} [delete]
func (controller CartController) Delete(c *gin.Context) {
	productId := c.Param("product_id")
	currentUser := c.MustGet("currentUser").(model.User)
	err := controller.Service.DeleteOne(currentUser.Id, productId)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Delete successfully"})
}

// DeleteMany  items in customer cart
//
//	@Summary		delete all items in cart
//	@Description	Delete all item in cart
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/cart [delete]
func (controller CartController) DeleteMany(c *gin.Context) {

	currentUser := c.MustGet("currentUser").(model.User)

	err := controller.Service.DeleteAll(currentUser.Id)

	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Delete successfully"})
}

// CheckOut items in cart
//
//	@Summary		Use to validate items of cart then modify it if invalid
//	@Description	Check out customer cart then delete invalid item(sold out item) return list of sold-out Items' ID.Use before make order to ensure order will valid
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
