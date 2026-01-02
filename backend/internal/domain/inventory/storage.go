package inventory

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"gorm.io/gorm"
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
	ensureInventorySearchIndexes(db)
	return st
}

func ensureInventorySearchIndexes(db *database.Database) {
	conn := db.GetDB()

	// Product search vector: name (weight A) + description (weight B)
	productExpr := "" +
		"setweight(to_tsvector('simple', coalesce(name,'')), 'A') || " +
		"setweight(to_tsvector('simple', coalesce(description,'')), 'B')"
	database.EnsureGeneratedTSVectorColumn(conn, ProductTable, "search_vector", productExpr)
	database.EnsureGinIndex(conn, "products_search_vector_gin_idx", ProductTable, "search_vector")

	// Variant search vector: name (weight A) + sku (weight A) + code (weight B)
	variantExpr := "" +
		"setweight(to_tsvector('simple', coalesce(name,'')), 'A') || " +
		"setweight(to_tsvector('simple', coalesce(sku,'')), 'A') || " +
		"setweight(to_tsvector('simple', coalesce(code,'')), 'B')"
	database.EnsureGeneratedTSVectorColumn(conn, VariantTable, "search_vector", variantExpr)
	database.EnsureGinIndex(conn, "variants_search_vector_gin_idx", VariantTable, "search_vector")

	// Trigram index for fast substring lookup on variant SKU
	database.EnsureTrigramGinIndex(conn, "variants_sku_trgm_idx", VariantTable, "sku")

	// Category search vector: name (weight A) + descriptor (weight B)
	categoryExpr := "" +
		"setweight(to_tsvector('simple', coalesce(name,'')), 'A') || " +
		"setweight(to_tsvector('simple', coalesce(descriptor,'')), 'B')"
	database.EnsureGeneratedTSVectorColumn(conn, CategoryTable, "search_vector", categoryExpr)
	database.EnsureGinIndex(conn, "categories_search_vector_gin_idx", CategoryTable, "search_vector")
}

func (s *Storage) ScopeLowStockVariants() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s <= %s", VariantSchema.StockQuantity.Column(), VariantSchema.StockQuantityAlert.Column()))
	}
}

func (s *Storage) ScopeOutOfStockVariants() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = 0", VariantSchema.StockQuantity.Column()))
	}
}

func (s *Storage) ScopeInStockVariants() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s > %s", VariantSchema.StockQuantity.Column(), VariantSchema.StockQuantityAlert.Column()))
	}
}
