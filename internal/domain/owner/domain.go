package owner

import (
	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type OwnerDomain struct {
	OwnerService      *OwnerService
	InvestmentService *InvestmentService
	OwnerDrawService  *OwnerDrawService
}

func NewDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, cache *db.Memcache, storeDomain *store.StoreDomain) *OwnerDomain {
	ownerRepo := newOwnerRepository(postgres)
	investmentRepo := newInvestmentRepository(postgres)
	ownerDrawRepo := newOwnerDrawRepository(postgres)
	postgres.AutoMigrate(&Owner{}, &Investment{}, &OwnerDraw{})

	return &OwnerDomain{
		OwnerService:      NewOwnerService(ownerRepo, storeDomain.StoreService),
		InvestmentService: NewInvestmentService(investmentRepo, storeDomain.StoreService),
		OwnerDrawService:  NewOwnerDrawService(ownerDrawRepo, storeDomain.StoreService),
	}
}
