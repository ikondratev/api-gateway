package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ikondratev/api-gateway/internal/logger"
	"github.com/ikondratev/api-gateway/internal/settings"
)

type HttpServer struct {
	engine   *http.Server
	settings *settings.Settings
	logger   logger.Logger
}

func NewHTTPServer(s *settings.Settings, logger logger.Logger, events EventCreator) *HttpServer {
	router := NewRouter(logger, events)
	handler := router.RegisterRoutes()

	server := &http.Server{
		Addr:              s.Server.Port,
		Handler:           handler,
		ReadHeaderTimeout: time.Duration(s.Server.HeaderTimeout) * time.Second,
		ReadTimeout:       time.Duration(s.Server.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(s.Server.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(s.Server.IdleTimeout) * time.Second,
	}

	return &HttpServer{
		engine:   server,
		settings: s,
		logger:   logger,
	}
}

func (s *HttpServer) Start(errCh chan error) {
	s.logger.Info("Http server started:", s.engine.Addr, s.settings.Environment)
	errCh <- s.engine.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.engine.Shutdown(ctx)
}
