package order

import (
	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type OrderDomain struct {
	OrderService *OrderService
}

func SetupOrderDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, storeService *store.StoreService) *OrderDomain {
	orderItemRepo := NewOrderItemRepository(postgres)
	orderRepo := NewOrderRepository(postgres)
	OrderNoteRepo := NewOrderNoteRepository(postgres)
	postgres.AutoMigrate(&Order{}, &OrderItem{}, &OrderNote{})
	return &OrderDomain{
		OrderService: NewOrderService(orderRepo, orderItemRepo, OrderNoteRepo, storeService, atomicProcess),
	}
}
