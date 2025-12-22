package analytics

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/platform/types/keyvalue"
	"github.com/abdelrahman146/kyora/internal/platform/types/timeseries"
	"github.com/shopspring/decimal"
)

type DashboardMetrics struct {
	BusinessID                 string                 `json:"businessID"`
	RevenueLast30Days          decimal.Decimal        `json:"revenueLast30Days"`          // This metric provides a snapshot of recent sales performance, helping to identify trends and measure growth over the past month.
	GrossProfitLast30Days      decimal.Decimal        `json:"grossProfitLast30Days"`      // This metric reflects the actual profit made from sales after accounting for the cost of goods sold, providing a clearer picture of financial health.
	OpenOrdersCount            int64                  `json:"openOrdersCount"`            // This metric provides a snapshot of current operational workload, indicating how many orders are still in progress.
	LowStockItemsCount         int64                  `json:"lowStockItemsCount"`         // This metric helps in proactive inventory management by highlighting items that are running low in stock.
	AllTimeRevenue             decimal.Decimal        `json:"allTimeRevenue"`             // The big-picture, motivational number showing the total revenue generated since the beginning.
	SafeToDrawAmount           decimal.Decimal        `json:"safeToDrawAmount"`           // This metric indicates the amount of money that can be safely withdrawn from the business without jeopardizing operational needs. (net profit - safety buffer)
	SalesPerformanceLast30Days *timeseries.TimeSeries `json:"salesPerformanceLast30Days"` //This chart provides the pulse of the business, showing daily activity over the past month.
	LiveOrderFunnel            []keyvalue.KeyValue    `json:"liveOrderFunnel"`            // This visualization gives an instant overview of your current operational workload. It shows where all non-completed orders are in the fulfillment process.
	TopSellingProducts         []*inventory.Product   `json:"topSellingProducts"`         // This chart highlights your best-performing products over the last 30 days, helping you identify trends and make informed inventory decisions.
	NewCustomersTimeSeries     *timeseries.TimeSeries `json:"newCustomersTimeSeries"`     // This chart tracks the number of new customers acquired each day over the past month, providing insights into customer growth trends.
}

type SalesAnalytics struct {
	BusinessID            string                 `json:"businessID"`
	From                  time.Time              `json:"from"`
	To                    time.Time              `json:"to"`
	TotalRevenue          decimal.Decimal        `json:"totalRevenue"`
	GrossProfit           decimal.Decimal        `json:"grossProfit"`
	TotalOrders           int64                  `json:"totalOrders"`
	AverageOrderValue     decimal.Decimal        `json:"averageOrderValue"`
	ItemsSold             int64                  `json:"itemsSold"`
	NumberOfSalesOverTime *timeseries.TimeSeries `json:"numberOfSalesOverTime"` // line chart
	RevenueOverTime       *timeseries.TimeSeries `json:"revenueOverTime"`       // line chart
	TopSellingProducts    []*inventory.Product   `json:"topSellingProducts"`    // pie chart
	OrderStatusBreakdown  []keyvalue.KeyValue    `json:"orderStatusBreakdown"`  // donut chart
	SalesByCountry        []keyvalue.KeyValue    `json:"salesByCountry"`        // Table
	SalesByChannel        []keyvalue.KeyValue    `json:"salesByChannel"`        // Table
}

type InventoryAnalytics struct {
	BusinessID                  string
	From                        time.Time
	To                          time.Time
	TotalInventoryValue         decimal.Decimal
	TotalInStock                int64
	LowStockItems               int64
	OutOfStockItems             int64
	InventoryTurnoverRatio      decimal.Decimal
	SellThroughRate             decimal.Decimal
	TopProductsByInventoryValue []*inventory.Product // Bar chart
}

type ExpenseAnalytics struct {
	BusinessID            string
	From                  time.Time
	To                    time.Time
	TotalExpenses         decimal.Decimal
	AveragExpenseAmount   decimal.Decimal
	TotalNumberOfEntries  int64
	ExpensesOverTime      *timeseries.TimeSeries // line chart
	TopExpensesByCategory []keyvalue.KeyValue    // Bar chart
}

type CustomerAnalytics struct {
	BusinessID                       string
	From                             time.Time
	To                               time.Time
	NewCustomers                     int64
	ReturningCustomers               int64
	RepeatCustomerRate               decimal.Decimal
	AverageRevenuePerCustomer        decimal.Decimal
	CustomerAcquisitionCost          decimal.Decimal
	CustomerLifetimeValue            decimal.Decimal
	AverageCustomerPurchaseFrequency decimal.Decimal
	NewCustomersOverTime             *timeseries.TimeSeries // line chart
	TopCustomersByRevenue            []*customer.Customer   // table 'customer id, customer name, country, total revenue
}

type AssetAnalytics struct {
	BusinessID              string
	From                    time.Time
	To                      time.Time
	TotalAssetsAquired      int64
	TotalAssetValue         decimal.Decimal
	AssetsByCategory        []keyvalue.KeyValue    // Bar chart
	AssetInvestmentOverTime *timeseries.TimeSeries // line chart
}

// Financial Reports Model (Balance Sheet, P&L, Cash Flow)
//-----------------------------------------------//

