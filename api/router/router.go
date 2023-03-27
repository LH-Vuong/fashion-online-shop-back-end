package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func SetUp(s *gin.Engine, container *dig.Container) {
	InitHealthRouter(s, container)
	InitProductRouter(s, container)
}
