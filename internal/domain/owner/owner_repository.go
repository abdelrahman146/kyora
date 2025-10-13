package owner

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OwnerRepository struct {
	db *db.Postgres
}

func NewOwnerRepository(db *db.Postgres) *OwnerRepository {
	return &OwnerRepository{db: db}
}

func (r *OwnerRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *OwnerRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *OwnerRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *OwnerRepository) createOne(ctx context.Context, owner *Owner, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(owner).Error
}

func (r *OwnerRepository) createMany(ctx context.Context, owners []*Owner, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&owners).Error
}

func (r *OwnerRepository) upsertMany(ctx context.Context, owners []*Owner, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name", "updated_at"}),
	}).Create(&owners).Error
}

func (r *OwnerRepository) updateOne(ctx context.Context, owner *Owner, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(owner).Error
}

func (r *OwnerRepository) updateMany(ctx context.Context, owners []*Owner, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&owners).Error
}

func (r *OwnerRepository) patchOne(ctx context.Context, updates *Owner, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Owner{}).Updates(updates).Error
}

func (r *OwnerRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Owner{}).Error
}

func (r *OwnerRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Owner{}).Error
}

func (r *OwnerRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Owner, error) {
	var owner Owner
	if err := r.db.Conn(ctx, opts...).First(&owner, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &owner, nil
}

func (r *OwnerRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Owner, error) {
	var owner Owner
	if err := r.db.Conn(ctx, opts...).First(&owner).Error; err != nil {
		return nil, err
	}
	return &owner, nil
}

func (r *OwnerRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Owner, error) {
	var owners []*Owner
	if err := r.db.Conn(ctx, opts...).Find(&owners).Error; err != nil {
		return nil, err
	}
	return owners, nil
}

func (r *OwnerRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Owner{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
