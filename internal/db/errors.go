package db

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/abdelrahman146/kyora/internal/utils"

	"gorm.io/gorm"
)

func HandleDBError(err error) *utils.ProblemDetails {
	if err == nil {
		return nil
	}
	if IsRecordNotFound(err) {
		return utils.Problem.NotFound("The requested resource was not found").WithError(err)
	}
	if IsUniqueViolation(err) {
		return utils.Problem.Conflict("A resource with the same unique attribute already exists").WithError(err)
	}
	return utils.Problem.InternalError().WithError(err)
}

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
