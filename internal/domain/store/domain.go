package store

import "github.com/abdelrahman146/kyora/internal/db"

type StoreDomain struct {
	StoreService *StoreService
}

func NewDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, cache *db.Memcache) *StoreDomain {
	storeRepo := newStoreRepository(postgres)
	postgres.AutoMigrate(&Store{})
	return &StoreDomain{
		StoreService: NewStoreService(storeRepo),
	}
}
