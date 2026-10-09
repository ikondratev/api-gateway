package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ikondratev/api-gateway/internal/eventservice"
	"github.com/ikondratev/api-gateway/internal/logger"
	"github.com/ikondratev/api-gateway/internal/server"
	"github.com/ikondratev/api-gateway/internal/settings"
)

type Server interface {
	Start(chan error)
	Stop(context.Context) error
}

type App struct {
	settings *settings.Settings
	server   Server
	logger   logger.Logger
	events   *eventservice.Client
}

func New(env string) (*App, error) {
	settings, err := settings.New(env)
	logger := logger.New(env)
	if err != nil {
		logger.Error("settings error:", "load", err)
		return nil, fmt.Errorf("settings error: %w", err)
	}

	eventts, err := eventservice.New(settings)
	if err != nil {
		logger.Error("event service client:", "error", err)
		return nil, fmt.Errorf("event service client: %w", err)
	}

	server := server.NewHTTPServer(settings, logger, eventts)

	return &App{
		server:   server,
		settings: settings,
		logger:   logger,
		events:   eventts,
	}, nil
}

func (a *App) Run() error {
	chErr := make(chan error, 1)
	go a.server.Start(chErr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	if err := a.awaitInterraption(chErr, quit); err != nil {
		if closeEventErr := a.events.Close(); closeEventErr != nil {
			a.logger.Error("event service close", "shutdown", closeEventErr)
		}
		return err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(a.settings.Server.WaitingShutdown)*time.Second)
	defer cancel()

	stopServerErr := a.server.Stop(ctx)
	stopEventsErr := a.events.Close()

	if stopServerErr != nil {
		a.logger.Error("server stopped with error", "shutdown", stopServerErr)
		return fmt.Errorf("server stopped with error: %w", stopServerErr)
	}
	if stopEventsErr != nil {
		a.logger.Error("event service close", "shutdown", stopEventsErr)
		return fmt.Errorf("event service close: %w", stopEventsErr)
	}

	a.logger.Info("Server stoppend gracefullty")
	return nil
}

func (a *App) awaitInterraption(chErr <-chan error, quit <-chan os.Signal) error {
	select {
	case err := <-chErr:
		if !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("internal error:", "shutdown", err)
			return fmt.Errorf("internal error: %w", err)
		}
		return nil
	case sig := <-quit:
		a.logger.Info("Server handled signal", "shutdown", sig)
		return nil
	}
}
