package database

import (
	"database/sql"
	"errors"
	"strings"

	"gorm.io/gorm"
)

func IsRecordNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound || err == sql.ErrNoRows
}

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	if strings.Contains(strings.ToLower(msg), "unique constraint") ||
		strings.Contains(strings.ToLower(msg), "duplicate key") ||
		strings.Contains(strings.ToLower(msg), "unique violation") {
		return true
	}
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil || unwrapped == err {
			break
		}
		err = unwrapped
		m := strings.ToLower(err.Error())
		if strings.Contains(m, "duplicate key") || strings.Contains(m, "unique") {
			return true
		}
	}
	return false
}
