package store

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StoreRepository struct {
	db *db.Postgres
}

func NewStoreRepository(db *db.Postgres) *StoreRepository {
	return &StoreRepository{db: db}
}

func (r *StoreRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *StoreRepository) ScopeName(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func (r *StoreRepository) ScopeSlug(slug string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("slug = ?", slug)
	}
}

func (r *StoreRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *StoreRepository) ScopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
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

func (r *StoreRepository) ScopeOrganizationID(organizationID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("organization_id = ?", organizationID)
	}
}

func (r *StoreRepository) ScopeCode(code string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("code = ?", code)
	}
}

func (r *StoreRepository) CreateOne(ctx context.Context, store *Store, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(store).Error
}

func (r *StoreRepository) CreateMany(ctx context.Context, stores []*Store, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&stores).Error
}

func (r *StoreRepository) UpsertMany(ctx context.Context, stores []*Store, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "slug", "updated_at"}),
	}).Create(&stores).Error
}

func (r *StoreRepository) UpdateOne(ctx context.Context, store *Store, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(store).Error
}

func (r *StoreRepository) UpdateMany(ctx context.Context, stores []*Store, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&stores).Error
}

func (r *StoreRepository) PatchOne(ctx context.Context, updates *Store, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Store{}).Updates(updates).Error
}

func (r *StoreRepository) DeleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Store{}).Error
}

func (r *StoreRepository) DeleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Store{}).Error
}

func (r *StoreRepository) FindByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Store, error) {
	var store Store
	if err := r.db.Conn(ctx, opts...).First(&store, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *StoreRepository) FindOne(ctx context.Context, opts ...db.PostgresOptions) (*Store, error) {
	var store Store
	if err := r.db.Conn(ctx, opts...).First(&store).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *StoreRepository) List(ctx context.Context, opts ...db.PostgresOptions) ([]*Store, error) {
	var stores []*Store
	if err := r.db.Conn(ctx, opts...).Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *StoreRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Store{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
