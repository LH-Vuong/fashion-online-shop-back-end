package controller

import (
	"github.com/gin-gonic/gin"
	"online_fashion_shop/api/service"
)

type ProductController struct {
	Service service.ProductService
}

func (cl ProductController) Get(c *gin.Context) {
	productId := c.Param("id")
	product := cl.Service.Get(productId)
	c.JSON(200, product)
	return
}
