package analytics

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/accounting"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/platform/types/keyvalue"
	"github.com/shopspring/decimal"
)

type ServiceParams struct {
	Inventory  *inventory.Service
	Orders     *order.Service
	Accounting *accounting.Service
	Customer   *customer.Service
}

type Service struct {
	inventory  *inventory.Service
	customer   *customer.Service
	orders     *order.Service
	accounting *accounting.Service
}

func NewService(params *ServiceParams) *Service {
	return &Service{
		inventory:  params.Inventory,
		orders:     params.Orders,
		accounting: params.Accounting,
		customer:   params.Customer,
	}
}

func (s *Service) ComputeDashboardMetrics(ctx context.Context, actor *account.User, biz *business.Business) (*DashboardMetrics, error) {
	dashboard := &DashboardMetrics{
		BusinessID: biz.ID,
	}
	var err error
	last30Days := time.Now().AddDate(0, 0, -30)
	// get revenue last 30 days
	dashboard.RevenueLast30Days, err = s.orders.SumOrdersTotal(ctx, actor, biz, last30Days, time.Now())
	if err != nil {
		return nil, err
	}
	cogsLast30Days, err := s.orders.SumOrdersCOGS(ctx, actor, biz, last30Days, time.Now())
	if err != nil {
		return nil, err
	}
	// get gross profit last 30 days
	dashboard.GrossProfitLast30Days = dashboard.RevenueLast30Days.Sub(cogsLast30Days)
	// open orders count
	dashboard.OpenOrdersCount, err = s.orders.CountOpenOrders(ctx, actor, biz)
	if err != nil {
		return nil, err
	}
	// LowStockItemsCount
	dashboard.LowStockItemsCount, err = s.inventory.CountLowStockVariants(ctx, actor, biz)
	if err != nil {
		return nil, err
	}
	// AllTimeRevenue
	dashboard.AllTimeRevenue, err = s.orders.SumOrdersTotal(ctx, actor, biz, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}
	// SafeToDrawAmount
	cogs, err := s.orders.SumOrdersCOGS(ctx, actor, biz, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}
	dashboard.SafeToDrawAmount, err = s.accounting.ComputeSafeToDrawAmount(ctx, actor, biz, dashboard.AllTimeRevenue, cogs, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}
	// SalesPerformanceLast30Days
	dashboard.SalesPerformanceLast30Days, err = s.orders.ComputeRevenueTimeSeries(ctx, actor, biz, last30Days, time.Now())
	if err != nil {
		return nil, err
	}
	// LiveOrderFunnel
	dashboard.LiveOrderFunnel, err = s.orders.ComputeLiveOrdersFunnel(ctx, actor, biz, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}
	// TopSellingProducts
	dashboard.TopSellingProducts, err = s.orders.ComputeTopSellingProducts(ctx, actor, biz, 5, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}
	// NewCustomersTimeSeries
	dashboard.NewCustomersTimeSeries, err = s.customer.ComputeCustomersTimeSeries(ctx, actor, biz, last30Days, time.Now())
	if err != nil {
		return nil, err
	}
	return dashboard, nil
}

