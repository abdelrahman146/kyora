package expense

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/govalues/decimal"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RecurringExpenseRepository struct {
	db *db.Postgres
}

func NewRecurringExpenseRepository(db *db.Postgres) *RecurringExpenseRepository {
	return &RecurringExpenseRepository{db: db}
}

func (r *RecurringExpenseRepository) ScopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *RecurringExpenseRepository) ScopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *RecurringExpenseRepository) ScopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *RecurringExpenseRepository) ScopeCategory(category ExpenseCategory) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category = ?", category)
	}
}

func (r *RecurringExpenseRepository) ScopeCategories(categories []ExpenseCategory) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category IN ?", categories)
	}
}

func (r *RecurringExpenseRepository) ScopeFrequency(frequency RecurringExpenseFrequency) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("frequency = ?", frequency)
	}
}

func (r *RecurringExpenseRepository) ScopeFrequencies(frequencies []RecurringExpenseFrequency) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("frequency IN ?", frequencies)
	}
}

func (r *RecurringExpenseRepository) ScopeFilter(filter *RecurringExpenseFilter) func(db *gorm.DB) *gorm.DB {
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
		if len(filter.Frequencies) > 0 {
			db = db.Where("frequency IN ?", filter.Frequencies)
		}
		return db
	}
}

func (r *RecurringExpenseRepository) ScopeStatus(status RecurringExpenseStatus) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db.Where("status = ?", status) }
}

func (r *RecurringExpenseRepository) ScopeStartDateLte(t time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db.Where("recurring_start_date <= ?", t) }
}

func (r *RecurringExpenseRepository) ScopeEndDateGteOrNull(t time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("recurring_end_date IS NULL OR recurring_end_date >= ?", t)
	}
}

func (r *RecurringExpenseRepository) CreateOne(ctx context.Context, expense *RecurringExpense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(expense).Error
}

func (r *RecurringExpenseRepository) CreateMany(ctx context.Context, expenses []*RecurringExpense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&expenses).Error
}

func (r *RecurringExpenseRepository) UpsertMany(ctx context.Context, expenses []*RecurringExpense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"frequency", "recurring_start_date", "recurring_end_date", "amount", "category", "type", "note", "updated_at"}),
	}).Create(&expenses).Error
}

func (r *RecurringExpenseRepository) UpdateOne(ctx context.Context, expense *RecurringExpense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(expense).Error
}

func (r *RecurringExpenseRepository) UpdateMany(ctx context.Context, expenses []*RecurringExpense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&expenses).Error
}

func (r *RecurringExpenseRepository) PatchOne(ctx context.Context, updates *RecurringExpense, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&RecurringExpense{}).Updates(updates).Error
}

func (r *RecurringExpenseRepository) DeleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&RecurringExpense{}).Error
}

func (r *RecurringExpenseRepository) DeleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&RecurringExpense{}).Error
}

func (r *RecurringExpenseRepository) FindByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*RecurringExpense, error) {
	var expense RecurringExpense
	if err := r.db.Conn(ctx, opts...).First(&expense, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *RecurringExpenseRepository) FindOne(ctx context.Context, opts ...db.PostgresOptions) (*RecurringExpense, error) {
	var expense RecurringExpense
	if err := r.db.Conn(ctx, opts...).First(&expense).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *RecurringExpenseRepository) List(ctx context.Context, opts ...db.PostgresOptions) ([]*RecurringExpense, error) {
	var expenses []*RecurringExpense
	if err := r.db.Conn(ctx, opts...).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *RecurringExpenseRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&RecurringExpense{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RecurringExpenseRepository) SumAmount(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&RecurringExpense{}).Select("COALESCE(SUM(amount), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}
