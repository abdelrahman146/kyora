package inventory

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/types"
	"github.com/shopspring/decimal"
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

func (r *variantRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
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

// ---- Analytics helpers ----

// sumStockQuantity returns the total units in stock across variants.
func (r *variantRepository) sumStockQuantity(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	type row struct{ Qty int64 }
	var res row
	if err := r.db.Conn(ctx, opts...).Model(&Variant{}).Select("COALESCE(SUM(stock_quantity),0) AS qty").Scan(&res).Error; err != nil {
		return 0, err
	}
	return res.Qty, nil
}

// sumInventoryValue returns the total inventory valuation using cost_price * stock_quantity.
func (r *variantRepository) sumInventoryValue(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&Variant{}).Select("COALESCE(SUM(cost_price * stock_quantity),0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}

// countLowStock returns number of variants where stock_quantity <= stock_alert and stock_alert > 0.
func (r *variantRepository) countLowStock(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Variant{}).Where("stock_alert > 0 AND stock_quantity > 0 AND stock_quantity <= stock_alert").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// countOutOfStock returns number of variants with zero quantity.
func (r *variantRepository) countOutOfStock(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Variant{}).Where("stock_quantity = 0").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// topProductsByInventoryValue aggregates inventory value per product (sum of variant cost*qty) and returns top-N rows.
// Returns key=product name, value=float (inventory value)
func (r *variantRepository) topProductsByInventoryValue(ctx context.Context, limit int, opts ...db.PostgresOptions) ([]types.KeyValue, error) {
	if limit <= 0 {
		limit = 10
	}
	rows := []types.KeyValue{}
	q := r.db.Conn(ctx, opts...).Table(VariantTable + " v").Joins("JOIN products p ON p.id = v.product_id")
	if err := q.Select("p.name AS key, COALESCE(SUM(v.cost_price * v.stock_quantity),0)::float AS value").Group("p.id, p.name").Order("value DESC").Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
