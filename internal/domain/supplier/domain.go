package supplier

import (
	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type SupplierDomain struct {
	SupplierService *SupplierService
}

func SetupSupplierDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, storeService *store.StoreService) *SupplierDomain {
	supplierRepo := NewSupplierRepository(postgres)
	postgres.AutoMigrate(&Supplier{})
	return &SupplierDomain{
		SupplierService: NewSupplierService(storeService, supplierRepo),
	}
}
