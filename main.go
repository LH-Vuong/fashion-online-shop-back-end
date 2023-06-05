package main

import (
	"log"
	container "online_fashion_shop/api"
	"online_fashion_shop/api/router"
	"online_fashion_shop/api/worker"
	_ "online_fashion_shop/docs"
	"online_fashion_shop/initializers"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {

}

func main() {
	// gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3001"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization", "refresh_token"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "PATCH", "DELETE"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	container := container.BuildContainer()

	server.Use(gin.Logger())

	router.SetUp(server, container)

	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	worker.Run()

	log.Fatal(server.Run(":" + config.ServerPort))
}

// @title			Online Shop API
// @version		1.0
// @description	Online Shop API
// @termsOfService	http://swagger.io/terms/

// @contact.name	vangxitrum
// @contact.url	http://www.swagger.io/support
// @contact.email	19522482@gm.uit.edu.vn

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	/api/

// @securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
