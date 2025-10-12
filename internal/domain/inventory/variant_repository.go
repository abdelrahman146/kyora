package inventory

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/govalues/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VariantRepository struct {
	db *db.Postgres
}

func NewVariantRepository(db *db.Postgres) *VariantRepository {
	return &VariantRepository{db: db}
}

func (r *VariantRepository) ScopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *VariantRepository) ScopeSKU(sku string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("sku = ?", sku)
	}
}

func (r *VariantRepository) ScopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *VariantRepository) ScopeSKUs(skus []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("sku IN ?", skus)
	}
}

func (r *VariantRepository) ScopeProductID(productID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("product_id = ?", productID)
	}
}

func (r *VariantRepository) ScopeProductIDs(productIDs []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("product_id IN ?", productIDs)
	}
}

func (r *VariantRepository) ScopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
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

func (r *VariantRepository) ScopeSearchQuery(query string) func(db *gorm.DB) *gorm.DB {
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

func (r *VariantRepository) ScopeFilter(filter *VariantFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		if len(filter.IDs) > 0 {
			db = db.Where("id IN ?", filter.IDs)
		}
		if len(filter.SKUs) > 0 {
			db = db.Where("sku IN ?", filter.SKUs)
		}
		if len(filter.ProductIDs) > 0 {
			db = db.Where("product_id IN ?", filter.ProductIDs)
		}
		if !filter.From.IsZero() || !filter.To.IsZero() {
			db = r.ScopeCreatedAt(filter.From, filter.To)(db)
		}
		if len(filter.ProductTags) > 0 {
			db = db.Joins("JOIN products ON products.id = variants.product_id").Where("products.tags && ?", filter.ProductTags)
		}
		if filter.SearchQuery != "" {
			db = r.ScopeSearchQuery(filter.SearchQuery)(db)
		}
		return db
	}
}

func (r *VariantRepository) CreateOne(ctx context.Context, variant *Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(variant).Error
}

func (r *VariantRepository) CreateMany(ctx context.Context, variants []*Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&variants).Error
}

// CreateManyStrict inserts all variants and returns an error on any conflict/violation.
// Useful when caller wants to handle unique conflicts explicitly.
func (r *VariantRepository) CreateManyStrict(ctx context.Context, variants []*Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(&variants).Error
}

func (r *VariantRepository) UpsertMany(ctx context.Context, variants []*Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "sku", "cost_price", "sale_price", "currency", "stock_quantity", "stock_alert", "updated_at"}),
	}).Create(&variants).Error
}

func (r *VariantRepository) UpdateOne(ctx context.Context, variant *Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(variant).Error
}

func (r *VariantRepository) UpdateMany(ctx context.Context, variants []*Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&variants).Error
}

func (r *VariantRepository) PatchOne(ctx context.Context, updates *Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Variant{}).Updates(updates).Error
}

func (r *VariantRepository) DeleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Variant{}).Error
}

func (r *VariantRepository) DeleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Variant{}).Error
}

func (r *VariantRepository) FindByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Variant, error) {
	var variant Variant
	if err := r.db.Conn(ctx, opts...).First(&variant, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *VariantRepository) FindOne(ctx context.Context, opts ...db.PostgresOptions) (*Variant, error) {
	var variant Variant
	if err := r.db.Conn(ctx, opts...).First(&variant).Error; err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *VariantRepository) List(ctx context.Context, opts ...db.PostgresOptions) ([]*Variant, error) {
	var variants []*Variant
	if err := r.db.Conn(ctx, opts...).Find(&variants).Error; err != nil {
		return nil, err
	}
	return variants, nil
}

func (r *VariantRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Variant{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *VariantRepository) SumCostPrice(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&Variant{}).Select("COALESCE(SUM(cost_price), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}

func (r *VariantRepository) SumSalePrice(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&Variant{}).Select("COALESCE(SUM(sale_price), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}
