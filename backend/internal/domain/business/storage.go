package business

import (
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
)

type Storage struct {
	cache    *cache.Cache
	business *database.Repository[Business]
}

func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
	return &Storage{
		cache:    cache,
		business: database.NewRepository[Business](db),
	}
}
