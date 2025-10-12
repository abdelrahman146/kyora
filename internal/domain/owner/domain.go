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

func SetupOwnerDomain(postgres *db.Postgres, storeService *store.StoreService, atomicProcess *db.AtomicProcess) *OwnerDomain {
	ownerRepo := NewOwnerRepository(postgres)
	investmentRepo := NewInvestmentRepository(postgres)
	ownerDrawRepo := NewOwnerDrawRepository(postgres)
	postgres.AutoMigrate(&Owner{}, &Investment{}, &OwnerDraw{})

	return &OwnerDomain{
		OwnerService:      NewOwnerService(ownerRepo, storeService),
		InvestmentService: NewInvestmentService(investmentRepo, storeService),
		OwnerDrawService:  NewOwnerDrawService(ownerDrawRepo, storeService),
	}
}
