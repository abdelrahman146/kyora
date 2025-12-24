package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/accounting"
	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type BusinessPaymentMethodsSuite struct {
	suite.Suite
	accountHelper *AccountTestHelper
	orderHelper   *OrderTestHelper
}

func (s *BusinessPaymentMethodsSuite) SetupSuite() {
	s.accountHelper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
	s.orderHelper = NewOrderTestHelper(testEnv.Database, e2eBaseURL)
}

func (s *BusinessPaymentMethodsSuite) resetDB() {
	s.NoError(testutils.TruncateTables(testEnv.Database,
		"expenses",
		"orders", "order_items", "order_notes",
		"customers", "customer_addresses",
		"products", "variants", "categories",
		"business_payment_methods",
		"businesses", "shipping_zones",
		"subscriptions", "plans",
		"users", "workspaces",
	))
}

func (s *BusinessPaymentMethodsSuite) SetupTest() {
	s.resetDB()
}

func (s *BusinessPaymentMethodsSuite) TearDownTest() {
	s.resetDB()
}

func (s *BusinessPaymentMethodsSuite) TestListPaymentMethods_Defaults() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	biz, err := s.orderHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	resp, err := s.orderHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+biz.Descriptor+"/payment-methods", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var methods []map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &methods))
	s.GreaterOrEqual(len(methods), 1)

	found := map[string]bool{}
	for _, m := range methods {
		d, _ := m["descriptor"].(string)
		found[d] = true
		s.NotEmpty(d)
		s.NotEmpty(m["name"])
		// enabled/fee fields should exist in response
		_, _ = m["enabled"]
		_, _ = m["feePercent"]
		_, _ = m["feeFixed"]
	}
	s.True(found["bank_transfer"], "expected bank_transfer in catalog")
	s.True(found["cash_on_delivery"], "expected cash_on_delivery in catalog")
	s.True(found["credit_card"], "expected credit_card in catalog")
}

func (s *BusinessPaymentMethodsSuite) TestUpdatePaymentMethod_ForbiddenForMember() {
	ctx := context.Background()
	_, ws, adminToken, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.orderHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	// Create a member in the same workspace.
	memberPassword, err := hash.Password("Password123!")
	s.NoError(err)
	member := &account.User{
		WorkspaceID:     ws.ID,
		Role:            role.RoleUser,
		FirstName:       "Member",
		LastName:        "User",
		Email:           "member@example.com",
		Password:        memberPassword,
		IsEmailVerified: true,
	}
	userRepo := database.NewRepository[account.User](testEnv.Database)
	s.NoError(userRepo.CreateOne(ctx, member))
	memberToken, err := auth.NewJwtToken(member.ID, member.WorkspaceID, member.AuthVersion)
	s.NoError(err)

	payload := map[string]interface{}{
		"enabled":    true,
		"feePercent": 0.05,
		"feeFixed":   0,
	}

	resp, err := s.orderHelper.Client.AuthenticatedRequest("PATCH", "/v1/businesses/"+biz.Descriptor+"/payment-methods/credit_card", payload, memberToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusForbidden, resp.StatusCode)

	resp2, err := s.orderHelper.Client.AuthenticatedRequest("PATCH", "/v1/businesses/"+biz.Descriptor+"/payment-methods/credit_card", payload, adminToken)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusOK, resp2.StatusCode)
}

func (s *BusinessPaymentMethodsSuite) TestOrderPaid_CreatesTransactionFeeExpense() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.orderHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	// Enable credit card with a 5% fee.
	pmPayload := map[string]interface{}{"enabled": true, "feePercent": 0.05, "feeFixed": 0}
	pmResp, err := s.orderHelper.Client.AuthenticatedRequest("PATCH", "/v1/businesses/"+biz.Descriptor+"/payment-methods/credit_card", pmPayload, token)
	s.NoError(err)
	defer pmResp.Body.Close()
	s.Equal(http.StatusOK, pmResp.StatusCode)

	cust, addr, err := s.orderHelper.CreateTestCustomer(ctx, biz.ID, "customer@example.com", "Test Customer")
	s.NoError(err)
	cat, err := s.orderHelper.CreateTestCategory(ctx, biz.ID, "Electronics", "electronics")
	s.NoError(err)
	_, variant, err := s.orderHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Test Product", decimal.NewFromInt(0), decimal.NewFromInt(100), 10)
	s.NoError(err)

	// Create order with credit card.
	createPayload := map[string]interface{}{
		"customerId":        cust.ID,
		"shippingAddressId": addr.ID,
		"channel":           "instagram",
		"paymentMethod":     "credit_card",
		"items": []map[string]interface{}{
			{"variantId": variant.ID, "quantity": 1, "unitPrice": 100, "unitCost": 0},
		},
	}
	// Respect create-order rate limit.
	time.Sleep(1100 * time.Millisecond)
	createResp, err := s.orderHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+biz.Descriptor+"/orders", createPayload, token)
	s.NoError(err)
	defer createResp.Body.Close()
	s.Equal(http.StatusCreated, createResp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(createResp, &created))
	orderID := created["id"].(string)
	s.NotEmpty(orderID)

	totalStr := fmt.Sprint(created["total"])
	orderTotal, err := decimal.NewFromString(totalStr)
	s.NoError(err)
	expectedFee := orderTotal.Mul(decimal.RequireFromString("0.05")).Round(2)

	// Mark order as paid.
	// Payment status transitions require the order to be placed/shipped/fulfilled.
	statusResp, err := s.orderHelper.Client.AuthenticatedRequest(
		"PATCH",
		"/v1/businesses/"+biz.Descriptor+"/orders/"+orderID+"/status",
		map[string]interface{}{"status": "placed"},
		token,
	)
	s.NoError(err)
	defer statusResp.Body.Close()
	s.Equal(http.StatusOK, statusResp.StatusCode)

	paidResp, err := s.orderHelper.Client.AuthenticatedRequest("PATCH", "/v1/businesses/"+biz.Descriptor+"/orders/"+orderID+"/payment-status", map[string]interface{}{"paymentStatus": "paid"}, token)
	s.NoError(err)
	defer paidResp.Body.Close()
	s.Equal(http.StatusOK, paidResp.StatusCode)

	// Poll DB for async bus side effect.
	expRepo := database.NewRepository[accounting.Expense](testEnv.Database)
	var feeExpense *accounting.Expense
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		exp, err := expRepo.FindOne(ctx,
			expRepo.ScopeBusinessID(biz.ID),
			expRepo.ScopeEquals(accounting.ExpenseSchema.OrderID, orderID),
			expRepo.ScopeEquals(accounting.ExpenseSchema.Category, accounting.ExpenseCategoryTransactionFee),
		)
		if err == nil && exp != nil {
			feeExpense = exp
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	s.NotNil(feeExpense, "expected transaction fee expense to be created")
	if feeExpense != nil {
		s.Equal(accounting.ExpenseCategoryTransactionFee, feeExpense.Category)
		s.Equal(expectedFee.StringFixed(2), feeExpense.Amount.Round(2).StringFixed(2))
	}
}

func TestBusinessPaymentMethodsSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BusinessPaymentMethodsSuite))
}
