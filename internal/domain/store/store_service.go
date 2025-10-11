package store

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
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
		Locale:         storeReq.Locale,
		Currency:       storeReq.Currency,
		Timezone:       storeReq.Timezone,
		VATRate:        storeReq.VATRate,
	}
	if err := s.storeRepo.CreateOne(ctx, store); err != nil {
		return nil, db.HandleDBError(err)
	}
	return store, nil
}

func (s *StoreService) UpdateStore(ctx context.Context, storeID string, storeReq *UpdateStoreRequest) (*Store, error) {
	if _, err := s.storeRepo.FindOne(ctx, s.storeRepo.ScopeID(storeID)); err != nil {
		return nil, db.HandleDBError(err)
	}
	store := &Store{
		Name:     storeReq.Name,
		Locale:   storeReq.Locale,
		Currency: storeReq.Currency,
		Timezone: storeReq.Timezone,
		VATRate:  storeReq.VATRate,
	}
	if err := s.storeRepo.PatchOne(ctx, store, s.storeRepo.ScopeID(storeID), db.WithReturning(&store)); err != nil {
		return nil, db.HandleDBError(err)
	}
	return store, nil
}

func (s *StoreService) GetStoreByID(ctx context.Context, id string) (*Store, error) {
	store, err := s.storeRepo.FindOne(ctx, db.WithScopes(s.storeRepo.ScopeID(id)))
	if err != nil {
		return nil, db.HandleDBError(err)
	}
	return store, nil
}

func (s *StoreService) ListOrganizationStores(ctx context.Context, orgID string) ([]*Store, error) {
	stores, err := s.storeRepo.List(ctx, s.storeRepo.ScopeOrganizationID(orgID))
	if err != nil {
		return nil, db.HandleDBError(err)
	}
	return stores, nil
}

func (s *StoreService) DeleteStore(ctx context.Context, storeID string) error {
	if err := s.storeRepo.DeleteOne(ctx, s.storeRepo.ScopeID(storeID)); err != nil {
		return db.HandleDBError(err)
	}
	return nil
}
