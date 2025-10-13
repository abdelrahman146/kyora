package inventory

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/govalues/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type variantRepository struct {
	db *db.Postgres
}

func newVariantRepository(db *db.Postgres) *variantRepository {
	return &variantRepository{db: db}
}

func (r *variantRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *variantRepository) scopeSKU(sku string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("sku = ?", sku)
	}
}

func (r *variantRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *variantRepository) scopeSKUs(skus []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("sku IN ?", skus)
	}
}

func (r *variantRepository) scopeProductID(productID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("product_id = ?", productID)
	}
}

func (r *variantRepository) scopeProductIDs(productIDs []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("product_id IN ?", productIDs)
	}
}

func (r *variantRepository) scopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
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

func (r *variantRepository) scopeSearchQuery(query string) func(db *gorm.DB) *gorm.DB {
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

func (r *variantRepository) scopeFilter(filter *VariantFilter) func(db *gorm.DB) *gorm.DB {
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
			db = r.scopeCreatedAt(filter.From, filter.To)(db)
		}
		if len(filter.ProductTags) > 0 {
			db = db.Joins("JOIN products ON products.id = variants.product_id").Where("products.tags && ?", filter.ProductTags)
		}
		if filter.SearchQuery != "" {
			db = r.scopeSearchQuery(filter.SearchQuery)(db)
		}
		return db
	}
}

func (r *variantRepository) createOne(ctx context.Context, variant *Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(variant).Error
}

func (r *variantRepository) createMany(ctx context.Context, variants []*Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&variants).Error
}

// CreateManyStrict inserts all variants and returns an error on any conflict/violation.
// Useful when caller wants to handle unique conflicts explicitly.
func (r *variantRepository) createManyStrict(ctx context.Context, variants []*Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(&variants).Error
}

func (r *variantRepository) upsertMany(ctx context.Context, variants []*Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "sku", "cost_price", "sale_price", "currency", "stock_quantity", "stock_alert", "updated_at"}),
	}).Create(&variants).Error
}

func (r *variantRepository) updateOne(ctx context.Context, variant *Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(variant).Error
}

func (r *variantRepository) updateMany(ctx context.Context, variants []*Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&variants).Error
}

func (r *variantRepository) patchOne(ctx context.Context, updates *Variant, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Variant{}).Updates(updates).Error
}

func (r *variantRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Variant{}).Error
}

func (r *variantRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Variant{}).Error
}

func (r *variantRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Variant, error) {
	var variant Variant
	if err := r.db.Conn(ctx, opts...).First(&variant, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *variantRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Variant, error) {
	var variant Variant
	if err := r.db.Conn(ctx, opts...).First(&variant).Error; err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *variantRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Variant, error) {
	var variants []*Variant
	if err := r.db.Conn(ctx, opts...).Find(&variants).Error; err != nil {
		return nil, err
	}
	return variants, nil
}

func (r *variantRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Variant{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *variantRepository) sumCostPrice(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&Variant{}).Select("COALESCE(SUM(cost_price), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}

func (r *variantRepository) sumSalePrice(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&Variant{}).Select("COALESCE(SUM(sale_price), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}
