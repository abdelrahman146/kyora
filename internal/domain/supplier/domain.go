package supplier

import (
	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type SupplierDomain struct {
	SupplierService *SupplierService
}

func NewDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, cache *db.Memcache, storeDomain *store.StoreDomain) *SupplierDomain {
	supplierRepo := newSupplierRepository(postgres)
	postgres.AutoMigrate(&Supplier{})
	return &SupplierDomain{
		SupplierService: NewSupplierService(storeDomain.StoreService, supplierRepo),
	}
}
