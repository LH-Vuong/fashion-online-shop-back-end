package controller

import (
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	Service service.ProductService
}

// ShowAccount godoc
//
//	@Summary		Show an account
//	@Description	get string by ID
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			id path	int	true	"Product ID"
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/products/:id [get]
func (cl ProductController) Get(c *gin.Context) {
	productId := c.Param("id")
	product := cl.Service.Get(productId)
	c.JSON(200, product)
}
