package main

import (
	"log"
	container "online_fashion_shop/api"
	"online_fashion_shop/api/controller"
	repository "online_fashion_shop/api/repository/user"
	"online_fashion_shop/api/router"
	"online_fashion_shop/api/service"
	_ "online_fashion_shop/docs"
	"online_fashion_shop/initializers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	server *gin.Engine

	userRepo            repository.UserRepository
	userService         service.UserService
	authController      controller.AuthController
	authRouteController router.AuthRouteController

	userController      controller.UserController
	userRouteController router.UserRouteController
)

func init() {

}

func main() {
	server = gin.Default()
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	corConfig := cors.DefaultConfig()
	corConfig.AllowOrigins = []string{"http://localhost:8081"}
	corConfig.AllowMethods = []string{"*"}
	corConfig.AllowCredentials = true
	corConfig.AllowHeaders = []string{"Origin", "Content-Type"}

	server.Use(cors.New(corConfig))

	r := server.Group("api")

	userRepo = repository.NewUserRepository(initializers.DB)
	userService = service.NewUserServiceImpl(userRepo)
	authController = controller.NewAuthController(userService)
	authRouteController = router.NewAuthRouteController(authController)
	authRouteController.AuthRoute(r)

	userController = controller.NewUserController(userService)
	userRouteController := router.NewUserRouteController(userController)
	userRouteController.UserRoute(r)

	container := container.BuildContainer()

	server.Use(gin.Logger())

	router.SetUp(server, container)

	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(server.Run(":" + config.ServerPort))
}

//	@title			Online Shop API
//	@version		1.0
//	@description	Online Shop API
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	vangxitrum
//	@contact.url	http://www.swagger.io/support
//	@contact.email	19522482@gm.uit.edu.vn

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		127.0.0.1:8081
//	@BasePath	/api/

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
