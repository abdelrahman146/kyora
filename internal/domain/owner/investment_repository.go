package owner

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/govalues/decimal"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type investmentRepository struct {
	db *db.Postgres
}

func newInvestmentRepository(db *db.Postgres) *investmentRepository {
	return &investmentRepository{db: db}
}

func (r *investmentRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *investmentRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *investmentRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *investmentRepository) scopeOwnerID(ownerID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("owner_id = ?", ownerID)
	}
}

func (r *investmentRepository) createOne(ctx context.Context, investment *Investment, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(investment).Error
}

func (r *investmentRepository) createMany(ctx context.Context, investments []*Investment, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&investments).Error
}

func (r *investmentRepository) upsertMany(ctx context.Context, investments []*Investment, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "amount", "currency", "note", "updated_at"}),
	}).Create(&investments).Error
}

func (r *investmentRepository) updateOne(ctx context.Context, investment *Investment, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(investment).Error
}

func (r *investmentRepository) updateMany(ctx context.Context, investments []*Investment, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&investments).Error
}

func (r *investmentRepository) patchOne(ctx context.Context, updates *Investment, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Investment{}).Updates(updates).Error
}

func (r *investmentRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Investment{}).Error
}

func (r *investmentRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Investment{}).Error
}

func (r *investmentRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Investment, error) {
	var investment Investment
	if err := r.db.Conn(ctx, opts...).First(&investment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &investment, nil
}

func (r *investmentRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Investment, error) {
	var investment Investment
	if err := r.db.Conn(ctx, opts...).First(&investment).Error; err != nil {
		return nil, err
	}
	return &investment, nil
}

func (r *investmentRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Investment, error) {
	var investments []*Investment
	if err := r.db.Conn(ctx, opts...).Find(&investments).Error; err != nil {
		return nil, err
	}
	return investments, nil
}

func (r *investmentRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Investment{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *investmentRepository) sumAmount(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&Investment{}).Select("COALESCE(SUM(amount), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}
