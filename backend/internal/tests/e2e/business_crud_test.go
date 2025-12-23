package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

type BusinessSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *BusinessSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *BusinessSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database,
		"users",
		"workspaces",
		"businesses",
		"subscriptions",
		"plans",
	)
	s.NoError(err)
}

func (s *BusinessSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database,
		"users",
		"workspaces",
		"businesses",
		"subscriptions",
		"plans",
	)
	s.NoError(err)
}

func (s *BusinessSuite) TestBusiness_CRUD_AndSecurity() {
	ctx := context.Background()

	adminUser, adminWS, adminToken, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NotEmpty(adminUser.ID)
	s.NotEmpty(adminWS.ID)
	s.NotEmpty(adminToken)

	s.NoError(s.helper.CreateTestSubscription(ctx, adminWS.ID))

	// Descriptor availability (true)
	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/descriptor/availability?descriptor=test-business", nil, adminToken)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	var avail map[string]any
	s.NoError(testutils.DecodeJSON(resp, &avail))
	s.Len(avail, 1)
	s.Equal(true, avail["available"])

	// Create
	createPayload := map[string]any{
		"name":        "Test Business",
		"descriptor":  "test-business",
		"countryCode": "eg",
		"currency":    "usd",
	}
	resp, err = s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", createPayload, adminToken)
	s.NoError(err)
	s.Equal(http.StatusCreated, resp.StatusCode)

	var createResult map[string]any
	s.NoError(testutils.DecodeJSON(resp, &createResult))
	s.Len(createResult, 1)
	s.Contains(createResult, "business")

	biz := createResult["business"].(map[string]any)
	s.Contains(biz, "id")
	s.Contains(biz, "workspaceId")
	s.Contains(biz, "descriptor")
	s.Contains(biz, "name")
	s.Contains(biz, "countryCode")
	s.Contains(biz, "currency")
	s.Contains(biz, "vatRate")
	s.Contains(biz, "safetyBuffer")
	s.Contains(biz, "establishedAt")
	s.Contains(biz, "createdAt")
	s.Contains(biz, "updatedAt")
	s.NotContains(biz, "workspace")

	businessID := biz["id"].(string)
	s.NotEmpty(businessID)

	// Descriptor availability (false)
	resp, err = s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/descriptor/availability?descriptor=test-business", nil, adminToken)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	avail = map[string]any{}
	s.NoError(testutils.DecodeJSON(resp, &avail))
	s.Equal(false, avail["available"])

	// Get by ID
	resp, err = s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+businessID, nil, adminToken)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	var getResult map[string]any
	s.NoError(testutils.DecodeJSON(resp, &getResult))
	s.Len(getResult, 1)
	s.Contains(getResult, "business")
	gotBiz := getResult["business"].(map[string]any)
	s.Equal("test-business", gotBiz["descriptor"])
	s.Equal("Test Business", gotBiz["name"])
	s.Equal(adminWS.ID, gotBiz["workspaceId"])

	// Get by descriptor
	resp, err = s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/descriptor/test-business", nil, adminToken)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	var getByDesc map[string]any
	s.NoError(testutils.DecodeJSON(resp, &getByDesc))
	s.Len(getByDesc, 1)
	gotBiz2 := getByDesc["business"].(map[string]any)
	s.Equal(businessID, gotBiz2["id"])

	// Update (PATCH)
	updatePayload := map[string]any{"name": "Updated Business"}
	resp, err = s.helper.Client.AuthenticatedRequest("PATCH", "/v1/businesses/"+businessID, updatePayload, adminToken)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	var updateResult map[string]any
	s.NoError(testutils.DecodeJSON(resp, &updateResult))
	s.Len(updateResult, 1)
	updatedBiz := updateResult["business"].(map[string]any)
	s.Equal("Updated Business", updatedBiz["name"])

	// Archive
	resp, err = s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+businessID+"/archive", nil, adminToken)
	s.NoError(err)
	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Unarchive
	resp, err = s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+businessID+"/unarchive", nil, adminToken)
	s.NoError(err)
	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Cross-tenant isolation: other workspace cannot access by ID
	_, otherWS, otherToken, err := s.helper.CreateTestUser(ctx, "other@example.com", "ValidPassword123!", "Other", "User", role.RoleAdmin)
	s.NoError(err)
	s.NotEmpty(otherWS.ID)
	s.NotEmpty(otherToken)

	resp, err = s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+businessID, nil, otherToken)
	s.NoError(err)
	s.Equal(http.StatusNotFound, resp.StatusCode)

	// BFLA: member cannot create
	_, memberWS, memberToken, err := s.helper.CreateTestUser(ctx, "member@example.com", "ValidPassword123!", "Member", "User", role.RoleUser)
	s.NoError(err)
	s.NotEmpty(memberWS.ID)
	s.NoError(s.helper.CreateTestSubscription(ctx, memberWS.ID))

	resp, err = s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", createPayload, memberToken)
	s.NoError(err)
	s.Equal(http.StatusForbidden, resp.StatusCode)

	// Delete
	resp, err = s.helper.Client.AuthenticatedRequest("DELETE", "/v1/businesses/"+businessID, nil, adminToken)
	s.NoError(err)
	s.Equal(http.StatusNoContent, resp.StatusCode)
}

func TestBusinessSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BusinessSuite))
}
