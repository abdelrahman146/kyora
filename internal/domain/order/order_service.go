package order

import (
	"context"
	"strings"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/govalues/decimal"
)

type OrderService struct {
	store         *store.StoreService
	orders        *OrderRepository
	orderItems    *OrderItemRepository
	orderNotes    *OrderNoteRepository
	atomicProcess *db.AtomicProcess
}

func NewOrderService(orders *OrderRepository, orderItems *OrderItemRepository, orderNotes *OrderNoteRepository, store *store.StoreService, atomicProcess *db.AtomicProcess) *OrderService {
	return &OrderService{
		store:         store,
		orders:        orders,
		orderItems:    orderItems,
		orderNotes:    orderNotes,
		atomicProcess: atomicProcess,
	}
}

func (s *OrderService) GetOrderByID(ctx context.Context, storeID string, id string) (*Order, error) {
	return s.orders.FindByID(ctx, id, s.orders.ScopeStoreID(storeID), db.WithPreload(OrderItemStruct), db.WithPreload(OrderNoteStruct), db.WithPreload(customer.CustomerStruct))
}

func (s *OrderService) GetOrderByOrderNumber(ctx context.Context, storeID string, orderNumber string) (*Order, error) {
	return s.orders.FindOne(ctx, s.orders.ScopeStoreID(storeID), s.orders.ScopeOrderNumber(orderNumber), db.WithPreload(OrderItemStruct), db.WithPreload(OrderNoteStruct), db.WithPreload(customer.CustomerStruct))
}

func (s *OrderService) ListOrders(ctx context.Context, storeID string, filter *OrderFilter, page, pageSize int, orderBy string, ascending bool) ([]*Order, error) {
	return s.orders.List(ctx, s.orders.ScopeStoreID(storeID), s.orders.ScopeOrderFilter(filter), db.WithPagination(page, pageSize), db.WithSorting(orderBy, ascending))
}

func (s *OrderService) CountOrders(ctx context.Context, storeID string, filter *OrderFilter) (int64, error) {
	return s.orders.Count(ctx, s.orders.ScopeStoreID(storeID), s.orders.ScopeOrderFilter(filter))
}

func (s *OrderService) calculateSubtotal(ctx context.Context, items []*CreateOrderItemRequest) decimal.Decimal {
	subtotal := decimal.MustNew(0, 0)
	for _, item := range items {
		quantity, err := decimal.NewFromInt64(int64(item.Quantity), 0, 0)
		if err != nil {
			utils.Log.FromContext(ctx).Error("Failed to calculate subtotal", "error", err)
			return decimal.MustNew(0, 0)
		}
		total, err := item.UnitPrice.Mul(quantity)
		if err != nil {
			utils.Log.FromContext(ctx).Error("Failed to calculate subtotal", "error", err)
			return decimal.MustNew(0, 0)
		}
		subtotal, err = subtotal.Add(total)
		if err != nil {
			utils.Log.FromContext(ctx).Error("Failed to calculate subtotal", "error", err)
			return decimal.MustNew(0, 0)
		}
	}
	return subtotal
}

func (s *OrderService) calculateTotal(ctx context.Context, subtotal, vat, shippingFee, discount decimal.Decimal) decimal.Decimal {
	total, err := subtotal.Add(vat)
	if err != nil {
		utils.Log.FromContext(ctx).Error("Failed to calculate total", "error", err)
		return decimal.MustNew(0, 0)
	}
	total, err = total.Add(shippingFee)
	if err != nil {
		utils.Log.FromContext(ctx).Error("Failed to calculate total", "error", err)
		return decimal.MustNew(0, 0)
	}
	total, err = total.Sub(discount)
	if err != nil {
		utils.Log.FromContext(ctx).Error("Failed to calculate total", "error", err)
		return decimal.MustNew(0, 0)
	}
	return total
}
func (s *OrderService) calculateVat(ctx context.Context, subtotal, vatRate decimal.Decimal) decimal.Decimal {
	vat, err := subtotal.Mul(vatRate)
	if err != nil {
		utils.Log.FromContext(ctx).Error("Failed to calculate VAT", "error", err)
		return decimal.MustNew(0, 0)
	}
	return vat
}

func (s *OrderService) generateOrderNumber() string {
	return utils.ID.NewBase62WithPrefix(strings.ToUpper(OrderAlias), 6)
}

