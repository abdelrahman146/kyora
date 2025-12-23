package e2e_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// ExpensesSuite tests /v1/businesses/:businessDescriptor/accounting/expenses endpoints.
type ExpensesSuite struct {
	suite.Suite
	helper *AccountingTestHelper
}

func (s *ExpensesSuite) SetupSuite() {
	s.helper = NewAccountingTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *ExpensesSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "expenses")
}

func (s *ExpensesSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "expenses")
}

func (s *ExpensesSuite) TestExpenses_CRUD_Admin() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	date := time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC)
	createPayload := map[string]interface{}{
		"category":   "supplies",
		"type":       "one_time",
		"amount":     "25.50",
		"occurredOn": date,
		"note":       "<img src=x onerror=alert(1)>",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/expenses", createPayload, ws.AdminToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &created))
	expenseID, ok := created["id"].(string)
	s.True(ok)
	s.Equal(ws.Business.ID, created["businessId"])

	getResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/expenses/"+expenseID, nil, ws.AdminToken)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusOK, getResp.StatusCode)

	listResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/expenses?page=1&pageSize=10", nil, ws.AdminToken)
	s.NoError(err)
	defer listResp.Body.Close()
	s.Equal(http.StatusOK, listResp.StatusCode)

	updatePayload := map[string]interface{}{"note": "updated", "amount": "30.00"}
	updResp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/expenses/"+expenseID, updatePayload, ws.AdminToken)
	s.NoError(err)
	defer updResp.Body.Close()
	s.Equal(http.StatusOK, updResp.StatusCode)

	delResp, err := s.helper.Client.AuthenticatedRequest("DELETE", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/expenses/"+expenseID, nil, ws.AdminToken)
	s.NoError(err)
	defer delResp.Body.Close()
	s.Equal(http.StatusNoContent, delResp.StatusCode)
}

func (s *ExpensesSuite) TestExpenses_Permissions_MemberCanViewButCannotManage() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	date := time.Date(2025, 4, 10, 0, 0, 0, 0, time.UTC)
	createPayload := map[string]interface{}{"category": "shipping", "type": "one_time", "amount": "10.00", "occurredOn": date}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/expenses", createPayload, ws.MemberToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusForbidden, resp.StatusCode)

	resp2, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/expenses", createPayload, ws.AdminToken)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Require().Equal(http.StatusCreated, resp2.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &created))
	expenseID, ok := created["id"].(string)
	s.True(ok)

	listResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/expenses", nil, ws.MemberToken)
	s.NoError(err)
	defer listResp.Body.Close()
	s.Equal(http.StatusOK, listResp.StatusCode)

	getResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/expenses/"+expenseID, nil, ws.MemberToken)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusOK, getResp.StatusCode)

	delResp, err := s.helper.Client.AuthenticatedRequest("DELETE", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/expenses/"+expenseID, nil, ws.MemberToken)
	s.NoError(err)
	defer delResp.Body.Close()
	s.Equal(http.StatusForbidden, delResp.StatusCode)
}

func TestExpensesSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(ExpensesSuite))
}
