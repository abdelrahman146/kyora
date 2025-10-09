package db

import (
	"context"
	"log/slog"
	"time"

	"github.com/abdelrahman146/kyora/internal/utils"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type SlogGormLogger struct {
	DefaultLogger *slog.Logger
	LogLevel      gormlogger.LogLevel
}

// NewSlogGormLogger creates a new GORM logger that uses slog.
func NewSlogGormLogger(level gormlogger.LogLevel) *SlogGormLogger {
	return &SlogGormLogger{
		DefaultLogger: slog.Default(),
		LogLevel:      level,
	}
}

func (l *SlogGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *SlogGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		utils.Log.FromContext(ctx).Info(msg, "data", data)
	}
}

func (l *SlogGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		utils.Log.FromContext(ctx).Warn(msg, "data", data)
	}
}

func (l *SlogGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		utils.Log.FromContext(ctx).Error(msg, "data", data)
	}
}

func (l *SlogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// Get the logger from context to include the trace_id
	log := utils.Log.FromContext(ctx)

	logAttrs := []any{
		slog.String("latency", elapsed.String()),
		slog.String("sql", sql),
		slog.Int64("rows_affected", rows),
	}

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && err != gorm.ErrRecordNotFound:
		log.Error("gorm query error", utils.Log.Err(err), slog.Group("query_info", logAttrs...))
	case l.LogLevel >= gormlogger.Warn && err == gorm.ErrRecordNotFound:
		log.Warn("gorm record not found", utils.Log.Err(err), slog.Group("query_info", logAttrs...))
	case l.LogLevel >= gormlogger.Info:
		log.Info("gorm query", slog.Group("query_info", logAttrs...))
	}
}
