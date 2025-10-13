package inventory

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	db *db.Postgres
}

func NewProductRepository(db *db.Postgres) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}
func (r *ProductRepository) ScopeSearchQuery(query string) func(db *gorm.DB) *gorm.DB {
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

func (r *ProductRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *ProductRepository) ScopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
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

func (r *ProductRepository) ScopeTags(tags []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(tags) > 0 {
			return db.Where("tags && ?", tags)
		}
		return db
	}
}

func (r *ProductRepository) ScopeFilter(filter *ProductFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		if len(filter.IDs) > 0 {
			db = db.Scopes(r.scopeIDs(filter.IDs))
		}
		if len(filter.Tags) > 0 {
			db = db.Scopes(r.ScopeTags(filter.Tags))
		}
		if filter.From.IsZero() && filter.To.IsZero() {
			db = db.Scopes(r.ScopeCreatedAt(filter.From, filter.To))
		}
		if filter.SearchQuery != "" {
			db = db.Scopes(r.ScopeSearchQuery(filter.SearchQuery))
		}
		return db
	}
}

func (r *ProductRepository) ScopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *ProductRepository) CreateOne(ctx context.Context, product *Product, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(product).Error
}

func (r *ProductRepository) CreateMany(ctx context.Context, products []*Product, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&products).Error
}

func (r *ProductRepository) UpsertMany(ctx context.Context, products []*Product, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "description", "tags", "updated_at"}),
	}).Create(&products).Error
}

func (r *ProductRepository) UpdateOne(ctx context.Context, product *Product, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(product).Error
}

func (r *ProductRepository) UpdateMany(ctx context.Context, products []*Product, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&products).Error
}

func (r *ProductRepository) PatchOne(ctx context.Context, updates *Product, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Product{}).Updates(updates).Error
}

func (r *ProductRepository) DeleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Product{}).Error
}

func (r *ProductRepository) DeleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Product{}).Error
}

func (r *ProductRepository) FindByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Product, error) {
	var product Product
	if err := r.db.Conn(ctx, opts...).First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindOne(ctx context.Context, opts ...db.PostgresOptions) (*Product, error) {
	var product Product
	if err := r.db.Conn(ctx, opts...).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) List(ctx context.Context, opts ...db.PostgresOptions) ([]*Product, error) {
	var products []*Product
	if err := r.db.Conn(ctx, opts...).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Product{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
