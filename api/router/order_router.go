package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func InitOrderRouter(s *gin.Engine, c *dig.Container) {
	//init an order
	s.PUT("/api/order")
	//list customer's order
	s.GET("/api/orders")

}
