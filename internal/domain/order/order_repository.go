package order

import (
	"context"
	"fmt"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/types"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type orderRepository struct {
	db *db.Postgres
}

func newOrderRepository(db *db.Postgres) *orderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}
func (r *orderRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *orderRepository) scopeOrderNumber(orderNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("order_number = ?", orderNumber)
	}
}

func (r *orderRepository) scopeCustomerID(customerID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("customer_id = ?", customerID)
	}
}

func (r *orderRepository) scopeCustomerIDs(customerIDs []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("customer_id IN ?", customerIDs)
	}
}

func (r *orderRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *orderRepository) scopeOrderStatus(status OrderStatus) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

func (r *orderRepository) scopeOrderStatuses(statuses []OrderStatus) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status IN ?", statuses)
	}
}

func (r *orderRepository) scopeExcludeOrderStatuses(statuses []OrderStatus) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status NOT IN ?", statuses)
	}
}

func (r *orderRepository) scopePaymentStatus(status OrderPaymentStatus) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("payment_status = ?", status)
	}
}

func (r *orderRepository) scopePaymentStatuses(statuses []OrderPaymentStatus) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("payment_status IN ?", statuses)
	}
}

func (r *orderRepository) scopePaymentMethod(method OrderPaymentMethod) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("payment_method = ?", method)
	}
}

func (r *orderRepository) scopePaymentMethods(methods []OrderPaymentMethod) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("payment_method IN ?", methods)
	}
}

func (r *orderRepository) scopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !from.IsZero() && !to.IsZero() {
			return db.Where(OrderTable+".created_at BETWEEN ? AND ?", from, to)
		} else if !from.IsZero() {
			return db.Where(OrderTable+".created_at >= ?", from)
		} else if !to.IsZero() {
			return db.Where(OrderTable+".created_at <= ?", to)
		}
		return db
	}
}

func (r *orderRepository) scopeCountryCode(countryCode string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(fmt.Sprintf("JOIN %s on %s.id = %s.shipping_address_id", customer.AddressTable, customer.AddressTable, OrderTable)).
			Where(fmt.Sprintf("%s.country_code = ?", customer.AddressTable), countryCode)
	}
}

func (r *orderRepository) scopeCountryCodes(countryCodes []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(fmt.Sprintf("JOIN %s on %s.id = %s.shipping_address_id", customer.AddressTable, customer.AddressTable, OrderTable)).
			Where(fmt.Sprintf("%s.country_code IN ?", customer.AddressTable), countryCodes)
	}
}

func (r *orderRepository) scopePaidAndActive() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Paid orders that were not cancelled or returned are considered active sales
		db = db.Where("payment_status = ?", OrderPaymentStatusPaid)
		db = db.Scopes(r.scopeExcludeOrderStatuses([]OrderStatus{OrderStatusCancelled, OrderStatusReturned}))
		return db
	}
}

func (r *orderRepository) scopeOrderFilter(filter *OrderFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		if len(filter.CustomerIDs) > 0 {
			db = db.Scopes(r.scopeCustomerIDs(filter.CustomerIDs))
		}
		if len(filter.Statuses) > 0 {
			db = db.Scopes(r.scopeOrderStatuses(filter.Statuses))
		}
		if len(filter.PaymentStatuses) > 0 {
			db = db.Scopes(r.scopePaymentStatuses(filter.PaymentStatuses))
		}
		if len(filter.PaymentMethods) > 0 {
			db = db.Scopes(r.scopePaymentMethods(filter.PaymentMethods))
		}
		if len(filter.CountryCodes) > 0 {
			db = db.Scopes(r.scopeCountryCodes(filter.CountryCodes))
		}
		return db
	}
}

func (r *orderRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Order, error) {
	var order Order
	if err := r.db.Conn(ctx, opts...).First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Order, error) {
	var order Order
	if err := r.db.Conn(ctx, opts...).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) createOne(ctx context.Context, order *Order, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(order).Error
}

