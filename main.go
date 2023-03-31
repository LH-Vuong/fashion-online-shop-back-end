package main

import (
	container "online_fashion_shop/api"
	"online_fashion_shop/api/router"

	"github.com/gin-gonic/gin"
)

// @title           Online Shop API
// @version         1.0
// @description     A online clothes shop service API in Go using Gin framework.
// @termsOfService  https://tos.santoshk.dev

// @contact.name   vangxitrum
// @contact.url    https://github.com/vangxitrum
// @contact.email  tuantq666@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
func main() {

	server := gin.Default()

	container := container.BuildContainer()
	router.SetUp(server, container)
	server.Run(":8080")

}
