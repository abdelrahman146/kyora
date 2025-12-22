package e2e_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// SummarySuite tests /v1/accounting/summary endpoint.
type SummarySuite struct {
	suite.Suite
	helper *AccountingTestHelper
}

func (s *SummarySuite) SetupSuite() {
	s.helper = NewAccountingTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *SummarySuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "investments", "withdrawals", "expenses")
}

func (s *SummarySuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "investments", "withdrawals", "expenses")
}

func (s *SummarySuite) TestSummary_EmptyWorkspace() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/summary", nil, ws.AdminToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &body))
	s.Contains(body, "totalAssetValue")
	s.Contains(body, "totalInvestments")
	s.Contains(body, "totalWithdrawals")
	s.Contains(body, "totalExpenses")
	s.Contains(body, "safeToDrawAmount")
	s.Contains(body, "currency")
}

func (s *SummarySuite) TestSummary_ComputedFields_WithInvestmentsWithdrawalsExpenses() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	now := time.Now().UTC()
	inv := map[string]interface{}{"investorId": ws.Admin.ID, "amount": "1000.00", "investedAt": now}
	exp := map[string]interface{}{"category": "software", "type": "one_time", "amount": "100.00", "occurredOn": now}
	wd := map[string]interface{}{"withdrawerId": ws.Admin.ID, "amount": "50.00", "withdrawnAt": now}

	resp1, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/investments", inv, ws.AdminToken)
	s.NoError(err)
	resp1.Body.Close()
	resp2, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/expenses", exp, ws.AdminToken)
	s.NoError(err)
	resp2.Body.Close()
	resp3, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/withdrawals", wd, ws.AdminToken)
	s.NoError(err)
	resp3.Body.Close()

	summaryResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/summary", nil, ws.AdminToken)
	s.NoError(err)
	defer summaryResp.Body.Close()
	s.Equal(http.StatusOK, summaryResp.StatusCode)

	var body map[string]interface{}
	s.NoError(testutils.DecodeJSON(summaryResp, &body))

	s.Equal("1000", body["totalInvestments"])
	s.Equal("50", body["totalWithdrawals"])
	s.Equal("100", body["totalExpenses"])
	// If SafetyBuffer is not set, it defaults to last-30-days expenses
	s.Equal("750", body["safeToDrawAmount"])
}

func (s *SummarySuite) TestSummary_DateRange_AppliesToSafeToDrawAndTotals() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	within := time.Date(2025, 12, 10, 12, 0, 0, 0, time.UTC)
	outside := time.Date(2025, 10, 10, 12, 0, 0, 0, time.UTC)

	invWithin := map[string]interface{}{"investorId": ws.Admin.ID, "amount": "1000.00", "investedAt": within}
	expWithin := map[string]interface{}{"category": "software", "type": "one_time", "amount": "100.00", "occurredOn": within}
	wdWithin := map[string]interface{}{"withdrawerId": ws.Admin.ID, "amount": "50.00", "withdrawnAt": within}

	expOutside := map[string]interface{}{"category": "software", "type": "one_time", "amount": "999.00", "occurredOn": outside}
	wdOutside := map[string]interface{}{"withdrawerId": ws.Admin.ID, "amount": "999.00", "withdrawnAt": outside}

	resp1, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/investments", invWithin, ws.AdminToken)
	s.NoError(err)
	resp1.Body.Close()
	resp2, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/expenses", expWithin, ws.AdminToken)
	s.NoError(err)
	resp2.Body.Close()
	resp3, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/withdrawals", wdWithin, ws.AdminToken)
	s.NoError(err)
	resp3.Body.Close()
	resp4, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/expenses", expOutside, ws.AdminToken)
	s.NoError(err)
	resp4.Body.Close()
	resp5, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/withdrawals", wdOutside, ws.AdminToken)
	s.NoError(err)
	resp5.Body.Close()

	from := "2025-12-01"
	to := "2025-12-31"
	summaryResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/summary?from="+from+"&to="+to, nil, ws.AdminToken)
	s.NoError(err)
	defer summaryResp.Body.Close()
	s.Equal(http.StatusOK, summaryResp.StatusCode)

	var body map[string]interface{}
	s.NoError(testutils.DecodeJSON(summaryResp, &body))

	s.Equal("1000", body["totalInvestments"])
	s.Equal("50", body["totalWithdrawals"])
	s.Equal("100", body["totalExpenses"])
	// SafetyBuffer defaults to last-30-days expenses anchored to `to`
	s.Equal("750", body["safeToDrawAmount"])
	s.Equal(from, body["from"])
	s.Equal(to, body["to"])
}

func (s *SummarySuite) TestSummary_Permissions_MemberCanView() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/summary", nil, ws.MemberToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)
}

func TestSummarySuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(SummarySuite))
}
