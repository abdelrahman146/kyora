package order

import (
	"context"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/types"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/shopspring/decimal"
)

type OrderService struct {
	store         *store.StoreService
	orders        *orderRepository
	orderItems    *orderItemRepository
	orderNotes    *orderNoteRepository
	atomicProcess *db.AtomicProcess
}

func NewOrderService(orders *orderRepository, orderItems *orderItemRepository, orderNotes *orderNoteRepository, store *store.StoreService, atomicProcess *db.AtomicProcess) *OrderService {
	return &OrderService{
		store:         store,
		orders:        orders,
		orderItems:    orderItems,
		orderNotes:    orderNotes,
		atomicProcess: atomicProcess,
	}
}

func (s *OrderService) GetOrderByID(ctx context.Context, storeID string, id string) (*Order, error) {
	return s.orders.findByID(ctx, id, s.orders.scopeStoreID(storeID), db.WithPreload(OrderItemStruct), db.WithPreload(OrderNoteStruct), db.WithPreload(customer.CustomerStruct))
}

func (s *OrderService) GetOrderByOrderNumber(ctx context.Context, storeID string, orderNumber string) (*Order, error) {
	return s.orders.findOne(ctx, s.orders.scopeStoreID(storeID), s.orders.scopeOrderNumber(orderNumber), db.WithPreload(OrderItemStruct), db.WithPreload(OrderNoteStruct), db.WithPreload(customer.CustomerStruct))
}

func (s *OrderService) ListOrders(ctx context.Context, storeID string, filter *OrderFilter, page, pageSize int, orderBy string, ascending bool) ([]*Order, error) {
	return s.orders.list(ctx, s.orders.scopeStoreID(storeID), s.orders.scopeOrderFilter(filter), db.WithPagination(page, pageSize), db.WithSorting(orderBy, ascending))
}

func (s *OrderService) CountOrders(ctx context.Context, storeID string, filter *OrderFilter) (int64, error) {
	return s.orders.count(ctx, s.orders.scopeStoreID(storeID), s.orders.scopeOrderFilter(filter))
}

func (s *OrderService) calculateSubtotal(ctx context.Context, items []*CreateOrderItemRequest) decimal.Decimal {
	subtotal := decimal.Zero
	for _, item := range items {
		quantity := decimal.NewFromInt(int64(item.Quantity))
		total := item.UnitPrice.Mul(quantity)
		subtotal = subtotal.Add(total)
	}
	return subtotal
}

func (s *OrderService) calculateCOGS(ctx context.Context, items []*CreateOrderItemRequest) decimal.Decimal {
	cogs := decimal.Zero
	for _, item := range items {
		quantity := decimal.NewFromInt(int64(item.Quantity))
		totalCost := item.UnitCost.Mul(quantity)
		cogs = cogs.Add(totalCost)
	}
	return cogs
}

