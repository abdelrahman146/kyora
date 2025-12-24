package e2e_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// InvestmentsSuite tests /v1/businesses/:businessDescriptor/accounting/investments endpoints.
type InvestmentsSuite struct {
	suite.Suite
	helper *AccountingTestHelper
}

func (s *InvestmentsSuite) SetupSuite() {
	s.helper = NewAccountingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *InvestmentsSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "investments")
}

func (s *InvestmentsSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "investments")
}

func (s *InvestmentsSuite) TestInvestments_CRUD_Admin() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	date := time.Date(2025, 3, 4, 0, 0, 0, 0, time.UTC)
	createPayload := map[string]interface{}{
		"investorId": ws.Admin.ID,
		"amount":     "500.00",
		"investedAt": date,
		"note":       "seed capital ' OR '1'='1",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/investments", createPayload, ws.AdminToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &created))
	s.NotEmpty(created["id"])
	s.Equal(ws.Business.ID, created["businessId"])

	investmentID, ok := created["id"].(string)
	s.True(ok)

	getResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/investments/"+investmentID, nil, ws.AdminToken)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusOK, getResp.StatusCode)

	listResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/investments?page=1&pageSize=10", nil, ws.AdminToken)
	s.NoError(err)
	defer listResp.Body.Close()
	s.Equal(http.StatusOK, listResp.StatusCode)

	updatePayload := map[string]interface{}{
		"note":   "updated",
		"amount": "600.00",
	}
	updResp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/investments/"+investmentID, updatePayload, ws.AdminToken)
	s.NoError(err)
	defer updResp.Body.Close()
	s.Equal(http.StatusOK, updResp.StatusCode)

	delResp, err := s.helper.Client.AuthenticatedRequest("DELETE", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/investments/"+investmentID, nil, ws.AdminToken)
	s.NoError(err)
	defer delResp.Body.Close()
	s.Equal(http.StatusNoContent, delResp.StatusCode)
}

func (s *InvestmentsSuite) TestInvestments_Permissions_MemberCanViewButCannotManage() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	date := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
	createPayload := map[string]interface{}{"investorId": ws.Admin.ID, "amount": "100.00", "investedAt": date}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/investments", createPayload, ws.MemberToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusForbidden, resp.StatusCode)

	resp2, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/investments", createPayload, ws.AdminToken)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Require().Equal(http.StatusCreated, resp2.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &created))
	investmentID, ok := created["id"].(string)
	s.True(ok)

	listResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/investments", nil, ws.MemberToken)
	s.NoError(err)
	defer listResp.Body.Close()
	s.Equal(http.StatusOK, listResp.StatusCode)

	getResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/investments/"+investmentID, nil, ws.MemberToken)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusOK, getResp.StatusCode)
}

func TestInvestmentsSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(InvestmentsSuite))
}
