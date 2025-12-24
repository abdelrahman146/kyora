package business

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
)

type Storage struct {
	cache    *cache.Cache
	business *database.Repository[Business]
	zone     *database.Repository[ShippingZone]
}

func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
	return &Storage{
		cache:    cache,
		business: database.NewRepository[Business](db),
		zone:     database.NewRepository[ShippingZone](db),
	}
}

func (s *Storage) CreateShippingZone(ctx context.Context, zone *ShippingZone) error {
	if zone == nil {
		return problem.BadRequest("shipping zone is required")
	}
	return s.zone.CreateOne(ctx, zone)
}

func (s *Storage) UpdateShippingZone(ctx context.Context, zone *ShippingZone) error {
	if zone == nil {
		return problem.BadRequest("shipping zone is required")
	}
	return s.zone.UpdateOne(ctx, zone)
}

func (s *Storage) DeleteShippingZone(ctx context.Context, zone *ShippingZone) error {
	if zone == nil {
		return problem.BadRequest("shipping zone is required")
	}
	return s.zone.DeleteOne(ctx, zone)
}

func (s *Storage) GetShippingZoneByID(ctx context.Context, businessID, zoneID string) (*ShippingZone, error) {
	return s.zone.FindByID(ctx, zoneID, s.zone.ScopeBusinessID(businessID))
}

func (s *Storage) ListShippingZones(ctx context.Context, businessID string) ([]*ShippingZone, error) {
	return s.zone.FindMany(ctx,
		s.zone.ScopeBusinessID(businessID),
		s.zone.WithOrderBy([]string{"name ASC"}),
	)
}

func (s *Storage) CountShippingZones(ctx context.Context, businessID string) (int64, error) {
	return s.zone.Count(ctx, s.zone.ScopeBusinessID(businessID))
}
