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
	auth   Auther
}

func NewRouter(logger logger.Logger, events EventCreator, auth Auther) *Router {
	router := mux.NewRouter()

	return &Router{
		engine: router,
		log:    logger,
		events: events,
		auth:   auth,
	}
}

func (r *Router) RegisterRoutes() *mux.Router {
	eventHandler := eventHandler{log: r.log, events: r.events}
	authHandler := authHandler{log: r.log, auth: r.auth}
	r.engine.HandleFunc("/ping", Pong).Methods(http.MethodGet)

	api := r.engine.PathPrefix("/api").Subrouter()
	api.Use(Trace())
	api.Use(Logging(r.log))

	v1 := api.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/event", eventHandler.create).Methods(http.MethodPost)

	v1_auth := v1.PathPrefix("/auth").Subrouter()
	v1_auth.HandleFunc("/register", authHandler.register).Methods(http.MethodPost)
	v1_auth.HandleFunc("/login", authHandler.login).Methods(http.MethodPost)
	v1_auth.HandleFunc("/refresh", authHandler.refresh).Methods(http.MethodPost)
	v1_auth.HandleFunc("/logout", authHandler.logout).Methods(http.MethodPost)

	return r.engine
}
