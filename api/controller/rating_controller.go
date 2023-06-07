package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online_fashion_shop/api/common/errs"
	"online_fashion_shop/api/model/rating"
	"online_fashion_shop/api/model/response"
	"online_fashion_shop/api/service"
)

type RatingController struct {
	Service service.RatingService
}

// Get Rating
//
//	@Summary		retrieve rating by rating id
//	@Description	retrieve rating by rating id
//	@Tags			rating
//	@Accept			json
//	@Produce		json
//	@Param          id     	path       string    true    "rating's id"
//	@Success		200				{object}	response.BaseResponse[rating.Rating]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/rating/{id} [get]
func (controller *RatingController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	rate, err := controller.Service.Get(id)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, response.BaseResponse[*rating.Rating]{
		Data:   rate,
		Status: "success",
	})
}

// List product's rating
//
//	@Summary		retrieve ratings by product id
//	@Description	retrieve ratings by product id
//	@Tags			rating
//	@Accept			json
//	@Produce		json
//	@Param          product_id     	path       string    true    "product's id"
//	@Success		200				{object}	response.BaseResponse[[]rating.Rating]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/ratings/{product_id} [get]
func (controller *RatingController) List(ctx *gin.Context) {
	productId := ctx.Param("product_id")
	rates, err := controller.Service.List(productId)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, response.BaseResponse[[]*rating.Rating]{
		Data:   rates,
		Status: "success",
	})
}

// Create new rating
//
//	@Summary		add a rating for product
//	@Description	rate to a product of the current Customer by using their access token.
//	@Tags			rating
//	@Accept			json
//	@Produce		json
//	@Param          rating_content     	body       rating.Rating    true    "customer's rating"
//	@Success		200				{object}	response.BaseResponse[[]rating.Rating]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/rating [put]
func (controller *RatingController) Create(ctx *gin.Context) {
	var rate rating.Rating
	err := ctx.BindJSON(&rate)
	if err != nil {
		return
	}
	err = controller.Service.Create(&rate)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, response.BaseResponse[*rating.Rating]{
		Status: "success",
		Data:   &rate,
	})
}

// Delete a rating
//
//	@Summary		Delete a rating by id
//	@Description	Delete a rating by id
//	@Tags			rating
//	@Accept			json
//	@Produce		json
//	@Param          id     	path       string    true    "customer's rating"
//	@Success		200				{object}	response.BaseResponse[string]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/rating/{id} [delete]
func (controller *RatingController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := controller.Service.Delete(id)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, response.BaseResponse[string]{
		Data:   id,
		Status: "success",
	})
}

// UpdateOne rating
//
//	@Summary		update a rating
//	@Description	update a rating, specify by id field
//	@Tags			rating
//	@Accept			json
//	@Produce		json
//	@Param          update_info     	body       rating.Rating    true    "update info"
//	@Success		200				{object}	response.BaseResponse[string]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/rating [post]
func (controller *RatingController) UpdateOne(ctx *gin.Context) {
	var updateInfo rating.Rating
	err := ctx.BindJSON(&updateInfo)
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
		return
	}
	err = controller.Service.Update(&updateInfo)
	if err != nil {
		errs.HandleFailStatus(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, response.BaseResponse[string]{
		Data:    updateInfo.Id,
		Message: "",
		Status:  "success",
	})
}
