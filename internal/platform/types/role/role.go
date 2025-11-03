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
	},
}

type Action string

const (
	ActionManage Action = "manage"
	ActionView   Action = "view"
)

type Resource string

const (
	ResourceAccount    Resource = "account"
	ResourceBilling    Resource = "billing"
	ResourceBusiness   Resource = "business"
	ResourceCustomer   Resource = "customer"
	ResourceOrder      Resource = "order"
	ResourceInventory  Resource = "inventory"
	ResourceExpense    Resource = "expense"
	ResourceAccounting Resource = "accounting"
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
