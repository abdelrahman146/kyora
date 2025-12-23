package e2e_test

import (
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
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
	testutils.TruncateTables(testEnv.Database,
		"stripe_events",
		"subscriptions",
		"plans",
		"users",
		"workspaces",
	)
}

func (s *BillingInvoicesSuite) TearDownTest() {
	s.SetupTest()
}

func (s *BillingInvoicesSuite) TestInvoices_CreateAndList_AndBOLAOnPay() {
	ctx := s.T().Context()
	_, _, token1 := s.helper.CreateTestUser(ctx, "admin1@example.com", role.RoleAdmin)
	_, _, token2 := s.helper.CreateTestUser(ctx, "admin2@example.com", role.RoleAdmin)

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

	// BOLA: second workspace must not be able to pay first workspace invoice
	respPay, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/invoices/"+invoiceID+"/pay", nil, token2)
	s.NoError(err)
	defer respPay.Body.Close()
	s.Equal(http.StatusNotFound, respPay.StatusCode)
}

func TestBillingInvoicesSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BillingInvoicesSuite))
}
