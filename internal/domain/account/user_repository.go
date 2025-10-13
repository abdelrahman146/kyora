package account

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *db.Postgres
}

func newUserRepository(db *db.Postgres) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *userRepository) scopeEmail(email string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	}
}

func (r *userRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *userRepository) scopeOrganizationID(orgID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("organization_id = ?", orgID)
	}
}

func (r *userRepository) scopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
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
func (r *userRepository) scopeUpdatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !from.IsZero() && !to.IsZero() {
			return db.Where("updated_at BETWEEN ? AND ?", from, to)
		} else if !from.IsZero() {
			return db.Where("updated_at >= ?", from)
		} else if !to.IsZero() {
			return db.Where("updated_at <= ?", to)
		}
		return db
	}
}

func (r *userRepository) createOne(ctx context.Context, user *User, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(user).Error
}

func (r *userRepository) createMany(ctx context.Context, users []*User, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&users).Error
}

func (r *userRepository) upsertMany(ctx context.Context, users []*User, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&users).Error
}

func (r *userRepository) updateOne(ctx context.Context, user *User, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(user).Error
}

func (r *userRepository) updateMany(ctx context.Context, users []*User, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&users).Error
}

func (r *userRepository) patchOne(ctx context.Context, updates *User, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&User{}).Updates(updates).Error
}

func (r *userRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&User{}).Error
}

func (r *userRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&User{}).Error
}

func (r *userRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*User, error) {
	var user User
	if err := r.db.Conn(ctx, opts...).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*User, error) {
	var user User
	if err := r.db.Conn(ctx, opts...).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*User, error) {
	var users []*User
	if err := r.db.Conn(ctx, opts...).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
