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
	"github.com/abdelrahman146/kyora/internal/types"
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
	NumberOfSalesOverTime *types.TimeSeries
	RevenueOverTime       *types.TimeSeries
	TopSellingProducts    []types.KeyValue
	OrderStatusBreakdown  []types.KeyValue
	SalesByLocation       []types.KeyValue
	SalesByChannel        []types.KeyValue
}

func (s *analyticsService) GenerateSalesAnalytics(ctx context.Context, storeId string, from, to time.Time) (*SalesAnalytics, error) {
	if to.Before(from) {
		from, to = to, from
	}
	bucket := types.BucketForRange(from, to)

	// Aggregate metrics
	totalRevenue, cogs, totalOrders, err := s.orderDomain.OrderService.AggregateSales(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}
	itemsSold, err := s.orderDomain.OrderService.ItemsSold(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}

	// Gross profit = revenue - COGS
	grossProfit := totalRevenue.Sub(cogs)

	// Average order value
	aov := decimal.Zero
	if totalOrders > 0 {
		aov = totalRevenue.Div(decimal.NewFromInt(totalOrders))
	}

	// Time series
	revRows, err := s.orderDomain.OrderService.RevenueTimeSeries(ctx, storeId, from, to, bucket)
	if err != nil {
		return nil, err
	}
	numRows, err := s.orderDomain.OrderService.CountTimeSeries(ctx, storeId, from, to, bucket)
	if err != nil {
		return nil, err
	}

	// Map to analytics.TimeSeries
	revenueOverTime := types.NewTimeSeries(ctx, revRows, from, to)
	numberOfSalesOverTime := types.NewTimeSeries(ctx, numRows, from, to)

	// Top selling products
	topProducts, err := s.orderDomain.OrderService.TopSellingProducts(ctx, storeId, from, to, 10)
	if err != nil {
		return nil, err
	}

	// Status breakdown (counts)
	statusKV, err := s.orderDomain.OrderService.BreakdownByStatus(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}

	// Sales by location (country revenue)
	locKV, err := s.orderDomain.OrderService.BreakdownByCountry(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}

	// Sales by channel (revenue)
	chKV, err := s.orderDomain.OrderService.BreakdownByChannel(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}

	result := &SalesAnalytics{
		StoreID:               storeId,
		From:                  from,
		To:                    to,
		TotalRevenue:          totalRevenue,
		GrossProfit:           grossProfit,
		TotalOrders:           totalOrders,
		AverageOrderValue:     aov,
		ItemsSold:             itemsSold,
		NumberOfSalesOverTime: numberOfSalesOverTime,
		RevenueOverTime:       revenueOverTime,
		TopSellingProducts:    topProducts,
		OrderStatusBreakdown:  statusKV,
		SalesByLocation:       locKV,
		SalesByChannel:        chKV,
	}
	return result, nil
}
