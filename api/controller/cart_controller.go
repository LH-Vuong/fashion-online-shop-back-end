package controller

import (
	"net/http"
	"online_fashion_shop/api/model/request"
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	Service service.CartService
}

func (controller CartController) Update(c *gin.Context) {
	var updateCartRequest request.CartRequest
	c.BindJSON(updateCartRequest)

	err := controller.Service.Update(updateCartRequest.CustomerId, updateCartRequest.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, "Update Cart successfully")
	}

}

func (controller CartController) Get(c *gin.Context) {
	response, err := controller.Service.Get(c.Param("customer_id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func (controller CartController) Add(c *gin.Context) {
	var request request.CartRequest
	err := c.BindJSON(&request)

	if err != nil {
		panic(err)
	}
	response, err := controller.Service.Add(request.CustomerId, request.Items[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, response)
	}

}

func (controller CartController) Delete(c *gin.Context) {

	customerId := c.Param("customer_id")
	productId := c.Param("product_id")
	err := controller.Service.DeleteOne(customerId, productId)

	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, "DeleteOne successfully")

}

func (controller CartController) CheckOut(c *gin.Context) {

	customerId := c.Param("customer_id")
	res, err := controller.Service.CheckOut(customerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, res)

}
