package owner

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type OwnerService struct {
	ownerRepo    *OwnerRepository
	storeService *store.StoreService
}

func NewOwnerService(ownerRepo *OwnerRepository, storeService *store.StoreService) *OwnerService {
	return &OwnerService{
		ownerRepo:    ownerRepo,
		storeService: storeService,
	}
}

func (s *OwnerService) CreateOwner(ctx context.Context, storeID string, req *CreateOwnerRequest) (*Owner, error) {
	store, err := s.storeService.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	owner := &Owner{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		StoreID:   store.ID,
	}

	if err := s.ownerRepo.createOne(ctx, owner); err != nil {
		return nil, err
	}

	return owner, nil
}

func (s *OwnerService) UpdateOwner(ctx context.Context, storeID string, ownerID string, req *UpdateOwnerRequest) (*Owner, error) {
	owner, err := s.ownerRepo.FindOne(ctx, s.ownerRepo.scopeID(ownerID), s.ownerRepo.scopeStoreID(storeID))
	if err != nil {
		return nil, err
	}

	if req.FirstName != "" {
		owner.FirstName = req.FirstName
	}
	if req.LastName != "" {
		owner.LastName = req.LastName
	}

	if err := s.ownerRepo.updateOne(ctx, owner); err != nil {
		return nil, err
	}

	return owner, nil
}

func (s *OwnerService) GetOwnerByID(ctx context.Context, storeID string, ownerID string) (*Owner, error) {
	return s.ownerRepo.FindOne(ctx, s.ownerRepo.scopeID(ownerID), s.ownerRepo.scopeStoreID(storeID))
}

func (s *OwnerService) ListOwners(ctx context.Context, storeID string) ([]*Owner, error) {
	return s.ownerRepo.List(ctx, s.ownerRepo.scopeStoreID(storeID))
}

func (s *OwnerService) DeleteOwner(ctx context.Context, storeID string, ownerID string) error {
	_, err := s.ownerRepo.FindOne(ctx, s.ownerRepo.scopeID(ownerID), s.ownerRepo.scopeStoreID(storeID))
	if err != nil {
		return err
	}

	return s.ownerRepo.deleteOne(ctx, s.ownerRepo.scopeID(ownerID))
}