// FinancialPosition represents a snapshot of a business's financial health at a specific point in time. (Balance Sheet)
type FinancialPosition struct {
	BusinessID string    `json:"businessID"`
	AsOf       time.Time `json:"asOf"` // The end date of the reporting period.
	// core totals
	TotalAssets      decimal.Decimal `json:"totalAssets"`      // The total value of everything the business owns. (CurrentAssets + FixedAssets)
	TotalLiabilities decimal.Decimal `json:"totalLiabilities"` // The total value of everything the business owes. For now, this is zero. because we are not tracking liabilities yet.
	TotalEquity      decimal.Decimal `json:"totalEquity"`      // The net value of the business (Assets - Liabilities). The value left over
	// breakdown of assets
	CashOnHand          decimal.Decimal `json:"cashOnHand"`          // Cash on Hand: The total cash business bank account (Revenue + Owner Investment) - (Expenses + Owner Draw + Asset Purchases)
	TotalInventoryValue decimal.Decimal `json:"totalInventoryValue"` // The total cost value of all products available for sale.
	CurrentAssets       decimal.Decimal `json:"currentAssets"`       // Short-term resources.  cashOnHand + totalInventoryValue
	FixedAssets         decimal.Decimal `json:"fixedAssets"`         // Long-term resources. The total cost value of all owned assets (e.g., equipment, property)
	// liabilities - for future use
	// equity breakdown
	OwnerInvestment  decimal.Decimal `json:"ownerInvestment"`  // The total amount of money the owner has invested into the business.
	RetainedEarnings decimal.Decimal `json:"retainedEarnings"` // The cumulative net profit that has been reinvested in the business rather than distributed to the owner. (All-Time Revenue - All-Time COGS - All-Time OPEX)
	OwnerDraws       decimal.Decimal `json:"ownerDraws"`       // The total amount of money the owner has withdrawn from the business for personal use.
}

// ProfitAndLossStatement represents the financial performance of a business over a specific period. (Income Statement)
type ProfitAndLossStatement struct {
	BusinessID         string              `json:"businessID"`
	AsOf               time.Time           `json:"asOf"`               // The end date of the reporting period
	GrossProfit        decimal.Decimal     `json:"grossProfit"`        // The profit made directly from selling your products, before any other business expenses. Calculation: Revenue - Cost of Goods Sold
	TotalExpenses      decimal.Decimal     `json:"totalExpenses"`      // The total of all operating expenses (OPEX) incurred in running the business.
	NetProfit          decimal.Decimal     `json:"netProfit"`          // The final profit after all expenses have been deducted from gross profit. Calculation: Gross Profit - Total Expenses
	Revenue            decimal.Decimal     `json:"revenue"`            // Total revenue generated from sales.
	COGS               decimal.Decimal     `json:"cogs"`               // Cost of Goods Sold: The direct costs attributable to the production of the goods sold by the business.
	ExpensesByCategory []keyvalue.KeyValue `json:"expensesByCategory"` // Breakdown of expenses by category
}

// CashFlowStatement represents the cash inflows and outflows of a business over a specific period.
type CashFlowStatement struct {
	BusinessID             string          `json:"businessID"`
	AsOf                   time.Time       `json:"asOf"`                   // The end date of the reporting period
	CashAtStart            decimal.Decimal `json:"cashAtStart"`            // The cash balance the business had at the beginning of the selected period.
	CashAtEnd              decimal.Decimal `json:"cashAtEnd"`              // The final cash balance the business had at the end of the period. This is the current runway.
	CashFromCustomers      decimal.Decimal `json:"cashFromCustomers"`      // Total cash received from customers (sales revenue) during the period.
	CashFromOwner          decimal.Decimal `json:"cashFromOwner"`          // Total cash invested into the business by the owner during the period.
	TotalCashIn            decimal.Decimal `json:"totalCashIn"`            // Total cash inflows (money coming into the business) during the period. Calculation: CashFromCustomers + CashFromOwner
	InventoryPurchases     decimal.Decimal `json:"inventoryPurchases"`     // Total cash spent on purchasing inventory during the period.
	OperatingExpenses      decimal.Decimal `json:"operatingExpenses"`      // Total cash spent on operating expenses (OPEX) during the period.
	TotalBusinessOperation decimal.Decimal `json:"totalBusinessOperation"` // Total cash outflows (money going out of the business) during the period. Calculation: InventoryPurchases + OperatingExpenses
	BusinessInvestments    decimal.Decimal `json:"businessInvestments"`    // Total cash spent on purchasing fixed assets (e.g., equipment, property) during the period.
	OwnerDraws             decimal.Decimal `json:"ownerDraws"`             // Total cash withdrawn from the business by the owner for personal use during the period.
	TotalCashOut           decimal.Decimal `json:"totalCashOut"`           // Total cash outflows (money going out of the business) during the period. Calculation: TotalBusinessOperation + BusinessInvestments + OwnerDraws
	NetCashFlow            decimal.Decimal `json:"netCashFlow"`            // The total change in the business's cash balance during the period (Total Cash In - Total Cash Out). This can be positive (✅) or negative (🔻)
}
