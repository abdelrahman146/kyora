package owner

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/govalues/decimal"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OwnerDrawRepository struct {
	db *db.Postgres
}

func NewOwnerDrawRepository(db *db.Postgres) *OwnerDrawRepository {
	return &OwnerDrawRepository{db: db}
}

func (r *OwnerDrawRepository) ScopeID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (r *OwnerDrawRepository) ScopeIDs(ids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (r *OwnerDrawRepository) ScopeStoreID(storeID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", storeID)
	}
}

func (r *OwnerDrawRepository) ScopeOwnerID(ownerID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("owner_id = ?", ownerID)
	}
}

func (r *OwnerDrawRepository) CreateOne(ctx context.Context, draw *OwnerDraw, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Create(draw).Error
}

func (r *OwnerDrawRepository) CreateMany(ctx context.Context, draws []*OwnerDraw, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{DoNothing: true}).Create(&draws).Error
}

func (r *OwnerDrawRepository) UpsertMany(ctx context.Context, draws []*OwnerDraw, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name", "updated_at"}),
	}).Create(&draws).Error
}

func (r *OwnerDrawRepository) UpdateOne(ctx context.Context, draw *OwnerDraw, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(draw).Error
}

func (r *OwnerDrawRepository) UpdateMany(ctx context.Context, draws []*OwnerDraw, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Save(&draws).Error
}

func (r *OwnerDrawRepository) PatchOne(ctx context.Context, updates *OwnerDraw, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Model(&OwnerDraw{}).Updates(updates).Error
}

func (r *OwnerDrawRepository) DeleteOne(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&OwnerDraw{}).Error
}

func (r *OwnerDrawRepository) DeleteMany(ctx context.Context, opts ...db.PostgresOptions) error {
	return r.db.Conn(ctx, opts...).Delete(&OwnerDraw{}).Error
}

func (r *OwnerDrawRepository) FindByID(ctx context.Context, id string, opts ...db.PostgresOptions) (*OwnerDraw, error) {
	var draw OwnerDraw
	if err := r.db.Conn(ctx, opts...).First(&draw, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &draw, nil
}

func (r *OwnerDrawRepository) FindOne(ctx context.Context, opts ...db.PostgresOptions) (*OwnerDraw, error) {
	var draw OwnerDraw
	if err := r.db.Conn(ctx, opts...).First(&draw).Error; err != nil {
		return nil, err
	}
	return &draw, nil
}

func (r *OwnerDrawRepository) List(ctx context.Context, opts ...db.PostgresOptions) ([]*OwnerDraw, error) {
	var draws []*OwnerDraw
	if err := r.db.Conn(ctx, opts...).Find(&draws).Error; err != nil {
		return nil, err
	}
	return draws, nil
}

func (r *OwnerDrawRepository) Count(ctx context.Context, opts ...db.PostgresOptions) (int64, error) {
	var count int64
	if err := r.db.Conn(ctx, opts...).Model(&OwnerDraw{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OwnerDrawRepository) SumAmount(ctx context.Context, opts ...db.PostgresOptions) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Conn(ctx, opts...).Model(&OwnerDraw{}).Select("COALESCE(SUM(amount), 0)").Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}
