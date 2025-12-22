package e2e_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// WithdrawalsSuite tests /v1/accounting/withdrawals endpoints.
type WithdrawalsSuite struct {
	suite.Suite
	helper *AccountingTestHelper
}

func (s *WithdrawalsSuite) SetupSuite() {
	s.helper = NewAccountingTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *WithdrawalsSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "withdrawals")
}

func (s *WithdrawalsSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "withdrawals")
}

func (s *WithdrawalsSuite) TestWithdrawals_CRUD_Admin() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	date := time.Date(2025, 3, 5, 0, 0, 0, 0, time.UTC)
	createPayload := map[string]interface{}{
		"amount":       "200.00",
		"withdrawerId": ws.Admin.ID,
		"withdrawnAt":  date,
		"note":         "owner draw",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/withdrawals", createPayload, ws.AdminToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &created))
	withdrawalID, ok := created["id"].(string)
	s.True(ok)

	getResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/withdrawals/"+withdrawalID, nil, ws.AdminToken)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusOK, getResp.StatusCode)

	listResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/withdrawals?page=1&pageSize=10", nil, ws.AdminToken)
	s.NoError(err)
	defer listResp.Body.Close()
	s.Equal(http.StatusOK, listResp.StatusCode)

	updatePayload := map[string]interface{}{"note": "updated", "amount": "250.00"}
	updResp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/accounting/withdrawals/"+withdrawalID, updatePayload, ws.AdminToken)
	s.NoError(err)
	defer updResp.Body.Close()
	s.Equal(http.StatusOK, updResp.StatusCode)

	delResp, err := s.helper.Client.AuthenticatedRequest("DELETE", "/v1/accounting/withdrawals/"+withdrawalID, nil, ws.AdminToken)
	s.NoError(err)
	defer delResp.Body.Close()
	s.Equal(http.StatusNoContent, delResp.StatusCode)
}

func (s *WithdrawalsSuite) TestWithdrawals_Permissions_MemberCanViewButCannotManage() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	date := time.Date(2025, 4, 2, 0, 0, 0, 0, time.UTC)
	createPayload := map[string]interface{}{"amount": "100.00", "withdrawerId": ws.Admin.ID, "withdrawnAt": date}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/withdrawals", createPayload, ws.MemberToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusForbidden, resp.StatusCode)

	resp2, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/withdrawals", createPayload, ws.AdminToken)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Require().Equal(http.StatusCreated, resp2.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &created))
	withdrawalID, ok := created["id"].(string)
	s.True(ok)

	listResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/withdrawals", nil, ws.MemberToken)
	s.NoError(err)
	defer listResp.Body.Close()
	s.Equal(http.StatusOK, listResp.StatusCode)

	getResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/withdrawals/"+withdrawalID, nil, ws.MemberToken)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusOK, getResp.StatusCode)
}

func TestWithdrawalsSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(WithdrawalsSuite))
}
