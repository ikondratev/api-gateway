package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ikondratev/api-gateway/internal/logger"
)

type Router struct {
	engine *mux.Router
}

func NewRouter(logger logger.Logger) *Router {
	router := mux.NewRouter()
	router.Use(Logging(logger))

	return &Router{
		engine: router,
	}
}

func (r *Router) RegisterRoutes() *mux.Router {
	r.engine.HandleFunc("/ping", Pong).Methods(http.MethodGet)
	
	return r.engine
}