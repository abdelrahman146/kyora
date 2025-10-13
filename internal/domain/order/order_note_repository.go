package order

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"gorm.io/gorm"
)

type OrderNoteRepository struct {
	db *db.Postgres
}

func NewOrderNoteRepository(db *db.Postgres) *OrderNoteRepository {
	return &OrderNoteRepository{db: db}
}

func (r *OrderNoteRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *OrderNoteRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *OrderNoteRepository) scopeOrderID(orderID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("order_id = ?", orderID)
	}
}

func (r *OrderNoteRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*OrderNote, error) {
	var orderNote OrderNote
	if err := r.db.Conn(ctx, opts...).First(&orderNote, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &orderNote, nil
}

func (r *OrderNoteRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*OrderNote, error) {
	var orderNote OrderNote
	if err := r.db.Conn(ctx, opts...).First(&orderNote).Error; err != nil {
		return nil, err
	}
	return &orderNote, nil
}

func (r *OrderNoteRepository) createOne(ctx context.Context, orderNote *OrderNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(orderNote).Error
}

func (r *OrderNoteRepository) createMany(ctx context.Context, orderNotes []*OrderNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(&orderNotes).Error
}

func (r *OrderNoteRepository) updateOne(ctx context.Context, orderNote *OrderNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(orderNote).Error
}

func (r *OrderNoteRepository) updateMany(ctx context.Context, orderNotes []*OrderNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&orderNotes).Error
}

func (r *OrderNoteRepository) deleteOne(ctx context.Context, orderNote *OrderNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(orderNote).Error
}

func (r *OrderNoteRepository) deleteMany(ctx context.Context, orderNotes []*OrderNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&orderNotes).Error
}

func (r *OrderNoteRepository) List(ctx context.Context, opts ...db.PostgresOptions) ([]*OrderNote, error) {
	var orderNotes []*OrderNote
	if err := r.db.Conn(ctx, opts...).Find(&orderNotes).Error; err != nil {
		return nil, err
	}
	return orderNotes, nil
}

func (r *OrderNoteRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&OrderNote{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
