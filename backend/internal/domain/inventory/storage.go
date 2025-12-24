package inventory

import (
	"fmt"
	"log/slog"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Storage struct {
	cache      *cache.Cache
	products   *database.Repository[Product]
	variants   *database.Repository[Variant]
	categories *database.Repository[Category]
}

func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
	st := &Storage{
		cache:      cache,
		products:   database.NewRepository[Product](db),
		variants:   database.NewRepository[Variant](db),
		categories: database.NewRepository[Category](db),
	}
	ensurePgTrgmAndIndexes(db)
	return st
}

func ensurePgTrgmAndIndexes(db *database.Database) {
	conn := db.GetDB()
	if err := conn.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error; err != nil {
		slog.Error("failed to ensure pg_trgm extension", "error", err)
		return
	}

	queries := []string{
		`CREATE INDEX IF NOT EXISTS "product_trgm_idx" ON "products" USING gin ("name" gin_trgm_ops)`,
		`CREATE INDEX IF NOT EXISTS "variant_trgm_idx" ON "variants" USING gin ("name" gin_trgm_ops)`,
		`CREATE INDEX IF NOT EXISTS "category_trgm_idx" ON "categories" USING gin ("name" gin_trgm_ops)`,
	}
	for _, q := range queries {
		if err := conn.Exec(q).Error; err != nil {
			slog.Error("failed to ensure trigram index", "error", err, "query", q)
		}
	}
}

func (s *Storage) ScopeSearchTermByName(term string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if term == "" {
			return db
		}
		return db.
			Where("name % ?", term).
			Order(clause.Expr{
				SQL:  "similarity(name, ?) DESC",
				Vars: []any{term},
			})
	}
}

func (s *Storage) ScopeLowStockVariants() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s <= %s", VariantSchema.StockQuantity.Column(), VariantSchema.StockQuantityAlert.Column()))
	}
}