func (s *Service) ComputeSalesAnalytics(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (*SalesAnalytics, error) {
	analytics := &SalesAnalytics{
		BusinessID: biz.ID,
		From:       from,
		To:         to,
	}
	// compute totals
	totalRevenue, err := s.orders.SumOrdersTotal(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	analytics.TotalRevenue = totalRevenue

	cogs, err := s.orders.SumOrdersCOGS(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	analytics.GrossProfit = totalRevenue.Sub(cogs)

	// orders metrics
	analytics.TotalOrders, err = s.orders.CountOrdersByDateRange(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	analytics.AverageOrderValue, err = s.orders.AvgOrdersTotal(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	analytics.ItemsSold, err = s.orders.SumItemsSold(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}

	// time series charts
	analytics.NumberOfSalesOverTime, err = s.orders.ComputeOrdersCountTimeSeries(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	analytics.RevenueOverTime, err = s.orders.ComputeRevenueTimeSeries(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}

	// breakdowns and top lists
	analytics.TopSellingProducts, err = s.orders.ComputeTopSellingProducts(ctx, actor, biz, 5, from, to)
	if err != nil {
		return nil, err
	}
	analytics.OrderStatusBreakdown, err = s.orders.CountOrdersByStatus(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	analytics.SalesByCountry, err = s.orders.SumOrdersTotalByCountry(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	analytics.SalesByChannel, err = s.orders.SumOrdersTotalByChannel(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}

	return analytics, nil
}

func (s *Service) ComputeInventoryAnalytics(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (*InventoryAnalytics, error) {
	analytics := &InventoryAnalytics{
		BusinessID: biz.ID,
		From:       from,
		To:         to,
	}
	// Totals and counts
	totalInvValue, err := s.inventory.SumInventoryValue(ctx, actor, biz)
	if err != nil {
		return nil, err
	}
	analytics.TotalInventoryValue = totalInvValue

	totalUnitsInStock, err := s.inventory.SumStockQuantity(ctx, actor, biz)
	if err != nil {
		return nil, err
	}
	analytics.TotalInStock = totalUnitsInStock

	lowStock, err := s.inventory.CountLowStockVariants(ctx, actor, biz)
	if err != nil {
		return nil, err
	}
	analytics.LowStockItems = lowStock

	outOfStock, err := s.inventory.CountOutOfStockVariants(ctx, actor, biz)
	if err != nil {
		return nil, err
	}
	analytics.OutOfStockItems = outOfStock

	// Ratios
	cogs, err := s.orders.SumOrdersCOGS(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	avgInventory := totalInvValue // best available proxy without historical ledger
	if !avgInventory.IsZero() {
		analytics.InventoryTurnoverRatio = cogs.Div(avgInventory)
	} else {
		analytics.InventoryTurnoverRatio = decimal.Zero
	}

	itemsSold, err := s.orders.SumItemsSold(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	denom := decimal.NewFromInt(itemsSold + totalUnitsInStock)
	if denom.GreaterThan(decimal.Zero) {
		analytics.SellThroughRate = decimal.NewFromInt(itemsSold).Div(denom)
	} else {
		analytics.SellThroughRate = decimal.Zero
	}

	topProducts, err := s.inventory.ComputeTopProductsByInventoryValueDetailed(ctx, actor, biz, 5)
	if err != nil {
		return nil, err
	}
	analytics.TopProductsByInventoryValue = make([]inventory.ProductResponse, len(topProducts))
	for i, tp := range topProducts {
		analytics.TopProductsByInventoryValue[i] = tp.Product
	}

	return analytics, nil
}

func (s *Service) ComputeCustomerAnalytics(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (*CustomerAnalytics, error) {
	analytics := &CustomerAnalytics{
		BusinessID: biz.ID,
		From:       from,
		To:         to,
	}
	// 1) New vs Returning customers (volume and rate)
	// New customers are those created (joined) in the selected period
	newCustomers, err := s.customer.CountCustomersByDateRange(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	analytics.NewCustomers = newCustomers

	// Unique purchasing customers and returning definition
	returningCustomers, err := s.orders.CountReturningCustomers(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	ordersByCustomer, err := s.orders.CountOrdersByCustomer(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	uniquePurchasers := int64(len(ordersByCustomer))
	analytics.ReturningCustomers = returningCustomers
	if uniquePurchasers > 0 {
		analytics.RepeatCustomerRate = decimal.NewFromInt(returningCustomers).Div(decimal.NewFromInt(uniquePurchasers))
	} else {
		analytics.RepeatCustomerRate = decimal.Zero
	}

	// 2) Revenue-centric metrics
	totalRevenue, err := s.orders.SumOrdersTotal(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	if uniquePurchasers > 0 {
		analytics.AverageRevenuePerCustomer = totalRevenue.Div(decimal.NewFromInt(uniquePurchasers))
	} else {
		analytics.AverageRevenuePerCustomer = decimal.Zero
	}

	// Customer Acquisition Cost (CAC) = Marketing Expenses in period / NewCustomers (if any)
	marketingSpend, err := s.accounting.SumExpensesAmountByCategory(ctx, actor, biz, accounting.ExpenseCategoryMarketing, from, to)
	if err != nil {
		return nil, err
	}
	if newCustomers > 0 {
		analytics.CustomerAcquisitionCost = marketingSpend.Div(decimal.NewFromInt(newCustomers))
	} else {
		analytics.CustomerAcquisitionCost = decimal.Zero
	}

	averageOrderValue, err := s.orders.AvgOrdersTotal(ctx, actor, biz, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}
	totalOrdersCount, err := s.orders.CountOrdersByDateRange(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}
	averageOrderFrequencey := decimal.Zero
	if uniquePurchasers > 0 {
		averageOrderFrequencey = decimal.NewFromInt(totalOrdersCount).Div(decimal.NewFromInt(uniquePurchasers))
	}
	// Customer Lifetime Value (CLV) = Average Order Value  x  Average Purchase Frequency
	analytics.CustomerLifetimeValue = averageOrderValue.Mul(averageOrderFrequencey)

	// Average purchase frequency in the period = total orders / unique purchasers
	analytics.AverageCustomerPurchaseFrequency = averageOrderFrequencey

	// 3) Time series
	analytics.NewCustomersOverTime, err = s.customer.ComputeCustomersTimeSeries(ctx, actor, biz, from, to)
	if err != nil {
		return nil, err
	}

	// 4) Top customers by revenue
	topKV, err := s.orders.SumOrdersTotalByCustomer(ctx, actor, biz, 5, from, to)
	if err != nil {
		return nil, err
	}
	ids := keyvalue.KeysFromKeyValueSlice(topKV)
	customers, err := s.customer.GetCustomersByIDs(ctx, actor, biz, ids)
	if err != nil {
		return nil, err
	}
	// Reorder fetched customers to match ids order
	byID := make(map[any]*customer.Customer, len(customers))
	for _, c := range customers {
		byID[c.ID] = c
	}
	ordered := make([]*customer.Customer, 0, len(ids))
	for _, id := range ids {
		if c, ok := byID[id]; ok {
			ordered = append(ordered, c)
		}
	}
	analytics.TopCustomersByRevenue = ordered

	return analytics, nil
}

func (s *Service) ComputeFinancialPosition(ctx context.Context, actor *account.User, biz *business.Business, asOf time.Time) (*FinancialPosition, error) {
	financialPosition := &FinancialPosition{
		BusinessID: biz.ID,
		AsOf:       asOf,
	}

	// Current inventory on hand (valued at cost)
	invValue, err := s.inventory.SumInventoryValue(ctx, actor, biz)
	if err != nil {
		return nil, err
	}
	financialPosition.TotalInventoryValue = invValue

	// Fixed assets purchased to date
	fixedAssets, err := s.accounting.SumAssetsValue(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}
	financialPosition.FixedAssets = fixedAssets

	// Owner equity movements
	ownerInvestment, err := s.accounting.SumInvestmentsAmount(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}
	financialPosition.OwnerInvestment = ownerInvestment

	ownerDraws, err := s.accounting.SumWithdrawalsAmount(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}
	financialPosition.OwnerDraws = ownerDraws

	// Operating metrics to compute retained earnings and cash
	totalRevenue, err := s.orders.SumOrdersTotal(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}
	totalCOGS, err := s.orders.SumOrdersCOGS(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}
	totalExpenses, err := s.accounting.SumExpensesAmount(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}

	// Retained Earnings = All-Time Revenue - All-Time COGS - All-Time OPEX
	financialPosition.RetainedEarnings = totalRevenue.Sub(totalCOGS).Sub(totalExpenses)

	// Cash on Hand approximation:
	// Cash = (Revenue + Owner Investment) - (Expenses + Owner Draws + Asset Purchases + Inventory Value)
	cashInflows := totalRevenue.Add(ownerInvestment)
	cashOutflows := totalExpenses.Add(ownerDraws).Add(fixedAssets).Add(invValue)
	financialPosition.CashOnHand = cashInflows.Sub(cashOutflows)

	// Current Assets = Cash + Inventory
	financialPosition.CurrentAssets = financialPosition.CashOnHand.Add(financialPosition.TotalInventoryValue)

	// Total Assets = Current Assets + Fixed Assets
	financialPosition.TotalAssets = financialPosition.CurrentAssets.Add(financialPosition.FixedAssets)

	// Liabilities are not tracked yet, set to zero
	financialPosition.TotalLiabilities = decimal.Zero

	// Total Equity = Assets - Liabilities
	financialPosition.TotalEquity = financialPosition.TotalAssets.Sub(financialPosition.TotalLiabilities)

	return financialPosition, nil
}

func (s *Service) ComputeProfitAndLoss(ctx context.Context, actor *account.User, biz *business.Business, asOf time.Time) (*ProfitAndLossStatement, error) {
	statement := &ProfitAndLossStatement{
		BusinessID: biz.ID,
		AsOf:       asOf,
	}
	// Revenue up to asOf
	revenue, err := s.orders.SumOrdersTotal(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}
	statement.Revenue = revenue

	// COGS up to asOf
	cogs, err := s.orders.SumOrdersCOGS(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}
	statement.COGS = cogs

	// Gross Profit = Revenue - COGS
	statement.GrossProfit = statement.Revenue.Sub(statement.COGS)

	// Operating Expenses up to asOf
	totalExpenses, err := s.accounting.SumExpensesAmount(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}
	statement.TotalExpenses = totalExpenses

	// Expense breakdown by category
	categories := accounting.ExpenseCategoriesList()
	breakdown := make([]keyvalue.KeyValue, 0, len(categories))
	for _, cat := range categories {
		amt, err := s.accounting.SumExpensesAmountByCategory(ctx, actor, biz, cat, time.Time{}, asOf)
		if err != nil {
			return nil, err
		}
		breakdown = append(breakdown, keyvalue.New(string(cat), amt))
	}
	statement.ExpensesByCategory = breakdown

	// Net Profit = Gross Profit - Total Expenses
	statement.NetProfit = statement.GrossProfit.Sub(statement.TotalExpenses)

	return statement, nil
}

func (s *Service) ComputeCashFlow(ctx context.Context, actor *account.User, biz *business.Business, asOf time.Time) (*CashFlowStatement, error) {
	statement := &CashFlowStatement{
		BusinessID: biz.ID,
		AsOf:       asOf,
	}
	// We currently don't track liabilities or an inventory purchase ledger.
	// To stay consistent with ComputeFinancialPosition, we approximate cash flows on a cash-basis using:
	// - Cash inflows: Revenue (cash from customers) + Owner investments
	// - Cash outflows: Operating expenses + Owner draws + Fixed asset purchases + Inventory on hand (as a proxy for historical inventory purchases)
	// This keeps CashAtEnd aligned with FinancialPosition.CashOnHand.

	// Inflows up to asOf (all-time to date)
	revenue, err := s.orders.SumOrdersTotal(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}
	statement.CashFromCustomers = revenue

	ownerInvestment, err := s.accounting.SumInvestmentsAmount(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}
	statement.CashFromOwner = ownerInvestment
	statement.TotalCashIn = statement.CashFromCustomers.Add(statement.CashFromOwner)

	// Operating outflows up to asOf
	totalExpenses, err := s.accounting.SumExpensesAmount(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}
	statement.OperatingExpenses = totalExpenses

	// Inventory purchases proxy: current inventory value (aligns with FinancialPosition approximation)
	invValue, err := s.inventory.SumInventoryValue(ctx, actor, biz)
	if err != nil {
		return nil, err
	}
	statement.InventoryPurchases = invValue
	statement.TotalBusinessOperation = statement.InventoryPurchases.Add(statement.OperatingExpenses)

	// Investing outflows (fixed assets) up to asOf
	fixedAssets, err := s.accounting.SumAssetsValue(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}
	statement.BusinessInvestments = fixedAssets

	// Financing outflows (owner draws) up to asOf
	ownerDraws, err := s.accounting.SumWithdrawalsAmount(ctx, actor, biz, time.Time{}, asOf)
	if err != nil {
		return nil, err
	}
	statement.OwnerDraws = ownerDraws

	// Totals and net change
	statement.TotalCashOut = statement.TotalBusinessOperation.
		Add(statement.BusinessInvestments).
		Add(statement.OwnerDraws)
	statement.NetCashFlow = statement.TotalCashIn.Sub(statement.TotalCashOut)

	// Since this statement aggregates from inception to asOf, the starting cash is assumed zero
	// and ending cash equals the net cash flow. This matches FinancialPosition.CashOnHand.
	statement.CashAtStart = decimal.Zero
	statement.CashAtEnd = statement.CashAtStart.Add(statement.NetCashFlow)

	return statement, nil
}
