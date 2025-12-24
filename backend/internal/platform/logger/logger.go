// Package logger provides logging utilities
// implements slog
package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/types/ctxkey"
	"github.com/spf13/viper"
)

var LoggerCtxKey = ctxkey.New("logger")

// Init configures the global slog logger from viper config values.
// It should be called early in application startup (e.g., Cobra PersistentPreRun).
func Init() {
	// Determine log level
	var level slog.Level
	logLevel := viper.GetString(config.LogLevel)
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO", "":
		level = slog.LevelInfo
	case "WARN", "WARNING":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Determine format: json (default) or text
	format := strings.ToLower(strings.TrimSpace(viper.GetString(config.LogFormat)))
	var handler slog.Handler
	if format == "text" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}
	slog.SetDefault(slog.New(handler))
}

func FromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(LoggerCtxKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, LoggerCtxKey, logger)
}
