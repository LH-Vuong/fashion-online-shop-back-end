package controller

import (
	"github.com/gin-gonic/gin"
	"online_fashion_shop/api/service"
)

type RatingController struct {
	service service.RatingService
}

func (controller *RatingController) Get(ctx *gin.Context) {
}

func (controller *RatingController) List(ctx *gin.Context) {

}

func (controller *RatingController) Create(ctx *gin.Context) {

}

func (controller *RatingController) delete(ctx *gin.Context) {

}
