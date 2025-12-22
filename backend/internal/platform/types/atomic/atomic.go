package atomic

import (
	"context"
	"database/sql"
)

type AtomicProcessOption func(*AtomicProcessOptions)

type AtomicProcessor interface {
	Exec(ctx context.Context, cb func(ctx context.Context) error, opts ...AtomicProcessOption) error
}

type IsolationLevel int

const (
	LevelDefault IsolationLevel = iota
	LevelReadUncommitted
	LevelReadCommitted
	LevelWriteCommitted
	LevelSnapshot
	LevelRepeatableRead
	LevelSerializable
	LevelLinearizable
)

func (l IsolationLevel) ToSQLIsolationLevel() sql.IsolationLevel {
	switch l {
	case LevelReadUncommitted:
		return sql.LevelReadUncommitted
	case LevelReadCommitted:
		return sql.LevelReadCommitted
	case LevelWriteCommitted:
		return sql.LevelWriteCommitted
	case LevelSnapshot:
		return sql.LevelSnapshot
	case LevelRepeatableRead:
		return sql.LevelRepeatableRead
	case LevelSerializable:
		return sql.LevelSerializable
	case LevelLinearizable:
		return sql.LevelLinearizable
	default:
		return sql.LevelDefault
	}
}

func (l IsolationLevel) String() string {
	switch l {
	case LevelReadUncommitted:
		return "READ UNCOMMITTED"
	case LevelReadCommitted:
		return "READ COMMITTED"
	case LevelWriteCommitted:
		return "WRITE COMMITTED"
	case LevelSnapshot:
		return "SNAPSHOT"
	case LevelRepeatableRead:
		return "REPEATABLE READ"
	case LevelSerializable:
		return "SERIALIZABLE"
	case LevelLinearizable:
		return "LINEARIZABLE"
	default:
		return "DEFAULT"
	}
}

type AtomicProcessOptions struct {
	Isolation IsolationLevel
	Retries   int
	ReadOnly  bool
}

func WithIsolationLevel(level IsolationLevel) AtomicProcessOption {
	return func(opts *AtomicProcessOptions) {
		opts.Isolation = level
	}
}

func WithReadOnly(readOnly bool) AtomicProcessOption {
	return func(opts *AtomicProcessOptions) {
		opts.ReadOnly = readOnly
	}
}

func WithRetries(retries int) AtomicProcessOption {
	return func(opts *AtomicProcessOptions) {
		opts.Retries = max(retries, 1)
	}
}
