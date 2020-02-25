package router

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/Sehsyha/crounch-back/configuration"
	"github.com/Sehsyha/crounch-back/handler"
	"github.com/Sehsyha/crounch-back/storage"
	"github.com/Sehsyha/crounch-back/util"
)

const (
	healthPath = "/health"

	userPath  = "/users"
	loginPath = "/users/login"

	listPath = "/lists"

	listProductPath = "/lists/:listID/products/:productID"

	productPath = "/products"
)

// Version represents the version of the application
var Version string

// Start launches the router which handle connection and execute the right functions
func Start(config *configuration.Config) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	hc := handler.NewContext(config)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	ConfigureRoutes(r, hc)

	r.Use(cors.New(corsConfig))
	log.SetLevel(log.DebugLevel)

	log.Info("Launching awesome server")
	err := r.Run(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func emptyHandler(c *gin.Context) {}

func ConfigureRoutes(r *gin.Engine, hc *handler.Context) {
	r.Use(otherMethodsHandler())

	// Health routes
	r.GET(healthPath, hc.Health)
	r.OPTIONS(healthPath, optionsHandler([]string{http.MethodGet}))

	// User routes
	r.POST(userPath, hc.Signup)
	r.OPTIONS(userPath, optionsHandler([]string{http.MethodPost}))
	r.POST(loginPath, hc.Login)
	r.OPTIONS(loginPath, optionsHandler([]string{http.MethodPost}))

	// List routes
	r.POST(listPath, checkAccess(hc.Storage), hc.CreateList)
	r.GET(listPath, checkAccess(hc.Storage), hc.GetOwnerLists)
	r.OPTIONS(listPath, optionsHandler([]string{http.MethodPost, http.MethodGet}))

	// Product routes
	r.POST(productPath, checkAccess(hc.Storage), hc.CreateProduct)
	r.OPTIONS(productPath, optionsHandler([]string{http.MethodPost}))

	// List product routes
	r.POST(listProductPath, checkAccess(hc.Storage), hc.AddProductToList)
	r.OPTIONS(listProductPath, optionsHandler([]string{http.MethodPost}))
}

func checkAccess(s storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			log.Info("Unauthorized - No token provided")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, err := s.GetUserIDByToken(token)

		if err != nil {
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
