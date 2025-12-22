package e2e_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/accounting"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// RecurringExpensesSuite tests /v1/accounting/recurring-expenses endpoints.
type RecurringExpensesSuite struct {
	suite.Suite
	helper *AccountingTestHelper
}

func (s *RecurringExpensesSuite) SetupSuite() {
	s.helper = NewAccountingTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *RecurringExpensesSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "recurring_expenses", "expenses")
}

func (s *RecurringExpensesSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "recurring_expenses", "expenses")
}

func (s *RecurringExpensesSuite) TestRecurringExpenses_Create_List_StatusTransition() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	payload := map[string]interface{}{
		"frequency":          "monthly",
		"recurringStartDate": startDate,
		"category":           "rent",
		"amount":             "100.00",
		"note":               "rent",
	}

	createResp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/recurring-expenses", payload, ws.AdminToken)
	s.NoError(err)
	defer createResp.Body.Close()
	s.Require().Equal(http.StatusCreated, createResp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(createResp, &created))
	recID, ok := created["id"].(string)
	s.True(ok)
	s.Equal("active", created["status"])

	// List should include it
	listResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/recurring-expenses?page=1&pageSize=10", nil, ws.AdminToken)
	s.NoError(err)
	defer listResp.Body.Close()
	s.Equal(http.StatusOK, listResp.StatusCode)

	// Transition to paused
	statusPayload := map[string]interface{}{"status": string(accounting.RecurringExpenseStatusPaused)}
	statusResp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/accounting/recurring-expenses/"+recID+"/status", statusPayload, ws.AdminToken)
	s.NoError(err)
	defer statusResp.Body.Close()
	s.Equal(http.StatusOK, statusResp.StatusCode)

	var statusBody map[string]interface{}
	s.NoError(testutils.DecodeJSON(statusResp, &statusBody))
	s.Equal("paused", statusBody["status"])

	// Invalid status value should be rejected by validation (400)
	invalidPayload := map[string]interface{}{"status": "invalid"}
	invalidResp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/accounting/recurring-expenses/"+recID+"/status", invalidPayload, ws.AdminToken)
	s.NoError(err)
	defer invalidResp.Body.Close()
	s.Equal(http.StatusBadRequest, invalidResp.StatusCode)
}

func (s *RecurringExpensesSuite) TestRecurringExpenses_Occurrences() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	startDate := time.Now().UTC().AddDate(0, -2, 0)
	payload := map[string]interface{}{
		"frequency":                    "monthly",
		"recurringStartDate":           startDate,
		"category":                     "software",
		"amount":                       "10.00",
		"autoCreateHistoricalExpenses": true,
	}

	createResp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/recurring-expenses", payload, ws.AdminToken)
	s.NoError(err)
	defer createResp.Body.Close()
	s.Require().Equal(http.StatusCreated, createResp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(createResp, &created))
	recID, ok := created["id"].(string)
	s.True(ok)

	occResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/recurring-expenses/"+recID+"/occurrences", nil, ws.AdminToken)
	s.NoError(err)
	defer occResp.Body.Close()
	s.Equal(http.StatusOK, occResp.StatusCode)

	var occBody []interface{}
	s.NoError(testutils.DecodeJSON(occResp, &occBody))
	s.NotEmpty(occBody)
}

func (s *RecurringExpensesSuite) TestRecurringExpenses_Permissions_MemberViewOnly() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	payload := map[string]interface{}{"frequency": "monthly", "recurringStartDate": startDate, "category": "rent", "amount": "10.00"}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/recurring-expenses", payload, ws.MemberToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusForbidden, resp.StatusCode)

	// member can list
	listResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/recurring-expenses", nil, ws.MemberToken)
	s.NoError(err)
	defer listResp.Body.Close()
	s.Equal(http.StatusOK, listResp.StatusCode)
}

func TestRecurringExpensesSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(RecurringExpensesSuite))
}
