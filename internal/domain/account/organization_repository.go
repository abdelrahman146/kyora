package account

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrganizationRepository struct {
	db *db.Postgres
}

func NewOrganizationRepository(db *db.Postgres) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) ScopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *OrganizationRepository) ScopeName(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func (r *OrganizationRepository) ScopeSlug(slug string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("slug = ?", slug)
	}
}

func (r *OrganizationRepository) ScopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *OrganizationRepository) ScopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
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

func (r *OrganizationRepository) ScopeUpdatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
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

func (r *OrganizationRepository) CreateOne(ctx context.Context, org *Organization, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(org).Error
}

func (r *OrganizationRepository) CreateMany(ctx context.Context, orgs []*Organization, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&orgs).Error
}

func (r *OrganizationRepository) UpsertMany(ctx context.Context, orgs []*Organization, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"slug", "name", "updated_at"}),
	}).Create(&orgs).Error
}

func (r *OrganizationRepository) UpdateOne(ctx context.Context, org *Organization, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(org).Error
}

func (r *OrganizationRepository) UpdateMany(ctx context.Context, orgs []*Organization, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&orgs).Error
}

func (r *OrganizationRepository) PatchOne(ctx context.Context, updates *Organization, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Organization{}).Updates(updates).Error
}

func (r *OrganizationRepository) DeleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Organization{}).Error
}

func (r *OrganizationRepository) DeleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Organization{}).Error
}

func (r *OrganizationRepository) FindByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Organization, error) {
	var org Organization
	if err := r.db.Conn(ctx, opts...).First(&org, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepository) FindOne(ctx context.Context, opts ...db.PostgresOptions) (*Organization, error) {
	var org Organization
	if err := r.db.Conn(ctx, opts...).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepository) List(ctx context.Context, opts ...db.PostgresOptions) ([]*Organization, error) {
	var orgs []*Organization
	if err := r.db.Conn(ctx, opts...).Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *OrganizationRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Organization{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
