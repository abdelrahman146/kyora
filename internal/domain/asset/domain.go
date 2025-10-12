package asset

import (
	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type AssetDomain struct {
	AssetService *AssetService
}

func SetupAssetDomain(postgres *db.Postgres, storeService *store.StoreService, atomicProcess *db.AtomicProcess) *AssetDomain {
	assetRepo := NewAssetRepository(postgres)
	postgres.AutoMigrate(&Asset{})

	return &AssetDomain{
		AssetService: NewAssetService(assetRepo, storeService),
	}
}
