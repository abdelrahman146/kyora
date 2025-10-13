package supplier

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type supplierRepository struct {
	db *db.Postgres
}

func newSupplierRepository(db *db.Postgres) *supplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *supplierRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *supplierRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *supplierRepository) scopeCountryCode(countryCode string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("country_code = ?", countryCode)
	}
}

func (r *supplierRepository) createOne(ctx context.Context, supplier *Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(supplier).Error
}

func (r *supplierRepository) createMany(ctx context.Context, suppliers []*Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&suppliers).Error
}

func (r *supplierRepository) upsertMany(ctx context.Context, suppliers []*Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "contact", "email", "phone", "website", "updated_at"}),
	}).Create(&suppliers).Error
}

func (r *supplierRepository) updateOne(ctx context.Context, supplier *Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(supplier).Error
}

func (r *supplierRepository) updateMany(ctx context.Context, suppliers []*Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&suppliers).Error
}

func (r *supplierRepository) patchOne(ctx context.Context, updates *Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Supplier{}).Updates(updates).Error
}

func (r *supplierRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Supplier{}).Error
}

func (r *supplierRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Supplier{}).Error
}

func (r *supplierRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Supplier, error) {
	var supplier Supplier
	if err := r.db.Conn(ctx, opts...).First(&supplier, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *supplierRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Supplier, error) {
	var supplier Supplier
	if err := r.db.Conn(ctx, opts...).First(&supplier).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *supplierRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Supplier, error) {
	var suppliers []*Supplier
	if err := r.db.Conn(ctx, opts...).Find(&suppliers).Error; err != nil {
		return nil, err
	}
	return suppliers, nil
}

func (r *supplierRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Supplier{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
