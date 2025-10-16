package expense

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type recurringExpenseRepository struct {
	db *db.Postgres
}

func newRecurringExpenseRepository(db *db.Postgres) *recurringExpenseRepository {
	return &recurringExpenseRepository{db: db}
}

func (r *recurringExpenseRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *recurringExpenseRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *recurringExpenseRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *recurringExpenseRepository) scopeCategory(category ExpenseCategory) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category = ?", category)
	}
}

func (r *recurringExpenseRepository) scopeCategories(categories []ExpenseCategory) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category IN ?", categories)
	}
}

func (r *recurringExpenseRepository) scopeFrequency(frequency RecurringExpenseFrequency) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("frequency = ?", frequency)
	}
}

func (r *recurringExpenseRepository) scopeFrequencies(frequencies []RecurringExpenseFrequency) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("frequency IN ?", frequencies)
	}
}

func (r *recurringExpenseRepository) scopeStatus(status RecurringExpenseStatus) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db.Where("status = ?", status) }
}

func (r *recurringExpenseRepository) scopeStartDateLte(t time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db.Where("recurring_start_date <= ?", t) }
}

func (r *recurringExpenseRepository) scopeEndDateGteOrNull(t time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("recurring_end_date IS NULL OR recurring_end_date >= ?", t)
	}
}

func (r *recurringExpenseRepository) createOne(ctx context.Context, expense *RecurringExpense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(expense).Error
}

func (r *recurringExpenseRepository) createMany(ctx context.Context, expenses []*RecurringExpense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&expenses).Error
}

func (r *recurringExpenseRepository) upsertMany(ctx context.Context, expenses []*RecurringExpense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"frequency", "recurring_start_date", "recurring_end_date", "amount", "category", "type", "note", "updated_at"}),
	}).Create(&expenses).Error
}

func (r *recurringExpenseRepository) updateOne(ctx context.Context, expense *RecurringExpense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(expense).Error
}

func (r *recurringExpenseRepository) updateMany(ctx context.Context, expenses []*RecurringExpense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&expenses).Error
}

func (r *recurringExpenseRepository) patchOne(ctx context.Context, updates *RecurringExpense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&RecurringExpense{}).Updates(updates).Error
}

func (r *recurringExpenseRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&RecurringExpense{}).Error
}

func (r *recurringExpenseRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&RecurringExpense{}).Error
}

func (r *recurringExpenseRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*RecurringExpense, error) {
	var expense RecurringExpense
	if err := r.db.Conn(ctx, opts...).First(&expense, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *recurringExpenseRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*RecurringExpense, error) {
	var expense RecurringExpense
	if err := r.db.Conn(ctx, opts...).First(&expense).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *recurringExpenseRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*RecurringExpense, error) {
	var expenses []*RecurringExpense
	if err := r.db.Conn(ctx, opts...).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *recurringExpenseRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&RecurringExpense{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *recurringExpenseRepository) sumAmount(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&RecurringExpense{}).Select("COALESCE(SUM(amount), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}
