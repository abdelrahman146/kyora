package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/types/keyvalue"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/types/timeseries"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DatabaseOption func(db *gorm.DB) *gorm.DB

type Repository[T any] struct {
	db *Database
}

func NewRepository[T any](db *Database) *Repository[T] {
	db.AutoMigrate(new(T))
	return &Repository[T]{db: db}
}

type LockingStrength string

const (
	LockingStrengthUpdate    LockingStrength = "UPDATE"
	LockingStrengthShare     LockingStrength = "SHARE"
	LockingOptionsSkipLocked LockingStrength = "SKIP LOCKED"
	LockingOptionsNoWait     LockingStrength = "NOWAIT"
)

func (r *Repository[T]) WithLockingStrength(strength LockingStrength) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(clause.Locking{Strength: string(strength)})
	}
}

func (r *Repository[T]) WithPreload(associations ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, association := range associations {
			db = db.Preload(association)
		}
		return db
	}
}

func (r *Repository[T]) WithJoins(joins ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, join := range joins {
			db = db.Joins(join)
		}
		return db
	}
}

func (r *Repository[T]) WithPagination(offset, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
}

// WithReturning applies a RETURNING clause to the query and scans the result into the provided value.
func (r *Repository[T]) WithReturning(value any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(clause.Returning{}).Scan(value)
	}
}

func (r *Repository[T]) WithLimit(limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

// WithOrderBy applies ordering based on the provided orderBy slice.
// Each element in orderBy should a column name and direction e.g., "name ASC" or "createdAt DESC".
func (r *Repository[T]) WithOrderBy(orderBy []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, ob := range orderBy {
			if ob != "" && len(ob) > 1 {
				db = db.Order(ob)
			}
		}
		return db
	}
}

func (r *Repository[T]) ScopeID(id any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *Repository[T]) ScopeIDs(ids []any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *Repository[T]) ScopeIn(field schema.Field, values []any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field.Column()+" IN ?", values)
	}
}

func (r *Repository[T]) ScopeNotIn(field schema.Field, values []any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field.Column()+" NOT IN ?", values)
	}
}

func (r *Repository[T]) ScopeEquals(field schema.Field, value any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field.Column()+" = ?", value)
	}
}

func (r *Repository[T]) ScopeIsNull(field schema.Field) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field.Column() + " IS NULL")
	}
}

func (r *Repository[T]) ScopeNotEquals(field schema.Field, value any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field.Column()+" <> ?", value)
	}
}

func (r *Repository[T]) ScopeGreaterThan(field schema.Field, threshold any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field.Column()+" > ?", threshold)
	}
}

func (r *Repository[T]) ScopeHavingGreaterThan(field schema.Field, threshold any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Having(field.Column()+" > ?", threshold)
	}
}

func (r *Repository[T]) ScopeLessThan(field schema.Field, threshold any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field.Column()+" < ?", threshold)
	}
}

func (r *Repository[T]) ScopeHavingLessThan(field schema.Field, threshold any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Having(field.Column()+" < ?", threshold)
	}
}

func (r *Repository[T]) ScopeGreaterThanOrEqual(field schema.Field, threshold any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field.Column()+" >= ?", threshold)
	}
}

func (r *Repository[T]) ScopeHavingGreaterThanOrEqual(field schema.Field, threshold any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Having(field.Column()+" >= ?", threshold)
	}
}

func (r *Repository[T]) ScopeLessThanOrEqual(field schema.Field, threshold any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field.Column()+" <= ?", threshold)
	}
}

func (r *Repository[T]) ScopeHavingLessThanOrEqual(field schema.Field, threshold any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Having(field.Column()+" <= ?", threshold)
	}
}

func (r *Repository[T]) ScopeBetween(field schema.Field, from, to any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field.Column()+" BETWEEN ? AND ?", from, to)
	}
}

func (r *Repository[T]) ScopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !from.IsZero() && !to.IsZero() {
			return db.Where("created_at BETWEEN ? AND ?", from, to)
		} else if !from.IsZero() {
			return db.Where("created_at >= ?", from)
		} else if !to.IsZero() {
			return db.Where("created_at <= ?", to)
		}
		return db
	}
}

