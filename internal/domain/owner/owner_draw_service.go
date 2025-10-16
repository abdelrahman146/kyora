package owner

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/types"
	"github.com/shopspring/decimal"
)

type OwnerDrawService struct {
	ownerDrawRepo *ownerDrawRepository
	storeService  *store.StoreService
}

func NewOwnerDrawService(ownerDrawRepo *ownerDrawRepository, storeService *store.StoreService) *OwnerDrawService {
	return &OwnerDrawService{
		ownerDrawRepo: ownerDrawRepo,
		storeService:  storeService,
	}
}

func (s *OwnerDrawService) CreateOwnerDraw(ctx context.Context, storeID string, req *CreateOwnerDrawRequest) (*OwnerDraw, error) {
	store, err := s.storeService.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	ownerDraw := &OwnerDraw{
		OwnerID:     req.OwnerID,
		Amount:      req.Amount,
		Note:        req.Note,
		WithdrawnAt: req.WithdrawnAt,
		Currency:    store.Currency,
		StoreID:     store.ID,
	}

	if err := s.ownerDrawRepo.createOne(ctx, ownerDraw); err != nil {
		return nil, err
	}

	return ownerDraw, nil
}

func (s *OwnerDrawService) UpdateOwnerDraw(ctx context.Context, storeID, ownerDrawID string, req *UpdateOwnerDrawRequest) (*OwnerDraw, error) {
	ownerDraw, err := s.ownerDrawRepo.findOne(ctx, s.ownerDrawRepo.scopeID(ownerDrawID), s.ownerDrawRepo.scopeStoreID(storeID))
	if err != nil {
		return nil, err
	}

	if req.OwnerID != "" {
		ownerDraw.OwnerID = req.OwnerID
	}
	if !req.Amount.IsZero() {
		ownerDraw.Amount = req.Amount
	}
	if req.Note != "" {
		ownerDraw.Note = req.Note
	}
	if !req.WithdrawnAt.IsZero() {
		ownerDraw.WithdrawnAt = req.WithdrawnAt
	}

	if err := s.ownerDrawRepo.updateOne(ctx, ownerDraw); err != nil {
		return nil, err
	}

	return ownerDraw, nil
}

func (s *OwnerDrawService) GetOwnerDrawByID(ctx context.Context, storeID, ownerDrawID string) (*OwnerDraw, error) {
	return s.ownerDrawRepo.findOne(ctx, s.ownerDrawRepo.scopeID(ownerDrawID), s.ownerDrawRepo.scopeStoreID(storeID))
}

func (s *OwnerDrawService) ListOwnerDraws(ctx context.Context, storeID string, listReq *types.ListRequest) ([]*OwnerDraw, error) {
	return s.ownerDrawRepo.list(ctx, s.ownerDrawRepo.scopeStoreID(storeID), db.WithPagination(listReq.Page, listReq.PageSize), db.WithOrderBy(listReq.OrderBy))
}

func (s *OwnerDrawService) CountOwnerDraws(ctx context.Context, storeID string) (int64, error) {
	return s.ownerDrawRepo.count(ctx, s.ownerDrawRepo.scopeStoreID(storeID))
}

func (s *OwnerDrawService) DeleteOwnerDraw(ctx context.Context, storeID, ownerDrawID string) error {
	_, err := s.ownerDrawRepo.findOne(ctx, s.ownerDrawRepo.scopeID(ownerDrawID), s.ownerDrawRepo.scopeStoreID(storeID))
	if err != nil {
		return err
	}

	return s.ownerDrawRepo.deleteOne(ctx, s.ownerDrawRepo.scopeID(ownerDrawID))
}

func (s *OwnerDrawService) SumTotalOwnerDraws(ctx context.Context, storeID string) (decimal.Decimal, error) {
	return s.ownerDrawRepo.sumAmount(ctx, s.ownerDrawRepo.scopeStoreID(storeID))
}

func (s *OwnerDrawService) SumTotalOwnerDrawsByOwner(ctx context.Context, storeID, ownerID string) (decimal.Decimal, error) {
	return s.ownerDrawRepo.sumAmount(ctx, s.ownerDrawRepo.scopeStoreID(storeID), s.ownerDrawRepo.scopeOwnerID(ownerID))
}

// SumOwnerDraws returns the total amount of owner draws in the given date range.
// When from or to is zero, it behaves as an open-ended bound (all-time if both are zero).
func (s *OwnerDrawService) SumOwnerDraws(ctx context.Context, storeID string, from, to time.Time) (decimal.Decimal, error) {
	return s.ownerDrawRepo.sumAmount(ctx, s.ownerDrawRepo.scopeStoreID(storeID), s.ownerDrawRepo.scopeWithDrawnAt(from, to))
}
