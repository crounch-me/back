package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/crounch-me/back/configuration"

	// Import documentations for swagger endpoint
	_ "github.com/crounch-me/back/docs"
	accountAdapters "github.com/crounch-me/back/internal/account/adapters"
	accountApp "github.com/crounch-me/back/internal/account/app"
	accountPorts "github.com/crounch-me/back/internal/account/ports"
	"github.com/crounch-me/back/internal/common/server"
	"github.com/crounch-me/back/internal/common/utils"
	listingAdapters "github.com/crounch-me/back/internal/listing/adapters"
	listingApp "github.com/crounch-me/back/internal/listing/app"
	listingPorts "github.com/crounch-me/back/internal/listing/ports"
	productsAdapters "github.com/crounch-me/back/internal/products/adapters"
	productsApp "github.com/crounch-me/back/internal/products/app"
	productPorts "github.com/crounch-me/back/internal/products/ports"
)

// @title Crounch Me API
// @version 1.0
// @description API serving the grocery application.

// @host localhost:3000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// Start launches the router which handle connection and execute the right functions
func Start(config *configuration.Config) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	validator := utils.NewValidator()
	generationLibrary := utils.NewGeneration()
	hashLibrary := utils.NewHash()

	authorizationsRepository := accountAdapters.NewAuthorizationsMemoryRepository()
	listsRepository := listingAdapters.NewListsMemoryRepository()
	usersRepository := accountAdapters.NewUsersMemoryRepository()
	productsRepository := productsAdapters.NewProductsMemoryRepository()

	listService, err := listingApp.NewListService(listsRepository, productsRepository, generationLibrary)
	if err != nil {
		log.Fatal(err)
	}

	productService, err := productsApp.NewProductService(generationLibrary, productsRepository)
	if err != nil {
		log.Fatal(err)
	}

	accountService, err := accountApp.NewAccountService(authorizationsRepository, generationLibrary, hashLibrary, usersRepository)
	if err != nil {
		log.Fatal(err)
	}

	listServer, err := listingPorts.NewGinServer(listService, accountService, validator)
	if err != nil {
		log.Fatal(err)
	}

	userServer, err := accountPorts.NewGinServer(accountService, validator)
	if err != nil {
		log.Fatal(err)
	}

	productServer, err := productPorts.NewGinServer(productService, accountService, validator)
	if err != nil {
		log.Fatal(err)
	}

	r.Use(accessControlAllowOriginHandler())

	userServer.ConfigureRoutes(r)
	listServer.ConfigureRoutes(r)
	productServer.ConfigureRoutes(r)

	r.Use(cors.New(corsConfig))
	r.Use(gin.Recovery())
	log.SetLevel(log.DebugLevel)

	url := ginSwagger.URL("http://localhost:3000/swagger/doc.json")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	log.Info("Launching awesome server")
	err = r.Run(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func accessControlAllowOriginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set(server.HeaderAccessControlAllowOrigin, "*")
		c.Next()
	}
}
