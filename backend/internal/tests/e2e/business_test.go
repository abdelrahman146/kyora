package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
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
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "subscriptions", "plans")
	s.NoError(err)
}

func (s *BusinessSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "subscriptions", "plans")
	s.NoError(err)
}

func (s *BusinessSuite) createBusiness(ctx context.Context, workspaceID, descriptor string) string {
	bizRepo := database.NewRepository[business.Business](testEnv.Database)
	biz := &business.Business{
		WorkspaceID:  workspaceID,
		Descriptor:   descriptor,
		Name:         "Test Business",
		CountryCode:  "EG",
		Currency:     "USD",
		VatRate:      decimal.NewFromFloat(0.14),
		SafetyBuffer: decimal.NewFromFloat(100),
	}
	s.NoError(bizRepo.CreateOne(ctx, biz))
	return biz.ID
}

func (s *BusinessSuite) TestCreateBusiness_Success() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))

	payload := map[string]interface{}{
		"name":          "Test Business",
		"descriptor":    "test-business",
		"countryCode":   "eg",
		"currency":      "usd",
		"vatRate":       "0.14",
		"safetyBuffer":  "100.50",
		"establishedAt": "2020-01-01",
	}
	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Len(result, 1)
	s.Contains(result, "business")

	biz := result["business"].(map[string]interface{})
	s.NotEmpty(biz["id"])
	s.Equal(ws.ID, biz["workspaceId"])
	s.Equal("test-business", biz["descriptor"])
	s.Equal("Test Business", biz["name"])
	s.Equal("EG", biz["countryCode"])
	s.Equal("USD", biz["currency"])
	s.Equal("0.14", biz["vatRate"])
	s.Equal("100.50", biz["safetyBuffer"])
	s.Contains(biz, "establishedAt")
	s.NotContains(biz, "workspace")

	// Verify DB state
	bizRepo := database.NewRepository[business.Business](testEnv.Database)
	dbBiz, err := bizRepo.FindByID(ctx, biz["id"].(string))
	s.NoError(err)
	s.Equal("test-business", dbBiz.Descriptor)
	s.Equal("EG", dbBiz.CountryCode)
	s.Equal("USD", dbBiz.Currency)
}

func (s *BusinessSuite) TestCreateBusiness_NormalizesInputs() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))

	tests := []struct {
		descriptor  string
		countryCode string
		currency    string
		expected    struct{ d, c, cur string }
	}{
		{"TEST-Business", "eg", "usd", struct{ d, c, cur string }{"test-business", "EG", "USD"}},
		{"  test-BIZ  ", "US", "eur", struct{ d, c, cur string }{"test-biz", "US", "EUR"}},
	}

	for i, tt := range tests {
		payload := map[string]interface{}{
			"name":        fmt.Sprintf("Business %d", i),
			"descriptor":  tt.descriptor,
			"countryCode": tt.countryCode,
			"currency":    tt.currency,
		}
		resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload, token)
		s.NoError(err)
		defer resp.Body.Close()
		s.Equal(http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		s.NoError(testutils.DecodeJSON(resp, &result))
		biz := result["business"].(map[string]interface{})
		s.Equal(tt.expected.d, biz["descriptor"])
		s.Equal(tt.expected.c, biz["countryCode"])
		s.Equal(tt.expected.cur, biz["currency"])
	}
}

func (s *BusinessSuite) TestCreateBusiness_ValidationErrors() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))

	tests := []struct {
		name    string
		payload map[string]interface{}
	}{
		{"missing name", map[string]interface{}{"descriptor": "test", "countryCode": "eg", "currency": "usd"}},
		{"missing descriptor", map[string]interface{}{"name": "Test", "countryCode": "eg", "currency": "usd"}},
		{"missing countryCode", map[string]interface{}{"name": "Test", "descriptor": "test", "currency": "usd"}},
		{"missing currency", map[string]interface{}{"name": "Test", "descriptor": "test", "countryCode": "eg"}},
		{"invalid countryCode", map[string]interface{}{"name": "Test", "descriptor": "test1", "countryCode": "e", "currency": "usd"}},
		{"invalid currency", map[string]interface{}{"name": "Test", "descriptor": "test2", "countryCode": "eg", "currency": "us"}},
		{"invalid descriptor format", map[string]interface{}{"name": "Test", "descriptor": "test_business", "countryCode": "eg", "currency": "usd"}},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", tt.payload, token)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func (s *BusinessSuite) TestCreateBusiness_DescriptorUniquenessInWorkspace() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))

	payload := map[string]interface{}{"name": "Business 1", "descriptor": "test-business", "countryCode": "eg", "currency": "usd"}
	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload, token)
	s.NoError(err)
	resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	// Attempt duplicate
	payload2 := map[string]interface{}{"name": "Business 2", "descriptor": "test-business", "countryCode": "eg", "currency": "usd"}
	resp2, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload2, token)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusConflict, resp2.StatusCode)
}

