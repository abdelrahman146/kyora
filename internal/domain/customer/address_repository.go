package customer

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AddressRepository struct {
	db *db.Postgres
}

func NewAddressRepository(db *db.Postgres) *AddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *AddressRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *AddressRepository) scopeCustomerID(customerID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("customer_id = ?", customerID)
	}
}

func (r *AddressRepository) createOne(ctx context.Context, address *Address, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(address).Error
}

func (r *AddressRepository) createMany(ctx context.Context, addresses []*Address, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&addresses).Error
}

func (r *AddressRepository) upsertMany(ctx context.Context, addresses []*Address, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"street", "city", "state", "country_code", "phone", "zip_code", "updated_at"}),
	}).Create(&addresses).Error
}

func (r *AddressRepository) updateOne(ctx context.Context, address *Address, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(address).Error
}

func (r *AddressRepository) updateMany(ctx context.Context, addresses []*Address, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&addresses).Error
}

func (r *AddressRepository) patchOne(ctx context.Context, updates *Address, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&Address{}).Updates(updates).Error
}

func (r *AddressRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Address{}).Error
}

func (r *AddressRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&Address{}).Error
}

func (r *AddressRepository) FindByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*Address, error) {
	var address Address
	if err := r.db.Conn(ctx, opts...).First(&address, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *AddressRepository) FindOne(ctx context.Context, opts ...db.PostgresOptions) (*Address, error) {
	var address Address
	if err := r.db.Conn(ctx, opts...).First(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *AddressRepository) List(ctx context.Context, opts ...db.PostgresOptions) ([]*Address, error) {
	var addresses []*Address
	if err := r.db.Conn(ctx, opts...).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *AddressRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&Address{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
