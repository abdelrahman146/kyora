package order

import (
	"context"
	"fmt"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/keyvalue"
	"github.com/abdelrahman146/kyora/internal/platform/types/list"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/timeseries"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Service struct {
	storage         *Storage
	atomicProcessor atomic.AtomicProcessor
	bus             *bus.Bus
	inventory       *inventory.Service
	customer        *customer.Service
}

func NewService(storage *Storage, atomicProcessor atomic.AtomicProcessor, bus *bus.Bus, inventory *inventory.Service, customer *customer.Service) *Service {
	return &Service{
		storage:         storage,
		atomicProcessor: atomicProcessor,
		bus:             bus,
		inventory:       inventory,
		customer:        customer,
	}
}

func (s *Service) CreateOrder(ctx context.Context, actor *account.User, biz *business.Business, req *CreateOrderRequest) (*Order, error) {
	// Basic abuse/double-submit protection (best-effort, cache-backed).
	if err := s.throttle(fmt.Sprintf("rl:order:create:%s:%s", biz.ID, actor.ID), time.Minute, 30, 1*time.Second); err != nil {
		return nil, err
	}

	var order *Order
	err := s.atomicProcessor.Exec(ctx, func(tctx context.Context) error {
		// basic validation
		if req == nil || len(req.Items) == 0 {
			return ErrEmptyOrderItems()
		}
		// ownership validation: ensure customer + address belong to this business
		if _, err := s.customer.GetCustomerByID(tctx, actor, biz, req.CustomerID); err != nil {
			return err
		}
		if _, err := s.customer.GetCustomerAddressByID(tctx, actor, biz, req.CustomerID, req.ShippingAddressID); err != nil {
			return err
		}

		orderItems, adjustments, err := s.prepareOrderItems(tctx, actor, biz, req.Items)
		if err != nil {
			return err
		}

		// calculate totals from items and fees and vat
		vatRate := biz.VatRate
		subtotal := s.calculateSubtotal(orderItems)
		cogs := s.calculateCOGS(orderItems)
		vat := s.calculateVAT(subtotal, vatRate)
		total := s.calculateTotal(subtotal, vat, req.ShippingFee, req.Discount)

		// generate order number with retry on conflict
		var orderNumber string
		const maxRetries = 5
		for range maxRetries {
			orderNumber = id.Base62(6)
			order = &Order{
				BusinessID:        biz.ID,
				CustomerID:        req.CustomerID,
				ShippingAddressID: req.ShippingAddressID,
				Channel:           req.Channel,
				Subtotal:          subtotal,
				VAT:               vat,
				VATRate:           vatRate,
				ShippingFee:       req.ShippingFee,
				Discount:          req.Discount,
				COGS:              cogs,
				Total:             total,
				Currency:          biz.Currency,
				Status:            OrderStatusPending,
				PaymentStatus:     OrderPaymentStatusPending,
				PaymentMethod:     req.PaymentMethod,
				PaymentReference:  req.PaymentReference,
				OrderNumber:       orderNumber,
			}
			if !req.OrderedAt.IsZero() {
				order.OrderedAt = req.OrderedAt
			} else {
				order.OrderedAt = time.Now()
			}
			// attempt create
			if err := s.storage.order.CreateOne(tctx, order); err != nil {
				if database.IsUniqueViolation(err) {
					// retry with a new order number
					continue
				}
				return err
			}
			// created successfully
			break
		}
		if order == nil || order.ID == "" {
			return ErrOrderNumberGenerationFailed(nil)
		}

		// create order items records
		for _, oi := range orderItems {
			oi.OrderID = order.ID
		}
		if err := s.storage.orderItem.CreateMany(tctx, orderItems); err != nil {
			return err
		}

		// adjust inventory levels
		if err := s.adjustInventoryLevels(tctx, actor, biz, adjustments); err != nil {
			return err
		}

		// attach items to order for return
		order.Items = orderItems
		return nil
	}, atomic.WithIsolationLevel(atomic.LevelSerializable), atomic.WithRetries(3))
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) UpdateOrder(ctx context.Context, actor *account.User, biz *business.Business, id string, req *UpdateOrderRequest) (*Order, error) {
	var updated *Order
	err := s.atomicProcessor.Exec(ctx, func(tctx context.Context) error {
		// load order with items scoped by business
		ord, err := s.storage.order.FindByID(tctx, id,
			s.storage.order.ScopeBusinessID(biz.ID),
			s.storage.order.WithPreload(OrderItemStruct),
			s.storage.order.WithLockingStrength(database.LockingStrengthUpdate),
		)
		if err != nil {
			return ErrOrderNotFound(id, err)
		}

		// Apply simple field updates
		if req.Channel != "" {
			ord.Channel = req.Channel
		}
		if req.ShippingFee.Valid {
			ord.ShippingFee = transformer.FromNullDecimal(req.ShippingFee)
		}
		if req.Discount.Valid {
			ord.Discount = transformer.FromNullDecimal(req.Discount)
		}
		if !req.OrderedAt.IsZero() {
			ord.OrderedAt = req.OrderedAt
		}

		// If items are provided, ensure status allows modification
		if req.Items != nil {
			switch ord.Status {
			case OrderStatusShipped, OrderStatusFulfilled, OrderStatusCancelled, OrderStatusReturned:
				return ErrOrderItemsUpdateNotAllowed(ord.ID, ord.Status)
			}
			if len(req.Items) == 0 {
				return ErrEmptyOrderItems()
			}

			// delete existing items and restock inventory to prepare for new ones
			if err := s.deleteOrderItems(tctx, actor, biz, ord.ID); err != nil {
				return err
			}
			// create new items
			orderItems, adjustments, err := s.prepareOrderItems(tctx, actor, biz, req.Items)
			if err != nil {
				return err
			}
			for _, oi := range orderItems {
				oi.OrderID = ord.ID
			}
			if err := s.storage.orderItem.CreateMany(tctx, orderItems); err != nil {
				return err
			}
			ord.Items = orderItems
			if err := s.adjustInventoryLevels(tctx, actor, biz, adjustments); err != nil {
				return err
			}
			// recalculate totals
			ord.Subtotal = s.calculateSubtotal(orderItems)
			ord.COGS = s.calculateCOGS(orderItems)
			ord.VAT = s.calculateVAT(ord.Subtotal, biz.VatRate)
			ord.Total = s.calculateTotal(ord.Subtotal, ord.VAT, ord.ShippingFee, ord.Discount)

		}

		if err := s.storage.order.UpdateOne(tctx, ord); err != nil {
			return err
		}

		updated = ord
		return nil
	}, atomic.WithIsolationLevel(atomic.LevelSerializable), atomic.WithRetries(3))
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// AddOrderPaymentDetails updates payment method/reference without changing payment status.
// This is intentionally separate from UpdateOrderPaymentStatus to avoid accidental status changes.
func (s *Service) AddOrderPaymentDetails(ctx context.Context, actor *account.User, biz *business.Business, id string, req *AddOrderPaymentDetailsRequest) (*Order, error) {
	if req == nil {
		return nil, problem.BadRequest("payment details are required")
	}
	var updated *Order
	err := s.atomicProcessor.Exec(ctx, func(tctx context.Context) error {
		ord, err := s.storage.order.FindByID(tctx, id,
			s.storage.order.ScopeBusinessID(biz.ID),
			s.storage.order.WithLockingStrength(database.LockingStrengthUpdate),
		)
		if err != nil {
			return ErrOrderNotFound(id, err)
		}
		// Prevent changing payment details for finalized states
		switch ord.Status {
		case OrderStatusCancelled, OrderStatusReturned:
			return ErrOrderPaymentStatusUpdateNotAllowedForOrderStatus(ord.ID, ord.Status, ord.PaymentStatus)
		}
		ord.PaymentMethod = req.PaymentMethod
		ord.PaymentReference = req.PaymentReference
		if err := s.storage.order.UpdateOne(tctx, ord); err != nil {
			return err
		}
		updated = ord
		return nil
	}, atomic.WithIsolationLevel(atomic.LevelSerializable), atomic.WithRetries(3))
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// itemVariant represents a variant and the quantity to adjust
type itemVariant struct {
	variant *inventory.Variant
	qty     int
}

// adjustInventoryLevels decreases stock for each variant in adjustments.
// It guards against negative stock and persists the new quantity.
func (s *Service) adjustInventoryLevels(ctx context.Context, actor *account.User, biz *business.Business, adjustments []itemVariant) error {
	return s.applyInventoryAdjustments(ctx, actor, biz, adjustments, -1)
}

// restockInventoryLevels increases stock for each variant in adjustments.
// Use this when order items are removed/cancelled and stock must be returned.
func (s *Service) restockInventoryLevels(ctx context.Context, actor *account.User, biz *business.Business, adjustments []itemVariant) error {
	return s.applyInventoryAdjustments(ctx, actor, biz, adjustments, +1)
}

// applyInventoryAdjustments applies a signed delta to stock quantity:
// sign -1 to decrement (allocate), +1 to increment (restock).
func (s *Service) applyInventoryAdjustments(ctx context.Context, actor *account.User, biz *business.Business, adjustments []itemVariant, sign int) error {
	for _, adj := range adjustments {
		delta := adj.qty * sign
		newStock := adj.variant.StockQuantity + delta
		if newStock < 0 {
			return ErrInsufficientStock(adj.variant, adj.qty)
		}
		adj.variant.StockQuantity = newStock
		stock := newStock
		if err := s.inventory.UpdateVariant(ctx, actor, biz, adj.variant.ID, &inventory.UpdateVariantRequest{
			StockQuantity: &stock,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) calculateSubtotal(items []*OrderItem) decimal.Decimal {
	subtotal := decimal.Zero
	for _, it := range items {
		subtotal = subtotal.Add(it.Total)
	}
	return subtotal
}

func (s *Service) calculateCOGS(items []*OrderItem) decimal.Decimal {
	cogs := decimal.Zero
	for _, it := range items {
		cogs = cogs.Add(it.TotalCost)
	}
	return cogs
}

func (s *Service) calculateTotal(subtotal, vat, shippingFee, discount decimal.Decimal) decimal.Decimal {
	return subtotal.Add(vat).Add(shippingFee).Sub(discount)
}

func (s *Service) calculateVAT(subtotal, vatRate decimal.Decimal) decimal.Decimal {
	return subtotal.Mul(vatRate)
}

func (s *Service) deleteOrderItems(ctx context.Context, actor *account.User, biz *business.Business, orderID string) error {
	return s.atomicProcessor.Exec(ctx, func(tctx context.Context) error {
		orderItems, err := s.storage.orderItem.FindMany(tctx, s.storage.orderItem.ScopeEquals(OrderItemSchema.OrderID, orderID), s.storage.orderItem.WithPreload(inventory.VariantStruct))
		if err != nil {
			return err
		}
		adjustments := make([]itemVariant, 0, len(orderItems))
		for _, oi := range orderItems {
			adjustments = append(adjustments, itemVariant{
				variant: oi.Variant,
				qty:     oi.Quantity,
			})
		}
		if err := s.storage.orderItem.DeleteMany(tctx,
			s.storage.orderItem.ScopeEquals(OrderItemSchema.OrderID, orderID),
		); err != nil {
			return err
		}
		return s.restockInventoryLevels(tctx, actor, biz, adjustments)
	})
}

func (s *Service) prepareOrderItems(ctx context.Context, actor *account.User, biz *business.Business, reqItems []*CreateOrderItemRequest) ([]*OrderItem, []itemVariant, error) {
	orderItems := make([]*OrderItem, 0, len(reqItems))
	adjustments := make([]itemVariant, 0, len(reqItems))

	for _, reqItem := range reqItems {
		variant, err := s.inventory.GetVariantByID(ctx, actor, biz, reqItem.VariantID)
		if err != nil {
			return nil, nil, ErrVariantNotFound(reqItem.VariantID, err)
		}
		if reqItem.Quantity <= 0 {
			return nil, nil, ErrInvalidOrderItemQuantity(reqItem.VariantID, reqItem.Quantity)
		}
		// Create order item
		orderItem := &OrderItem{
			VariantID: reqItem.VariantID,
			ProductID: variant.ProductID,
			Currency:  biz.Currency,
			Quantity:  reqItem.Quantity,
			UnitPrice: reqItem.UnitPrice,
			UnitCost:  reqItem.UnitCost,
			Total:     reqItem.UnitPrice.Mul(decimal.NewFromInt(int64(reqItem.Quantity))),
			TotalCost: reqItem.UnitCost.Mul(decimal.NewFromInt(int64(reqItem.Quantity))),
		}
		orderItems = append(orderItems, orderItem)

		// Prepare inventory adjustment
		adjustments = append(adjustments, itemVariant{
			variant: variant,
			qty:     reqItem.Quantity,
		})
	}

	return orderItems, adjustments, nil
}

func (s *Service) UpdateOrderStatus(ctx context.Context, actor *account.User, biz *business.Business, id string, status OrderStatus) (*Order, error) {
	order, err := s.storage.order.FindByID(ctx, id, s.storage.order.ScopeBusinessID(biz.ID))
	if err != nil {
		return nil, ErrOrderNotFound(id, err)
	}
	sm := newOrderStateMachine(order)

	if err := sm.transitionStateTo(status); err != nil {
		return nil, err
	}
	if err := s.storage.order.UpdateOne(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) UpdateOrderPaymentStatus(ctx context.Context, actor *account.User, biz *business.Business, id string, paymentStatus OrderPaymentStatus) (*Order, error) {
	order, err := s.storage.order.FindByID(ctx, id, s.storage.order.ScopeBusinessID(biz.ID))
	if err != nil {
		return nil, ErrOrderNotFound(id, err)
	}
	sm := newOrderStateMachine(order)

	if err := sm.transitionPaymentStatusTo(paymentStatus); err != nil {
		return nil, err
	}
	if err := s.storage.order.UpdateOne(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) GetOrderByID(ctx context.Context, actor *account.User, biz *business.Business, id string) (*Order, error) {
	return s.storage.order.FindByID(ctx, id,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.WithPreload(OrderItemStruct),
		s.storage.order.WithPreload(OrderNoteStruct),
	)
}

func (s *Service) GetOrderByOrderNumber(ctx context.Context, actor *account.User, biz *business.Business, orderNumber string) (*Order, error) {
	return s.storage.order.FindOne(ctx,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeEquals(OrderSchema.OrderNumber, orderNumber),
		s.storage.order.WithPreload(OrderItemStruct),
		s.storage.order.WithPreload(OrderNoteStruct),
	)
}

type ListOrdersFilters struct {
	Statuses        []OrderStatus
	PaymentStatuses []OrderPaymentStatus
	CustomerID      string
	OrderNumber     string
	From            time.Time
	To              time.Time
}

func (s *Service) ListOrders(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest, filters *ListOrdersFilters) ([]*Order, int64, error) {
	opts := []func(db *gorm.DB) *gorm.DB{
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.WithPagination(req.Offset(), req.Limit()),
		s.storage.order.WithOrderBy(req.ParsedOrderBy(OrderSchema)),
	}
	if req.SearchTerm() != "" {
		opts = append(opts, s.storage.order.ScopeSearchTerm(req.SearchTerm(), OrderSchema.OrderNumber, OrderSchema.Channel))
	}
	if filters != nil {
		if len(filters.Statuses) > 0 {
			vals := make([]any, 0, len(filters.Statuses))
			for _, st := range filters.Statuses {
				vals = append(vals, st)
			}
			opts = append(opts, s.storage.order.ScopeIn(OrderSchema.Status, vals))
		}
		if len(filters.PaymentStatuses) > 0 {
			vals := make([]any, 0, len(filters.PaymentStatuses))
			for _, st := range filters.PaymentStatuses {
				vals = append(vals, st)
			}
			opts = append(opts, s.storage.order.ScopeIn(OrderSchema.PaymentStatus, vals))
		}
		if filters.CustomerID != "" {
			opts = append(opts, s.storage.order.ScopeEquals(OrderSchema.CustomerID, filters.CustomerID))
		}
		if filters.OrderNumber != "" {
			opts = append(opts, s.storage.order.ScopeEquals(OrderSchema.OrderNumber, filters.OrderNumber))
		}
		if !filters.From.IsZero() || !filters.To.IsZero() {
			opts = append(opts, s.storage.order.ScopeTime(OrderSchema.OrderedAt, filters.From, filters.To))
		}
	}

	items, err := s.storage.order.FindMany(ctx, opts...)
	if err != nil {
		return nil, 0, err
	}
	// Count without pagination/order
	countOpts := []func(db *gorm.DB) *gorm.DB{
		s.storage.order.ScopeBusinessID(biz.ID),
	}
	if req.SearchTerm() != "" {
		countOpts = append(countOpts, s.storage.order.ScopeSearchTerm(req.SearchTerm(), OrderSchema.OrderNumber, OrderSchema.Channel))
	}
	if filters != nil {
		if len(filters.Statuses) > 0 {
			vals := make([]any, 0, len(filters.Statuses))
			for _, st := range filters.Statuses {
				vals = append(vals, st)
			}
			countOpts = append(countOpts, s.storage.order.ScopeIn(OrderSchema.Status, vals))
		}
		if len(filters.PaymentStatuses) > 0 {
			vals := make([]any, 0, len(filters.PaymentStatuses))
			for _, st := range filters.PaymentStatuses {
				vals = append(vals, st)
			}
			countOpts = append(countOpts, s.storage.order.ScopeIn(OrderSchema.PaymentStatus, vals))
		}
		if filters.CustomerID != "" {
			countOpts = append(countOpts, s.storage.order.ScopeEquals(OrderSchema.CustomerID, filters.CustomerID))
		}
		if filters.OrderNumber != "" {
			countOpts = append(countOpts, s.storage.order.ScopeEquals(OrderSchema.OrderNumber, filters.OrderNumber))
		}
		if !filters.From.IsZero() || !filters.To.IsZero() {
			countOpts = append(countOpts, s.storage.order.ScopeTime(OrderSchema.OrderedAt, filters.From, filters.To))
		}
	}
	count, err := s.storage.order.Count(ctx, countOpts...)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}

func (s *Service) CountOrders(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	return s.storage.order.Count(ctx,
		s.storage.order.ScopeBusinessID(biz.ID),
	)
}

// CountOrdersByDateRange returns the number of orders in the provided date range (by OrderedAt)
func (s *Service) CountOrdersByDateRange(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (int64, error) {
	return s.storage.order.Count(ctx,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to),
	)
}

func (s *Service) DeleteOrder(ctx context.Context, actor *account.User, biz *business.Business, id string) error {
	return s.atomicProcessor.Exec(ctx, func(tctx context.Context) error {
		order, err := s.storage.order.FindByID(tctx, id,
			s.storage.order.ScopeBusinessID(biz.ID),
			s.storage.order.WithPreload(OrderItemStruct),
			s.storage.order.WithLockingStrength(database.LockingStrengthUpdate),
		)
		if err != nil {
			return ErrOrderNotFound(id, err)
		}
		// Deletion is a destructive action; restrict to safe states only.
		if order.Status != OrderStatusPending && order.Status != OrderStatusCancelled {
			return ErrOrderCannotBeDeleted(order.ID, order.Status)
		}
		// delete order items and restock inventory
		if err := s.deleteOrderItems(tctx, actor, biz, order.ID); err != nil {
			return err
		}
		// delete order
		return s.storage.order.DeleteOne(tctx, order)
	})
}

func (s *Service) CreateOrderNote(ctx context.Context, actor *account.User, biz *business.Business, orderID string, req *CreateOrderNoteRequest) (*OrderNote, error) {
	// Basic abuse protection for note creation.
	if err := s.throttle(fmt.Sprintf("rl:order:note:create:%s:%s:%s", biz.ID, actor.ID, orderID), 5*time.Minute, 60, 1*time.Second); err != nil {
		return nil, err
	}

	if _, err := s.storage.order.FindByID(ctx, orderID, s.storage.order.ScopeBusinessID(biz.ID)); err != nil {
		return nil, ErrOrderNotFound(orderID, err)
	}
	note := &OrderNote{
		OrderID: orderID,
		Content: req.Content,
	}
	if err := s.storage.orderNote.CreateOne(ctx, note); err != nil {
		return nil, err
	}
	return note, nil
}

// throttle implements a simple token bucket using cache with JSON state.
// This is best-effort and becomes a no-op when cache isn't configured.
type throttleState struct {
	Count int   `json:"count"`
	Last  int64 `json:"last"` // unix seconds
}

func (s *Service) throttle(key string, window time.Duration, max int, minInterval time.Duration) error {
	if s.storage.cache == nil {
		return nil
	}
	now := time.Now()
	// read state
	var st throttleState
	if data, err := s.storage.cache.Get(key); err == nil && len(data) > 0 {
		_ = s.storage.cache.Unmarshal(data, &st)
	}
	// enforce min interval
	if st.Last != 0 && now.Sub(time.Unix(st.Last, 0)) < minInterval {
		return ErrOrderRateLimited()
	}
	// increment and cap
	st.Count++
	st.Last = now.Unix()
	if st.Count > max {
		// write back to ensure TTL stays
		if b, err := s.storage.cache.Marshal(st); err == nil {
			_ = s.storage.cache.SetX(key, b, int32(window.Seconds()))
		}
		return ErrOrderRateLimited()
	}
	// persist with TTL window
	if b, err := s.storage.cache.Marshal(st); err == nil {
		_ = s.storage.cache.SetX(key, b, int32(window.Seconds()))
	}
	return nil
}

func (s *Service) UpdateOrderNote(ctx context.Context, actor *account.User, biz *business.Business, orderID string, noteID string, req *UpdateOrderNoteRequest) (*OrderNote, error) {
	if _, err := s.storage.order.FindByID(ctx, orderID, s.storage.order.ScopeBusinessID(biz.ID)); err != nil {
		return nil, ErrOrderNotFound(orderID, err)
	}
	note, err := s.storage.orderNote.FindByID(ctx, noteID)
	if err != nil {
		return nil, ErrOrderNoteNotFound(noteID, err)
	}
	if note.OrderID != orderID {
		return nil, ErrOrderNoteNotFound(noteID, gorm.ErrRecordNotFound)
	}
	if req.Content != "" {
		note.Content = req.Content
	}
	err = s.storage.orderNote.UpdateOne(ctx, note)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (s *Service) DeleteOrderNote(ctx context.Context, actor *account.User, biz *business.Business, orderID string, noteID string) error {
	if _, err := s.storage.order.FindByID(ctx, orderID, s.storage.order.ScopeBusinessID(biz.ID)); err != nil {
		return ErrOrderNotFound(orderID, err)
	}
	note, err := s.storage.orderNote.FindByID(ctx, noteID)
	if err != nil {
		return ErrOrderNoteNotFound(noteID, err)
	}
	if note.OrderID != orderID {
		return ErrOrderNoteNotFound(noteID, gorm.ErrRecordNotFound)
	}
	return s.storage.orderNote.DeleteOne(ctx, note)
}

func (s *Service) SumOrdersTotal(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (decimal.Decimal, error) {
	return s.storage.order.Sum(ctx, OrderSchema.Total, s.storage.order.ScopeBusinessID(biz.ID), s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to))
}

func (s *Service) CountOpenOrders(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	return s.storage.order.Count(ctx,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeIn(OrderSchema.Status, []any{
			OrderStatusPending,
			OrderStatusPlaced,
			OrderStatusReadyForShipment,
		}),
	)
}

func (s *Service) AvgOrdersTotal(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (decimal.Decimal, error) {
	return s.storage.order.Avg(ctx, OrderSchema.Total, s.storage.order.ScopeBusinessID(biz.ID), s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to))
}

func (s *Service) SumOrdersCOGS(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (decimal.Decimal, error) {
	return s.storage.order.Sum(ctx, OrderSchema.COGS, s.storage.order.ScopeBusinessID(biz.ID), s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to))
}

func (s *Service) AvgOrdersCOGS(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (decimal.Decimal, error) {
	return s.storage.order.Avg(ctx, OrderSchema.COGS, s.storage.order.ScopeBusinessID(biz.ID), s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to))
}

func (s *Service) TopOrdersByTotal(ctx context.Context, actor *account.User, biz *business.Business, limit int, from, to time.Time) ([]*Order, error) {
	return s.storage.order.FindMany(ctx,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to),
		s.storage.order.WithLimit(limit),
		s.storage.order.WithOrderBy([]string{OrderSchema.Total.Column() + " DESC"}),
	)
}

func (s *Service) TopOrdersByCOGS(ctx context.Context, actor *account.User, biz *business.Business, limit int, from, to time.Time) ([]*Order, error) {
	return s.storage.order.FindMany(ctx,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to),
		s.storage.order.WithLimit(limit),
		s.storage.order.WithOrderBy([]string{OrderSchema.COGS.Column() + " DESC"}),
	)
}

func (s *Service) ComputeRevenueTimeSeries(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (*timeseries.TimeSeries, error) {
	granularity := timeseries.GetTimeGranularityByDateRange(from, to)
	return s.storage.order.TimeSeriesSum(ctx, OrderSchema.Total, OrderSchema.OrderedAt, granularity,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to),
	)
}

// ComputeOrdersCountTimeSeries returns a time series of order counts over time (bucketed by date granularity) within the given range.
func (s *Service) ComputeOrdersCountTimeSeries(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (*timeseries.TimeSeries, error) {
	granularity := timeseries.GetTimeGranularityByDateRange(from, to)
	return s.storage.order.TimeSeriesCount(ctx, OrderSchema.OrderedAt, granularity,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to),
	)
}

func (s *Service) ComputeLiveOrdersFunnel(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) ([]keyvalue.KeyValue, error) {
	return s.storage.order.CountBy(ctx, OrderSchema.Status,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to),
		s.storage.order.ScopeNotIn(OrderSchema.Status, []any{OrderStatusCancelled, OrderStatusReturned, OrderStatusFulfilled}),
	)
}

func (s *Service) ComputeTopSellingProducts(ctx context.Context, actor *account.User, biz *business.Business, limit int, from, to time.Time) ([]*inventory.Product, error) {
	joinOrders := s.storage.orderItem.WithJoins("JOIN orders ON orders.id = order_items.order_id")
	scopeOrdersBusiness := func(db *gorm.DB) *gorm.DB {
		return db.Where("orders.business_id = ?", biz.ID)
	}
	scopeOrdersTime := func(db *gorm.DB) *gorm.DB {
		if !from.IsZero() && !to.IsZero() {
			return db.Where("orders.ordered_at BETWEEN ? AND ?", from, to)
		}
		if !from.IsZero() {
			return db.Where("orders.ordered_at >= ?", from)
		}
		if !to.IsZero() {
			return db.Where("orders.ordered_at <= ?", to)
		}
		return db
	}

	res, err := s.storage.orderItem.SumBy(ctx, OrderItemSchema.ProductID, OrderItemSchema.Quantity,
		joinOrders,
		scopeOrdersBusiness,
		scopeOrdersTime,
		s.storage.orderItem.WithLimit(limit),
		s.storage.orderItem.WithOrderBy([]string{fmt.Sprintf("%s DESC", keyvalue.Schema.Value.Column())}),
		s.storage.orderItem.WithPreload(inventory.ProductStruct),
	)
	if err != nil {
		return nil, err
	}
	products := make([]*inventory.Product, 0, len(res))
	for _, kv := range res {
		if p, ok := kv.Key.(*inventory.Product); ok {
			products = append(products, p)
		}
	}
	return products, nil
}

// CountOrdersByStatus returns a breakdown of order counts by status over the given date range.
func (s *Service) CountOrdersByStatus(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) ([]keyvalue.KeyValue, error) {
	return s.storage.order.CountBy(ctx, OrderSchema.Status,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to),
	)
}

// SumOrdersTotalByChannel returns revenue grouped by sales channel for the given range.
func (s *Service) SumOrdersTotalByChannel(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) ([]keyvalue.KeyValue, error) {
	return s.storage.order.SumBy(ctx, OrderSchema.Channel, OrderSchema.Total,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to),
		s.storage.order.WithOrderBy([]string{fmt.Sprintf("%s DESC", keyvalue.Schema.Value.Column())}),
	)
}

// SumOrdersTotalByCountry returns revenue grouped by destination country using the order's shipping address for the given range.
func (s *Service) SumOrdersTotalByCountry(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) ([]keyvalue.KeyValue, error) {
	// Load orders with ShippingAddress to sum in-memory for correctness and simplicity
	orders, err := s.storage.order.FindMany(ctx,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to),
		s.storage.order.WithPreload("ShippingAddress"),
	)
	if err != nil {
		return nil, err
	}
	sums := map[string]decimal.Decimal{}
	for _, o := range orders {
		cc := ""
		if o.ShippingAddress != nil {
			cc = o.ShippingAddress.CountryCode
		}
		sums[cc] = sums[cc].Add(o.Total)
	}
	out := make([]keyvalue.KeyValue, 0, len(sums))
	for k, v := range sums {
		out = append(out, keyvalue.New(k, v))
	}
	return out, nil
}

// SumItemsSold returns the total number of items sold in the provided range, summing order item quantities.
func (s *Service) SumItemsSold(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (int64, error) {
	// Load orders with items and sum quantities to properly scope by business and order time
	orders, err := s.storage.order.FindMany(ctx,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to),
		s.storage.order.WithPreload(OrderItemStruct),
	)
	if err != nil {
		return 0, err
	}
	var total int64
	for _, o := range orders {
		for _, it := range o.Items {
			total += int64(it.Quantity)
		}
	}
	return total, nil
}

// CountOrdersByCustomer returns a breakdown of order counts grouped by CustomerID within the given date range.
func (s *Service) CountOrdersByCustomer(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) ([]keyvalue.KeyValue, error) {
	return s.storage.order.CountBy(ctx, OrderSchema.CustomerID,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to),
	)
}

func (s *Service) CountReturningCustomers(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (int64, error) {
	// Find customers with more than 1 order in the given date range
	results, err := s.storage.order.CountBy(ctx, OrderSchema.CustomerID,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to),
		s.storage.order.ScopeHavingGreaterThan(keyvalue.Schema.Value, 1),
	)
	if err != nil {
		return 0, err
	}
	return int64(len(results)), nil
}

// SumOrdersTotalByCustomer returns revenue grouped by CustomerID within the given range ordered by total DESC with an optional limit.
func (s *Service) SumOrdersTotalByCustomer(ctx context.Context, actor *account.User, biz *business.Business, limit int, from, to time.Time) ([]keyvalue.KeyValue, error) {
	return s.storage.order.SumBy(ctx, OrderSchema.CustomerID, OrderSchema.Total,
		s.storage.order.ScopeBusinessID(biz.ID),
		s.storage.order.ScopeTime(OrderSchema.OrderedAt, from, to),
		s.storage.order.WithOrderBy([]string{fmt.Sprintf("%s DESC", keyvalue.Schema.Value.Column())}),
		s.storage.order.WithLimit(limit),
	)
}
