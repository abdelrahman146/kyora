package customer

import (
	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type CustomerDomain struct {
	CustomerService *CustomerService
}

func SetupCustomerDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, storeService *store.StoreService) *CustomerDomain {
	addressRepo := NewAddressRepository(postgres)
	customerRepo := NewCustomerRepository(postgres)
	customerNotesRepo := NewCustomerNoteRepository(postgres)
	postgres.AutoMigrate(&Address{}, &Customer{}, &CustomerNote{})
	return &CustomerDomain{
		CustomerService: NewCustomerService(customerRepo, addressRepo, customerNotesRepo, atomicProcess),
	}
}
