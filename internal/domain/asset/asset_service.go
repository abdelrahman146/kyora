package asset

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type AssetService struct {
	assetRepo    *AssetRepository
	storeService *store.StoreService
}

func NewAssetService(assetRepo *AssetRepository, storeService *store.StoreService) *AssetService {
	return &AssetService{assetRepo: assetRepo, storeService: storeService}
}

func (s *AssetService) GetAssetByID(ctx context.Context, storeID, assetID string) (*Asset, error) {
	asset, err := s.assetRepo.FindOne(ctx, s.assetRepo.scopeID(assetID), s.assetRepo.scopeStoreID(storeID))
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *AssetService) ListAssets(ctx context.Context, storeID string, filter *AssetFilter, page, pageSize int, orderBy string, ascending bool) ([]*Asset, error) {
	assets, err := s.assetRepo.List(ctx, s.assetRepo.scopeStoreID(storeID), s.assetRepo.scopeFilter(filter), db.WithPagination(page, pageSize), db.WithSorting(orderBy, ascending))
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (s *AssetService) CreateAsset(ctx context.Context, storeID string, req *CreateAssetRequest) (*Asset, error) {
	store, err := s.storeService.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, err
	}
	asset := &Asset{
		StoreID:  storeID,
		Name:     req.Name,
		Type:     req.Type,
		Currency: store.Currency,
		Value:    req.Value,
		Note:     req.Note,
	}
	if err := s.assetRepo.createOne(ctx, asset); err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *AssetService) UpdateAsset(ctx context.Context, storeID, assetID string, req *UpdateAssetRequest) (*Asset, error) {
	asset, err := s.assetRepo.FindOne(ctx, s.assetRepo.scopeID(assetID), s.assetRepo.scopeStoreID(storeID))
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		asset.Name = req.Name
	}
	if req.Type != "" {
		asset.Type = req.Type
	}
	if !req.Value.IsZero() {
		asset.Value = req.Value
	}
	if req.Currency != "" {
		asset.Currency = req.Currency
	}
	if req.PurchasedAt != nil {
		asset.PurchasedAt = *req.PurchasedAt
	}
	if req.Note != "" {
		asset.Note = req.Note
	}
	if err := s.assetRepo.updateOne(ctx, asset); err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *AssetService) DeleteAsset(ctx context.Context, storeID, assetID string) error {
	if err := s.assetRepo.deleteOne(ctx, s.assetRepo.scopeID(assetID), s.assetRepo.scopeStoreID(storeID)); err != nil {
		return err
	}
	return nil
}
