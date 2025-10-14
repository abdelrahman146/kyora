package owner

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ownerRepository struct {
	db *db.Postgres
}

func newOwnerRepository(db *db.Postgres) *ownerRepository {
	return &ownerRepository{db: db}
}

func (r *ownerRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *ownerRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *ownerRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *ownerRepository) scopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !from.IsZero() && !to.IsZero() {
			return db.Where("created_at BETWEEN ? AND ?", from, to)
		} else if !from.IsZero() {
			return db.Where("created_at >= ?", from)
		} else if !to.IsZero() {
			return db.Where("created_at <= ?", to)
		}
		return db
	}
}

func (r *ownerRepository) createOne(ctx context.Context, owner *Owner, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(owner).Error
}

func (r *ownerRepository) createMany(ctx context.Context, owners []*Owner, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&owners).Error
}

func (r *ownerRepository) upsertMany(ctx context.Context, owners []*Owner, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name", "updated_at"}),
	}).Create(&owners).Error
}

func (r *ownerRepository) updateOne(ctx context.Context, owner *Owner, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(owner).Error
}

func (r *ownerRepository) updateMany(ctx context.Context, owners []*Owner, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&owners).Error
}

func (r *ownerRepository) patchOne(ctx context.Context, updates *Owner, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Owner{}).Updates(updates).Error
}

func (r *ownerRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Owner{}).Error
}

func (r *ownerRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Owner{}).Error
}

func (r *ownerRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Owner, error) {
	var owner Owner
	if err := r.db.Conn(ctx, opts...).First(&owner, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &owner, nil
}

func (r *ownerRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Owner, error) {
	var owner Owner
	if err := r.db.Conn(ctx, opts...).First(&owner).Error; err != nil {
		return nil, err
	}
	return &owner, nil
}

func (r *ownerRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Owner, error) {
	var owners []*Owner
	if err := r.db.Conn(ctx, opts...).Find(&owners).Error; err != nil {
		return nil, err
	}
	return owners, nil
}

func (r *ownerRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Owner{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
