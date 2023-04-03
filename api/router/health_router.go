package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"online_fashion_shop/configs"
)

func InitHealthRouter(s *gin.Engine, _ *dig.Container) {

	configs.GetString("hoolo")
	s.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
