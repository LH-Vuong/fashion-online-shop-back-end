package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func InitHealthRouter(s *gin.Engine, _ *dig.Container) {
	v1 := s.Group("/api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
}
