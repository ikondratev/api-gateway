package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Logger
}

type slogLogger struct {
	engine *slog.Logger
}

func New(env string) Logger {
	level := slog.LevelInfo
	if env == "dev" {
		level = slog.LevelDebug
	}

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:level,
	})

	return &slogLogger{engine: slog.New(h)}
}

func (l *slogLogger) Debug(msg string, args ...any) { l.engine.Debug(msg, args...) }

func (l *slogLogger) Info(msg string, args ...any)  { l.engine.Info(msg, args...) }

func (l *slogLogger) Warn(msg string, args ...any)  { l.engine.Warn(msg, args...) }

func (l *slogLogger) Error(msg string, args ...any) { l.engine.Error(msg, args...) }

func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{engine: l.engine.With(args...)}
}