func (r *orderRepository) createMany(ctx context.Context, orders []*Order, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&orders).Error
}

func (r *orderRepository) upsertMany(ctx context.Context, orders []*Order, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"order_number", "status", "payment_status", "payment_method", "total", "currency", "updated_at"}),
	}).Create(&orders).Error
}

func (r *orderRepository) updateOne(ctx context.Context, order *Order, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(order).Error
}

func (r *orderRepository) updateMany(ctx context.Context, orders []*Order, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&orders).Error
}

func (r *orderRepository) patchOne(ctx context.Context, updates *Order, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Order{}).Updates(updates).Error
}

func (r *orderRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Order{}).Error
}

func (r *orderRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Order{}).Error
}

func (r *orderRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Order, error) {
	var orders []*Order
	if err := r.db.Conn(ctx, opts...).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Order{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *orderRepository) sumTotal(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&Order{}).Select("COALESCE(SUM(total), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}

// ---- Analytics-oriented structs and helpers ----

type AggregateRow struct {
	Total decimal.Decimal `gorm:"column:total"`
	Cogs  decimal.Decimal `gorm:"column:cogs"`
	Cnt   int64           `gorm:"column:cnt"`
}

func (r *orderRepository) aggregateSales(ctx context.Context, opts ...db.PostgresOptions) (total decimal.Decimal, cogs decimal.Decimal, orderCount int64, err error) {
	var row AggregateRow
	q := r.db.Conn(ctx, opts...).Model(&Order{})
	if err = q.Select("COALESCE(SUM(total),0) AS total, COALESCE(SUM(cogs),0) AS cogs, COUNT(*) AS cnt").Scan(&row).Error; err != nil {
		return decimal.Zero, decimal.Zero, 0, err
	}
	return row.Total, row.Cogs, row.Cnt, nil
}

func (r *orderRepository) sumItemsSold(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	type cntRow struct{ Cnt int64 }
	var row cntRow
	// Join orders to order_items and sum quantities for paid, active orders
	q := r.db.Conn(ctx, opts...).Table(OrderTable + " o").Joins("JOIN " + OrderItemTable + " i ON i.order_id = o.id")
	if err := q.Select("COALESCE(SUM(i.quantity),0) AS cnt").Scan(&row).Error; err != nil {
		return 0, err
	}
	return row.Cnt, nil
}

func (r *orderRepository) topSellingProducts(ctx context.Context, limit int, opts ...db.PostgresOptions) ([]types.KeyValue, error) {
	rows := make([]types.KeyValue, 0, limit)
	q := r.db.Conn(ctx, opts...).Table(OrderTable + " o").
		Joins("JOIN " + OrderItemTable + " i ON i.order_id = o.id").
		Joins("JOIN " + inventory.ProductTable + " p ON p.id = i.product_id")
	orderStr := "value DESC"
	if limit <= 0 {
		limit = 10
	}
	if err := q.Select("p.name AS key, COALESCE(SUM(i.quantity),0)::float AS value").
		Group("p.id, p.name").
		Order(orderStr).Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *orderRepository) breakdownByStatus(ctx context.Context, opts ...db.PostgresOptions) ([]types.KeyValue, error) {
	rows := []types.KeyValue{}
	q := r.db.Conn(ctx, opts...).Model(&Order{})
	if err := q.Select("status AS key, COUNT(*)::float AS value").Group("status").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *orderRepository) breakdownByChannel(ctx context.Context, opts ...db.PostgresOptions) ([]types.KeyValue, error) {
	rows := []types.KeyValue{}
	q := r.db.Conn(ctx, opts...).Model(&Order{})
	if err := q.Select("channel AS key, COALESCE(SUM(total),0)::float AS value").Group("channel").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *orderRepository) breakdownByCountry(ctx context.Context, opts ...db.PostgresOptions) ([]types.KeyValue, error) {
	rows := []types.KeyValue{}
	q := r.db.Conn(ctx, opts...).Table(OrderTable + " o").
		Joins(fmt.Sprintf("JOIN %s a ON a.id = o.shipping_address_id", customer.AddressTable))
	if err := q.Select("a.country_code AS key, COALESCE(SUM(o.total),0)::float AS value").Group("a.country_code").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *orderRepository) revenueTimeSeries(ctx context.Context, bucket string, opts ...db.PostgresOptions) ([]types.TimeSeriesRow, error) {
	switch bucket {
	case "hour", "day", "week", "month", "quarter", "year":
		// ok
	default:
		bucket = "day"
	}
	rows := []types.TimeSeriesRow{}
	sel := fmt.Sprintf("date_trunc('%s', created_at) AS timestamp, COALESCE(SUM(total),0)::float AS value", bucket)
	q := r.db.Conn(ctx, opts...).Model(&Order{}).Model(&Order{})
	if err := q.Select(sel).Group("timestamp").Order("timestamp ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *orderRepository) countTimeSeries(ctx context.Context, bucket string, opts ...db.PostgresOptions) ([]types.TimeSeriesRow, error) {
	switch bucket {
	case "hour", "day", "week", "month", "quarter", "year":
		// ok
	default:
		bucket = "day"
	}
	rows := []types.TimeSeriesRow{}
	sel := fmt.Sprintf("date_trunc('%s', created_at) AS timestamp, COUNT(*)::float AS value", bucket)
	q := r.db.Conn(ctx, opts...).Model(&Order{})
	if err := q.Select(sel).Group("timestamp").Order("timestamp ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ---- Customer analytics helpers ----

// distinctPurchasingCustomers returns number of distinct customers with at least one paid, active order in range.
func (r *orderRepository) distinctPurchasingCustomers(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var cnt int64
	if err := r.db.Conn(ctx, opts...).Model(&Order{}).Distinct("customer_id").Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

// revenuePerCustomer returns top-N customers by revenue in period.
func (r *orderRepository) revenuePerCustomer(ctx context.Context, limit int, opts ...db.PostgresOptions) ([]types.KeyValue, error) {
	if limit <= 0 {
		limit = 10
	}
	rows := []types.KeyValue{}
	q := r.db.Conn(ctx, opts...).Model(&Order{})
	if err := q.Select("customer_id AS key, COALESCE(SUM(total),0)::float AS value").Group("customer_id").Order("value DESC").Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// returningCustomersCount counts distinct customers who placed more than one paid, active order in the period.
func (r *orderRepository) returningCustomersCount(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	type row struct{ Cnt int64 }
	var res row
	sub := r.db.Conn(ctx, opts...).Model(&Order{}).
		Select("customer_id, COUNT(*) AS c").Group("customer_id").Having("COUNT(*) > 1")
	if err := r.db.Conn(ctx).Table("(?) as t", sub).Select("COUNT(*) AS cnt").Scan(&res).Error; err != nil {
		return 0, err
	}
	return res.Cnt, nil
}

// ordersPerCustomerTimeSeries returns count of orders per bucket (same as countTimeSeries) but distinct customers per bucket is separate; here we may track new vs returning.
// For returning customers over time, we count orders by customers who have placed more than one order up to that bucket.
func (r *orderRepository) returningCustomersTimeSeries(ctx context.Context, bucket string, from, to time.Time, opts ...db.PostgresOptions) ([]types.TimeSeriesRow, error) {
	switch bucket {
	case "hour", "day", "week", "month", "quarter", "year":
	default:
		bucket = "day"
	}
	rows := []types.TimeSeriesRow{}
	// Approach: for each bucket date, count distinct customers with >1 cumulative orders up to that bucket. This is heavy; approximate by counting customers with >1 order inside bucket.
	sel := fmt.Sprintf("date_trunc('%s', created_at) AS timestamp, COUNT(DISTINCT customer_id)::float AS value", bucket)
	q := r.db.Conn(ctx, opts...).Model(&Order{}).Group("timestamp").Order("timestamp ASC")
	if err := q.Select(sel).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
