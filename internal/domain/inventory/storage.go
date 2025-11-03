package inventory

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Storage struct {
	cache      *cache.Cache
	products   *database.Repository[Product]
	variants   *database.Repository[Variant]
	categories *database.Repository[Category]
}

func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
	return &Storage{
		cache:      cache,
		products:   database.NewRepository[Product](db),
		variants:   database.NewRepository[Variant](db),
		categories: database.NewRepository[Category](db),
	}
}

func (s *Storage) ScopeSearchTermByName(term string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if term == "" {
			return db
		}
		return db.
			Where("name % ?", term).
			Order(clause.Expr{
				SQL:  "similarity(name, ?) DESC",
				Vars: []any{term},
			})
	}
}

func (s *Storage) ScopeLowStockVariants() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s <= %s", VariantSchema.StockQuantity, VariantSchema.StockQuantityAlert))
	}
}
