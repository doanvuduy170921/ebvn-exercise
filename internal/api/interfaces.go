package api

import "net/http"

// Engine defines the HTTP server interface, decoupling the routing framework from main.
type Engine interface {
	Start() error
	InitRoutes()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
