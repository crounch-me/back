package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/Sehsyha/crounch-back/configuration"
	"github.com/Sehsyha/crounch-back/handler"
	"github.com/Sehsyha/crounch-back/storage"
)

const (
	healthPath = "/health"

	userPath  = "/users"
	loginPath = "/users/login"

	listPath = "/lists"
)

// Version represents the version of the application
var Version string

// Start launches the router which handle connection and execute the right functions
func Start(config *configuration.Config) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(cors.Default())

	hc := handler.NewContext(config)
	configureRoutes(r, hc)

	log.SetLevel(log.DebugLevel)

	log.Info("Launching awesome server")
	err := r.Run(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func configureRoutes(r *gin.Engine, hc *handler.Context) {
	// Health routes
	r.GET(healthPath, hc.Health)

	// User routes
	r.POST(userPath, hc.Signup)
	r.POST(loginPath, hc.Login)

	// List routes
	r.POST(listPath, checkAccess(hc.Storage), hc.CreateList)
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
