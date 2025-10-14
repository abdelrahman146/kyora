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
	NumberOfSalesOverTime *types.TimeSeries // line chart
	RevenueOverTime       *types.TimeSeries // line chart
	TopSellingProducts    []types.KeyValue  // pie chart
	OrderStatusBreakdown  []types.KeyValue  // donut chart
	SalesByCountry        []types.KeyValue  // Table
	SalesByChannel        []types.KeyValue  // Table
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
		SalesByCountry:        locKV,
		SalesByChannel:        chKV,
	}
	return result, nil
}

type InventoryAnalytics struct {
	StoreID                     string
	From                        time.Time
	To                          time.Time
	TotalInventoryValue         decimal.Decimal
	TotalInStock                int64
	LowStockItems               int64
	OutOfStockItems             int64
	InventoryTurnoverRatio      decimal.Decimal
	SellThroughRate             decimal.Decimal
	InventoryValueOverTime      *types.TimeSeries // line chart
	TopProductsByInventoryValue []types.KeyValue  // Bar chart
}

func (s *analyticsService) GenerateInventoryAnalytics(ctx context.Context, storeId string, from, to time.Time) (*InventoryAnalytics, error) {
	return nil, nil
}

type ExpenseAnalytics struct {
	StoreID               string
	From                  time.Time
	To                    time.Time
	TotalExpenses         decimal.Decimal
	AveragExpenseAmount   decimal.Decimal
	TotalNumberOfEntries  int64
	ExpensesOverTime      *types.TimeSeries // line chart
	TopExpensesByCategory []types.KeyValue  // Bar chart
}

func (s *analyticsService) GenerateExpenseAnalytics(ctx context.Context, storeId string, from, to time.Time) (*ExpenseAnalytics, error) {
	return nil, nil
}

type CustomerAnalytics struct {
	StoreID                          string
	From                             time.Time
	To                               time.Time
	TotalCustomers                   int64
	NewCustomers                     int64
	ReturningCustomers               int64
	RepeatCustomerRate               decimal.Decimal
	AverageRevenuePerCustomer        decimal.Decimal
	CustomerAcquisitionCost          decimal.Decimal
	CustomerLifetimeValue            decimal.Decimal
	AverageCustomerPurchaseFrequency decimal.Decimal
	NewCustomersOverTime             *types.TimeSeries // line chart
	ReturningCustomersOverTime       *types.TimeSeries // line chart
	TopCustomersByRevenue            []any             // table 'cusomter id, custoemr name, country, total revenue
}

func (s *analyticsService) GenerateCustomerAnalytics(ctx context.Context, storeId string, from, to time.Time) (*CustomerAnalytics, error) {
	return nil, nil
}

type AssetAnalytics struct {
	StoreID                 string
	From                    time.Time
	To                      time.Time
	TotalAssetsAquired      int64
	TotalAssetValue         decimal.Decimal
	AssetsByCategory        []types.KeyValue  // Bar chart
	AssetInvestmentOverTime *types.TimeSeries // line chart
}

func (s *analyticsService) GenerateAssetAnalytics(ctx context.Context, storeId string, from, to time.Time) (*AssetAnalytics, error) {
	return nil, nil
}
