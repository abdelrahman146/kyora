package role

import (
	"fmt"
	"slices"

	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// static permissions
var RolePermissions = map[Role][]string{
	RoleUser: {
		"view:account",
		"view:business",
		"view:customer",
		"view:order",
		"view:inventory",
		"view:expense",
		"view:accounting",
		"view:basic_analytics",
		"view:basic_financial_reports",
	},
	RoleAdmin: {
		"view:account",
		"manage:account",
		"view:billing",
		"manage:billing",
		"view:business",
		"manage:business",
		"view:customer",
		"manage:customer",
		"view:order",
		"manage:order",
		"view:inventory",
		"manage:inventory",
		"view:expense",
		"manage:expense",
		"view:accounting",
		"manage:accounting",
		"view:basic_analytics",
		"view:basic_financial_reports",
	},
}

type Action string

const (
	ActionManage Action = "manage"
	ActionView   Action = "view"
)

type Resource string

const (
	ResourceAccount                  Resource = "account"
	ResourceBilling                  Resource = "billing"
	ResourceBusiness                 Resource = "business"
	ResourceCustomer                 Resource = "customer"
	ResourceOrder                    Resource = "order"
	ResourceInventory                Resource = "inventory"
	ResourceExpense                  Resource = "expense"
	ResourceAccounting               Resource = "accounting"
	ResourceBasicAnalytics           Resource = "basic_analytics"
	ResourceAdvancedAnalytics        Resource = "advanced_analytics"
	ResourceFinancialReports         Resource = "basic_financial_reports"
	ResourceAdvancedFinancialReports Resource = "advanced_financial_reports"
	ResourceOrderPaymentLinks        Resource = "order_payment_links"
	ResourceOrderInvoiceGeneration   Resource = "order_invoice_generation"
	ResourceAIBusinessAssistant      Resource = "ai_business_assistant"
	ResourceExportAnalyticsData      Resource = "export_analytics_data"
	ResourceDataImport               Resource = "data_import"
	ResourceDataExport               Resource = "data_export"
)

func (r Role) HasPermission(action Action, resource Resource) error {
	permissions, ok := RolePermissions[r]
	if !ok {
		return UnauthorizedError(action, resource)
	}
	permissionToCheck := string(action) + ":" + string(resource)
	if !slices.Contains(permissions, permissionToCheck) {
		return UnauthorizedError(action, resource)
	}
	return nil
}

func UnauthorizedError(action Action, resource Resource) error {
	return problem.Forbidden(fmt.Sprintf("unauthorized to %s %s", action, resource)).With("action", action).With("resource", resource)
}