func (r *Repository[T]) ScopeTime(field schema.Field, from, to time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !from.IsZero() && !to.IsZero() {
			return db.Where(field.Column()+" BETWEEN ? AND ?", from, to)
		} else if !from.IsZero() {
			return db.Where(field.Column()+" >= ?", from)
		} else if !to.IsZero() {
			return db.Where(field.Column()+" <= ?", to)
		}
		return db
	}
}

func (r *Repository[T]) ScopeBusinessID(businessID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("business_id = ?", businessID)
	}
}

func (r *Repository[T]) ScopeWorkspaceID(workspaceID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("workspace_id = ?", workspaceID)
	}
}

func (r *Repository[T]) ScopeSearchTerm(searchTerm string, fields ...schema.Field) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if searchTerm == "" || len(fields) == 0 {
			return db
		}
		likePattern := "%" + searchTerm + "%"
		conditions := make([]string, len(fields))
		values := make([]any, len(fields))
		for i, field := range fields {
			conditions[i] = field.Column() + " ILIKE ?"
			values[i] = likePattern
		}
		return db.Where("("+strings.Join(conditions, " OR ")+")", values...)
	}
}

func (r *Repository[T]) CreateOne(ctx context.Context, entity *T, opts ...func(db *gorm.DB) *gorm.DB) error {
	return r.db.Conn(ctx).Scopes(opts...).Create(entity).Error
}

func (r *Repository[T]) CreateMany(ctx context.Context, entities []*T, opts ...func(db *gorm.DB) *gorm.DB) error {
	return r.db.Conn(ctx).Scopes(opts...).Create(&entities).Error
}

func (r *Repository[T]) UpdateOne(ctx context.Context, entity *T, opts ...func(db *gorm.DB) *gorm.DB) error {
	return r.db.Conn(ctx).Scopes(opts...).Save(entity).Error
}

func (r *Repository[T]) UpdateMany(ctx context.Context, entities []*T, opts ...func(db *gorm.DB) *gorm.DB) error {
	return r.db.Conn(ctx).Scopes(opts...).Save(&entities).Error
}

func (r *Repository[T]) DeleteOne(ctx context.Context, entity *T, opts ...func(db *gorm.DB) *gorm.DB) error {
	return r.db.Conn(ctx).Scopes(opts...).Delete(entity).Error
}

func (r *Repository[T]) DeleteMany(ctx context.Context, opts ...func(db *gorm.DB) *gorm.DB) error {
	return r.db.Conn(ctx).Scopes(opts...).Delete(new(T)).Error
}

func (r *Repository[T]) FindByID(ctx context.Context, id any, opts ...func(db *gorm.DB) *gorm.DB) (*T, error) {
	var entity T
	err := r.db.Conn(ctx).Scopes(append(opts, r.ScopeID(id))...).First(&entity).Error
	return &entity, err
}

func (r *Repository[T]) FindOne(ctx context.Context, opts ...func(db *gorm.DB) *gorm.DB) (*T, error) {
	var entity T
	err := r.db.Conn(ctx).Scopes(opts...).First(&entity).Error
	return &entity, err
}

func (r *Repository[T]) FindMany(ctx context.Context, opts ...func(db *gorm.DB) *gorm.DB) ([]*T, error) {
	var entities []*T
	err := r.db.Conn(ctx).Scopes(opts...).Find(&entities).Error
	return entities, err
}

func (r *Repository[T]) Count(ctx context.Context, opts ...func(db *gorm.DB) *gorm.DB) (int64, error) {
	var count int64
	err := r.db.Conn(ctx).Scopes(opts...).Model(new(T)).Count(&count).Error
	return count, err
}

func (r *Repository[T]) Sum(ctx context.Context, column schema.Field, opts ...func(db *gorm.DB) *gorm.DB) (decimal.Decimal, error) {
	var sum decimal.Decimal
	err := r.db.Conn(ctx).Scopes(opts...).Model(new(T)).Select("COALESCE(SUM(" + column.Column() + "),0)::decimal").Scan(&sum).Error
	if err != nil {
		return decimal.Zero, err
	}
	return sum, nil
}

