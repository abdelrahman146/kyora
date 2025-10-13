package order

import (
	"context"
	"fmt"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/govalues/decimal"
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
			return db.Where("created_at BETWEEN ? AND ?", from, to)
		} else if !from.IsZero() {
			return db.Where("created_at >= ?", from)
		} else if !to.IsZero() {
			return db.Where("created_at <= ?", to)
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