func (s *BusinessSuite) TestCreateBusiness_DescriptorAvailableAcrossWorkspaces() {
	ctx := context.Background()
	_, ws1, token1, err := s.helper.CreateTestUser(ctx, "admin1@example.com", "ValidPassword123!", "Admin", "One", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws1.ID))

	_, ws2, token2, err := s.helper.CreateTestUser(ctx, "admin2@example.com", "ValidPassword123!", "Admin", "Two", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws2.ID))

	payload := map[string]interface{}{"name": "Test Business", "descriptor": "shared", "countryCode": "eg", "currency": "usd"}
	
	resp1, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload, token1)
	s.NoError(err)
	resp1.Body.Close()
	s.Equal(http.StatusCreated, resp1.StatusCode)

	resp2, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload, token2)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusCreated, resp2.StatusCode)
}

func (s *BusinessSuite) TestCreateBusiness_RequiresAuth() {
	payload := map[string]interface{}{"name": "Test", "descriptor": "test", "countryCode": "eg", "currency": "usd"}
	resp, err := s.helper.Client.Post("/v1/businesses", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *BusinessSuite) TestCreateBusiness_RequiresManagePermission() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "member@example.com", "ValidPassword123!", "Member", "User", role.RoleUser)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))

	payload := map[string]interface{}{"name": "Test", "descriptor": "test", "countryCode": "eg", "currency": "usd"}
	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusForbidden, resp.StatusCode)
}

func (s *BusinessSuite) TestCreateBusiness_EnforcesPlanLimit() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)

	// Create plan with MaxBusinesses=1
	plan, err := testutils.CreateTestPlan(ctx, testEnv.Database, "limited-plan")
	s.NoError(err)
	plan.Limits.MaxBusinesses = 1
	
	// Update plan limits
	planRepo := database.NewRepository[billing.Plan](testEnv.Database)
	s.NoError(planRepo.UpdateOne(ctx, plan))
	
	// Create subscription with limited plan
	subRepo := database.NewRepository[billing.Subscription](testEnv.Database)
	sub := &billing.Subscription{
		WorkspaceID: ws.ID,
		PlanID:      plan.ID,
		Status:      billing.SubscriptionStatusActive,
	}
	s.NoError(subRepo.CreateOne(ctx, sub))

	// First business succeeds
	payload1 := map[string]interface{}{"name": "Business 1", "descriptor": "business-1", "countryCode": "eg", "currency": "usd"}
	resp1, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload1, token)
	s.NoError(err)
	resp1.Body.Close()
	s.Equal(http.StatusCreated, resp1.StatusCode)

	// Second business fails due to plan limit
	payload2 := map[string]interface{}{"name": "Business 2", "descriptor": "business-2", "countryCode": "eg", "currency": "usd"}
	resp2, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload2, token)
	s.NoError(err)
	defer resp2.Body.Close()
	s.GreaterOrEqual(resp2.StatusCode, 400)
}

func (s *BusinessSuite) TestGetBusiness_ById() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)

	bizID := s.createBusiness(ctx, ws.ID, "test-business")

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+bizID, nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	biz := result["business"].(map[string]interface{})
	s.Equal(bizID, biz["id"])
	s.Equal("test-business", biz["descriptor"])
}

