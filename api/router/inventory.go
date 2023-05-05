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
		s.PUT("/api/inventory", c.Create)
		s.POST("/api/inventory/", c.Update)
		s.GET("/api/inventory/:quantity_id", c.Get)
		s.GET("/api/inventory/with_detail/:detail_id", c.ListByDetailId)
		s.DELETE("/api/inventory/:quantity_id", c.DeleteById)
		s.DELETE("/api/inventory/with_detail/:detail_id", c.DeleteByDetailId)
	})

}
