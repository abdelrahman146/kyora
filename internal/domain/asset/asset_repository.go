package asset

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/govalues/decimal"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type assetRepository struct {
	db *db.Postgres
}

func newAssetRepository(db *db.Postgres) *assetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *assetRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *assetRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *assetRepository) scopeType(assetType AssetType) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("type = ?", assetType)
	}
}

func (r *assetRepository) scopeTypes(assetTypes []AssetType) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("type IN ?", assetTypes)
	}
}

func (r *assetRepository) scopeFilter(filter *AssetFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		if len(filter.Types) > 0 {
			db = db.Where("type IN ?", filter.Types)
		}
		return db
	}
}

func (r *assetRepository) createOne(ctx context.Context, asset *Asset, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(asset).Error
}

func (r *assetRepository) createMany(ctx context.Context, assets []*Asset, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&assets).Error
}

func (r *assetRepository) upsertMany(ctx context.Context, assets []*Asset, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "type", "amount", "currency", "purchased_at", "note", "updated_at"}),
	}).Create(&assets).Error
}

func (r *assetRepository) updateOne(ctx context.Context, asset *Asset, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(asset).Error
}

func (r *assetRepository) updateMany(ctx context.Context, assets []*Asset, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&assets).Error
}

func (r *assetRepository) patchOne(ctx context.Context, updates *Asset, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Asset{}).Updates(updates).Error
}

func (r *assetRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Asset{}).Error
}

func (r *assetRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Asset{}).Error
}

func (r *assetRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Asset, error) {
	var asset Asset
	if err := r.db.Conn(ctx, opts...).First(&asset, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Asset, error) {
	var asset Asset
	if err := r.db.Conn(ctx, opts...).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Asset, error) {
	var assets []*Asset
	if err := r.db.Conn(ctx, opts...).Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *assetRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Asset{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *assetRepository) sumValue(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&Asset{}).Select("COALESCE(SUM(value), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}
