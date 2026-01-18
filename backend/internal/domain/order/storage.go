package order

import (
	"strings"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"gorm.io/gorm"
)

type Storage struct {
	cache     *cache.Cache
	order     *database.Repository[Order]
	orderItem *database.Repository[OrderItem]
	orderNote *database.Repository[OrderNote]
}

func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
	st := &Storage{
		cache:     cache,
		order:     database.NewRepository[Order](db),
		orderItem: database.NewRepository[OrderItem](db),
		orderNote: database.NewRepository[OrderNote](db),
	}
	ensureOrderSearchIndexes(db)
	return st
}

func ensureOrderSearchIndexes(db *database.Database) {
	conn := db.GetDB()

	expr := "" +
		"setweight(to_tsvector('simple', coalesce(order_number,'')), 'A') || " +
		"setweight(to_tsvector('simple', coalesce(channel,'')), 'B') || " +
		"setweight(to_tsvector('simple', coalesce(currency,'')), 'C') || " +
		"setweight(to_tsvector('simple', coalesce(payment_reference,'')), 'C')"

	database.EnsureGeneratedTSVectorColumn(conn, OrderTable, "search_vector", expr)
	database.EnsureGinIndex(conn, "orders_search_vector_gin_idx", OrderTable, "search_vector")

	// Trigram index for fast substring lookup on order numbers.
	database.EnsureTrigramGinIndex(conn, "orders_order_number_trgm_idx", OrderTable, "order_number")
}

// ScopeOrderSearch adds search conditions for orders, including customer search vector.
// Requires a LEFT JOIN on customers table to be called before this scope.
func (s *Storage) ScopeOrderSearch(term string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		like := "%" + term + "%"
		return db.Where(
			"(orders.search_vector @@ websearch_to_tsquery('simple', ?) OR customers.search_vector @@ websearch_to_tsquery('simple', ?) OR orders.order_number ILIKE ? OR customers.name ILIKE ? OR customers.email ILIKE ?)",
			term, term, like, like, like,
		)
	}
}

// ScopeChannels filters orders by channel (case-insensitive).
func (s *Storage) ScopeChannels(channels []string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(channels) == 0 {
			return db
		}

		normalized := make([]string, 0, len(channels))
		seen := map[string]struct{}{}
		for _, ch := range channels {
			c := strings.ToLower(strings.TrimSpace(ch))
			if c == "" || seen[c] != struct{}{} {
				continue
			}
			seen[c] = struct{}{}
			normalized = append(normalized, c)
		}

		if len(normalized) > 0 {
			return db.Where("LOWER(orders.channel) IN ?", normalized)
		}
		return db
	}
}

// WithOrderCustomerJoin adds a LEFT JOIN with customers table for search.
func (s *Storage) WithOrderCustomerJoin() func(*gorm.DB) *gorm.DB {
	return s.order.WithJoins("LEFT JOIN customers ON customers.id = orders.customer_id")
}
