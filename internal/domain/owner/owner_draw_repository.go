package owner

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ownerDrawRepository struct {
	db *db.Postgres
}

func newOwnerDrawRepository(db *db.Postgres) *ownerDrawRepository {
	return &ownerDrawRepository{db: db}
}

func (r *ownerDrawRepository) scopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *ownerDrawRepository) scopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *ownerDrawRepository) scopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *ownerDrawRepository) scopeOwnerID(ownerID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("owner_id = ?", ownerID)
	}
}

func (r *ownerDrawRepository) scopeCreatedAt(from, to time.Time) func(db *gorm.DB) *gorm.DB {
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

func (r *ownerDrawRepository) createOne(ctx context.Context, draw *OwnerDraw, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(draw).Error
}

func (r *ownerDrawRepository) createMany(ctx context.Context, draws []*OwnerDraw, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&draws).Error
}

func (r *ownerDrawRepository) upsertMany(ctx context.Context, draws []*OwnerDraw, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name", "updated_at"}),
	}).Create(&draws).Error
}

func (r *ownerDrawRepository) updateOne(ctx context.Context, draw *OwnerDraw, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(draw).Error
}

func (r *ownerDrawRepository) updateMany(ctx context.Context, draws []*OwnerDraw, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&draws).Error
}

func (r *ownerDrawRepository) patchOne(ctx context.Context, updates *OwnerDraw, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&OwnerDraw{}).Updates(updates).Error
}

func (r *ownerDrawRepository) deleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&OwnerDraw{}).Error
}

func (r *ownerDrawRepository) deleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&OwnerDraw{}).Error
}

func (r *ownerDrawRepository) findByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*OwnerDraw, error) {
	var draw OwnerDraw
	if err := r.db.Conn(ctx, opts...).First(&draw, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &draw, nil
}

func (r *ownerDrawRepository) findOne(ctx context.Context, opts ...db.PostgresOptions) (*OwnerDraw, error) {
	var draw OwnerDraw
	if err := r.db.Conn(ctx, opts...).First(&draw).Error; err != nil {
		return nil, err
	}
	return &draw, nil
}

func (r *ownerDrawRepository) list(ctx context.Context, opts ...db.PostgresOptions) ([]*OwnerDraw, error) {
	var draws []*OwnerDraw
	if err := r.db.Conn(ctx, opts...).Find(&draws).Error; err != nil {
		return nil, err
	}
	return draws, nil
}

func (r *ownerDrawRepository) count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&OwnerDraw{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ownerDrawRepository) sumAmount(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&OwnerDraw{}).Select("COALESCE(SUM(amount), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}
