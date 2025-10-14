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
	"github.com/shopspring/decimal"
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
type analyticsDeps struct {
	storeDomain     *store.StoreDomain
	orderDomain     *order.OrderDomain
	assetDomain     *asset.AssetDomain
	customerDomain  *customer.CustomerDomain
	expenseDomain   *expense.ExpenseDomain
	inventoryDomain *inventory.InventoryDomain
	ownerDomain     *owner.OwnerDomain
	supplierDomain  *supplier.SupplierDomain
}

func newAnalyticsService(d analyticsDeps) *analyticsService {
	return &analyticsService{
		storeDomain:     d.storeDomain,
		orderDomain:     d.orderDomain,
		assetDomain:     d.assetDomain,
		customerDomain:  d.customerDomain,
		expenseDomain:   d.expenseDomain,
		inventoryDomain: d.inventoryDomain,
		ownerDomain:     d.ownerDomain,
		supplierDomain:  d.supplierDomain,
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
