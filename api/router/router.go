package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func SetUp(s *gin.Engine, container *dig.Container) {
	InitHealthRouter(s, container)
	InitProductDetailRouter(s, container)
	InitCartRouter(s, container)
	InitOrderRouter(s, container)
	InitAuthRouter(s, container)
	InitUserRouter(s, container)
	InitInventoryRouter(s, container)
	InitPhotoRouter(s, container)
	InitChatRouter(s, container)
}
