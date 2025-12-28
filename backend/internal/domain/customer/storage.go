package customer

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
)

type Storage struct {
	db              *database.Database
	cache           *cache.Cache
	customer        *database.Repository[Customer]
	customerNote    *database.Repository[CustomerNote]
	customerAddress *database.Repository[CustomerAddress]
}

func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
	st := &Storage{
		db:              db,
		cache:           cache,
		customer:        database.NewRepository[Customer](db),
		customerNote:    database.NewRepository[CustomerNote](db),
		customerAddress: database.NewRepository[CustomerAddress](db),
	}
	ensureCustomerSearchIndexes(db)
	return st
}

// CustomerAggregation holds order count and total spent for a customer
type CustomerAggregation struct {
	CustomerID  string  `gorm:"column:customer_id"`
	OrdersCount int     `gorm:"column:orders_count"`
	TotalSpent  float64 `gorm:"column:total_spent"`
}

// GetCustomerAggregations fetches order count and total spent for given customer IDs
func (s *Storage) GetCustomerAggregations(ctx context.Context, businessID string, customerIDs []string) (map[string]CustomerAggregation, error) {
	if len(customerIDs) == 0 {
		return make(map[string]CustomerAggregation), nil
	}

	var results []CustomerAggregation
	err := s.db.Conn(ctx).
		Table("orders").
		Select(
			"customer_id",
			"COUNT(DISTINCT id)::int as orders_count",
			"COALESCE(SUM(total), 0)::numeric as total_spent",
		).
		Where("customer_id IN ?", customerIDs).
		Where("deleted_at IS NULL").
		Where("business_id = ?", businessID).
		Where("status NOT IN ?", []string{"cancelled", "returned", "failed"}).
		Group("customer_id").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert to map for easy lookup
	aggMap := make(map[string]CustomerAggregation)
	for _, agg := range results {
		aggMap[agg.CustomerID] = agg
	}

	return aggMap, nil
}

func ensureCustomerSearchIndexes(db *database.Database) {
	conn := db.GetDB()

	// Weighted full-text document (most relevant fields first).
	expr := "" +
		"setweight(to_tsvector('simple', coalesce(name,'')), 'A') || " +
		"setweight(to_tsvector('simple', coalesce(email,'')), 'A') || " +
		"setweight(to_tsvector('simple', coalesce(phone_number,'')), 'B') || " +
		"setweight(to_tsvector('simple', coalesce(whatsapp_number,'')), 'B') || " +
		"setweight(to_tsvector('simple', coalesce(instagram_username,'')), 'C') || " +
		"setweight(to_tsvector('simple', coalesce(tiktok_username,'')), 'C')"

	database.EnsureGeneratedTSVectorColumn(conn, CustomerTable, "search_vector", expr)
	database.EnsureGinIndex(conn, "customers_search_vector_gin_idx", CustomerTable, "search_vector")

	// Trigram indexes speed up substring matches (ILIKE) for user-friendly partial search.
	database.EnsureTrigramGinIndex(conn, "customers_name_trgm_idx", CustomerTable, "name")
	database.EnsureTrigramGinIndex(conn, "customers_email_trgm_idx", CustomerTable, "email")
}