func (r *Repository[T]) Avg(ctx context.Context, column schema.Field, opts ...func(db *gorm.DB) *gorm.DB) (decimal.Decimal, error) {
	var avg decimal.Decimal
	err := r.db.Conn(ctx).Scopes(opts...).Model(new(T)).Select("COALESCE(AVG(" + column.Column() + "),0)::decimal").Scan(&avg).Error
	if err != nil {
		return decimal.Zero, err
	}
	return avg, nil
}

func (r *Repository[T]) TimeSeriesSum(ctx context.Context, valueColumn schema.Field, timeColumn schema.Field, granularity timeseries.Granularity, opts ...func(db *gorm.DB) *gorm.DB) (*timeseries.TimeSeries, error) {
	var rows []timeseries.TimeSeriesRow
	sel := fmt.Sprintf("date_trunc('%s', %s) AS timestamp, COALESCE(SUM(%s),0)::decimal AS value", granularity.Bucket(), timeColumn.Column(), valueColumn.Column())
	q := r.db.Conn(ctx).Scopes(opts...).Model(new(T))
	if err := q.Select(sel).Group("timestamp").Order("timestamp ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return timeseries.New(rows, granularity), nil
}

func (r *Repository[T]) TimeSeriesCount(ctx context.Context, timeColumn schema.Field, granularity timeseries.Granularity, opts ...func(db *gorm.DB) *gorm.DB) (*timeseries.TimeSeries, error) {
	var rows []timeseries.TimeSeriesRow
	sel := fmt.Sprintf("date_trunc('%s', %s) AS timestamp, COUNT(*) AS value", granularity.Bucket(), timeColumn.Column())
	q := r.db.Conn(ctx).Scopes(opts...).Model(new(T))
	if err := q.Select(sel).Group("timestamp").Order("timestamp ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return timeseries.New(rows, granularity), nil
}

func (r *Repository[T]) CountBy(ctx context.Context, column schema.Field, opts ...func(db *gorm.DB) *gorm.DB) ([]keyvalue.KeyValue, error) {
	var rows []map[string]any
	err := r.db.Conn(ctx).Scopes(opts...).Model(new(T)).Select(column.Column() + " AS key, COUNT(*) AS value").Group(column.Column()).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	results := make([]keyvalue.KeyValue, 0, len(rows))
	for _, row := range rows {
		results = append(results, keyvalue.New(row["key"], row["value"]))
	}
	return results, nil
}

func (r *Repository[T]) SumBy(ctx context.Context, groupBy schema.Field, sumColumn schema.Field, opts ...func(db *gorm.DB) *gorm.DB) ([]keyvalue.KeyValue, error) {
	var rows []map[string]any
	err := r.db.Conn(ctx).Scopes(opts...).Model(new(T)).Select(groupBy.Column() + " AS key, COALESCE(SUM(" + sumColumn.Column() + "),0)::decimal AS value").Group(groupBy.Column()).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	results := make([]keyvalue.KeyValue, 0, len(rows))
	for _, row := range rows {
		results = append(results, keyvalue.New(row["key"], row["value"]))
	}
	return results, nil
}

func (r *Repository[T]) AvgBy(ctx context.Context, groupBy schema.Field, avgColumn schema.Field, opts ...func(db *gorm.DB) *gorm.DB) ([]keyvalue.KeyValue, error) {
	var rows []map[string]any
	err := r.db.Conn(ctx).Scopes(opts...).Model(new(T)).Select(groupBy.Column() + " AS key, COALESCE(AVG(" + avgColumn.Column() + "),0)::decimal AS value").Group(groupBy.Column()).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	results := make([]keyvalue.KeyValue, 0, len(rows))
	for _, row := range rows {
		results = append(results, keyvalue.New(row["key"], row["value"]))
	}
	return results, nil
}
