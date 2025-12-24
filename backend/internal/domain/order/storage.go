package order

import (
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
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
