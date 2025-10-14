package analytics

import (
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

func (s *analyticsService) GenerateDashboardAnalytics(storeID string) (*Dashboard, error) {
	return nil, nil
}
