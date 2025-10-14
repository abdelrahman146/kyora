package analytics

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/types"
	"github.com/shopspring/decimal"
)

type FinancialPosition struct {
	StoreID string    `json:"storeId"`
	From    time.Time `json:"from"` // The start date of the reporting period.
	To      time.Time `json:"to"`   // The end date of the reporting period.
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

func (s *analyticsService) GenerateFinancialPositionReport(ctx context.Context, storeID string, from, to time.Time) (*FinancialPosition, error) {
	revenue, cogs, _, err := s.orderDomain.OrderService.AggregateSales(ctx, storeID, from, to)
	if err != nil {
		return nil, err
	}
	expenses, _, err := s.expenseDomain.ExpenseService.ExpenseTotals(ctx, storeID, from, to)
	if err != nil {
		return nil, err
	}
	ownerInvestment, ownerDraws, err := s.sumEquityMovements(ctx, storeID, from, to)
	if err != nil {
		return nil, err
	}
	// Inventory valuation is a point-in-time snapshot; we use current value
	totalInventoryValue, _, _, _, err := s.inventoryDomain.InventoryService.InventoryTotals(ctx, storeID, time.Time{}, to)
	if err != nil {
		return nil, err
	}
	fixedAssets, _, err := s.assetDomain.AssetService.AssetTotals(ctx, storeID, from, to)
	if err != nil {
		return nil, err
	}

	// Retained earnings and cash approximation
	retainedEarnings := revenue.Sub(cogs).Sub(expenses)
	cashOnHand := revenue.Add(ownerInvestment).Sub(expenses).Sub(ownerDraws).Sub(fixedAssets)
	currentAssets := cashOnHand.Add(totalInventoryValue)
	totalAssets := currentAssets.Add(fixedAssets)
	totalLiabilities := decimal.Zero // Not tracked yet
	totalEquity := totalAssets.Sub(totalLiabilities)

	fp := &FinancialPosition{
		StoreID:             storeID,
		From:                from,
		To:                  to,
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

func (s *analyticsService) GenerateProfitAndLossReport(ctx context.Context, storeID string, from, to time.Time) (*ProfitAndLossStatement, error) {
	// Revenue, COGS and orders in range
	revenue, cogs, _, err := s.orderDomain.OrderService.AggregateSales(ctx, storeID, from, to)
	if err != nil {
		return nil, err
	}
	gross := revenue.Sub(cogs)

	// OPEX and by-category breakdown in range
	totalOpex, _, err := s.expenseDomain.ExpenseService.ExpenseTotals(ctx, storeID, from, to)
	if err != nil {
		return nil, err
	}
	breakdown, err := s.expenseDomain.ExpenseService.ExpenseBreakdownByCategory(ctx, storeID, from, to)
	if err != nil {
		return nil, err
	}

	net := gross.Sub(totalOpex)
	pl := &ProfitAndLossStatement{
		StoreID:            storeID,
		From:               from,
		To:                 to,
		GrossProfit:        gross,
		TotalExpenses:      totalOpex,
		NetProfit:          net,
		Revenue:            revenue,
		COGS:               cogs,
		ExpensesByCategory: breakdown,
	}
	return pl, nil
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

func (s *analyticsService) GenerateCashFlowReport(ctx context.Context, storeID string, from, to time.Time) (*CashFlowStatement, error) {
	// Cash inflows
	revenue, _, _, err := s.orderDomain.OrderService.AggregateSales(ctx, storeID, from, to)
	if err != nil {
		return nil, err
	}
	ownerIn, err := s.ownerDomain.InvestmentService.SumInvestments(ctx, storeID, from, to)
	if err != nil {
		return nil, err
	}
	totalIn := revenue.Add(ownerIn)

	// Cash outflows
	opex, _, err := s.expenseDomain.ExpenseService.ExpenseTotals(ctx, storeID, from, to)
	if err != nil {
		return nil, err
	}
	// Inventory purchases not tracked explicitly; set to zero until purchase orders are modeled
	inventoryPurchases := decimal.Zero
	// Fixed asset purchases in range
	assetValue, _, err := s.assetDomain.AssetService.AssetTotals(ctx, storeID, from, to)
	if err != nil {
		return nil, err
	}
	ownerDraws, err := s.ownerDomain.OwnerDrawService.SumOwnerDraws(ctx, storeID, from, to)
	if err != nil {
		return nil, err
	}
	businessOperations := inventoryPurchases.Add(opex)
	totalOut := businessOperations.Add(assetValue).Add(ownerDraws)

	// Net cash flow in period
	net := totalIn.Sub(totalOut)

	// Cash at end approximated using financial position as of 'to'
	// Use same approximation as FinancialPosition cashOnHand
	endFP, err := s.GenerateFinancialPositionReport(ctx, storeID, time.Time{}, to)
	if err != nil {
		return nil, err
	}
	cashEnd := endFP.CashOnHand

	// Cash at start = cashEnd - net
	cashStart := cashEnd.Sub(net)

	cs := &CashFlowStatement{
		StoreID:                storeID,
		From:                   from,
		To:                     to,
		CashAtStart:            cashStart,
		CashAtEnd:              cashEnd,
		CashFromCustomers:      revenue,
		CashFromOwner:          ownerIn,
		TotalCashIn:            totalIn,
		InventoryPurchases:     inventoryPurchases,
		OperatingExpenses:      opex,
		TotalBusinessOperation: businessOperations,
		BusinessInvestments:    assetValue,
		OwnerDraws:             ownerDraws,
		TotalCashOut:           totalOut,
		NetCashFlow:            net,
	}
	return cs, nil
}

func (s *analyticsService) sumEquityMovements(ctx context.Context, storeID string, from, to time.Time) (investment decimal.Decimal, draws decimal.Decimal, err error) {
	inv, e1 := s.ownerDomain.InvestmentService.SumInvestments(ctx, storeID, from, to)
	if e1 != nil {
		return decimal.Zero, decimal.Zero, e1
	}
	dr, e2 := s.ownerDomain.OwnerDrawService.SumOwnerDraws(ctx, storeID, from, to)
	if e2 != nil {
		return decimal.Zero, decimal.Zero, e2
	}
	return inv, dr, nil
}
