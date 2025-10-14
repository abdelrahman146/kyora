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
	if to.Before(from) {
		from, to = to, from
	}
	bucket := types.BucketForRange(from, to)

	// Aggregate current inventory state
	totalValue, totalUnits, lowStock, outOfStock, err := s.inventoryDomain.InventoryService.InventoryTotals(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}

	// Sales + COGS for turnover & sell-through approximations
	revenue, cogs, orderCount, err := s.orderDomain.OrderService.AggregateSales(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}
	_ = revenue // currently unused in inventory metrics but could be exposed later
	itemsSold, err := s.orderDomain.OrderService.ItemsSold(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}

	// Inventory Turnover Ratio (approx): COGS / Average Inventory Value.
	// Lacking historical daily valuation snapshots we approximate average with ending inventory value.
	// Safe-guard division by zero.
	invTurnover := decimal.Zero
	if totalValue.Sign() > 0 {
		invTurnover = cogs.Div(totalValue)
	}

	// Sell-through rate (approx): items sold / (items sold + ending inventory units)
	sellThrough := decimal.Zero
	denom := itemsSold + totalUnits
	if denom > 0 {
		sellThrough = decimal.NewFromInt(itemsSold).Div(decimal.NewFromInt(denom))
	}

	// Inventory value over time & top products by inventory value
	// For lack of movement history we only compute current top products.
	topProducts, err := s.inventoryDomain.InventoryService.TopProductsByInventoryValue(ctx, storeId, 10)
	if err != nil {
		return nil, err
	}

	// Placeholder empty time series until movement tracking implemented
	inventoryValueTS := types.NewTimeSeries(ctx, []types.TimeSeriesRow{}, from, to)

	result := &InventoryAnalytics{
		StoreID:                     storeId,
		From:                        from,
		To:                          to,
		TotalInventoryValue:         totalValue,
		TotalInStock:                totalUnits,
		LowStockItems:               lowStock,
		OutOfStockItems:             outOfStock,
		InventoryTurnoverRatio:      invTurnover,
		SellThroughRate:             sellThrough,
		InventoryValueOverTime:      inventoryValueTS,
		TopProductsByInventoryValue: topProducts,
	}
	_ = orderCount // reserved for future metrics
	_ = bucket
	return result, nil
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
	if to.Before(from) {
		from, to = to, from
	}
	bucket := types.BucketForRange(from, to)

	total, count, err := s.expenseDomain.ExpenseService.ExpenseTotals(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}
	avg := decimal.Zero
	if count > 0 {
		avg = total.Div(decimal.NewFromInt(count))
	}

	// Time series & breakdown
	rows, err := s.expenseDomain.ExpenseService.ExpenseAmountTimeSeries(ctx, storeId, from, to, bucket)
	if err != nil {
		return nil, err
	}
	ts := types.NewTimeSeries(ctx, rows, from, to)
	breakdown, err := s.expenseDomain.ExpenseService.ExpenseBreakdownByCategory(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}

	result := &ExpenseAnalytics{
		StoreID:               storeId,
		From:                  from,
		To:                    to,
		TotalExpenses:         total,
		AveragExpenseAmount:   avg,
		TotalNumberOfEntries:  count,
		ExpensesOverTime:      ts,
		TopExpensesByCategory: breakdown,
	}
	return result, nil
}

