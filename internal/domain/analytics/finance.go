package analytics

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/types"
	"github.com/shopspring/decimal"
)

type FinancialPosition struct {
	StoreID string     `json:"storeId"`
	From    *time.Time `json:"from"` // The start date of the reporting period.
	To      *time.Time `json:"to"`   // The end date of the reporting period.
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

func (s *analyticsService) GenerateFinancialPositionReport(ctx context.Context, storeID string, startDate, endDate *time.Time) (*FinancialPosition, error) {
	// Reuse existing domain aggregates to build a snapshot-style financial position
	// All-Time Revenue, COGS
	allTimeRevenue, allTimeCOGS, _, err := s.orderDomain.OrderService.AllTimeSalesAggregate(ctx, storeID)
	if err != nil {
		return nil, err
	}

	// All-Time Operating Expenses (OPEX)
	allTimeExpenses, err := s.expenseDomain.ExpenseService.TotalExpensesAllTime(ctx, storeID)
	if err != nil {
		return nil, err
	}

	// Owner investments and draws (equity movements)
	ownerInvestment, err := s.ownerDomain.InvestmentService.CalculateTotalInvestedAmount(ctx, storeID)
	if err != nil {
		return nil, err
	}
	ownerDraws, err := s.ownerDomain.OwnerDrawService.SumTotalOwnerDraws(ctx, storeID)
	if err != nil {
		return nil, err
	}

	// Inventory valuation (current, at cost)
	totalInventoryValue, _, _, _, err := s.inventoryDomain.InventoryService.InventoryTotals(ctx, storeID)
	if err != nil {
		return nil, err
	}

	// Fixed assets value (sum of purchased asset values, all-time)
	// Pass zero time range to aggregate for all-time based on repository implementation
	fixedAssets, _, err := s.assetDomain.AssetService.AssetTotals(ctx, storeID, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}

	// Retained earnings: revenue - COGS - OPEX
	retainedEarnings := allTimeRevenue.Sub(allTimeCOGS).Sub(allTimeExpenses)

	// Cash on hand: (Revenue + Owner Investment) - (Expenses + Owner Draws + Asset Purchases)
	cashOnHand := allTimeRevenue.Add(ownerInvestment).Sub(allTimeExpenses).Sub(ownerDraws).Sub(fixedAssets)

	// Current assets: cash + inventory (simplified model)
	currentAssets := cashOnHand.Add(totalInventoryValue)

	// Totals
	totalAssets := currentAssets.Add(fixedAssets)
	totalLiabilities := decimal.Zero // Not tracked yet
	totalEquity := totalAssets.Sub(totalLiabilities)

	fp := &FinancialPosition{
		StoreID:             storeID,
		TotalAssets:         totalAssets,
		TotalLiabilities:    totalLiabilities,
		TotalEquity:         totalEquity,
		CashOnHand:          cashOnHand,
		TotalInventoryValue: totalInventoryValue,
		CurrentAssets:       currentAssets,
		FixedAssets:         fixedAssets,
		OwnerInvestment:     ownerInvestment,
		RetainedEarnings:    retainedEarnings,
		OwnerDraws:          ownerDraws,
	}
	return fp, nil
}

type ProfitAndLossStatement struct {
	StoreID            string           `json:"storeId"`
	From               time.Time        `json:"from"`               // The start date of the reporting period.
	To                 time.Time        `json:"to"`                 // The end date of the reporting period.
	GrossProfit        decimal.Decimal  `json:"grossProfit"`        // The profit made directly from selling your products, before any other business expenses. Calculation: Revenue - Cost of Goods Sold
	TotalExpenses      decimal.Decimal  `json:"totalExpenses"`      // The total of all operating expenses (OPEX) incurred in running the business.
	NetProfit          decimal.Decimal  `json:"netProfit"`          // The final profit after all expenses have been deducted from gross profit. Calculation: Gross Profit - Total Expenses
	Revenue            decimal.Decimal  `json:"revenue"`            // Total revenue generated from sales.
	COGS               decimal.Decimal  `json:"cogs"`               // Cost of Goods Sold: The direct costs attributable to the production of the goods sold by the business.
	ExpensesByCategory []types.KeyValue `json:"expensesByCategory"` // Breakdown of expenses by category
}

func (s *analyticsService) GenerateProfitAndLossReport(ctx context.Context, storeID string, startDate, endDate *time.Time) (*ProfitAndLossStatement, error) {
	return nil, nil
}

type CashFlowStatement struct {
	StoreID                string          `json:"storeId"`
	From                   time.Time       `json:"from"`                   // The start date of the reporting period.
	To                     time.Time       `json:"to"`                     // The end date of the reporting period.
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

func (s *analyticsService) GenerateCashFlowReport(ctx context.Context, storeID string, startDate, endDate *time.Time) (*CashFlowStatement, error) {
	return nil, nil
}