func (s *BusinessSuite) TestGetBusiness_NotFound() {
	ctx := context.Background()
	_, _, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/bus_nonexistent", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *BusinessSuite) TestGetBusiness_CrossWorkspaceIsolation() {
	ctx := context.Background()
	_, ws1, _, err := s.helper.CreateTestUser(ctx, "admin1@example.com", "ValidPassword123!", "Admin", "One", role.RoleAdmin)
	s.NoError(err)
	bizID := s.createBusiness(ctx, ws1.ID, "business-1")

	_, _, token2, err := s.helper.CreateTestUser(ctx, "admin2@example.com", "ValidPassword123!", "Admin", "Two", role.RoleAdmin)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+bizID, nil, token2)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *BusinessSuite) TestGetBusinessByDescriptor_Success() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	bizID := s.createBusiness(ctx, ws.ID, "test-business")

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/descriptor/test-business", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	biz := result["business"].(map[string]interface{})
	s.Equal(bizID, biz["id"])
}

func (s *BusinessSuite) TestGetBusinessByDescriptor_NormalizesQuery() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.createBusiness(ctx, ws.ID, "test-business")

	// All variations should find the business
	for _, desc := range []string{"TEST-BUSINESS", "Test-Business", "  test-business  "} {
		resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/descriptor/"+desc, nil, token)
		s.NoError(err)
		resp.Body.Close()
		s.Equal(http.StatusOK, resp.StatusCode)
	}
}

func (s *BusinessSuite) TestListBusinesses_ReturnsWorkspaceBusinessesOnly() {
	ctx := context.Background()
	_, ws1, token1, err := s.helper.CreateTestUser(ctx, "admin1@example.com", "ValidPassword123!", "Admin", "One", role.RoleAdmin)
	s.NoError(err)
	s.createBusiness(ctx, ws1.ID, "business-1")
	s.createBusiness(ctx, ws1.ID, "business-2")

	_, ws2, token2, err := s.helper.CreateTestUser(ctx, "admin2@example.com", "ValidPassword123!", "Admin", "Two", role.RoleAdmin)
	s.NoError(err)
	s.createBusiness(ctx, ws2.ID, "business-3")

	resp1, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses", nil, token1)
	s.NoError(err)
	defer resp1.Body.Close()
	var result1 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp1, &result1))
	businesses1 := result1["businesses"].([]interface{})
	s.Len(businesses1, 2)

	resp2, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses", nil, token2)
	s.NoError(err)
	defer resp2.Body.Close()
	var result2 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &result2))
	businesses2 := result2["businesses"].([]interface{})
	s.Len(businesses2, 1)
}

func (s *BusinessSuite) TestUpdateBusiness_PartialUpdate() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	bizID := s.createBusiness(ctx, ws.ID, "test-business")

	payload := map[string]interface{}{"name": "Updated Name"}
	resp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/businesses/"+bizID, payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	biz := result["business"].(map[string]interface{})
	s.Equal("Updated Name", biz["name"])
	s.Equal("test-business", biz["descriptor"]) // Unchanged
}

func (s *BusinessSuite) TestUpdateBusiness_ChangesDescriptor() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	bizID := s.createBusiness(ctx, ws.ID, "test-business")

	payload := map[string]interface{}{"descriptor": "NEW-Descriptor"}
	resp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/businesses/"+bizID, payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	biz := result["business"].(map[string]interface{})
	s.Equal("new-descriptor", biz["descriptor"]) // Normalized
}

func (s *BusinessSuite) TestUpdateBusiness_DescriptorConflict() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.createBusiness(ctx, ws.ID, "business-1")
	bizID2 := s.createBusiness(ctx, ws.ID, "business-2")

	payload := map[string]interface{}{"descriptor": "business-1"}
	resp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/businesses/"+bizID2, payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusConflict, resp.StatusCode)
}

func (s *BusinessSuite) TestUpdateBusiness_CrossWorkspaceIsolation() {
	ctx := context.Background()
	_, ws1, _, err := s.helper.CreateTestUser(ctx, "admin1@example.com", "ValidPassword123!", "Admin", "One", role.RoleAdmin)
	s.NoError(err)
	bizID := s.createBusiness(ctx, ws1.ID, "business-1")

	_, _, token2, err := s.helper.CreateTestUser(ctx, "admin2@example.com", "ValidPassword123!", "Admin", "Two", role.RoleAdmin)
	s.NoError(err)

	payload := map[string]interface{}{"name": "Hacked"}
	resp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/businesses/"+bizID, payload, token2)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *BusinessSuite) TestArchiveBusiness_SetsTimestamp() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	bizID := s.createBusiness(ctx, ws.ID, "test-business")

	resp, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/archive", bizID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify archived state
	bizRepo := database.NewRepository[business.Business](testEnv.Database)
	dbBiz, err := bizRepo.FindByID(ctx, bizID)
	s.NoError(err)
	s.NotNil(dbBiz.ArchivedAt)
}

