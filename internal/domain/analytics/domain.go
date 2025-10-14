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

type AnalyticsDomain struct {
	Service *analyticsService
}

func NewDomain(storeDomain *store.StoreDomain, orderDomain *order.OrderDomain, assetDomain *asset.AssetDomain, customerDomain *customer.CustomerDomain, expenseDomain *expense.ExpenseDomain, inventoryDomain *inventory.InventoryDomain, ownerDomain *owner.OwnerDomain, supplierDomain *supplier.SupplierDomain) *AnalyticsDomain {
	return &AnalyticsDomain{
		Service: newAnalyticsService(storeDomain, orderDomain, assetDomain, customerDomain, expenseDomain, inventoryDomain, ownerDomain, supplierDomain),
	}
}
