package utils

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-errors/errors" // For stack traces
)

type loggerHelper struct{}

type loggerKeyType struct{}

var loggerKey = loggerKeyType{}

func (loggerHelper) Setup(logLevel string) {
	var level slog.Level

	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create a new JSON handler.
	// The ReplaceAttr function is key for adding stack traces.
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	// Create a new logger and set it as the default.
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func (loggerHelper) Err(err error) slog.Attr {
	// Use go-errors to wrap the error and capture the stack trace.
	return slog.Any("error", errors.Wrap(err, 1))
}

func (loggerHelper) FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

func (loggerHelper) Default() *slog.Logger {
	return slog.Default()
}

func (loggerHelper) WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

var Log = loggerHelper{}
