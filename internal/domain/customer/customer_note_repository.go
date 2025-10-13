package customer

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type customerNoteRepository struct {
	db *db.Postgres
}

func newCustomerNoteRepository(db *db.Postgres) *customerNoteRepository {
	return &customerNoteRepository{db: db}
}

func (r *customerNoteRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *customerNoteRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *customerNoteRepository) scopeCustomerID(customerID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("customer_id = ?", customerID)
	}
}

func (r *customerNoteRepository) createOne(ctx context.Context, customerNote *CustomerNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(customerNote).Error
}

func (r *customerNoteRepository) createMany(ctx context.Context, customerNotes []*CustomerNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&customerNotes).Error
}

func (r *customerNoteRepository) upsertMany(ctx context.Context, customerNotes []*CustomerNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"note", "updated_at"}),
	}).Create(&customerNotes).Error
}

func (r *customerNoteRepository) updateOne(ctx context.Context, customerNote *CustomerNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(customerNote).Error
}

func (r *customerNoteRepository) updateMany(ctx context.Context, customerNotes []*CustomerNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&customerNotes).Error
}

func (r *customerNoteRepository) patchOne(ctx context.Context, updates *CustomerNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&CustomerNote{}).Updates(updates).Error
}

func (r *customerNoteRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&CustomerNote{}).Error
}

func (r *customerNoteRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&CustomerNote{}).Error
}

func (r *customerNoteRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*CustomerNote, error) {
	var customerNote CustomerNote
	if err := r.db.Conn(ctx, opts...).First(&customerNote, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customerNote, nil
}

func (r *customerNoteRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*CustomerNote, error) {
	var customerNote CustomerNote
	if err := r.db.Conn(ctx, opts...).First(&customerNote).Error; err != nil {
		return nil, err
	}
	return &customerNote, nil
}

func (r *customerNoteRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*CustomerNote, error) {
	var customerNotes []*CustomerNote
	if err := r.db.Conn(ctx, opts...).Find(&customerNotes).Error; err != nil {
		return nil, err
	}
	return customerNotes, nil
}

func (r *customerNoteRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&CustomerNote{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
