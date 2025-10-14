package analytics

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/types"
	"github.com/shopspring/decimal"
)

type Dashboard struct {
	StoreID                    string            `json:"storeId"`
	RevenueLast30Days          decimal.Decimal   `json:"revenueLast30Days"`          // This metric provides a snapshot of recent sales performance, helping to identify trends and measure growth over the past month.
	GrossProfitLast30Days      decimal.Decimal   `json:"grossProfitLast30Days"`      // This metric reflects the actual profit made from sales after accounting for the cost of goods sold, providing a clearer picture of financial health.
	OpenOrdersCount            int64             `json:"openOrdersCount"`            // This metric provides a snapshot of current operational workload, indicating how many orders are still in progress.
	LowStockItemsCount         int64             `json:"lowStockItemsCount"`         // This metric helps in proactive inventory management by highlighting items that are running low in stock.
	AllTimeRevenue             decimal.Decimal   `json:"allTimeRevenue"`             // The big-picture, motivational number showing the total revenue generated since the beginning.
	SafeToDrawAmount           decimal.Decimal   `json:"safeToDrawAmount"`           // This metric indicates the amount of money that can be safely withdrawn from the business without jeopardizing operational needs. (net profit - safety buffer)
	SalesPerformanceLast30Days *types.TimeSeries `json:"salesPerformanceLast30Days"` //This chart provides the pulse of the business, showing daily activity over the past month.
	LiveOrderFunnel            []types.KeyValue  `json:"liveOrderFunnel"`            // This visualization gives an instant overview of your current operational workload. It shows where all non-completed orders are in the fulfillment process.
	TopSellingProducts         []types.KeyValue  `json:"topSellingProducts"`         // This chart highlights your best-performing products over the last 30 days, helping you identify trends and make informed inventory decisions.
	NewCustomersTimeSeries     *types.TimeSeries `json:"newCustomersTimeSeries"`     // This chart tracks the number of new customers acquired each day over the past month, providing insights into customer growth trends.
}

// GenerateDashboardAnalytics compiles fast-loading dashboard snapshot metrics for the
// provided store. It intentionally limits the look-back window for recent
// performance metrics to the last 30 days while using all‑time aggregates for
// motivational / cumulative values. Expensive multi-range queries are kept
// minimal and reused computations are avoided.
//
// SafeToDrawAmount heuristic:
//
//	safe = max( (allTimeRevenue - allTimeCOGS - allTimeExpenses) - totalOwnerDraws - safetyBuffer, 0 )
//
// Rationale: preserve one month of operating expense runway before distributing
// additional profits; owner draws already taken reduce remaining distributable
// profit. This is a conservative, easily explainable buffer.
func (s *analyticsService) GenerateDashboardAnalytics(ctx context.Context, storeID string) (*Dashboard, error) {
	// Define core time windows
	to := time.Now().UTC()
	from30 := to.AddDate(0, 0, -30)
	bucket := "day" // daily granularity for 30‑day series

	// 1. Last 30 days revenue & gross profit / open orders / low stock / top products / sales performance / new customers.
	revenue30, cogs30, _, err := s.orderDomain.OrderService.AggregateSales(ctx, storeID, from30, to)
	if err != nil {
		return nil, err
	}
	grossProfit30 := revenue30.Sub(cogs30)

	openOrdersCount, err := s.orderDomain.OrderService.OpenOrdersCount(ctx, storeID)
	if err != nil {
		return nil, err
	}

	// Inventory totals (low stock value only needed but call once to reuse).
	_, _, lowStock, _, err := s.inventoryDomain.InventoryService.InventoryTotals(ctx, storeID, from30, to)
	if err != nil {
		return nil, err
	}

	// Revenue time series (sales performance) – daily revenue last 30 days
	revRows, err := s.orderDomain.OrderService.RevenueTimeSeries(ctx, storeID, from30, to, bucket)
	if err != nil {
		return nil, err
	}
	salesPerformanceTS := types.NewTimeSeries(ctx, revRows, from30, to)

	// Live order funnel (non-completed)
	funnel, err := s.orderDomain.OrderService.OpenOrdersFunnel(ctx, storeID)
	if err != nil {
		return nil, err
	}

	// Top selling products (qty) last 30 days
	topProducts, err := s.orderDomain.OrderService.TopSellingProducts(ctx, storeID, from30, to, 5)
	if err != nil {
		return nil, err
	}

	// New customers time series last 30 days
	newCustRows, err := s.customerDomain.CustomerService.NewCustomersTimeSeries(ctx, storeID, from30, to, bucket)
	if err != nil {
		return nil, err
	}
	newCustomersTS := types.NewTimeSeries(ctx, newCustRows, from30, to)

	// 2. All‑time aggregates: revenue & expenses & COGS & owner draws for SafeToDraw.
	allTimeRevenue, allTimeCogs, _, err := s.orderDomain.OrderService.AllTimeSalesAggregate(ctx, storeID)
	if err != nil {
		return nil, err
	}
	allTimeExpenses, err := s.expenseDomain.ExpenseService.TotalExpensesAllTime(ctx, storeID)
	if err != nil {
		return nil, err
	}
	totalOwnerDraws, err := s.ownerDomain.OwnerDrawService.SumTotalOwnerDraws(ctx, storeID)
	if err != nil {
		return nil, err
	}
	last30Expenses, _, err := s.expenseDomain.ExpenseService.ExpenseTotals(ctx, storeID, from30, to)
	if err != nil {
		return nil, err
	}

	netProfitAllTime := allTimeRevenue.Sub(allTimeCogs).Sub(allTimeExpenses)
	// conservative buffer: one month of expenses
	safetyBuffer := last30Expenses
	safeToDraw := netProfitAllTime.Sub(totalOwnerDraws).Sub(safetyBuffer)
	if safeToDraw.IsNegative() {
		safeToDraw = decimal.Zero
	}

	dashboard := &Dashboard{
		StoreID:                    storeID,
		RevenueLast30Days:          revenue30,
		GrossProfitLast30Days:      grossProfit30,
		OpenOrdersCount:            openOrdersCount,
		LowStockItemsCount:         lowStock,
		AllTimeRevenue:             allTimeRevenue,
		SafeToDrawAmount:           safeToDraw,
		SalesPerformanceLast30Days: salesPerformanceTS,
		LiveOrderFunnel:            funnel,
		TopSellingProducts:         topProducts,
		NewCustomersTimeSeries:     newCustomersTS,
	}
	return dashboard, nil
}
