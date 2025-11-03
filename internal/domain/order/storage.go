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
	return &Storage{
		cache:     cache,
		order:     database.NewRepository[Order](db),
		orderItem: database.NewRepository[OrderItem](db),
		orderNote: database.NewRepository[OrderNote](db),
	}
}
