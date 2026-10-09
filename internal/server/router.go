package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ikondratev/api-gateway/internal/logger"
)

type Router struct {
	engine *mux.Router
	log    logger.Logger
	events EventCreator
}

func NewRouter(logger logger.Logger, events EventCreator) *Router {
	router := mux.NewRouter()
	router.Use(Logging(logger))

	return &Router{
		engine: router,
		log:    logger,
		events: events,
	}
}

func (r *Router) RegisterRoutes() *mux.Router {
	eventHandler := eventHandler{log: r.log, events: r.events}
	r.engine.HandleFunc("/ping", Pong).Methods(http.MethodGet)

	api := r.engine.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/event", eventHandler.create).Methods(http.MethodPost)

	return r.engine
}
