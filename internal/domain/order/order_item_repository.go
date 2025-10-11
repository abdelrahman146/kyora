package order

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"gorm.io/gorm"
)

type OrderItemRepository struct {
	db *db.Postgres
}

func NewOrderItemRepository(db *db.Postgres) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (r *OrderItemRepository) ScopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *OrderItemRepository) ScopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *OrderItemRepository) ScopeOrderID(orderID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("order_id = ?", orderID)
	}
}

func (r *OrderItemRepository) FindByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*OrderItem, error) {
	var orderItem OrderItem
	if err := r.db.Conn(ctx, opts...).First(&orderItem, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &orderItem, nil
}

func (r *OrderItemRepository) FindOne(ctx context.Context, opts ...db.PostgresOptions) (*OrderItem, error) {
	var orderItem OrderItem
	if err := r.db.Conn(ctx, opts...).First(&orderItem).Error; err != nil {
		return nil, err
	}
	return &orderItem, nil
}

func (r *OrderItemRepository) CreateOne(ctx context.Context, orderItem *OrderItem, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(orderItem).Error
}

func (r *OrderItemRepository) CreateMany(ctx context.Context, orderItems []*OrderItem, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(&orderItems).Error
}

func (r *OrderItemRepository) UpdateOne(ctx context.Context, orderItem *OrderItem, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(orderItem).Error
}

func (r *OrderItemRepository) UpdateMany(ctx context.Context, orderItems []*OrderItem, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&orderItems).Error
}

func (r *OrderItemRepository) DeleteOne(ctx context.Context, orderItem *OrderItem, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(orderItem).Error
}

func (r *OrderItemRepository) DeleteMany(ctx context.Context, orderItems []*OrderItem, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&orderItems).Error
}

func (r *OrderItemRepository) List(ctx context.Context, opts ...db.PostgresOptions) ([]*OrderItem, error) {
	var orderItems []*OrderItem
	if err := r.db.Conn(ctx, opts...).Find(&orderItems).Error; err != nil {
		return nil, err
	}
	return orderItems, nil
}

func (r *OrderItemRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&OrderItem{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
