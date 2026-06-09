package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "lesson01-ebvn/docs"
	"lesson01-ebvn/internal/config"
	"lesson01-ebvn/internal/handler"
	"lesson01-ebvn/internal/service"
	"net/http"
)

type engine struct {
	app *gin.Engine
	cfg *config.Config
}

// NewEngine initializes a new instance of the HTTP server engine.
func NewEngine(config *config.Config) Engine {
	app := &engine{
		app: gin.Default(),
		cfg: config,
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

// InitRoutes wires up clean architecture layers and registers HTTP endpoints.
func (e *engine) initRoutes() {
	bookMarkSvc := service.NewBookMarkService(e.cfg)
	bookMarkHdl := handler.NewBookMarkHandler(bookMarkSvc)
	e.app.GET("/health-check", bookMarkHdl.HealthCheck)
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
