package logger

import (
	"log/slog"
	"os"
)

const (
	DEVELOPMENT = "DEVELOPMENT"
	DEBUG       = "DEBUG"
	INFO        = "INFO"
)

type Interface interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Logger struct {
	logger *slog.Logger
}

func New(env string) *Logger {
	var level slog.Level
	switch env {
	case DEVELOPMENT, DEBUG:
		level = slog.LevelDebug
	default:
		level = slog.LevelInfo
	}

	return &Logger{slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})),
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}
