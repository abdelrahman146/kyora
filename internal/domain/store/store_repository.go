package store

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type storeRepository struct {
	db *db.Postgres
}

func newStoreRepository(db *db.Postgres) *storeRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *storeRepository) scopeName(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func (r *storeRepository) scopeSlug(slug string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("slug = ?", slug)
	}
}

func (r *storeRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *storeRepository) scopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
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

func (r *storeRepository) scopeOrganizationID(organizationID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("organization_id = ?", organizationID)
	}
}

func (r *storeRepository) scopeCode(code string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("code = ?", code)
	}
}

func (r *storeRepository) createOne(ctx context.Context, store *Store, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(store).Error
}

func (r *storeRepository) createMany(ctx context.Context, stores []*Store, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&stores).Error
}

func (r *storeRepository) upsertMany(ctx context.Context, stores []*Store, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "slug", "updated_at"}),
	}).Create(&stores).Error
}

func (r *storeRepository) updateOne(ctx context.Context, store *Store, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(store).Error
}

func (r *storeRepository) updateMany(ctx context.Context, stores []*Store, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&stores).Error
}

func (r *storeRepository) patchOne(ctx context.Context, updates *Store, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Store{}).Updates(updates).Error
}

func (r *storeRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Store{}).Error
}

func (r *storeRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Store{}).Error
}

func (r *storeRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Store, error) {
	var store Store
	if err := r.db.Conn(ctx, opts...).First(&store, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Store, error) {
	var store Store
	if err := r.db.Conn(ctx, opts...).First(&store).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Store, error) {
	var stores []*Store
	if err := r.db.Conn(ctx, opts...).Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *storeRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Store{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
