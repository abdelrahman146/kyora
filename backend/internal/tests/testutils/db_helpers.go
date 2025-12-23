package testutils

import (
	"fmt"
	"strings"

	"github.com/abdelrahman146/kyora/internal/platform/database"
)

// TruncateTables truncates the specified database tables for test isolation
func TruncateTables(db *database.Database, tables ...string) error {
	if len(tables) == 0 {
		return nil
	}

	conn := db.GetDB()
	stmt := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", strings.Join(tables, ", "))
	if err := conn.Exec(stmt).Error; err != nil {
		return fmt.Errorf("failed to truncate tables %v: %w", tables, err)
	}

	return nil
}

// ClearTable is an alias for TruncateTables for backward compatibility
func ClearTable(db *database.Database, tables ...string) error {
	return TruncateTables(db, tables...)
}
