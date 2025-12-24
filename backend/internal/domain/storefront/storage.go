package storefront

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"gorm.io/gorm"
)

type Storage struct {
	requests *database.Repository[StorefrontRequest]
	cache    *cache.Cache
}

func NewStorage(db *database.Database, c *cache.Cache) *Storage {
	return &Storage{
		requests: database.NewRepository[StorefrontRequest](db),
		cache:    c,
	}
}

func (s *Storage) Cache() *cache.Cache { return s.cache }

func (s *Storage) CreateRequest(ctx context.Context, req *StorefrontRequest, opts ...func(db *gorm.DB) *gorm.DB) error {
	return s.requests.CreateOne(ctx, req, opts...)
}

func (s *Storage) UpdateRequest(ctx context.Context, req *StorefrontRequest, opts ...func(db *gorm.DB) *gorm.DB) error {
	return s.requests.UpdateOne(ctx, req, opts...)
}

func (s *Storage) GetRequestByKey(ctx context.Context, businessID, key string) (*StorefrontRequest, error) {
	return s.requests.FindOne(ctx,
		s.requests.ScopeEquals(StorefrontRequestSchema.BusinessID, businessID),
		s.requests.ScopeEquals(StorefrontRequestSchema.IdempotencyKey, key),
	)
}
