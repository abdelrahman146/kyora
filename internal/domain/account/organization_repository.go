package account

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type organizationRepository struct {
	db *db.Postgres
}

func newOrganizationRepository(db *db.Postgres) *organizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *organizationRepository) scopeName(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func (r *organizationRepository) scopeSlug(slug string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("slug = ?", slug)
	}
}

func (r *organizationRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *organizationRepository) scopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
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

func (r *organizationRepository) scopeUpdatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !from.IsZero() && !to.IsZero() {
			return db.Where("updated_at BETWEEN ? AND ?", from, to)
		} else if !from.IsZero() {
			return db.Where("updated_at >= ?", from)
		} else if !to.IsZero() {
			return db.Where("updated_at <= ?", to)
		}
		return db
	}
}

func (r *organizationRepository) createOne(ctx context.Context, org *Organization, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(org).Error
}

func (r *organizationRepository) createMany(ctx context.Context, orgs []*Organization, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&orgs).Error
}

func (r *organizationRepository) upsertMany(ctx context.Context, orgs []*Organization, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"slug", "name", "updated_at"}),
	}).Create(&orgs).Error
}

func (r *organizationRepository) updateOne(ctx context.Context, org *Organization, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(org).Error
}

func (r *organizationRepository) updateMany(ctx context.Context, orgs []*Organization, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&orgs).Error
}

func (r *organizationRepository) patchOne(ctx context.Context, updates *Organization, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Organization{}).Updates(updates).Error
}

func (r *organizationRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Organization{}).Error
}

func (r *organizationRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Organization{}).Error
}

func (r *organizationRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Organization, error) {
	var org Organization
	if err := r.db.Conn(ctx, opts...).First(&org, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Organization, error) {
	var org Organization
	if err := r.db.Conn(ctx, opts...).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Organization, error) {
	var orgs []*Organization
	if err := r.db.Conn(ctx, opts...).Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *organizationRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Organization{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
