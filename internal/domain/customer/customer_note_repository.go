package customer

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CustomerNoteRepository struct {
	db *db.Postgres
}

func NewCustomerNoteRepository(db *db.Postgres) *CustomerNoteRepository {
	return &CustomerNoteRepository{db: db}
}

func (r *CustomerNoteRepository) ScopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *CustomerNoteRepository) ScopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *CustomerNoteRepository) ScopeCustomerID(customerID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("customer_id = ?", customerID)
	}
}

func (r *CustomerNoteRepository) CreateOne(ctx context.Context, customerNote *CustomerNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(customerNote).Error
}

func (r *CustomerNoteRepository) CreateMany(ctx context.Context, customerNotes []*CustomerNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&customerNotes).Error
}

func (r *CustomerNoteRepository) UpsertMany(ctx context.Context, customerNotes []*CustomerNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"note", "updated_at"}),
	}).Create(&customerNotes).Error
}

func (r *CustomerNoteRepository) UpdateOne(ctx context.Context, customerNote *CustomerNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(customerNote).Error
}

func (r *CustomerNoteRepository) UpdateMany(ctx context.Context, customerNotes []*CustomerNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&customerNotes).Error
}

func (r *CustomerNoteRepository) PatchOne(ctx context.Context, updates *CustomerNote, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&CustomerNote{}).Updates(updates).Error
}

func (r *CustomerNoteRepository) DeleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&CustomerNote{}).Error
}

func (r *CustomerNoteRepository) DeleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&CustomerNote{}).Error
}

func (r *CustomerNoteRepository) FindByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*CustomerNote, error) {
	var customerNote CustomerNote
	if err := r.db.Conn(ctx, opts...).First(&customerNote, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customerNote, nil
}

func (r *CustomerNoteRepository) FindOne(ctx context.Context, opts ...db.PostgresOptions) (*CustomerNote, error) {
	var customerNote CustomerNote
	if err := r.db.Conn(ctx, opts...).First(&customerNote).Error; err != nil {
		return nil, err
	}
	return &customerNote, nil
}

func (r *CustomerNoteRepository) List(ctx context.Context, opts ...db.PostgresOptions) ([]*CustomerNote, error) {
	var customerNotes []*CustomerNote
	if err := r.db.Conn(ctx, opts...).Find(&customerNotes).Error; err != nil {
		return nil, err
	}
	return customerNotes, nil
}

func (r *CustomerNoteRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&CustomerNote{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
