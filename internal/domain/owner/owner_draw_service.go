package owner

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/govalues/decimal"
)

type OwnerDrawService struct {
	ownerDrawRepo *OwnerDrawRepository
	storeService  *store.StoreService
}

func NewOwnerDrawService(ownerDrawRepo *OwnerDrawRepository, storeService *store.StoreService) *OwnerDrawService {
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
		OwnerID: req.OwnerID,
		Amount:  req.Amount,
		Note:    req.Note,
		StoreID: store.ID,
	}

	if err := s.ownerDrawRepo.CreateOne(ctx, ownerDraw); err != nil {
		return nil, err
	}

	return ownerDraw, nil
}

func (s *OwnerDrawService) UpdateOwnerDraw(ctx context.Context, storeID, ownerDrawID string, req *UpdateOwnerDrawRequest) (*OwnerDraw, error) {
	ownerDraw, err := s.ownerDrawRepo.FindOne(ctx, s.ownerDrawRepo.ScopeID(ownerDrawID), s.ownerDrawRepo.ScopeStoreID(storeID))
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

	if err := s.ownerDrawRepo.UpdateOne(ctx, ownerDraw); err != nil {
		return nil, err
	}

	return ownerDraw, nil
}

func (s *OwnerDrawService) GetOwnerDrawByID(ctx context.Context, storeID, ownerDrawID string) (*OwnerDraw, error) {
	return s.ownerDrawRepo.FindOne(ctx, s.ownerDrawRepo.ScopeID(ownerDrawID), s.ownerDrawRepo.ScopeStoreID(storeID))
}

func (s *OwnerDrawService) ListOwnerDraws(ctx context.Context, storeID string, page, pageSize int, orderBy string, ascending bool) ([]*OwnerDraw, error) {
	return s.ownerDrawRepo.List(ctx, s.ownerDrawRepo.ScopeStoreID(storeID), db.WithPagination(page, pageSize), db.WithSorting(orderBy, ascending))
}

func (s *OwnerDrawService) CountOwnerDraws(ctx context.Context, storeID string) (int64, error) {
	return s.ownerDrawRepo.Count(ctx, s.ownerDrawRepo.ScopeStoreID(storeID))
}

func (s *OwnerDrawService) DeleteOwnerDraw(ctx context.Context, storeID, ownerDrawID string) error {
	_, err := s.ownerDrawRepo.FindOne(ctx, s.ownerDrawRepo.ScopeID(ownerDrawID), s.ownerDrawRepo.ScopeStoreID(storeID))
	if err != nil {
		return err
	}

	return s.ownerDrawRepo.DeleteOne(ctx, s.ownerDrawRepo.ScopeID(ownerDrawID))
}

func (s *OwnerDrawService) SumTotalOwnerDraws(ctx context.Context, storeID string) (decimal.Decimal, error) {
	return s.ownerDrawRepo.SumAmount(ctx, s.ownerDrawRepo.ScopeStoreID(storeID))
}

func (s *OwnerDrawService) SumTotalOwnerDrawsByOwner(ctx context.Context, storeID, ownerID string) (decimal.Decimal, error) {
	return s.ownerDrawRepo.SumAmount(ctx, s.ownerDrawRepo.ScopeStoreID(storeID), s.ownerDrawRepo.ScopeOwnerID(ownerID))
}
