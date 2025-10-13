package store

import (
	"context"
	"fmt"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/utils"
)

type StoreService struct {
	storeRepo *StoreRepository
}

func NewStoreService(storeRepo *StoreRepository) *StoreService {
	return &StoreService{storeRepo: storeRepo}
}

func (s *StoreService) CreateStore(ctx context.Context, organizationID string, storeReq *CreateStoreRequest) (*Store, error) {
	store := &Store{
		OrganizationID: organizationID,
		Name:           storeReq.Name,
		Currency:       storeReq.Currency,
		CountryCode:    storeReq.CountryCode,
		VatRate:        storeReq.VatRate,
	}
	maxAttempts := 5
	for i := range maxAttempts {
		store.Code = utils.ID.NewBase62(6)
		available, err := s.IsStoreCodeAvailable(ctx, organizationID, store.Code)
		if err != nil {
			return nil, err
		}
		if available {
			break
		}
		if i == maxAttempts-1 {
			return nil, utils.Problem.InternalError().WithError(fmt.Errorf("max attempts reached for generating unique store code"))
		}
	}
	if err := s.storeRepo.CreateOne(ctx, store); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *StoreService) IsStoreCodeAvailable(ctx context.Context, organizationID, code string) (bool, error) {
	existingStore, err := s.storeRepo.FindOne(ctx, s.storeRepo.ScopeOrganizationID(organizationID), s.storeRepo.ScopeCode(code))
	if err != nil {
		return false, err
	}
	return existingStore == nil, nil
}

func (s *StoreService) ValidateStoreID(ctx context.Context, storeID string) error {
	_, err := s.storeRepo.FindOne(ctx, s.storeRepo.ScopeID(storeID))
	if err != nil {
		return err
	}
	return nil
}

func (s *StoreService) UpdateStore(ctx context.Context, storeID string, storeReq *UpdateStoreRequest) (*Store, error) {
	if _, err := s.storeRepo.FindOne(ctx, s.storeRepo.ScopeID(storeID)); err != nil {
		return nil, err
	}
	store := &Store{
		Name:         storeReq.Name,
		Currency:     storeReq.Currency,
		VatRate:      storeReq.VatRate,
		CountryCode:  storeReq.CountryCode,
		SafetyBuffer: storeReq.SafetyBuffer,
	}
	if err := s.storeRepo.PatchOne(ctx, store, s.storeRepo.ScopeID(storeID), db.WithReturning(&store)); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *StoreService) GetStoreByID(ctx context.Context, id string) (*Store, error) {
	store, err := s.storeRepo.FindOne(ctx, db.WithScopes(s.storeRepo.ScopeID(id)))
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (s *StoreService) ListOrganizationStores(ctx context.Context, orgID string) ([]*Store, error) {
	stores, err := s.storeRepo.List(ctx, s.storeRepo.ScopeOrganizationID(orgID))
	if err != nil {
		return nil, err
	}
	return stores, nil
}

func (s *StoreService) DeleteStore(ctx context.Context, storeID string) error {
	if err := s.storeRepo.DeleteOne(ctx, s.storeRepo.ScopeID(storeID)); err != nil {
		return err
	}
	return nil
}
