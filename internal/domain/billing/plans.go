package billing

import "github.com/shopspring/decimal"

// predefined plans
// these will be upserted into the database on storage initialization
// StripePlanID will be populated by SyncPlansToStripe service method
var plans = []Plan{
	{
		Descriptor:   "starter",
		Name:         "Starter Plan",
		Description:  "Perfect for sole owners just getting started with their business",
		Price:        decimal.NewFromInt(0),
		Currency:     "aed",
		StripePlanID: "", // Will be populated by Stripe sync
		BillingCycle: BillingCycleMonthly,
		Features: PlanFeature{
			CustomerManagement:       true,
			InventoryManagement:      true,
			OrderManagement:          true,
			ExpenseManagement:        true,
			Accounting:               true,
			BasicAnalytics:           true,
			FinancialReports:         true,
			DataImport:               false,
			DataExport:               false,
			AdvancedAnalytics:        false,
			AdvancedFinancialReports: false,
			OrderPaymentLinks:        false,
			InvoiceGeneration:        false,
			ExportAnalyticsData:      false,
			AIBusinessAssistant:      false,
		},
		Limits: PlanLimit{
			MaxOrdersPerMonth: 25,
			MaxTeamMembers:    1,
			MaxBusinesses:     1,
		},
	},
	{
		Descriptor:   "professional",
		Name:         "Professional Plan",
		Description:  "For growing businesses that need advanced features and higher limits",
		Price:        decimal.NewFromFloat(54.99),
		Currency:     "aed",
		StripePlanID: "", // Will be populated by Stripe sync
		BillingCycle: BillingCycleMonthly,
		Features: PlanFeature{
			CustomerManagement:       true,
			InventoryManagement:      true,
			OrderManagement:          true,
			ExpenseManagement:        true,
			Accounting:               true,
			BasicAnalytics:           true,
			FinancialReports:         true,
			DataImport:               true,
			DataExport:               true,
			AdvancedAnalytics:        true,
			AdvancedFinancialReports: true,
			OrderPaymentLinks:        true,
			InvoiceGeneration:        true,
			ExportAnalyticsData:      false,
			AIBusinessAssistant:      false,
		},
		Limits: PlanLimit{
			MaxOrdersPerMonth: 500,
			MaxTeamMembers:    5,
			MaxBusinesses:     3,
		},
	},
	{
		Descriptor:   "enterprise",
		Name:         "Enterprise Plan",
		Description:  "For large businesses with advanced needs and unlimited usage",
		Price:        decimal.NewFromFloat(155.99),
		Currency:     "aed",
		StripePlanID: "", // Will be populated by Stripe sync
		BillingCycle: BillingCycleMonthly,
		Features: PlanFeature{
			CustomerManagement:       true,
			InventoryManagement:      true,
			OrderManagement:          true,
			ExpenseManagement:        true,
			Accounting:               true,
			BasicAnalytics:           true,
			FinancialReports:         true,
			DataImport:               true,
			DataExport:               true,
			AdvancedAnalytics:        true,
			AdvancedFinancialReports: true,
			OrderPaymentLinks:        true,
			InvoiceGeneration:        true,
			ExportAnalyticsData:      true,
			AIBusinessAssistant:      true,
		},
		Limits: PlanLimit{
			MaxOrdersPerMonth: -1, // Unlimited
			MaxTeamMembers:    -1, // Unlimited
			MaxBusinesses:     -1, // Unlimited
		},
	},
}
