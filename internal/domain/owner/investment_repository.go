package owner

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/govalues/decimal"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InvestmentRepository struct {
	db *db.Postgres
}

func NewInvestmentRepository(db *db.Postgres) *InvestmentRepository {
	return &InvestmentRepository{db: db}
}

func (r *InvestmentRepository) ScopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *InvestmentRepository) ScopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *InvestmentRepository) ScopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *InvestmentRepository) ScopeOwnerID(ownerID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("owner_id = ?", ownerID)
	}
}

func (r *InvestmentRepository) CreateOne(ctx context.Context, investment *Investment, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(investment).Error
}

func (r *InvestmentRepository) CreateMany(ctx context.Context, investments []*Investment, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&investments).Error
}

func (r *InvestmentRepository) UpsertMany(ctx context.Context, investments []*Investment, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "amount", "currency", "note", "updated_at"}),
	}).Create(&investments).Error
}

func (r *InvestmentRepository) UpdateOne(ctx context.Context, investment *Investment, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(investment).Error
}

func (r *InvestmentRepository) UpdateMany(ctx context.Context, investments []*Investment, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&investments).Error
}

func (r *InvestmentRepository) PatchOne(ctx context.Context, updates *Investment, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Investment{}).Updates(updates).Error
}

func (r *InvestmentRepository) DeleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Investment{}).Error
}

func (r *InvestmentRepository) DeleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Investment{}).Error
}

func (r *InvestmentRepository) FindByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Investment, error) {
	var investment Investment
	if err := r.db.Conn(ctx, opts...).First(&investment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &investment, nil
}

func (r *InvestmentRepository) FindOne(ctx context.Context, opts ...db.PostgresOptions) (*Investment, error) {
	var investment Investment
	if err := r.db.Conn(ctx, opts...).First(&investment).Error; err != nil {
		return nil, err
	}
	return &investment, nil
}

func (r *InvestmentRepository) List(ctx context.Context, opts ...db.PostgresOptions) ([]*Investment, error) {
	var investments []*Investment
	if err := r.db.Conn(ctx, opts...).Find(&investments).Error; err != nil {
		return nil, err
	}
	return investments, nil
}

func (r *InvestmentRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Investment{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *InvestmentRepository) SumAmount(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&Investment{}).Select("COALESCE(SUM(amount), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}
