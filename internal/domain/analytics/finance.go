package analytics

import (
	"context"

	"github.com/shopspring/decimal"
)

type FinancialPosition struct {
	StoreID string `json:"storeId"`
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

func (s *analyticsService) GenerateFinancialPositionReport(ctx context.Context, storeID string) (*FinancialPosition, error) {
	return nil, nil
}
