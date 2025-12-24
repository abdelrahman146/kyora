package database

import (
	"context"
	"math/rand/v2"
	"time"

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
	if options.Retries < 1 {
		options.Retries = 1
	}

	var lastErr error
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
		lastErr = err
		if !IsRetryableTxError(err) {
			return err
		}

		// Exponential backoff with jitter to avoid thundering herd.
		// Base: 50ms, 100ms, 200ms, ... capped at 2s.
		backoff := 50 * time.Millisecond * time.Duration(1<<i)
		if backoff > 2*time.Second {
			backoff = 2 * time.Second
		}
		jitter := time.Duration(rand.Int64N(int64(backoff / 2)))
		delay := backoff + jitter

		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
			// retry
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return ctx.Err()
		}
	}
	return lastErr
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
