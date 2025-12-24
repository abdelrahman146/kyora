package database

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"

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

// IsRetryableTxError returns true for Postgres transaction errors that are safe to retry.
// This is primarily used for SERIALIZABLE/REPEATABLE READ transactions where
// serialization failures and deadlocks can occur.
func IsRetryableTxError(err error) bool {
	if err == nil {
		return false
	}

	// Postgres SQLSTATEs
	// 40001 = serialization_failure
	// 40P01 = deadlock_detected
	// 55P03 = lock_not_available
	// 57014 = query_canceled (often due to statement_timeout) -> not safe to blindly retry
	const (
		serializationFailure = "40001"
		deadlockDetected     = "40P01"
		lockNotAvailable     = "55P03"
	)

	// Unwrap chain
	for e := err; e != nil; e = errors.Unwrap(e) {
		var pgErr *pgconn.PgError
		if errors.As(e, &pgErr) {
			switch pgErr.Code {
			case serializationFailure, deadlockDetected, lockNotAvailable:
				return true
			}
		}
		var pqErr *pq.Error
		if errors.As(e, &pqErr) {
			switch string(pqErr.Code) {
			case serializationFailure, deadlockDetected, lockNotAvailable:
				return true
			}
		}
	}

	// Fallback heuristic (covers some wrapped gorm/driver errors)
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "sqlstate 40001") || strings.Contains(msg, "sqlstate 40p01") || strings.Contains(msg, "sqlstate 55p03") {
		return true
	}
	return false
}
