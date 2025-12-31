package asset

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
)

// Storage provides data access for assets.
type Storage struct {
	cache *cache.Cache
	asset *database.Repository[Asset]
}

// NewStorage creates a new asset storage instance.
func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
	return &Storage{
		cache: cache,
		asset: database.NewRepository[Asset](db),
	}
}

// Cache returns the cache instance.
func (s *Storage) Cache() *cache.Cache { return s.cache }

// Create creates a new asset record.
func (s *Storage) Create(ctx context.Context, a *Asset) error {
	return s.asset.CreateOne(ctx, a)
}

// Update updates an existing asset record.
func (s *Storage) Update(ctx context.Context, a *Asset) error {
	return s.asset.UpdateOne(ctx, a)
}

// Delete deletes an asset record.
func (s *Storage) Delete(ctx context.Context, a *Asset) error {
	return s.asset.DeleteOne(ctx, a)
}

// GetByID returns an asset by ID scoped to a business.
func (s *Storage) GetByID(ctx context.Context, businessID, assetID string) (*Asset, error) {
	return s.asset.FindOne(ctx, s.asset.ScopeBusinessID(businessID), s.asset.ScopeID(assetID))
}

// FindByID returns an asset by ID without applying business/workspace scoping.
// Intended for internal operations (e.g., public reads, GC jobs).
func (s *Storage) FindByID(ctx context.Context, assetID string) (*Asset, error) {
	return s.asset.FindByID(ctx, assetID)
}

// MarkUploadComplete marks a multipart upload as complete.
// TODO: Implement proper JSONB update when Repository supports raw queries.
func (s *Storage) MarkUploadComplete(ctx context.Context, assetID string) error {
	// For now, just return nil - this will be implemented in the GC rewrite
	return nil
}
