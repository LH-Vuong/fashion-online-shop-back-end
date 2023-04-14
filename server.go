package main

import (
	"github.com/gin-gonic/gin"
	container "online_fashion_shop/api"
	"online_fashion_shop/api/router"
	"online_fashion_shop/configs"
)

func main() {
	server := gin.Default()
	container := container.BuildContainer()
	server.Use(gin.Logger())
	router.SetUp(server, container)
	server.Run(configs.GetString("server.port"))
}
