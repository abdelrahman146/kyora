package analytics

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/asset"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/expense"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/domain/owner"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/domain/supplier"
	"github.com/govalues/decimal"
)

type analyticsService struct {
	storeDomain     *store.StoreDomain
	orderDomain     *order.OrderDomain
	assetDomain     *asset.AssetDomain
	customerDomain  *customer.CustomerDomain
	expenseDomain   *expense.ExpenseDomain
	inventoryDomain *inventory.InventoryDomain
	ownerDomain     *owner.OwnerDomain
	supplierDomain  *supplier.SupplierDomain
}

func newAnalyticsService(storeDomain *store.StoreDomain, orderDomain *order.OrderDomain, assetDomain *asset.AssetDomain, customerDomain *customer.CustomerDomain, expenseDomain *expense.ExpenseDomain, inventoryDomain *inventory.InventoryDomain, ownerDomain *owner.OwnerDomain, supplierDomain *supplier.SupplierDomain) *analyticsService {
	return &analyticsService{
		storeDomain:     storeDomain,
		orderDomain:     orderDomain,
		assetDomain:     assetDomain,
		customerDomain:  customerDomain,
		expenseDomain:   expenseDomain,
		inventoryDomain: inventoryDomain,
		ownerDomain:     ownerDomain,
		supplierDomain:  supplierDomain,
	}
}

type SalesAnalytics struct {
	StoreID               string
	From                  time.Time
	To                    time.Time
	TotalRevenue          decimal.Decimal
	GrossProfit           decimal.Decimal
	TotalOrders           int64
	AverageOrderValue     decimal.Decimal
	ItemsSold             int64
	NumberOfSalesOverTime *TimeSeries
	RevenueOverTime       *TimeSeries
	TopSellingProducts    []*KeyValue
	OrderStatusBreakdown  []*KeyValue
	SalesByLocation       []*KeyValue
	SalesByChannel        []*KeyValue
}

func (s *analyticsService) GenerateSalesAnalytics(ctx context.Context, storeId string, from, to time.Time) (*SalesAnalytics, error) {
	return nil, nil
}
