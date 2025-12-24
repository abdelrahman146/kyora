package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

type InventoryCategoriesSuite struct {
	suite.Suite
	accountHelper   *AccountTestHelper
	inventoryHelper *InventoryTestHelper
}

func (s *InventoryCategoriesSuite) SetupSuite() {
	s.accountHelper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
	s.inventoryHelper = NewInventoryTestHelper(testEnv.Database, e2eBaseURL)
}

func (s *InventoryCategoriesSuite) resetDB() {
	s.NoError(testutils.TruncateTables(testEnv.Database, inventoryTables...))
}

func (s *InventoryCategoriesSuite) SetupTest() {
	s.resetDB()
}

func (s *InventoryCategoriesSuite) TearDownTest() {
	s.resetDB()
}

func (s *InventoryCategoriesSuite) TestCreateCategory_Success_NormalizesDescriptor() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	_, err = s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	payload := map[string]interface{}{"name": "T-Shirts", "descriptor": "  TShirts  "}
	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/inventory/categories", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.NotEmpty(result["id"])
	s.Equal("T-Shirts", result["name"])
	s.Equal("tshirts", result["descriptor"])
	s.Contains(result, "businessId")

	catID := result["id"].(string)
	cat, err := s.inventoryHelper.GetCategory(ctx, catID)
	s.NoError(err)
	s.Equal("tshirts", cat.Descriptor)
}

func (s *InventoryCategoriesSuite) TestCreateCategory_DuplicateDescriptorSameBusiness() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	_, err = s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "First", "dupe")
	s.NoError(err)

	payload := map[string]interface{}{"name": "Second", "descriptor": "dupe"}
	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/inventory/categories", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusConflict, resp.StatusCode)
}

func (s *InventoryCategoriesSuite) TestCreateCategory_SameDescriptorDifferentBusiness() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz1, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "biz-1")
	s.NoError(err)
	biz2, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "biz-2")
	s.NoError(err)

	_, err = s.inventoryHelper.CreateTestCategory(ctx, biz1.ID, "Cat", "shared")
	s.NoError(err)

	payload := map[string]interface{}{"name": "Cat", "descriptor": "shared"}
	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/biz-2/inventory/categories", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Equal(biz2.ID, result["businessId"])
}

func (s *InventoryCategoriesSuite) TestUpdateCategory_Success_NormalizesDescriptor() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "old")
	s.NoError(err)

	payload := map[string]interface{}{"name": "New Name", "descriptor": "  NEW  "}
	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/businesses/test-biz/inventory/categories/%s", cat.ID), payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Equal("New Name", result["name"])
	s.Equal("new", result["descriptor"])

	dbCat, err := s.inventoryHelper.GetCategory(ctx, cat.ID)
	s.NoError(err)
	s.Equal("new", dbCat.Descriptor)
}

func (s *InventoryCategoriesSuite) TestListCategories_Success_ViewPermission() {
	ctx := context.Background()
	ws, users, err := testutils.CreateWorkspaceWithUsers(ctx, testEnv.Database, "admin@example.com", "Password123!", []struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
		Role      role.Role
	}{
		{Email: "viewer@example.com", Password: "Password123!", FirstName: "Viewer", LastName: "User", Role: role.RoleUser},
	})
	s.NoError(err)
	viewerToken, err := auth.NewJwtToken(users[1].ID, ws.ID, users[1].AuthVersion)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	_, err = s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "A", "a")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "B", "b")
	s.NoError(err)

	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/inventory/categories", nil, viewerToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result []map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Len(result, 2)
}

func (s *InventoryCategoriesSuite) TestGetCategory_NotFound() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	_, err = s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/inventory/categories/cat_does_not_exist", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *InventoryCategoriesSuite) TestDeleteCategory_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "cat")
	s.NoError(err)

	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/test-biz/inventory/categories/%s", cat.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNoContent, resp.StatusCode)

	getResp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/test-biz/inventory/categories/%s", cat.ID), nil, token)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusNotFound, getResp.StatusCode)
}

func (s *InventoryCategoriesSuite) TestManageEndpoints_RequireManagePermission() {
	ctx := context.Background()
	ws, users, err := testutils.CreateWorkspaceWithUsers(ctx, testEnv.Database, "admin@example.com", "Password123!", []struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
		Role      role.Role
	}{
		{Email: "member@example.com", Password: "Password123!", FirstName: "Member", LastName: "User", Role: role.RoleUser},
	})
	s.NoError(err)
	memberToken, err := auth.NewJwtToken(users[1].ID, ws.ID, users[1].AuthVersion)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "cat")
	s.NoError(err)

	createPayload := map[string]interface{}{"name": "X", "descriptor": "x"}
	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/inventory/categories", createPayload, memberToken)
	s.NoError(err)
	resp.Body.Close()
	s.Equal(http.StatusForbidden, resp.StatusCode)

	updatePayload := map[string]interface{}{"name": "Y"}
	resp2, err := s.inventoryHelper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/businesses/test-biz/inventory/categories/%s", cat.ID), updatePayload, memberToken)
	s.NoError(err)
	resp2.Body.Close()
	s.Equal(http.StatusForbidden, resp2.StatusCode)

	resp3, err := s.inventoryHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/test-biz/inventory/categories/%s", cat.ID), nil, memberToken)
	s.NoError(err)
	resp3.Body.Close()
	s.Equal(http.StatusForbidden, resp3.StatusCode)
}

func (s *InventoryCategoriesSuite) TestRequiresAuth() {
	payload := map[string]interface{}{"name": "Cat", "descriptor": "cat"}
	resp, err := s.inventoryHelper.Client.Post("/v1/businesses/test-biz/inventory/categories", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func TestInventoryCategoriesSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(InventoryCategoriesSuite))
}
