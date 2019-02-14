package router

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/Sehsyha/crounch-back/configuration"
	"github.com/Sehsyha/crounch-back/handler"
)

const (
	healthPath = "/health"

	userPath  = "/users"
	loginPath = "/users/login"
)

var Version string

func Start(config *configuration.Config) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.HandleMethodNotAllowed = true

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
}
