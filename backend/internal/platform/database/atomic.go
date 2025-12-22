package database

import (
	"context"
	"database/sql"

	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"gorm.io/gorm"
)

type AtomicProcess struct {
	tx *gorm.DB
}

func NewAtomicProcess(db *Database) *AtomicProcess {
	return &AtomicProcess{db.GetDB()}
}

func (u *AtomicProcess) Exec(ctx context.Context, cb func(ctx context.Context) error, opts ...atomic.AtomicProcessOption) error {
	options := &atomic.AtomicProcessOptions{
		Isolation: atomic.LevelDefault,
		Retries:   3,
		ReadOnly:  false,
	}
	for _, opt := range opts {
		opt(options)
	}
	for i := 0; i < options.Retries; i++ {
		err := u.tx.Transaction(func(tx *gorm.DB) error {
			if err := u.setupTransaction(tx, options); err != nil {
				return err
			}
			return cb(context.WithValue(ctx, TxKey, tx))
		})
		if err == nil {
			return nil
		}
		if err != sql.ErrTxDone {
			return err
		}
	}
	return sql.ErrTxDone
}

func (u *AtomicProcess) setupTransaction(tx *gorm.DB, options *atomic.AtomicProcessOptions) error {
	if options.Isolation != atomic.LevelDefault {
		if err := tx.Exec("SET TRANSACTION ISOLATION LEVEL " + options.Isolation.String()).Error; err != nil {
			return err
		}
	}
	if options.ReadOnly {
		if err := tx.Exec("SET TRANSACTION READ ONLY").Error; err != nil {
			return err
		}
	}
	return nil
}
