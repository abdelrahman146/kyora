package billing

import (
	"context"
	"log/slog"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
)

type Storage struct {
	cache        *cache.Cache
	db           *database.Database
	plan         *database.Repository[Plan]
	subscription *database.Repository[Subscription]
}

func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
	s := &Storage{
		cache:        cache,
		db:           db,
		plan:         database.NewRepository[Plan](db),
		subscription: database.NewRepository[Subscription](db),
	}
	s.init()
	return s
}

func (s *Storage) init() {
	// upsert plans into the database
	ctx := context.Background()
	for _, plan := range plans {
		existing, err := s.plan.FindOne(ctx, s.plan.ScopeEquals(PlanSchema.Descriptor, plan.Descriptor))
		if err != nil && !database.IsRecordNotFound(err) {
			slog.Default().Error("failed to fetch existing plan", "error", err, "descriptor", plan.Descriptor)
			continue
		}
		if existing == nil {
			if err := s.plan.CreateOne(ctx, &plan); err != nil {
				slog.Default().Error("failed to create plan", "error", err, "descriptor", plan.Descriptor)
				continue
			}
		} else {
			plan.ID = existing.ID
			if err := s.plan.UpdateOne(ctx, &plan); err != nil {
				slog.Default().Error("failed to update plan", "error", err, "descriptor", plan.Descriptor)
				continue
			}
		}
	}
}

// CountBusinessesByWorkspace returns count of businesses for a workspace
func (s *Storage) CountBusinessesByWorkspace(ctx context.Context, workspaceID string) (int64, error) {
	var count int64
	err := s.db.Conn(ctx).Table("businesses").Where("workspace_id = ?", workspaceID).Count(&count).Error
	return count, err
}

// CountMonthlyOrdersByWorkspace counts orders within date range across all businesses in a workspace
func (s *Storage) CountMonthlyOrdersByWorkspace(ctx context.Context, workspaceID string, from, to time.Time) (int64, error) {
	var count int64
	q := s.db.Conn(ctx).Table("orders as o").
		Joins("join businesses b on b.id = o.business_id").
		Where("b.workspace_id = ?", workspaceID)
	if !from.IsZero() && !to.IsZero() {
		q = q.Where("o.ordered_at BETWEEN ? AND ?", from, to)
	} else if !from.IsZero() {
		q = q.Where("o.ordered_at >= ?", from)
	} else if !to.IsZero() {
		q = q.Where("o.ordered_at <= ?", to)
	}
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
