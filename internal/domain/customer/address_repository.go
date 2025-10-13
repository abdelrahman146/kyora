package customer

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type addressRepository struct {
	db *db.Postgres
}

func newAddressRepository(db *db.Postgres) *addressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *addressRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *addressRepository) scopeCustomerID(customerID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("customer_id = ?", customerID)
	}
}

func (r *addressRepository) createOne(ctx context.Context, address *Address, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(address).Error
}

func (r *addressRepository) createMany(ctx context.Context, addresses []*Address, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&addresses).Error
}

func (r *addressRepository) upsertMany(ctx context.Context, addresses []*Address, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"street", "city", "state", "country_code", "phone", "zip_code", "updated_at"}),
	}).Create(&addresses).Error
}

func (r *addressRepository) updateOne(ctx context.Context, address *Address, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(address).Error
}

func (r *addressRepository) updateMany(ctx context.Context, addresses []*Address, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&addresses).Error
}

func (r *addressRepository) patchOne(ctx context.Context, updates *Address, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Address{}).Updates(updates).Error
}

func (r *addressRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Address{}).Error
}

func (r *addressRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Address{}).Error
}

func (r *addressRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Address, error) {
	var address Address
	if err := r.db.Conn(ctx, opts...).First(&address, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*Address, error) {
	var address Address
	if err := r.db.Conn(ctx, opts...).First(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*Address, error) {
	var addresses []*Address
	if err := r.db.Conn(ctx, opts...).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Address{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