type CustomerAnalytics struct {
	StoreID                          string
	From                             time.Time
	To                               time.Time
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
	if to.Before(from) {
		from, to = to, from
	}
	bucket := types.BucketForRange(from, to)

	// New customers in range
	newCustomers, err := s.customerDomain.CustomerService.CountNewCustomersInRange(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}

	// Distinct purchasing customers & returning customers
	distinctPurchasers, err := s.orderDomain.OrderService.DistinctPurchasingCustomers(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}
	returningCustomers, err := s.orderDomain.OrderService.ReturningCustomersCount(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}

	repeatRate := decimal.Zero
	if distinctPurchasers > 0 {
		repeatRate = decimal.NewFromInt(returningCustomers).Div(decimal.NewFromInt(distinctPurchasers))
	}

	// Revenue & gross profit in period for ARPC and margin
	totalRevenue, cogs, orderCount, err := s.orderDomain.OrderService.AggregateSales(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}
	grossProfit := totalRevenue.Sub(cogs)
	avgRevenuePerCustomer := decimal.Zero
	if distinctPurchasers > 0 {
		avgRevenuePerCustomer = totalRevenue.Div(decimal.NewFromInt(distinctPurchasers))
	}

	// Average purchase frequency (orders per purchasing customer)
	avgPurchaseFreq := decimal.Zero
	if distinctPurchasers > 0 {
		avgPurchaseFreq = decimal.NewFromInt(orderCount).Div(decimal.NewFromInt(distinctPurchasers))
	}

	// Customer Lifetime Value (simple heuristic): (avg revenue per customer * gross margin %) * repeat purchase rate
	grossMarginPct := decimal.Zero
	if !totalRevenue.IsZero() {
		grossMarginPct = grossProfit.Div(totalRevenue)
	}
	clv := avgRevenuePerCustomer.Mul(grossMarginPct).Mul(repeatRate)

	// Acquisition cost attempts: marketing expenses in period / new customers (if available)
	marketingExpenses, err := s.expenseDomain.ExpenseService.MarketingExpensesInRange(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}
	cac := decimal.Zero
	if newCustomers > 0 {
		cac = marketingExpenses.Div(decimal.NewFromInt(newCustomers))
	}

	// Time series placeholders
	newCustRows, err := s.customerDomain.CustomerService.NewCustomersTimeSeries(ctx, storeId, from, to, bucket)
	if err != nil {
		return nil, err
	}
	newCustTS := types.NewTimeSeries(ctx, newCustRows, from, to)
	returningCustTS, err := s.orderDomain.OrderService.ReturningCustomersTimeSeries(ctx, storeId, from, to, bucket)
	if err != nil {
		return nil, err
	}
	returningTS := types.NewTimeSeries(ctx, returningCustTS, from, to)

	// Top customers by revenue (id only; enrichment left for future join)
	topCust, err := s.orderDomain.OrderService.RevenuePerCustomer(ctx, storeId, from, to, 10)
	if err != nil {
		return nil, err
	}
	topAny := make([]any, len(topCust))
	for i, kv := range topCust {
		topAny[i] = kv
	}

	result := &CustomerAnalytics{
		StoreID:                          storeId,
		From:                             from,
		To:                               to,
		NewCustomers:                     newCustomers,
		ReturningCustomers:               returningCustomers,
		RepeatCustomerRate:               repeatRate,
		AverageRevenuePerCustomer:        avgRevenuePerCustomer,
		CustomerAcquisitionCost:          cac,
		CustomerLifetimeValue:            clv,
		AverageCustomerPurchaseFrequency: avgPurchaseFreq,
		NewCustomersOverTime:             newCustTS,
		ReturningCustomersOverTime:       returningTS,
		TopCustomersByRevenue:            topAny,
	}
	return result, nil
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
	if to.Before(from) {
		from, to = to, from
	}
	bucket := types.BucketForRange(from, to)

	totalValue, count, err := s.assetDomain.AssetService.AssetTotals(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}
	breakdown, err := s.assetDomain.AssetService.AssetBreakdownByType(ctx, storeId, from, to)
	if err != nil {
		return nil, err
	}
	tsRows, err := s.assetDomain.AssetService.AssetValueTimeSeries(ctx, storeId, from, to, bucket)
	if err != nil {
		return nil, err
	}
	ts := types.NewTimeSeries(ctx, tsRows, from, to)
	result := &AssetAnalytics{
		StoreID:                 storeId,
		From:                    from,
		To:                      to,
		TotalAssetsAquired:      count,
		TotalAssetValue:         totalValue,
		AssetsByCategory:        breakdown,
		AssetInvestmentOverTime: ts,
	}
	return result, nil
}
