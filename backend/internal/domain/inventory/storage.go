package inventory

import (
	"fmt"
	"strings"

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

// ScopeProductSearch applies search filter across products, variants, categories, and SKU
func (s *Storage) ScopeProductSearch(searchTerm string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if searchTerm == "" {
			return db
		}
		like := "%" + searchTerm + "%"
		return db.
			Joins("LEFT JOIN categories ON categories.id = products.category_id AND categories.deleted_at IS NULL").
			Joins("LEFT JOIN variants ON variants.product_id = products.id AND variants.deleted_at IS NULL").
			Where(
				"(products.search_vector @@ websearch_to_tsquery('simple', ?) OR variants.search_vector @@ websearch_to_tsquery('simple', ?) OR categories.search_vector @@ websearch_to_tsquery('simple', ?) OR variants.sku ILIKE ?)",
				searchTerm,
				searchTerm,
				searchTerm,
				like,
			)
	}
}

// ScopeProductStockStatus applies stock status filter (requires variants join)
func (s *Storage) ScopeProductStockStatus(status StockStatus) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch status {
		case StockStatusInStock:
			return db.Where("variants.stock_quantity > variants.stock_alert")
		case StockStatusLowStock:
			return db.Where("variants.stock_quantity > 0 AND variants.stock_quantity <= variants.stock_alert")
		case StockStatusOutOfStock:
			return db.Where("variants.stock_quantity = 0")
		default:
			return db
		}
	}
}

// WithProductVariantJoin adds a join to variants table (LEFT or INNER)
func (s *Storage) WithProductVariantJoin(useInnerJoin bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		joinType := "LEFT"
		if useInnerJoin {
			joinType = "INNER"
		}
		return db.Joins(fmt.Sprintf("%s JOIN variants ON variants.product_id = products.id AND variants.deleted_at IS NULL", joinType))
	}
}

// WithProductAggregation adds LATERAL JOIN for product aggregations based on order fields
func (s *Storage) WithProductAggregation(orderBy []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(orderBy) == 0 {
			return db
		}

		needsVariantsCount := false
		needsCostPrice := false
		needsStock := false

		for _, ob := range orderBy {
			switch {
			case ob == "variantsCount" || ob == "-variantsCount":
				needsVariantsCount = true
			case ob == "costPrice" || ob == "-costPrice":
				needsCostPrice = true
			case ob == "stock" || ob == "-stock":
				needsStock = true
			}
		}

		if !needsVariantsCount && !needsCostPrice && !needsStock {
			return db
		}

		// Build aggregation SQL with all needed fields
		var aggFields []string
		if needsVariantsCount {
			aggFields = append(aggFields, "COUNT(*)::int as variants_count")
		}
		if needsCostPrice {
			aggFields = append(aggFields, "AVG(cost_price)::numeric as avg_cost_price")
		}
		if needsStock {
			aggFields = append(aggFields, "SUM(stock_quantity)::int as total_stock")
		}

		joinSQL := fmt.Sprintf(`
			LEFT JOIN LATERAL (
				SELECT %s
				FROM variants
				WHERE variants.product_id = products.id 
					AND variants.deleted_at IS NULL
			) AS product_agg ON true
		`, strings.Join(aggFields, ", "))

		return db.Joins(joinSQL)
	}
}

// WithProductGroupBy adds GROUP BY clause for products with aggregations
func (s *Storage) WithProductGroupBy(needsAggregation bool, aggregationFields []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		groupByColumns := []string{"products.id"}

		if needsAggregation && len(aggregationFields) > 0 {
			groupByColumns = append(groupByColumns, aggregationFields...)
		}

		return db.Group(strings.Join(groupByColumns, ", "))
	}
}

// ParseProductCustomOrdering returns custom ORDER BY clauses for aggregated fields
func (s *Storage) ParseProductCustomOrdering(orderBy []string) []string {
	customOrders := []string{}
	for _, ob := range orderBy {
		switch ob {
		case "variantsCount":
			customOrders = append(customOrders, "product_agg.variants_count ASC")
		case "-variantsCount":
			customOrders = append(customOrders, "product_agg.variants_count DESC")
		case "costPrice":
			customOrders = append(customOrders, "product_agg.avg_cost_price ASC")
		case "-costPrice":
			customOrders = append(customOrders, "product_agg.avg_cost_price DESC")
		case "stock":
			customOrders = append(customOrders, "product_agg.total_stock ASC")
		case "-stock":
			customOrders = append(customOrders, "product_agg.total_stock DESC")
		}
	}
	return customOrders
}
