package customer

import (
	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type CustomerDomain struct {
	CustomerService *CustomerService
}

func NewDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, cache *db.Memcache, storeDomain *store.StoreDomain) *CustomerDomain {
	addressRepo := newAddressRepository(postgres)
	customerRepo := newCustomerRepository(postgres)
	customerNotesRepo := newCustomerNoteRepository(postgres)
	postgres.AutoMigrate(&Address{}, &Customer{}, &CustomerNote{})
	return &CustomerDomain{
		CustomerService: NewCustomerService(customerRepo, addressRepo, customerNotesRepo, atomicProcess),
	}
}
