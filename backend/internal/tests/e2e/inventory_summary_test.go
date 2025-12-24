package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type InventorySummarySuite struct {
	suite.Suite
	accountHelper   *AccountTestHelper
	inventoryHelper *InventoryTestHelper
}

func (s *InventorySummarySuite) SetupSuite() {
	s.accountHelper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
	s.inventoryHelper = NewInventoryTestHelper(testEnv.Database, e2eBaseURL)
}

func (s *InventorySummarySuite) resetDB() {
	s.NoError(testutils.TruncateTables(testEnv.Database, inventoryTables...))
}

func (s *InventorySummarySuite) SetupTest() {
	s.resetDB()
}

func (s *InventorySummarySuite) TearDownTest() {
	s.resetDB()
}

func (s *InventorySummarySuite) TestSummary_ComputesAllMetrics() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "cat")
	s.NoError(err)

	p1, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "P1", "")
	s.NoError(err)
	p2, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "P2", "")
	s.NoError(err)

	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, p1.ID, "A", "SKU-A", "USD", decimal.NewFromInt(10), decimal.NewFromInt(15), 2, 5)
	s.NoError(err) // low stock
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, p1.ID, "B", "SKU-B", "USD", decimal.NewFromInt(1), decimal.NewFromInt(2), 0, 0)
	s.NoError(err) // low stock + out of stock
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, p2.ID, "C", "SKU-C", "USD", decimal.NewFromInt(3), decimal.NewFromInt(4), 7, 2)
	s.NoError(err)

	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/inventory/summary?topLimit=2", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Equal(float64(2), result["productsCount"])
	s.Equal(float64(3), result["variantsCount"])
	s.Equal(float64(1), result["categoriesCount"])
	s.Equal(float64(2), result["lowStockVariantsCount"])
	s.Equal(float64(1), result["outOfStockVariantsCount"])
	s.Equal(float64(9), result["totalStockUnits"])
	s.Equal("41", result["inventoryValue"])

	top := result["topProductsByInventoryValue"].([]interface{})
	s.Len(top, 2)
	s.Equal(p2.ID, top[0].(map[string]interface{})["id"])
	s.Equal(p1.ID, top[1].(map[string]interface{})["id"])
}

func (s *InventorySummarySuite) TestTopProductsByInventoryValue_DetailedIncludesValuesAndOrdering() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "cat")
	s.NoError(err)

	p1, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "P1", "")
	s.NoError(err)
	p2, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "P2", "")
	s.NoError(err)

	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, p1.ID, "A", "SKU-A", "USD", decimal.NewFromInt(10), decimal.NewFromInt(15), 2, 5)
	s.NoError(err) // value 20
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, p2.ID, "B", "SKU-B", "USD", decimal.NewFromInt(3), decimal.NewFromInt(4), 7, 2)
	s.NoError(err) // value 21

	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/inventory/top-products?limit=2", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result []map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Len(result, 2)

	first := result[0]
	second := result[1]
	firstProduct := first["product"].(map[string]interface{})
	secondProduct := second["product"].(map[string]interface{})

	s.Equal(p2.ID, firstProduct["id"])
	s.Equal("21", first["inventoryValue"])
	s.Equal(p1.ID, secondProduct["id"])
	s.Equal("20", second["inventoryValue"])
}

func (s *InventorySummarySuite) TestSummary_ViewAllowed_ForUserRole() {
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
	viewerToken, err := auth.NewJwtToken(users[1].ID, ws.ID)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	_, err = s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/inventory/summary", nil, viewerToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)
}

func TestInventorySummarySuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(InventorySummarySuite))
}
