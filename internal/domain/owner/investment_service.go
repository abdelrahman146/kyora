package owner

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/shopspring/decimal"
)

type InvestmentService struct {
	investmentRepo *investmentRepository
	storeService   *store.StoreService
}

func NewInvestmentService(investmentRepo *investmentRepository, storeService *store.StoreService) *InvestmentService {
	return &InvestmentService{
		investmentRepo: investmentRepo,
		storeService:   storeService,
	}
}

func (s *InvestmentService) CreateInvestment(ctx context.Context, storeID string, req *CreateInvestmentRequest) (*Investment, error) {
	store, err := s.storeService.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	investment := &Investment{
		Name:     req.Name,
		Amount:   req.Amount,
		Currency: store.Currency,
		Note:     req.Note,
		StoreID:  store.ID,
	}

	if err := s.investmentRepo.createOne(ctx, investment); err != nil {
		return nil, err
	}

	return investment, nil
}

func (s *InvestmentService) UpdateInvestment(ctx context.Context, storeID, investmentID string, req *UpdateInvestmentRequest) (*Investment, error) {
	investment, err := s.investmentRepo.findOne(ctx, s.investmentRepo.scopeID(investmentID), s.investmentRepo.scopeStoreID(storeID))
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		investment.Name = req.Name
	}
	if !req.Amount.IsZero() {
		investment.Amount = req.Amount
	}
	if req.Note != "" {
		investment.Note = req.Note
	}

	if err := s.investmentRepo.updateOne(ctx, investment); err != nil {
		return nil, err
	}

	return investment, nil
}

func (s *InvestmentService) DeleteInvestment(ctx context.Context, storeID, investmentID string) error {
	_, err := s.investmentRepo.findOne(ctx, s.investmentRepo.scopeID(investmentID), s.investmentRepo.scopeStoreID(storeID))
	if err != nil {
		return err
	}

	return s.investmentRepo.deleteOne(ctx, s.investmentRepo.scopeID(investmentID))
}

func (s *InvestmentService) GetInvestmentByID(ctx context.Context, storeID string, investmentID string) (*Investment, error) {
	return s.investmentRepo.findOne(ctx, s.investmentRepo.scopeID(investmentID), s.investmentRepo.scopeStoreID(storeID))
}

func (s *InvestmentService) ListInvestments(ctx context.Context, storeID string, page, pageSize int, orderBy string, ascending bool) ([]*Investment, error) {
	return s.investmentRepo.list(ctx,
		s.investmentRepo.scopeStoreID(storeID),
		db.WithPagination(page, pageSize),
		db.WithSorting(orderBy, ascending),
	)
}

func (s *InvestmentService) CountInvestments(ctx context.Context, storeID string) (int64, error) {
	return s.investmentRepo.count(ctx, s.investmentRepo.scopeStoreID(storeID))
}

func (s *InvestmentService) CalculateTotalInvestedAmount(ctx context.Context, storeID string) (decimal.Decimal, error) {
	total, err := s.investmentRepo.sumAmount(ctx, s.investmentRepo.scopeStoreID(storeID))
	if err != nil {
		return decimal.Zero, err
	}
	return total, nil
}

func (s *InvestmentService) CalculateTotalInvestedAmountByOwner(ctx context.Context, storeID, ownerID string) (decimal.Decimal, error) {
	total, err := s.investmentRepo.sumAmount(ctx,
		s.investmentRepo.scopeStoreID(storeID),
		s.investmentRepo.scopeOwnerID(ownerID),
	)
	if err != nil {
		return decimal.Zero, err
	}
	return total, nil
}
