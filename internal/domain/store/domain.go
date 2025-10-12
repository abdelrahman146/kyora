package store

import "github.com/abdelrahman146/kyora/internal/db"

type StoreDomain struct {
	StoreService *StoreService
}

func SetupStoreDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess) *StoreDomain {
	storeRepo := NewStoreRepository(postgres)
	postgres.AutoMigrate(&Store{})
	return &StoreDomain{
		StoreService: NewStoreService(storeRepo),
	}
}
