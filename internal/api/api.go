package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lesson01-ebvn/docs"
	"time"

	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "lesson01-ebvn/docs"
	"lesson01-ebvn/internal/config"
	"lesson01-ebvn/internal/handler"
	"lesson01-ebvn/internal/repository"
	"lesson01-ebvn/internal/service"
	redispkg "lesson01-ebvn/pkg/redis"
	"net/http"
)

type engine struct {
	app *gin.Engine
	cfg *config.Config
	re  redispkg.RedisClient
}

// NewEngine initializes a new instance of the HTTP server engine.
func NewEngine(config *config.Config, re redispkg.RedisClient) Engine {
	app := &engine{
		app: gin.Default(),
		cfg: config,
		re:  re,
	}
	app.initRoutes()
	return app
}

// Start begins listening and serving HTTP requests on the configured port.
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.Port))
}

// ServeHTTP implements the http.Handler interface.
func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}

// initRoutes wires up clean architecture layers and registers HTTP endpoints.
func (e *engine) initRoutes() {
	// config CORS
	e.app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:9999", "http://localhost:8080", "*"}, // Mở cho FE port 9999
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	docs.SwaggerInfo.Host = e.cfg.HostName

	urlRepo := repository.NewUrlRepo(e.re)
	bookMarkSvc := service.NewBookMarkService(e.cfg, urlRepo)
	bookMarkHdl := handler.NewBookMarkHandler(bookMarkSvc)
	e.app.GET("/health-check", bookMarkHdl.HealthCheck)
	e.app.POST("/v1/links/shorten", bookMarkHdl.ShortenURL)
	e.app.GET("/v1/links/redirect/:code", bookMarkHdl.Redirect)
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
