package main

import (
	"github.com/gin-gonic/gin"
	container "online_fashion_shop/api"
	"online_fashion_shop/api/router"
)

func main() {

	server := gin.Default()
	container := container.BuildContainer()
	router.SetUp(server, container)
	server.Run(":8080")

}
