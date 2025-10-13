package inventory

import (
	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type InventoryDomain struct {
	InventoryService *InventoryService
}

func NewDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, cache *db.Memcache, storeDomain *store.StoreDomain) *InventoryDomain {
	productRepo := newProductRepository(postgres)
	variantRepo := newVariantRepository(postgres)
	postgres.AutoMigrate(&Product{}, &Variant{})
	return &InventoryDomain{
		InventoryService: NewInventoryService(productRepo, variantRepo, storeDomain.StoreService, atomicProcess),
	}
}
