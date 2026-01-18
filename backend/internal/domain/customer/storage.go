package customer

import (
	"context"
	"strings"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"gorm.io/gorm"
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

// ScopeHasOrders filters customers by whether they have any non-deleted orders.
func (s *Storage) ScopeHasOrders(hasOrders bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if hasOrders {
			return db.Where("EXISTS (SELECT 1 FROM orders WHERE orders.customer_id = customers.id AND orders.deleted_at IS NULL)")
		}
		return db.Where("NOT EXISTS (SELECT 1 FROM orders WHERE orders.customer_id = customers.id AND orders.deleted_at IS NULL)")
	}
}

// ScopeSocialPlatforms filters customers by social media platform presence.
// Only includes customers that have a non-empty username/number for at least one of the specified platforms.
func (s *Storage) ScopeSocialPlatforms(platforms []string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(platforms) == 0 {
			return db
		}

		conditions := []string{}
		for _, platform := range platforms {
			switch strings.ToLower(platform) {
			case "instagram":
				conditions = append(conditions, "customers.instagram_username IS NOT NULL AND customers.instagram_username != ''")
			case "tiktok":
				conditions = append(conditions, "customers.tiktok_username IS NOT NULL AND customers.tiktok_username != ''")
			case "facebook":
				conditions = append(conditions, "customers.facebook_username IS NOT NULL AND customers.facebook_username != ''")
			case "x":
				conditions = append(conditions, "customers.x_username IS NOT NULL AND customers.x_username != ''")
			case "snapchat":
				conditions = append(conditions, "customers.snapchat_username IS NOT NULL AND customers.snapchat_username != ''")
			case "whatsapp":
				conditions = append(conditions, "customers.whatsapp_number IS NOT NULL AND customers.whatsapp_number != ''")
			}
		}

		if len(conditions) > 0 {
			return db.Where("(" + strings.Join(conditions, " OR ") + ")")
		}
		return db
	}
}

// WithCustomerAggregation adds a LATERAL join to compute orders count and total spent per customer.
// This is used when sorting by ordersCount or totalSpent.
func (s *Storage) WithCustomerAggregation() func(*gorm.DB) *gorm.DB {
	return s.customer.WithJoins(`
		LEFT JOIN LATERAL (
			SELECT 
				COUNT(DISTINCT orders.id)::int as orders_count,
				COALESCE(SUM(orders.total), 0)::numeric as total_spent
			FROM orders
			WHERE orders.customer_id = customers.id 
				AND orders.deleted_at IS NULL
				AND orders.status NOT IN ('cancelled', 'returned', 'failed')
		) AS customer_agg ON true
	`)
}
