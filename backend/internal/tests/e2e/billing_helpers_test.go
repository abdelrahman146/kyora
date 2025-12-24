package e2e_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/nullable"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	stripelib "github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/paymentmethod"
)

type BillingTestHelper struct {
	db     *database.Database
	client *testutils.HTTPClient
}

func NewBillingTestHelper(db *database.Database, cacheAddr, baseURL string) *BillingTestHelper {
	_ = cache.NewConnection([]string{cacheAddr})
	return &BillingTestHelper{db: db, client: testutils.NewHTTPClient(baseURL)}
}

func (h *BillingTestHelper) Client() *testutils.HTTPClient { return h.client }

func (h *BillingTestHelper) UniqueSlug(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		prefix = "t"
	}
	return fmt.Sprintf("%s_%s", prefix, id.Base62(10))
}

func (h *BillingTestHelper) UniqueEmail(prefix string) string {
	return fmt.Sprintf("%s_%s@example.com", h.UniqueSlug(prefix), id.Base62(6))
}

func (h *BillingTestHelper) CreateTestUser(ctx context.Context, email string, userRole role.Role) (*account.User, *account.Workspace, string, error) {
	user, ws, token, err := testutils.CreateAuthenticatedUser(ctx, h.db, email, "Pass123!", "Test", "User", userRole)
	if err != nil {
		return nil, nil, "", err
	}
	return user, ws, token, nil
}

func (h *BillingTestHelper) CreateStripeCardPaymentMethod() (string, error) {
	params := &stripelib.PaymentMethodParams{
		Type: stripelib.String(string(stripelib.PaymentMethodTypeCard)),
		Card: &stripelib.PaymentMethodCardParams{
			Number:   stripelib.String("4242424242424242"),
			ExpMonth: stripelib.Int64(12),
			ExpYear:  stripelib.Int64(2030),
			CVC:      stripelib.String("123"),
		},
	}
	pm, err := paymentmethod.New(params)
	if err != nil {
		return "", err
	}
	return pm.ID, nil
}

func (h *BillingTestHelper) CreatePlan(ctx context.Context, descriptor string, price decimal.Decimal, limits billing.PlanLimit) (*billing.Plan, error) {
	planRepo := database.NewRepository[billing.Plan](h.db)
	if existing, err := planRepo.FindOne(ctx, planRepo.ScopeEquals(billing.PlanSchema.Descriptor, descriptor)); err == nil && existing != nil {
		return existing, nil
	}

	features := billing.PlanFeature{
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
	}

	stripePlanID := fmt.Sprintf("stripe_test_plan_%s_%d", descriptor, time.Now().UnixNano())
	plan := &billing.Plan{
		Descriptor:   descriptor,
		Name:         "Test " + descriptor,
		Description:  "E2E plan " + descriptor,
		StripePlanID: stripePlanID,
		Price:        price,
		Currency:     "usd",
		BillingCycle: billing.BillingCycleMonthly,
		Features:     features,
		Limits:       limits,
	}
	if err := planRepo.CreateOne(ctx, plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (h *BillingTestHelper) AddWorkspaceUser(ctx context.Context, wsID, email string, userRole role.Role) (*account.User, error) {
	hashed, err := hash.Password("Pass123!")
	if err != nil {
		return nil, err
	}
	user := &account.User{
		WorkspaceID:     wsID,
		Role:            userRole,
		FirstName:       "Extra",
		LastName:        "User",
		Email:           email,
		Password:        hashed,
		IsEmailVerified: true,
	}
	userRepo := database.NewRepository[account.User](h.db)
	if err := userRepo.CreateOne(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (h *BillingTestHelper) CreateBusiness(ctx context.Context, wsID string) (*business.Business, error) {
	bizRepo := database.NewRepository[business.Business](h.db)
	biz := &business.Business{
		WorkspaceID:  wsID,
		Descriptor:   fmt.Sprintf("biz_%d", time.Now().UnixNano()),
		Name:         "Test Biz",
		CountryCode:  "US",
		Currency:     "usd",
		VatRate:      decimal.Zero,
		SafetyBuffer: decimal.Zero,
	}
	if err := bizRepo.CreateOne(ctx, biz); err != nil {
		return nil, err
	}
	return biz, nil
}

func (h *BillingTestHelper) CreateCustomerWithAddress(ctx context.Context, businessID string) (*customer.Customer, *customer.CustomerAddress, error) {
	custRepo := database.NewRepository[customer.Customer](h.db)
	addrRepo := database.NewRepository[customer.CustomerAddress](h.db)

	cust := &customer.Customer{
		BusinessID:  businessID,
		Name:        "John Doe",
		CountryCode: "US",
		Gender:      customer.GenderOther,
		Email:       nullable.NewString(fmt.Sprintf("cust_%d@example.com", time.Now().UnixNano())),
	}
	if err := custRepo.CreateOne(ctx, cust); err != nil {
		return nil, nil, err
	}

	addr := &customer.CustomerAddress{
		CustomerID:  cust.ID,
		CountryCode: "US",
		State:       "CA",
		City:        "SF",
		Street:      nullable.NewString("Market"),
		PhoneCode:   "+1",
		PhoneNumber: "5551234",
	}
	if err := addrRepo.CreateOne(ctx, addr); err != nil {
		return nil, nil, err
	}
	return cust, addr, nil
}

func (h *BillingTestHelper) CreateOrder(ctx context.Context, businessID, customerID, addressID string) (*order.Order, error) {
	orderRepo := database.NewRepository[order.Order](h.db)
	o := &order.Order{
		OrderNumber:       fmt.Sprintf("ord_%d", time.Now().UnixNano()),
		BusinessID:        businessID,
		CustomerID:        customerID,
		ShippingAddressID: addressID,
		Channel:           "instagram",
		Subtotal:          decimal.Zero,
		VAT:               decimal.Zero,
		VATRate:           decimal.Zero,
		ShippingFee:       decimal.Zero,
		Discount:          decimal.Zero,
		COGS:              decimal.Zero,
		Total:             decimal.Zero,
		Currency:          "usd",
		Status:            order.OrderStatusPending,
		PaymentStatus:     order.OrderPaymentStatusPending,
		PaymentMethod:     order.OrderPaymentMethodBankTransfer,
		OrderedAt:         time.Now().UTC(),
	}
	if err := orderRepo.CreateOne(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}
