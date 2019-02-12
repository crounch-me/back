package router

import (
	"github.com/Sehsyha/crounch-back/handler"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	healthPath = "/health"

	userPath = "/users"
)

func Start() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.HandleMethodNotAllowed = true

	hc := handler.NewContext()
	configureRoutes(r, hc)

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
}
