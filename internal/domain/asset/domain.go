package asset

import (
	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type AssetDomain struct {
	AssetService *AssetService
}

func NewDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, cache *db.Memcache, storeService *store.StoreService) *AssetDomain {
	assetRepo := NewAssetRepository(postgres)
	postgres.AutoMigrate(&Asset{})

	return &AssetDomain{
		AssetService: NewAssetService(assetRepo, storeService),
	}
}
