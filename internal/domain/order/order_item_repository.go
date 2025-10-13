package order

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"gorm.io/gorm"
)

type orderItemRepository struct {
	db *db.Postgres
}

func newOrderItemRepository(db *db.Postgres) *orderItemRepository {
	return &orderItemRepository{db: db}
}

func (r *orderItemRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *orderItemRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *orderItemRepository) scopeOrderID(orderID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("order_id = ?", orderID)
	}
}

func (r *orderItemRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*OrderItem, error) {
	var orderItem OrderItem
	if err := r.db.Conn(ctx, opts...).First(&orderItem, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &orderItem, nil
}

func (r *orderItemRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*OrderItem, error) {
	var orderItem OrderItem
	if err := r.db.Conn(ctx, opts...).First(&orderItem).Error; err != nil {
		return nil, err
	}
	return &orderItem, nil
}

func (r *orderItemRepository) createOne(ctx context.Context, orderItem *OrderItem, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(orderItem).Error
}

func (r *orderItemRepository) createMany(ctx context.Context, orderItems []*OrderItem, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(&orderItems).Error
}

func (r *orderItemRepository) updateOne(ctx context.Context, orderItem *OrderItem, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(orderItem).Error
}

func (r *orderItemRepository) updateMany(ctx context.Context, orderItems []*OrderItem, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&orderItems).Error
}

func (r *orderItemRepository) deleteOne(ctx context.Context, orderItem *OrderItem, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(orderItem).Error
}

func (r *orderItemRepository) deleteMany(ctx context.Context, orderItems []*OrderItem, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&orderItems).Error
}

func (r *orderItemRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*OrderItem, error) {
	var orderItems []*OrderItem
	if err := r.db.Conn(ctx, opts...).Find(&orderItems).Error; err != nil {
		return nil, err
	}
	return orderItems, nil
}

func (r *orderItemRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&OrderItem{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
