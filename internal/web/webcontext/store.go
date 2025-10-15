package webcontext

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/web/middleware"
)

func GetStoreID(c context.Context) string {
	store := GetVal[*store.Store](c, middleware.StoreKey)
	return store.ID
}

func GetStoreName(c context.Context) string {
	store := GetVal[*store.Store](c, middleware.StoreKey)
	return store.Name
}

func GetStore(c context.Context) (*store.Store, bool) {
	store := GetVal[*store.Store](c, middleware.StoreKey)
	return store, store != nil
}

func GetOrganizationStores(c context.Context) []*store.Store {
	return GetVal[[]*store.Store](c, middleware.StoresListKey)
}
