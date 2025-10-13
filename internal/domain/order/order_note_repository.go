package order

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"gorm.io/gorm"
)

type orderNoteRepository struct {
	db *db.Postgres
}

func newOrderNoteRepository(db *db.Postgres) *orderNoteRepository {
	return &orderNoteRepository{db: db}
}

func (r *orderNoteRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *orderNoteRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *orderNoteRepository) scopeOrderID(orderID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("order_id = ?", orderID)
	}
}

func (r *orderNoteRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*OrderNote, error) {
	var orderNote OrderNote
	if err := r.db.Conn(ctx, opts...).First(&orderNote, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &orderNote, nil
}

func (r *orderNoteRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*OrderNote, error) {
	var orderNote OrderNote
	if err := r.db.Conn(ctx, opts...).First(&orderNote).Error; err != nil {
		return nil, err
	}
	return &orderNote, nil
}

func (r *orderNoteRepository) createOne(ctx context.Context, orderNote *OrderNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(orderNote).Error
}

func (r *orderNoteRepository) createMany(ctx context.Context, orderNotes []*OrderNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(&orderNotes).Error
}

func (r *orderNoteRepository) updateOne(ctx context.Context, orderNote *OrderNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(orderNote).Error
}

func (r *orderNoteRepository) updateMany(ctx context.Context, orderNotes []*OrderNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&orderNotes).Error
}

func (r *orderNoteRepository) deleteOne(ctx context.Context, orderNote *OrderNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(orderNote).Error
}

func (r *orderNoteRepository) deleteMany(ctx context.Context, orderNotes []*OrderNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&orderNotes).Error
}

func (r *orderNoteRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*OrderNote, error) {
	var orderNotes []*OrderNote
	if err := r.db.Conn(ctx, opts...).Find(&orderNotes).Error; err != nil {
		return nil, err
	}
	return orderNotes, nil
}

func (r *orderNoteRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&OrderNote{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
