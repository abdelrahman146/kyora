package testutils

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/platform/database"
)

// TruncateTables truncates the specified database tables for test isolation
func TruncateTables(db *database.Database, tables ...string) error {
	conn := db.GetDB()
	for _, table := range tables {
		if err := conn.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}
	return nil
}

// ClearTable is an alias for TruncateTables for backward compatibility
func ClearTable(db *database.Database, tables ...string) error {
	return TruncateTables(db, tables...)
}
