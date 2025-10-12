package inventory

import (
	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type InventoryDomain struct {
	InventoryService *InventoryService
}

func NewDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, cache *db.Memcache, storeService *store.StoreService) *InventoryDomain {
	productRepo := NewProductRepository(postgres)
	variantRepo := NewVariantRepository(postgres)
	postgres.AutoMigrate(&Product{}, &Variant{})
	return &InventoryDomain{
		InventoryService: NewInventoryService(productRepo, variantRepo, storeService, atomicProcess),
	}
}
