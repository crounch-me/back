package router

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/crounch-me/back/configuration"
	_ "github.com/crounch-me/back/docs"
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"
	"github.com/crounch-me/back/handler"
	"github.com/crounch-me/back/util"
)

const (
	healthPath = "/health"

	userPath  = "/users"
	loginPath = "/users/login"
  mePath    = "/me"
  logoutPath = "/logout"

	listPath       = "/lists"
	listWithIDPath = "/lists/:listID"

	listProductPath = "/lists/:listID/products/:productID"

	productPath       = "/products"
	productSearchPath = "/products/search"
)

// Version represents the version of the application
var Version string

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

  hc := handler.NewContext(config)
  corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	configureRoutes(r, hc)

	r.Use(cors.New(corsConfig))
	r.Use(gin.Recovery())
	log.SetLevel(log.DebugLevel)

  url := ginSwagger.URL("http://localhost:3000/swagger/doc.json")

  r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	log.Info("Launching awesome server")
	err := r.Run(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func emptyHandler(c *gin.Context) {}

func configureRoutes(r *gin.Engine, hc *handler.Context) {
	r.Use(otherMethodsHandler())

	// Health routes
	r.GET(healthPath, hc.Health)
	r.OPTIONS(healthPath, optionsHandler([]string{http.MethodGet}))

	// User routes
	r.POST(userPath, hc.Signup)
	r.OPTIONS(userPath, optionsHandler([]string{http.MethodPost}))
	r.POST(loginPath, hc.Login)
	r.OPTIONS(loginPath, optionsHandler([]string{http.MethodPost}))
	r.GET(mePath, checkAccess(hc.Storage), hc.Me)
  r.OPTIONS(mePath, optionsHandler([]string{http.MethodGet}))
  r.POST(logoutPath, checkAccess(hc.Storage), hc.Logout)
  r.OPTIONS(logoutPath, optionsHandler([]string{http.MethodPost}))

	// List routes
	r.POST(listPath, checkAccess(hc.Storage), hc.CreateList)
	r.GET(listPath, checkAccess(hc.Storage), hc.GetOwnerLists)
	r.OPTIONS(listPath, optionsHandler([]string{http.MethodPost, http.MethodGet}))

	r.GET(listWithIDPath, checkAccess(hc.Storage), hc.GetList)
	r.DELETE(listWithIDPath, checkAccess(hc.Storage), hc.DeleteList)
	r.OPTIONS(listWithIDPath, optionsHandler([]string{http.MethodDelete}))

	// Product routes
	r.POST(productPath, checkAccess(hc.Storage), hc.CreateProduct)
	r.OPTIONS(productPath, optionsHandler([]string{http.MethodPost}))

	r.POST(productSearchPath, checkAccess(hc.Storage), hc.SearchDefaultProducts)
	r.OPTIONS(productSearchPath, optionsHandler([]string{http.MethodPost}))

	// List product routes
	r.POST(listProductPath, checkAccess(hc.Storage), hc.AddProductToList)
	r.DELETE(listProductPath, checkAccess(hc.Storage), hc.DeleteProductFromList)
	r.PATCH(listProductPath, checkAccess(hc.Storage), hc.UpdateProductInList)
	r.OPTIONS(listProductPath, optionsHandler([]string{http.MethodPost, http.MethodPatch, http.MethodDelete}))
}

func checkAccess(us users.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			log.Info("Unauthorized - No token provided")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, err := us.GetUserIDByToken(token)

		if err != nil {
			if err.Code == users.UserNotFoundErrorCode {
				c.AbortWithStatusJSON(http.StatusUnauthorized, domain.NewError(domain.UnauthorizedErrorCode))
				return
			}
			log.WithError(err).Error("Unauthorized - Error while accessing database")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Set(handler.ContextUserID, *userID)
	}
}

func optionsHandler(allowedMethods []string) gin.HandlerFunc {
	allowedMethods = append(allowedMethods, http.MethodOptions)
	allowedHeaders := []string{util.HeaderContentType, util.HeaderAuthorization, util.HeaderAccept}
	return func(c *gin.Context) {
		c.Writer.Header().Set(util.HeaderAccessControlAllowOrigin, "*")
		c.Writer.Header().Set(util.HeaderAccessControlAllowMethods, strings.Join(allowedMethods, ","))
		c.Writer.Header().Set(util.HeaderAccessControlAllowHeaders, strings.Join(allowedHeaders, ","))
		c.AbortWithStatus(http.StatusOK)
	}
}

func otherMethodsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set(util.HeaderAccessControlAllowOrigin, "*")
		c.Next()
	}
}
