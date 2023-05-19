package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online_fashion_shop/api/common/errs"
	"online_fashion_shop/api/model/product"
	"online_fashion_shop/api/model/response"
	"online_fashion_shop/api/service"
)

type InventoryController struct {
	Service service.ProductQuantityService
}

// Create
//
//	@Summary		init product quantity for product size and color
//	@Description	New product quantity by size and color of product If the quantity of an existing one with the same size and color or not found product detail will throw error.
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Param          request     		body       product.ProductQuantity   true    "id of product detail"
//	@Success		200				{object}	response.BaseResponse[any]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/inventory [put]
func (ic InventoryController) Create(ctx *gin.Context) {
	var newQuantity product.ProductQuantity
	ctx.BindJSON(&newQuantity)
	rs, err := ic.Service.Create(newQuantity)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, response.BaseResponse[*product.ProductQuantity]{Data: rs, Message: "create new one", Status: "success"})
}

// DeleteById
//
//	@Summary		Delete many by quantity_id
//	@Description	Delete many by quantity_id
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Param          quantity_id     		path       string   true    "quantity_id represent quantity of size & color & detail_id(product_id)"
//	@Success		200				{object}	response.BaseResponse[string]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/inventory/{quantity_id} [delete]
func (ic InventoryController) DeleteById(ctx *gin.Context) {
	quantityId := ctx.Param("quantity_id")

	err := ic.Service.DeleteOne(quantityId)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, response.BaseResponse[string]{
		Data:    "deleted one",
		Message: "deleted one",
		Status:  "success",
	})
}

// DeleteByDetailId
//
//	@Summary		Delete many by detail_id
//	@Description	Delete many by detail_id
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Param          detail_id     		path       string   true    "id of product detail"
//	@Success		200				{object}	response.BaseResponse[string]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/inventory/with_detail/{detail_id} [delete]
func (ic InventoryController) DeleteByDetailId(ctx *gin.Context) {
	quantityId := ctx.Param("detail_id")
	err := ic.Service.DeleteManyByDetailId(quantityId)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, response.BaseResponse[string]{
		Data:    "deleted ones",
		Message: "deleted many by detail_id",
		Status:  "success",
	})
}

// Update
//
//	@Summary		add more product to inventory
//	@Description	Increase the amount of a product by product_quantity_id
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Param          request     	body       product.ProductQuantity   true    "quantity_id represent for amount of product's size & color "
//	@Success		200				{object}	response.BaseResponse[product.ProductQuantity]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/inventory [post]
func (ic InventoryController) Update(ctx *gin.Context) {
	var updateQuantity product.ProductQuantity
	err := ctx.BindJSON(&updateQuantity)
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "invalid param")
	}

	rs, err := ic.Service.Update(updateQuantity)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
	}

	ctx.JSON(http.StatusOK, response.BaseResponse[*product.ProductQuantity]{
		Data:    rs,
		Message: "updated one",
		Status:  "success",
	})

}

// ListByDetailId
//
//	@Summary		List by detail_id(product_id)
//	@Description	List product's quantity
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Param          detail_id     		path       string   true    "id of product detail"
//	@Success		200				{object}	response.BaseResponse[[]product.ProductQuantity]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/inventory/with_detail/{detail_id} [get]
func (ic InventoryController) ListByDetailId(ctx *gin.Context) {
	detailId := ctx.Param("detail_id")
	rs, err := ic.Service.ListByDetailId(detailId)

	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse[[]*product.ProductQuantity]{
		Data:    rs,
		Message: "list product by detail_id",
		Status:  "success",
	})
}

// Get
//
//	@Summary		Get one by quantity_id
//	@Description	Get one by quantity_id
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Param          quantity_id     		path       string   true    "id of product detail"
//	@Success		200				{object}	response.BaseResponse[product.ProductQuantity]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/inventory/{quantity_id} [get]
func (ic InventoryController) Get(ctx *gin.Context) {
	quantityId := ctx.Param("quantity_id")
	rs, err := ic.Service.Get(quantityId)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse[*product.ProductQuantity]{
		Data:    rs,
		Message: "list product by detail_id",
		Status:  "success",
	})

}
