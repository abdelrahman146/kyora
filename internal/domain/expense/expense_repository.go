package expense

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/govalues/decimal"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ExpenseRepository struct {
	db *db.Postgres
}

func NewExpenseRepository(db *db.Postgres) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *ExpenseRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *ExpenseRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *ExpenseRepository) scopeCategory(category ExpenseCategory) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category = ?", category)
	}
}

func (r *ExpenseRepository) scopeCategories(categories []ExpenseCategory) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category IN ?", categories)
	}
}

func (r *ExpenseRepository) scopeType(expenseType ExpenseType) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("type = ?", expenseType)
	}
}

func (r *ExpenseRepository) scopeRecurringExpenseID(recurringExpenseID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("recurring_expense_id = ?", recurringExpenseID)
	}
}

func (r *ExpenseRepository) scopeCreatedAtBetween(start, end time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at >= ? AND created_at < ?", start, end)
	}
}

func (r *ExpenseRepository) scopeOccurredOn(day time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Ensure we only match by date portion in case the DB stores timezone
		return db.Where("occurred_on = ?", day)
	}
}

func (r *ExpenseRepository) scopeFilter(filter *ExpenseFilter) func(db *gorm.DB) *gorm.DB {
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

func (r *ExpenseRepository) createOne(ctx context.Context, expense *Expense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(expense).Error
}

func (r *ExpenseRepository) createMany(ctx context.Context, expenses []*Expense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&expenses).Error
}

func (r *ExpenseRepository) upsertMany(ctx context.Context, expenses []*Expense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "amount", "category", "note", "updated_at"}),
	}).Create(&expenses).Error
}

func (r *ExpenseRepository) updateOne(ctx context.Context, expense *Expense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(expense).Error
}

func (r *ExpenseRepository) updateMany(ctx context.Context, expenses []*Expense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&expenses).Error
}

func (r *ExpenseRepository) patchOne(ctx context.Context, updates *Expense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Expense{}).Updates(updates).Error
}

func (r *ExpenseRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Expense{}).Error
}

func (r *ExpenseRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Expense{}).Error
}

func (r *ExpenseRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Expense, error) {
	var expense Expense
	if err := r.db.Conn(ctx, opts...).First(&expense, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *ExpenseRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Expense, error) {
	var expense Expense
	if err := r.db.Conn(ctx, opts...).First(&expense).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *ExpenseRepository) List(ctx context.Context, opts ...db.PostgresOptions) ([]*Expense, error) {
	var expenses []*Expense
	if err := r.db.Conn(ctx, opts...).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Expense{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ExpenseRepository) SumAmount(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&Expense{}).Select("COALESCE(SUM(amount), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}
