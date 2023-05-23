package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"online_fashion_shop/api/controller"
	"online_fashion_shop/api/service"
)

func InitInventoryRouter(s *gin.Engine, c *dig.Container) {

	c.Invoke(func(sv service.ProductQuantityService) {
		c := controller.InventoryController{Service: sv}
		inventoryRouter := s.Group("api/inventory")

		inventoryRouter.PUT("", c.Create)
		inventoryRouter.POST("", c.Update)
		inventoryRouter.GET("/:quantity_id", c.Get)
		inventoryRouter.GET("/with_detail/:detail_id", c.ListByDetailId)
		inventoryRouter.DELETE("/:quantity_id", c.DeleteById)
		inventoryRouter.DELETE("/with_detail/:detail_id", c.DeleteByDetailId)
	})

}
