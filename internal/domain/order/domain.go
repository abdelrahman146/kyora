package order

import (
	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type OrderDomain struct {
	OrderService *OrderService
}

func NewDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, cache *db.Memcache, storeDomain *store.StoreDomain) *OrderDomain {
	orderItemRepo := newOrderItemRepository(postgres)
	orderRepo := newOrderRepository(postgres)
	OrderNoteRepo := newOrderNoteRepository(postgres)
	postgres.AutoMigrate(&Order{}, &OrderItem{}, &OrderNote{})
	return &OrderDomain{
		OrderService: NewOrderService(orderRepo, orderItemRepo, OrderNoteRepo, storeDomain.StoreService, atomicProcess),
	}
}
