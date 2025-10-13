package customer

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type customerRepository struct {
	db *db.Postgres
}

func newCustomerRepository(db *db.Postgres) *customerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *customerRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *customerRepository) scopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
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

func (r *customerRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *customerRepository) createOne(ctx context.Context, customer *Customer, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(customer).Error
}

func (r *customerRepository) createMany(ctx context.Context, customers []*Customer, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&customers).Error
}

func (r *customerRepository) upsertMany(ctx context.Context, customers []*Customer, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name", "gender", "email", "phone", "tiktok_username", "instagram_username", "facebook_username", "x_username", "snapchat_username", "whatsapp_number", "updated_at"}),
	}).Create(&customers).Error
}

func (r *customerRepository) updateOne(ctx context.Context, customer *Customer, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(customer).Error
}

func (r *customerRepository) updateMany(ctx context.Context, customers []*Customer, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&customers).Error
}

func (r *customerRepository) patchOne(ctx context.Context, updates *Customer, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Customer{}).Updates(updates).Error
}

func (r *customerRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Customer{}).Error
}

func (r *customerRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Customer{}).Error
}

func (r *customerRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Customer, error) {
	var customer Customer
	if err := r.db.Conn(ctx, opts...).First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Customer, error) {
	var customer Customer
	if err := r.db.Conn(ctx, opts...).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Customer, error) {
	var customers []*Customer
	if err := r.db.Conn(ctx, opts...).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func (r *customerRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Customer{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