func (s *BusinessSuite) TestArchiveBusiness_Idempotent() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	bizID := s.createBusiness(ctx, ws.ID, "test-business")

	// Archive twice
	resp1, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/archive", bizID), nil, token)
	s.NoError(err)
	resp1.Body.Close()
	s.Equal(http.StatusNoContent, resp1.StatusCode)

	bizRepo := database.NewRepository[business.Business](testEnv.Database)
	dbBiz1, _ := bizRepo.FindByID(ctx, bizID)
	firstTime := dbBiz1.ArchivedAt

	resp2, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/archive", bizID), nil, token)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusNoContent, resp2.StatusCode)

	dbBiz2, _ := bizRepo.FindByID(ctx, bizID)
	s.Equal(firstTime, dbBiz2.ArchivedAt) // Timestamp unchanged
}

func (s *BusinessSuite) TestUnarchiveBusiness_ClearsTimestamp() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	bizID := s.createBusiness(ctx, ws.ID, "test-business")

	// Archive then unarchive
	s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/archive", bizID), nil, token)
	resp, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/unarchive", bizID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNoContent, resp.StatusCode)

	bizRepo := database.NewRepository[business.Business](testEnv.Database)
	dbBiz, err := bizRepo.FindByID(ctx, bizID)
	s.NoError(err)
	s.Nil(dbBiz.ArchivedAt)
}

func (s *BusinessSuite) TestDeleteBusiness_RemovesPermanently() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	bizID := s.createBusiness(ctx, ws.ID, "test-business")

	resp, err := s.helper.Client.AuthenticatedRequest("DELETE", "/v1/businesses/"+bizID, nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify deletion
	bizRepo := database.NewRepository[business.Business](testEnv.Database)
	_, err = bizRepo.FindByID(ctx, bizID)
	s.True(database.IsRecordNotFound(err))
}

func (s *BusinessSuite) TestCheckDescriptorAvailability_BeforeAndAfterCreation() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))

	// Check availability before creation
	resp1, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/descriptor/availability?descriptor=test-business", nil, token)
	s.NoError(err)
	defer resp1.Body.Close()
	var result1 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp1, &result1))
	s.Equal(true, result1["available"])

	// Create business
	payload := map[string]interface{}{"name": "Test", "descriptor": "test-business", "countryCode": "eg", "currency": "usd"}
	s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload, token)

	// Check availability after creation
	resp2, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/descriptor/availability?descriptor=test-business", nil, token)
	s.NoError(err)
	defer resp2.Body.Close()
	var result2 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &result2))
	s.Equal(false, result2["available"])
}

func (s *BusinessSuite) TestCheckDescriptorAvailability_CrossWorkspace() {
	ctx := context.Background()
	_, ws1, token1, err := s.helper.CreateTestUser(ctx, "admin1@example.com", "ValidPassword123!", "Admin", "One", role.RoleAdmin)
	s.NoError(err)
	s.createBusiness(ctx, ws1.ID, "shared")

	_, _, token2, err := s.helper.CreateTestUser(ctx, "admin2@example.com", "ValidPassword123!", "Admin", "Two", role.RoleAdmin)
	s.NoError(err)

	// Should be unavailable in ws1
	resp1, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/descriptor/availability?descriptor=shared", nil, token1)
	s.NoError(err)
	defer resp1.Body.Close()
	var result1 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp1, &result1))
	s.Equal(false, result1["available"])

	// Should be available in ws2
	resp2, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/descriptor/availability?descriptor=shared", nil, token2)
	s.NoError(err)
	defer resp2.Body.Close()
	var result2 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &result2))
	s.Equal(true, result2["available"])
}

func TestBusinessSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BusinessSuite))
}
