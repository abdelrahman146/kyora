package asset

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
)

type Storage struct {
	cache *cache.Cache
	asset *database.Repository[Asset]
}

func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
	return &Storage{
		cache: cache,
		asset: database.NewRepository[Asset](db),
	}
}

func (s *Storage) Cache() *cache.Cache { return s.cache }

func (s *Storage) Create(ctx context.Context, a *Asset) error {
	return s.asset.CreateOne(ctx, a)
}

func (s *Storage) Update(ctx context.Context, a *Asset) error {
	return s.asset.UpdateOne(ctx, a)
}

func (s *Storage) GetByID(ctx context.Context, businessID, assetID string) (*Asset, error) {
	return s.asset.FindOne(ctx, s.asset.ScopeBusinessID(businessID), s.asset.ScopeID(assetID))
}

func (s *Storage) FindByBusinessAndIdempotencyKey(ctx context.Context, businessID, idemKey string) (*Asset, error) {
	return s.asset.FindOne(ctx, s.asset.ScopeBusinessID(businessID), s.asset.ScopeEquals(AssetSchema.IdempotencyKey, idemKey))
}
