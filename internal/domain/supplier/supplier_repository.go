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

func (r *SupplierRepository) ScopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *SupplierRepository) ScopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *SupplierRepository) ScopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *SupplierRepository) CreateOne(ctx context.Context, supplier *Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(supplier).Error
}

func (r *SupplierRepository) CreateMany(ctx context.Context, suppliers []*Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&suppliers).Error
}

func (r *SupplierRepository) UpsertMany(ctx context.Context, suppliers []*Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "contact", "email", "phone", "website", "updated_at"}),
	}).Create(&suppliers).Error
}

func (r *SupplierRepository) UpdateOne(ctx context.Context, supplier *Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(supplier).Error
}

func (r *SupplierRepository) UpdateMany(ctx context.Context, suppliers []*Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&suppliers).Error
}

func (r *SupplierRepository) PatchOne(ctx context.Context, updates *Supplier, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Supplier{}).Updates(updates).Error
}

func (r *SupplierRepository) DeleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Supplier{}).Error
}

func (r *SupplierRepository) DeleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Supplier{}).Error
}

func (r *SupplierRepository) FindByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Supplier, error) {
	var supplier Supplier
	if err := r.db.Conn(ctx, opts...).First(&supplier, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *SupplierRepository) FindOne(ctx context.Context, opts ...db.PostgresOptions) (*Supplier, error) {
	var supplier Supplier
	if err := r.db.Conn(ctx, opts...).First(&supplier).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *SupplierRepository) List(ctx context.Context, opts ...db.PostgresOptions) ([]*Supplier, error) {
	var suppliers []*Supplier
	if err := r.db.Conn(ctx, opts...).Find(&suppliers).Error; err != nil {
		return nil, err
	}
	return suppliers, nil
}

func (r *SupplierRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Supplier{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
