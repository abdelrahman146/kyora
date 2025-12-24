// Package database provides database utilities
// implements gorm
package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/types/ctxkey"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

var TxKey = ctxkey.New("transaction")

func NewConnection(dsn string, logLevel string) (*Database, error) {
	maxAttempts := 5
	var db *gorm.DB
	var err error
	logger := NewSlogGormLogger(logLevel)
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger,
		})
		if err == nil {
			break
		}
		slog.Warn("Failed to connect to database, retrying...", "attempt", attempts, "error", err)
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	if err != nil {
		slog.Error("Could not connect to the database", "error", err)
		return nil, fmt.Errorf("connect database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Could not get database instance", "error", err)
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	// Ensure required Postgres extensions exist (used by search scopes and indexes).
	// This is safe to run multiple times and avoids test/dev failures on fresh databases.
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error; err != nil {
		slog.Warn("Failed to ensure pg_trgm extension", "error", err)
	}
	maxOpenConns := viper.GetInt(config.DatabaseMaxOpenConns)
	maxIdleConns := viper.GetInt(config.DatabaseMaxIdleConns)
	maxIdleTime := viper.GetDuration(config.DatabaseMaxIdleTime)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)
	return &Database{db: db}, nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}

func (d *Database) CloseConnection() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) Conn(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(TxKey).(*gorm.DB); ok {
		return tx
	}
	return d.db.WithContext(ctx)
}

func (d *Database) ApplyOptions(db *gorm.DB, opts ...DatabaseOption) *gorm.DB {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (d *Database) AutoMigrate(models ...interface{}) error {
	return d.db.AutoMigrate(models...)
}
