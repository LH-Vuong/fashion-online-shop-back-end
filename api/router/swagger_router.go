package router

import (
	_ "online_fashion_shop/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/dig"
)

func InitSwaggerRouter(s *gin.Engine, _ *dig.Container) {
	s.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
