package supplier

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SupplierRepository struct {
	db *db.Postgres
}

func NewSupplierRepository(db *db.Postgres) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *SupplierRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *SupplierRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *SupplierRepository) scopeCountryCode(countryCode string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("country_code = ?", countryCode)
	}
}

func (r *SupplierRepository) createOne(ctx context.Context, supplier *Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(supplier).Error
}

func (r *SupplierRepository) createMany(ctx context.Context, suppliers []*Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&suppliers).Error
}

func (r *SupplierRepository) upsertMany(ctx context.Context, suppliers []*Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "contact", "email", "phone", "website", "updated_at"}),
	}).Create(&suppliers).Error
}

func (r *SupplierRepository) updateOne(ctx context.Context, supplier *Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(supplier).Error
}

func (r *SupplierRepository) updateMany(ctx context.Context, suppliers []*Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&suppliers).Error
}

func (r *SupplierRepository) patchOne(ctx context.Context, updates *Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Supplier{}).Updates(updates).Error
}

func (r *SupplierRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Supplier{}).Error
}

func (r *SupplierRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Supplier{}).Error
}

func (r *SupplierRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Supplier, error) {
	var supplier Supplier
	if err := r.db.Conn(ctx, opts...).First(&supplier, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *SupplierRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Supplier, error) {
	var supplier Supplier
	if err := r.db.Conn(ctx, opts...).First(&supplier).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *SupplierRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Supplier, error) {
	var suppliers []*Supplier
	if err := r.db.Conn(ctx, opts...).Find(&suppliers).Error; err != nil {
		return nil, err
	}
	return suppliers, nil
}

func (r *SupplierRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Supplier{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