func (s *OrderService) calculateTotal(ctx context.Context, subtotal, vat, shippingFee, discount decimal.Decimal) decimal.Decimal {
	total := subtotal.Add(vat)
	total = total.Add(shippingFee)
	total = total.Sub(discount)
	return total
}
func (s *OrderService) calculateVat(ctx context.Context, subtotal, vatRate decimal.Decimal) decimal.Decimal {
	vat := subtotal.Mul(vatRate)
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
	cogs := s.calculateCOGS(ctx, order.Items)
	// if cogs.GreaterThan(subtotal) {
	// 	return nil, utils.Problem.BadRequest("COGS cannot be greater than subtotal")
	// }
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
		Channel:           order.Channel,
		COGS:              cogs,
		VAT:               vat,
		VATRate:           store.VatRate,
		ShippingFee:       order.ShippingFee,
		Discount:          order.Discount,
		Total:             total,
		ShippingAddressID: order.ShippingAddressID,
	}

	err = s.atomicProcess.Exec(ctx, func(txCtx context.Context) error {
		if err := s.orders.createOne(txCtx, newOrder); err != nil {
			return err
		}

		orderItems := make([]*OrderItem, len(order.Items))
		for i, item := range order.Items {
			quantity := decimal.NewFromInt(int64(item.Quantity))
			total := item.UnitPrice.Mul(quantity)
			orderItems[i] = &OrderItem{
				OrderID:   newOrder.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
				UnitCost:  item.UnitCost,
				Total:     total,
				Currency:  store.Currency,
			}
		}

		if err := s.orderItems.createMany(txCtx, orderItems); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	newOrder.Items, err = s.orderItems.list(ctx, s.orderItems.scopeOrderID(newOrder.ID))
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
	if err := s.orderNotes.createOne(ctx, note); err != nil {
		return nil, err
	}
	return note, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, storeID, orderID string, newStatus OrderStatus) (*Order, error) {
	order, err := s.GetOrderByID(ctx, storeID, orderID)
	if err != nil {
		return nil, err
	}
	sm := newOrderStateMachine(order)
	if err := sm.transitionStateTo(newStatus); err != nil {
		return nil, err
	}
	if err := s.orders.updateOne(ctx, order, s.orders.scopeID(order.ID)); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) PayOrder(ctx context.Context, storeID, orderID string, paymentDetails *AddOrderPaymentDetailsRequest) (*Order, error) {
	order, err := s.GetOrderByID(ctx, storeID, orderID)
	if err != nil {
		return nil, err
	}
	sm := newOrderStateMachine(order)
	if err := sm.transitionPaymentStatusTo(OrderPaymentStatusPaid); err != nil {
		return nil, err
	}
	order.PaymentMethod = paymentDetails.PaymentMethod
	order.PaymentReference = paymentDetails.PaymentReference
	if err := s.orders.updateOne(ctx, order, s.orders.scopeID(order.ID)); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) RefundOrder(ctx context.Context, storeID, orderID string) (*Order, error) {
	order, err := s.GetOrderByID(ctx, storeID, orderID)
	if err != nil {
		return nil, err
	}
	sm := newOrderStateMachine(order)
	if err := sm.transitionPaymentStatusTo(OrderPaymentStatusRefunded); err != nil {
		return nil, err
	}
	if err := s.orders.updateOne(ctx, order, s.orders.scopeID(order.ID)); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, storeID, orderID string) error {
	order, err := s.GetOrderByID(ctx, storeID, orderID)
	if err != nil {
		return err
	}
	if err := s.orders.deleteOne(ctx, s.orders.scopeID(order.ID)); err != nil {
		return err
	}
	return nil
}

func (s *OrderService) AggregateSales(ctx context.Context, storeID string, from, to time.Time) (total decimal.Decimal, cogs decimal.Decimal, orderCount int64, err error) {
	return s.orders.aggregateSales(ctx, s.orders.scopeStoreID(storeID), s.orders.scopeCreatedAt(from, to), s.orders.scopePaidAndActive())
}

// ItemsSold returns the sum of quantities across order items of paid, active orders.
func (s *OrderService) ItemsSold(ctx context.Context, storeID string, from, to time.Time) (int64, error) {
	return s.orders.sumItemsSold(ctx, s.orders.scopeStoreID(storeID), s.orders.scopeCreatedAt(from, to), s.orders.scopePaidAndActive())
}

// TopSellingProducts returns top-N products by quantity for paid, active orders.
func (s *OrderService) TopSellingProducts(ctx context.Context, storeID string, from, to time.Time, limit int) ([]types.KeyValue, error) {
	rows, err := s.orders.topSellingProducts(ctx, limit, s.orders.scopeStoreID(storeID), s.orders.scopeCreatedAt(from, to), s.orders.scopePaidAndActive())
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *OrderService) BreakdownByStatus(ctx context.Context, storeID string, from, to time.Time) ([]types.KeyValue, error) {
	rows, err := s.orders.breakdownByStatus(ctx, s.orders.scopeStoreID(storeID), s.orders.scopeCreatedAt(from, to))
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *OrderService) BreakdownByChannel(ctx context.Context, storeID string, from, to time.Time) ([]types.KeyValue, error) {
	rows, err := s.orders.breakdownByChannel(ctx, s.orders.scopeStoreID(storeID), s.orders.scopeCreatedAt(from, to), s.orders.scopePaidAndActive())
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *OrderService) BreakdownByCountry(ctx context.Context, storeID string, from, to time.Time) ([]types.KeyValue, error) {
	rows, err := s.orders.breakdownByCountry(ctx, s.orders.scopeStoreID(storeID), s.orders.scopeCreatedAt(from, to), s.orders.scopePaidAndActive())
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// RevenueTimeSeries returns revenue per bucket defined by PostgreSQL date_trunc bucket.
func (s *OrderService) RevenueTimeSeries(ctx context.Context, storeID string, from, to time.Time, bucket string) ([]types.TimeSeriesRow, error) {
	rows, err := s.orders.revenueTimeSeries(ctx, bucket, s.orders.scopeStoreID(storeID), s.orders.scopeCreatedAt(from, to), s.orders.scopePaidAndActive())
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// CountTimeSeries returns number of orders per bucket for paid, active orders.
func (s *OrderService) CountTimeSeries(ctx context.Context, storeID string, from, to time.Time, bucket string) ([]types.TimeSeriesRow, error) {
	rows, err := s.orders.countTimeSeries(ctx, bucket, s.orders.scopeStoreID(storeID), s.orders.scopeCreatedAt(from, to), s.orders.scopePaidAndActive())
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ---- Customer analytics wrappers ----

func (s *OrderService) DistinctPurchasingCustomers(ctx context.Context, storeID string, from, to time.Time) (int64, error) {
	return s.orders.distinctPurchasingCustomers(ctx, s.orders.scopeStoreID(storeID), s.orders.scopeCreatedAt(from, to), s.orders.scopePaidAndActive())
}

func (s *OrderService) RevenuePerCustomer(ctx context.Context, storeID string, from, to time.Time, limit int) ([]types.KeyValue, error) {
	return s.orders.revenuePerCustomer(ctx, limit, s.orders.scopeStoreID(storeID), s.orders.scopeCreatedAt(from, to), s.orders.scopePaidAndActive())
}

func (s *OrderService) ReturningCustomersCount(ctx context.Context, storeID string, from, to time.Time) (int64, error) {
	return s.orders.returningCustomersCount(ctx, s.orders.scopeStoreID(storeID), s.orders.scopeCreatedAt(from, to), s.orders.scopePaidAndActive())
}

func (s *OrderService) ReturningCustomersTimeSeries(ctx context.Context, storeID string, from, to time.Time, bucket string) ([]types.TimeSeriesRow, error) {
	return s.orders.returningCustomersTimeSeries(ctx, bucket, from, to, s.orders.scopeStoreID(storeID), s.orders.scopeCreatedAt(from, to), s.orders.scopePaidAndActive())
}

// ---- Dashboard helpers ----

// OpenOrdersCount returns the number of orders which are still in-flight
// (not fulfilled, not cancelled, not returned) regardless of payment status.
func (s *OrderService) OpenOrdersCount(ctx context.Context, storeID string) (int64, error) {
	return s.orders.count(ctx,
		s.orders.scopeStoreID(storeID),
		s.orders.scopeExcludeOrderStatuses([]OrderStatus{OrderStatusFulfilled, OrderStatusCancelled, OrderStatusReturned}),
	)
}

// OpenOrdersFunnel returns a breakdown of open orders by status for funnel visualization.
func (s *OrderService) OpenOrdersFunnel(ctx context.Context, storeID string) ([]types.KeyValue, error) {
	return s.orders.breakdownByStatus(ctx,
		s.orders.scopeStoreID(storeID),
		s.orders.scopeExcludeOrderStatuses([]OrderStatus{OrderStatusFulfilled, OrderStatusCancelled, OrderStatusReturned}),
	)
}

// AllTimeSalesAggregate returns total revenue, COGS and order count for all time
// considering paid, active orders only.
func (s *OrderService) AllTimeSalesAggregate(ctx context.Context, storeID string) (total decimal.Decimal, cogs decimal.Decimal, orderCount int64, err error) {
	return s.orders.aggregateSales(ctx, s.orders.scopeStoreID(storeID), s.orders.scopePaidAndActive())
}
