package analytics

import (
	"github.com/abdelrahman146/kyora/internal/domain/asset"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/expense"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/domain/owner"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/domain/supplier"
)

type AnalyticsDomain struct{ Service *analyticsService }

type DomainDeps struct {
	Store     *store.StoreDomain
	Order     *order.OrderDomain
	Asset     *asset.AssetDomain
	Customer  *customer.CustomerDomain
	Expense   *expense.ExpenseDomain
	Inventory *inventory.InventoryDomain
	Owner     *owner.OwnerDomain
	Supplier  *supplier.SupplierDomain
}

func NewDomain(d DomainDeps) *AnalyticsDomain {
	return &AnalyticsDomain{Service: newAnalyticsService(analyticsDeps{
		storeDomain:     d.Store,
		orderDomain:     d.Order,
		assetDomain:     d.Asset,
		customerDomain:  d.Customer,
		expenseDomain:   d.Expense,
		inventoryDomain: d.Inventory,
		ownerDomain:     d.Owner,
		supplierDomain:  d.Supplier,
	})}
}