func (s *OrderService) CreateOrder(ctx context.Context, storeID string, order *CreateOrderRequest) (*Order, error) {
	store, err := s.store.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, err
	}
	subtotal := s.calculateSubtotal(ctx, order.Items)
	vat := s.calculateVat(ctx, subtotal, store.VatRate)
	total := s.calculateTotal(ctx, subtotal, vat, order.ShippingFee, order.Discount)
	orderNumber := s.generateOrderNumber()
	newOrder := &Order{
		StoreID:           storeID,
		CustomerID:        order.CustomerID,
		OrderNumber:       orderNumber,
		Status:            OrderStatusPending,
		PaymentStatus:     OrderPaymentStatusPending,
		PaymentMethod:     order.PaymentMethod,
		Currency:          store.Currency,
		Subtotal:          subtotal,
		VAT:               vat,
		VATRate:           store.VatRate,
		ShippingFee:       order.ShippingFee,
		Discount:          order.Discount,
		Total:             total,
		ShippingAddressID: order.ShippingAddressID,
	}

	err = s.atomicProcess.Exec(ctx, func(txCtx context.Context) error {
		if err := s.orders.CreateOne(txCtx, newOrder); err != nil {
			return err
		}

		orderItems := make([]*OrderItem, len(order.Items))
		for i, item := range order.Items {
			quantity, err := decimal.NewFromInt64(int64(item.Quantity), 0, 0)
			if err != nil {
				return err
			}
			total, err := item.UnitPrice.Mul(quantity)
			if err != nil {
				return err
			}
			orderItems[i] = &OrderItem{
				OrderID:   newOrder.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
				Total:     total,
				Currency:  store.Currency,
			}
		}

		if err := s.orderItems.CreateMany(txCtx, orderItems); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	newOrder.Items, err = s.orderItems.List(ctx, s.orderItems.ScopeOrderID(newOrder.ID))
	if err != nil {
		return nil, err
	}
	return newOrder, nil
}

func (s *OrderService) AddOrderNote(ctx context.Context, storeID, orderID string, content string) (*OrderNote, error) {
	order, err := s.GetOrderByID(ctx, storeID, orderID)
	if err != nil {
		return nil, err
	}
	note := &OrderNote{
		OrderID: order.ID,
		Note:    content,
	}
	if err := s.orderNotes.CreateOne(ctx, note); err != nil {
		return nil, err
	}
	return note, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, storeID, orderID string, newStatus OrderStatus) (*Order, error) {
	order, err := s.GetOrderByID(ctx, storeID, orderID)
	if err != nil {
		return nil, err
	}
	sm := NewOrderStateMachine(order)
	if err := sm.TransitionStateTo(newStatus); err != nil {
		return nil, err
	}
	if err := s.orders.UpdateOne(ctx, order, s.orders.ScopeID(order.ID)); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) PayOrder(ctx context.Context, storeID, orderID string, paymentDetails *AddOrderPaymentDetailsRequest) (*Order, error) {
	order, err := s.GetOrderByID(ctx, storeID, orderID)
	if err != nil {
		return nil, err
	}
	sm := NewOrderStateMachine(order)
	if err := sm.TransitionPaymentStatusTo(OrderPaymentStatusPaid); err != nil {
		return nil, err
	}
	order.PaymentMethod = paymentDetails.PaymentMethod
	order.PaymentReference = paymentDetails.PaymentReference
	if err := s.orders.UpdateOne(ctx, order, s.orders.ScopeID(order.ID)); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) RefundOrder(ctx context.Context, storeID, orderID string) (*Order, error) {
	order, err := s.GetOrderByID(ctx, storeID, orderID)
	if err != nil {
		return nil, err
	}
	sm := NewOrderStateMachine(order)
	if err := sm.TransitionPaymentStatusTo(OrderPaymentStatusRefunded); err != nil {
		return nil, err
	}
	if err := s.orders.UpdateOne(ctx, order, s.orders.ScopeID(order.ID)); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, storeID, orderID string) error {
	order, err := s.GetOrderByID(ctx, storeID, orderID)
	if err != nil {
		return err
	}
	if err := s.orders.DeleteOne(ctx, s.orders.ScopeID(order.ID)); err != nil {
		return err
	}
	return nil
}
