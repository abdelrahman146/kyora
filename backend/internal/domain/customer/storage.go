package customer

import (
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
)

type Storage struct {
	cache           *cache.Cache
	customer        *database.Repository[Customer]
	customerNote    *database.Repository[CustomerNote]
	customerAddress *database.Repository[CustomerAddress]
}

func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
	st := &Storage{
		cache:           cache,
		customer:        database.NewRepository[Customer](db),
		customerNote:    database.NewRepository[CustomerNote](db),
		customerAddress: database.NewRepository[CustomerAddress](db),
	}
	ensureCustomerSearchIndexes(db)
	return st
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
