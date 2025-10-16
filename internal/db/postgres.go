package db

import (
	"context"
	"database/sql"

	"github.com/samber/lo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TxKey struct{}

type Postgres struct {
	db *gorm.DB
}

func NewPostgres(dsn string, cfg *gorm.Config) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, err
	}
	return &Postgres{db: db}, nil
}

func (p *Postgres) Conn(ctx context.Context, opts ...PostgresOptions) *gorm.DB {
	if tx, ok := ctx.Value(TxKey{}).(*gorm.DB); ok {
		return p.applyOptions(tx, opts...)
	}
	return p.applyOptions(p.db.WithContext(ctx), opts...)
}

func (p *Postgres) DB() *gorm.DB {
	return p.db
}

func (p *Postgres) applyOptions(db *gorm.DB, opts ...PostgresOptions) *gorm.DB {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (p *Postgres) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *Postgres) AutoMigrate(models ...interface{}) error {
	return p.db.AutoMigrate(models...)
}

type PostgresOptions func(db *gorm.DB) *gorm.DB

type LockingStrength string

const (
	LockingStrengthUpdate    LockingStrength = "UPDATE"
	LockingStrengthShare     LockingStrength = "SHARE"
	LockingOptionsSkipLocked LockingStrength = "SKIP LOCKED"
	LockingOptionsNoWait     LockingStrength = "NOWAIT"
)

func WithLock(strength LockingStrength) PostgresOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(clause.Locking{Strength: string(strength)})
	}
}

func WithReturning(value any) PostgresOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(value).Clauses(clause.Returning{})
	}
}

func WithPreload(associations ...string) PostgresOptions {
	return func(db *gorm.DB) *gorm.DB {
		for _, association := range associations {
			db = db.Preload(association)
		}
		return db
	}
}

func WithJoins(joins ...string) PostgresOptions {
	return func(db *gorm.DB) *gorm.DB {
		for _, join := range joins {
			db = db.Joins(join)
		}
		return db
	}
}

func WithScopes(scopes ...func(*gorm.DB) *gorm.DB) PostgresOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(scopes...)
	}
}

func WithPagination(page, pageSize int) PostgresOptions {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 30
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func WithLimit(limit int) PostgresOptions {
	return func(db *gorm.DB) *gorm.DB {
		if limit > 0 {
			return db.Limit(limit)
		}
		return db
	}
}

func WithOrderBy(orderBy []string) PostgresOptions {
	return func(db *gorm.DB) *gorm.DB {
		for _, ob := range orderBy {
			if ob != "" && len(ob) > 1 {
				if ob[0] == '-' {
					db = db.Order(lo.SnakeCase(ob[1:]) + " DESC")
				} else {
					db = db.Order(lo.SnakeCase(ob) + " ASC")
				}
			}
		}
		return db
	}
}

type AtomicProcess struct {
	tx *gorm.DB
}

func NewAtomicProcess(db *gorm.DB) *AtomicProcess {
	return &AtomicProcess{db}
}

func (u *AtomicProcess) Exec(ctx context.Context, cb func(ctx context.Context) error, opts ...AtomicProccessOption) error {
	options := &AtomicProcessOptions{
		Isolation: sql.LevelDefault,
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
			return cb(context.WithValue(ctx, TxKey{}, tx))
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

func (u *AtomicProcess) setupTransaction(tx *gorm.DB, options *AtomicProcessOptions) error {
	if options.Isolation != sql.LevelDefault {
		if err := tx.Exec("SET TRANSACTION ISOLATION LEVEL " + u.isolationLevelToString(options.Isolation)).Error; err != nil {
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

func (u *AtomicProcess) isolationLevelToString(level sql.IsolationLevel) string {
	switch level {
	case sql.LevelReadUncommitted:
		return "READ UNCOMMITTED"
	case sql.LevelReadCommitted:
		return "READ COMMITTED"
	case sql.LevelWriteCommitted:
		return "WRITE COMMITTED"
	case sql.LevelRepeatableRead:
		return "REPEATABLE READ"
	case sql.LevelSnapshot:
		return "SNAPSHOT"
	case sql.LevelSerializable:
		return "SERIALIZABLE"
	case sql.LevelLinearizable:
		return "LINEARIZABLE"
	default:
		return "DEFAULT"
	}
}

type AtomicProcessOptions struct {
	Isolation sql.IsolationLevel
	Retries   int
	ReadOnly  bool
}

type AtomicProccessOption func(*AtomicProcessOptions)

func WithIsolationLevel(level sql.IsolationLevel) AtomicProccessOption {
	return func(opts *AtomicProcessOptions) {
		opts.Isolation = level
	}
}

func WithReadOnly(readOnly bool) AtomicProccessOption {
	return func(opts *AtomicProcessOptions) {
		opts.ReadOnly = readOnly
	}
}

func WithRetries(retries int) AtomicProccessOption {
	return func(opts *AtomicProcessOptions) {
		if retries < 1 {
			opts.Retries = 1
		} else {
			opts.Retries = retries
		}
	}
}
