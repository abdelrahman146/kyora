package db

import (
	"database/sql"

	"github.com/abdelrahman146/kyora/internal/utils"

	"gorm.io/gorm"
)

func HandleDBError(err error) *utils.ProblemDetails {
	if err == nil {
		return nil
	}
	if err == gorm.ErrRecordNotFound || err == sql.ErrNoRows {
		return utils.Problem.NotFound("The requested resource was not found").WithError(err)
	}
	return utils.Problem.InternalError().WithError(err)
}
