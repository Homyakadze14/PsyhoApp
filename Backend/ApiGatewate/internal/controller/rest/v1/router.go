// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"log/slog"
	"net/http"

	_ "github.com/Homyakadze14/PsyhoApp/ApiGatewate/docs"
	"github.com/Homyakadze14/PsyhoApp/ApiGatewate/internal/lib/s3"

	authv1 "github.com/Homyakadze14/PsyhoApp/ApiGatewate/proto/gen/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Clients struct {
	Auth authv1.AuthClient
}

// Swagger spec:
// @title       API Gatewate
// @description API Gatewate for orbit of success services
// @version     1.0
// @schemes 	https
// @host        cookhub.space
// @BasePath    /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRouter(handler *gin.Engine, c Clients, log *slog.Logger, s3 *s3.S3Storage) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Set cors
	corsConf := cors.DefaultConfig()
	corsConf.AllowOrigins = []string{"http://localhost:5173", "http://147.45.235.14:5173"}
	corsConf.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
	corsConf.AllowCredentials = true
	handler.Use(cors.New(corsConf))

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	g := handler.Group("/api/v1")
	{
		NewAuthRoutes(log, g, c.Auth)
		NewMediaRoutes(log, g, s3)
	}
}
