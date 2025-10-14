package expense

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/types"
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type expenseRepository struct {
	db *db.Postgres
}

func NewExpenseRepository(db *db.Postgres) *expenseRepository {
	return &expenseRepository{db: db}
}

func (r *expenseRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *expenseRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *expenseRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *expenseRepository) scopeCategory(category ExpenseCategory) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category = ?", category)
	}
}

func (r *expenseRepository) scopeCategories(categories []ExpenseCategory) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category IN ?", categories)
	}
}

func (r *expenseRepository) scopeType(expenseType ExpenseType) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("type = ?", expenseType)
	}
}

func (r *expenseRepository) scopeRecurringExpenseID(recurringExpenseID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("recurring_expense_id = ?", recurringExpenseID)
	}
}

func (r *expenseRepository) scopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
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

func (r *expenseRepository) scopeOccurredOn(day time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Ensure we only match by date portion in case the DB stores timezone
		return db.Where("occurred_on = ?", day)
	}
}

func (r *expenseRepository) scopeFilter(filter *ExpenseFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		if len(filter.IDs) > 0 {
			db = db.Where("id IN ?", filter.IDs)
		}
		if len(filter.Categories) > 0 {
			db = db.Where("category IN ?", filter.Categories)
		}
		if len(filter.Types) > 0 {
			db = db.Where("type IN ?", filter.Types)
		}
		return db
	}
}

func (r *expenseRepository) createOne(ctx context.Context, expense *Expense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(expense).Error
}

func (r *expenseRepository) createMany(ctx context.Context, expenses []*Expense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&expenses).Error
}

func (r *expenseRepository) upsertMany(ctx context.Context, expenses []*Expense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "amount", "category", "note", "updated_at"}),
	}).Create(&expenses).Error
}

func (r *expenseRepository) updateOne(ctx context.Context, expense *Expense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(expense).Error
}

func (r *expenseRepository) updateMany(ctx context.Context, expenses []*Expense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&expenses).Error
}

func (r *expenseRepository) patchOne(ctx context.Context, updates *Expense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Expense{}).Updates(updates).Error
}

func (r *expenseRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Expense{}).Error
}

func (r *expenseRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Expense{}).Error
}

func (r *expenseRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Expense, error) {
	var expense Expense
	if err := r.db.Conn(ctx, opts...).First(&expense, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *expenseRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Expense, error) {
	var expense Expense
	if err := r.db.Conn(ctx, opts...).First(&expense).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *expenseRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Expense, error) {
	var expenses []*Expense
	if err := r.db.Conn(ctx, opts...).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *expenseRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Expense{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *expenseRepository) sumAmount(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&Expense{}).Select("COALESCE(SUM(amount), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}

// ---- Analytics helpers ----

// breakdownByCategory returns total expense amount per category.
func (r *expenseRepository) breakdownByCategory(ctx context.Context, opts ...db.PostgresOptions) ([]types.KeyValue, error) {
	rows := []types.KeyValue{}
	q := r.db.Conn(ctx, opts...).Model(&Expense{})
	if err := q.Select("category AS key, COALESCE(SUM(amount),0)::float AS value").Group("category").Order("value DESC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// amountTimeSeries returns sum(amount) per date_trunc bucket for occurred_on (preferred) falling back to created_at.
func (r *expenseRepository) amountTimeSeries(ctx context.Context, bucket string, from, to time.Time, opts ...db.PostgresOptions) ([]types.TimeSeriesRow, error) {
	switch bucket {
	case "hour", "day", "week", "month", "quarter", "year":
	default:
		bucket = "day"
	}
	rows := []types.TimeSeriesRow{}
	sel := "date_trunc('" + bucket + "', occurred_on) AS timestamp, COALESCE(SUM(amount),0)::float AS value"
	q := r.db.Conn(ctx, opts...).Model(&Expense{}).Scopes(r.scopeCreatedAt(from, to))
	if err := q.Select(sel).Group("timestamp").Order("timestamp ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
