package e2e_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// AssetsSuite tests /v1/accounting/assets endpoints.
type AssetsSuite struct {
	suite.Suite
	helper *AccountingTestHelper
}

func (s *AssetsSuite) SetupSuite() {
	s.helper = NewAccountingTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *AssetsSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "assets")
}

func (s *AssetsSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "assets")
}

func (s *AssetsSuite) TestAssets_CRUD_Admin() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	purchasedAt := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	createPayload := map[string]interface{}{
		"name":        "MacBook Pro",
		"type":        "equipment",
		"value":       "1200.50",
		"purchasedAt": purchasedAt,
		"note":        "<script>alert('x')</script>' OR '1'='1",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/assets", createPayload, ws.AdminToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &created))
	s.Contains(created, "id")
	s.NotEmpty(created["id"])
	s.Equal(ws.Business.ID, created["businessId"])
	s.Equal("MacBook Pro", created["name"])
	s.Equal("equipment", created["type"])
	s.Equal("aed", created["currency"])

	assetID := created["id"].(string)

	// Get
	getResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/assets/"+assetID, nil, ws.AdminToken)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusOK, getResp.StatusCode)

	var got map[string]interface{}
	s.NoError(testutils.DecodeJSON(getResp, &got))
	s.Equal(assetID, got["id"])

	// List
	listResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/assets?page=1&pageSize=10", nil, ws.AdminToken)
	s.NoError(err)
	defer listResp.Body.Close()
	s.Equal(http.StatusOK, listResp.StatusCode)

	var listBody map[string]interface{}
	s.NoError(testutils.DecodeJSON(listResp, &listBody))
	s.Contains(listBody, "items")
	items, ok := listBody["items"].([]interface{})
	s.True(ok)
	s.Len(items, 1)

	// Update
	updatePayload := map[string]interface{}{
		"name":  "MacBook Pro Updated",
		"value": "999.99",
		"note":  "updated",
	}
	updResp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/accounting/assets/"+assetID, updatePayload, ws.AdminToken)
	s.NoError(err)
	defer updResp.Body.Close()
	s.Equal(http.StatusOK, updResp.StatusCode)

	var updated map[string]interface{}
	s.NoError(testutils.DecodeJSON(updResp, &updated))
	s.Equal("MacBook Pro Updated", updated["name"])

	// Delete
	delResp, err := s.helper.Client.AuthenticatedRequest("DELETE", "/v1/accounting/assets/"+assetID, nil, ws.AdminToken)
	s.NoError(err)
	defer delResp.Body.Close()
	s.Equal(http.StatusNoContent, delResp.StatusCode)

	// Get after delete => 404
	get2Resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/assets/"+assetID, nil, ws.AdminToken)
	s.NoError(err)
	defer get2Resp.Body.Close()
	s.Equal(http.StatusNotFound, get2Resp.StatusCode)
}

func (s *AssetsSuite) TestAssets_Permissions_MemberCanViewButCannotManage() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	purchasedAt := time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC)
	createPayload := map[string]interface{}{
		"name":        "Camera",
		"type":        "equipment",
		"value":       "200.00",
		"purchasedAt": purchasedAt,
	}

	// member cannot create
	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/assets", createPayload, ws.MemberToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusForbidden, resp.StatusCode)

	// admin creates
	resp2, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/assets", createPayload, ws.AdminToken)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusCreated, resp2.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &created))
	assetID := created["id"].(string)

	// member can list
	listResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/assets", nil, ws.MemberToken)
	s.NoError(err)
	defer listResp.Body.Close()
	s.Equal(http.StatusOK, listResp.StatusCode)

	// member can get
	getResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/assets/"+assetID, nil, ws.MemberToken)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusOK, getResp.StatusCode)

	// member cannot delete
	delResp, err := s.helper.Client.AuthenticatedRequest("DELETE", "/v1/accounting/assets/"+assetID, nil, ws.MemberToken)
	s.NoError(err)
	defer delResp.Body.Close()
	s.Equal(http.StatusForbidden, delResp.StatusCode)
}

func (s *AssetsSuite) TestAssets_Validation() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/accounting/assets", map[string]interface{}{}, ws.AdminToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)

	// invalid query param
	resp2, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/accounting/assets?page=-1", nil, ws.AdminToken)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusBadRequest, resp2.StatusCode)
}

func TestAssetsSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(AssetsSuite))
}
