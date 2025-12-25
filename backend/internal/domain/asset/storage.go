package asset

import (
	"context"
	"strings"
	"time"

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

func (s *Storage) Delete(ctx context.Context, a *Asset) error {
	return s.asset.DeleteOne(ctx, a)
}

func (s *Storage) GetByID(ctx context.Context, businessID, assetID string) (*Asset, error) {
	return s.asset.FindOne(ctx, s.asset.ScopeBusinessID(businessID), s.asset.ScopeID(assetID))
}

// FindByID returns an asset by ID without applying business/workspace scoping.
// Intended for internal operations (e.g., public reads, maintenance jobs).
func (s *Storage) FindByID(ctx context.Context, assetID string) (*Asset, error) {
	return s.asset.FindByID(ctx, assetID)
}

func (s *Storage) FindByBusinessAndIdempotencyKey(ctx context.Context, businessID, idemKey string) (*Asset, error) {
	return s.asset.FindOne(ctx, s.asset.ScopeBusinessID(businessID), s.asset.ScopeEquals(AssetSchema.IdempotencyKey, idemKey))
}

func (s *Storage) ListExpiredPending(ctx context.Context, now time.Time, limit int) ([]*Asset, error) {
	if limit <= 0 {
		limit = 200
	}
	return s.asset.FindMany(ctx,
		s.asset.ScopeEquals(AssetSchema.Status, StatusPending),
		s.asset.ScopeWhere("upload_expires_at IS NOT NULL AND upload_expires_at < ?", now),
		s.asset.WithOrderBy([]string{"upload_expires_at ASC"}),
		s.asset.WithLimit(limit),
	)
}

// ListReadyOrphans returns ready assets that are not referenced by known URL-only fields.
// We only consider assets with a non-empty public_url.
func (s *Storage) ListReadyOrphans(ctx context.Context, now time.Time, minAge time.Duration, limit int) ([]*Asset, error) {
	if limit <= 0 {
		limit = 200
	}
	cutoff := now.Add(-minAge)

	// Note: this is Postgres-specific (jsonb_build_array) and intentionally uses
	// parameterized SQL to avoid injection issues.
	where := strings.Join([]string{
		"public_url <> ''",
		"status = 'ready'",
		"completed_at IS NOT NULL AND completed_at < ?",
		"NOT EXISTS (SELECT 1 FROM businesses b WHERE b.deleted_at IS NULL AND b.logo_url = uploaded_assets.public_url)",
		"NOT EXISTS (SELECT 1 FROM products p WHERE p.deleted_at IS NULL AND p.photos @> jsonb_build_array(uploaded_assets.public_url))",
		"NOT EXISTS (SELECT 1 FROM variants v WHERE v.deleted_at IS NULL AND v.photos @> jsonb_build_array(uploaded_assets.public_url))",
	}, " AND ")

	return s.asset.FindMany(ctx,
		s.asset.ScopeWhere(where, cutoff),
		s.asset.WithOrderBy([]string{"completed_at ASC"}),
		s.asset.WithLimit(limit),
	)
}
