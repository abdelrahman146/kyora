package inventory

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type productRepository struct {
	db *db.Postgres
}

func newProductRepository(db *db.Postgres) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}
func (r *productRepository) scopeSearchQuery(query string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if query == "" {
			return db
		}
		return db.
			Where("name % ?", query).
			Order(clause.Expr{
				SQL:  "similarity(name, ?) DESC",
				Vars: []any{query},
			})
	}
}

func (r *productRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *productRepository) scopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
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

func (r *productRepository) scopeTags(tags []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(tags) > 0 {
			return db.Where("tags && ?", tags)
		}
		return db
	}
}

func (r *productRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *productRepository) createOne(ctx context.Context, product *Product, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(product).Error
}

func (r *productRepository) createMany(ctx context.Context, products []*Product, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&products).Error
}

func (r *productRepository) upsertMany(ctx context.Context, products []*Product, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "description", "tags", "updated_at"}),
	}).Create(&products).Error
}

func (r *productRepository) updateOne(ctx context.Context, product *Product, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(product).Error
}

func (r *productRepository) updateMany(ctx context.Context, products []*Product, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&products).Error
}

func (r *productRepository) patchOne(ctx context.Context, updates *Product, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Product{}).Updates(updates).Error
}

func (r *productRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Product{}).Error
}

func (r *productRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Product{}).Error
}

func (r *productRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Product, error) {
	var product Product
	if err := r.db.Conn(ctx, opts...).First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Product, error) {
	var product Product
	if err := r.db.Conn(ctx, opts...).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Product, error) {
	var products []*Product
	if err := r.db.Conn(ctx, opts...).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Product{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
