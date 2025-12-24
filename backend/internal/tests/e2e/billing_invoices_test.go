package e2e_test

import (
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type BillingInvoicesSuite struct {
	suite.Suite
	helper *BillingTestHelper
}

func (s *BillingInvoicesSuite) SetupSuite() {
	s.helper = NewBillingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *BillingInvoicesSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database,
		"plans",
		"subscriptions",
		"billing_invoice_records",
		"stripe_events",
		"users",
		"workspaces",
	)
	s.NoError(err)
}

func (s *BillingInvoicesSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database,
		"plans",
		"subscriptions",
		"billing_invoice_records",
		"stripe_events",
		"users",
		"workspaces",
	)
	s.NoError(err)
}

func (s *BillingInvoicesSuite) TestInvoices_CreateAndList_AndBOLAOnPay() {
	ctx := s.T().Context()
	_, _, token1, err := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("admin1"), role.RoleAdmin)
	s.NoError(err)
	_, _, token2, err := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("admin2"), role.RoleAdmin)
	s.NoError(err)

	respCreate, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/invoices", map[string]interface{}{
		"description": "Manual invoice",
		"amount":      500,
		"currency":    "usd",
	}, token1)
	s.NoError(err)
	defer respCreate.Body.Close()
	s.Equal(http.StatusCreated, respCreate.StatusCode)

	var inv map[string]interface{}
	s.NoError(testutils.DecodeJSON(respCreate, &inv))
	s.Contains(inv, "hosted_invoice_url")
	s.Contains(inv, "invoice_pdf")
	s.Contains(inv, "id")
	invoiceID, ok := inv["id"].(string)
	s.True(ok)
	s.NotEmpty(invoiceID)

	respList, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/invoices?page=1&pageSize=10", nil, token1)
	s.NoError(err)
	defer respList.Body.Close()
	s.Equal(http.StatusOK, respList.StatusCode)

	var listResp map[string]interface{}
	s.NoError(testutils.DecodeJSON(respList, &listResp))
	s.Contains(listResp, "items")
	items, ok := listResp["items"].([]interface{})
	s.True(ok)
	s.NotEmpty(items)
	first, ok := items[0].(map[string]interface{})
	s.True(ok)
	s.Contains(first, "id")
	s.Contains(first, "status")
	s.Contains(first, "currency")
	s.Contains(first, "amountDue")
	s.Contains(first, "amountPaid")
	s.Contains(first, "createdAt")

	// Download should redirect for owner
	respDownload, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/invoices/"+invoiceID+"/download", nil, token1)
	s.NoError(err)
	defer respDownload.Body.Close()
	s.Equal(http.StatusConflict, respDownload.StatusCode)
	var prob map[string]interface{}
	s.NoError(testutils.DecodeJSON(respDownload, &prob))
	s.Equal(float64(http.StatusConflict), prob["status"])
	s.Equal("Conflict", prob["title"])
	s.Equal("invoice is not ready for download", prob["detail"])

	// BOLA: second workspace must not be able to pay first workspace invoice
	respPay, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/invoices/"+invoiceID+"/pay", nil, token2)
	s.NoError(err)
	defer respPay.Body.Close()
	s.Equal(http.StatusNotFound, respPay.StatusCode)

	// BOLA: second workspace must not be able to download first workspace invoice
	respDownload2, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/invoices/"+invoiceID+"/download", nil, token2)
	s.NoError(err)
	defer respDownload2.Body.Close()
	s.Equal(http.StatusNotFound, respDownload2.StatusCode)
}

func (s *BillingInvoicesSuite) TestInvoices_Pay_HappyPath_WithPaymentMethod() {
	ctx := s.T().Context()
	descriptor := s.helper.UniqueSlug("paid")
	_, err := s.helper.CreatePlan(ctx, descriptor, decimal.NewFromInt(10), billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)

	_, _, token, err := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("admin"), role.RoleAdmin)
	s.NoError(err)

	// Ensure subscription + customer exist
	respSub, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": descriptor}, token)
	s.NoError(err)
	defer respSub.Body.Close()
	s.Require().Equal(http.StatusOK, respSub.StatusCode)

	// Attach a card payment method so invoice.pay can succeed.
	pmID, err := s.helper.CreateStripeCardPaymentMethod()
	s.NoError(err)
	respAttach, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/payment-methods/attach", map[string]interface{}{"paymentMethodId": pmID}, token)
	s.NoError(err)
	defer respAttach.Body.Close()
	s.Equal(http.StatusOK, respAttach.StatusCode)

	respCreate, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/invoices", map[string]interface{}{
		"description": "Manual invoice",
		"amount":      500,
		"currency":    "usd",
	}, token)
	s.NoError(err)
	defer respCreate.Body.Close()
	s.Require().Equal(http.StatusCreated, respCreate.StatusCode)
	var inv map[string]interface{}
	s.NoError(testutils.DecodeJSON(respCreate, &inv))
	invoiceID, _ := inv["id"].(string)
	s.NotEmpty(invoiceID)

	respPay, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/invoices/"+invoiceID+"/pay", nil, token)
	s.NoError(err)
	defer respPay.Body.Close()
	s.Equal(http.StatusNoContent, respPay.StatusCode)
}

func TestBillingInvoicesSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BillingInvoicesSuite))
}
