package accounting

import (
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
)

type Storage struct {
	cache            *cache.Cache
	investment       *database.Repository[Investment]
	withdrawal       *database.Repository[Withdrawal]
	asset            *database.Repository[Asset]
	expense          *database.Repository[Expense]
	recurringExpense *database.Repository[RecurringExpense]
}

func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
	return &Storage{
		cache:            cache,
		investment:       database.NewRepository[Investment](db),
		withdrawal:       database.NewRepository[Withdrawal](db),
		asset:            database.NewRepository[Asset](db),
		expense:          database.NewRepository[Expense](db),
		recurringExpense: database.NewRepository[RecurringExpense](db),
	}
}
