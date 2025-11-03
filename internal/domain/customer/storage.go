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
	return &Storage{
		cache:           cache,
		customer:        database.NewRepository[Customer](db),
		customerNote:    database.NewRepository[CustomerNote](db),
		customerAddress: database.NewRepository[CustomerAddress](db),
	}
}
