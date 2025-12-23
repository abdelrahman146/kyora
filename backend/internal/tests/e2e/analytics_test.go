package e2e_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type AnalyticsSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *AnalyticsSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *AnalyticsSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database,
		"users",
		"workspaces",
		"businesses",
		"products",
		"variants",
		"categories",
		"customers",
		"orders",
		"order_items",
		"assets",
		"investments",
		"withdrawals",
		"expenses",
		"recurring_expenses",
	)
	s.NoError(err)
}

func (s *AnalyticsSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database,
		"users",
		"workspaces",
		"businesses",
		"products",
		"variants",
		"categories",
		"customers",
		"orders",
		"order_items",
		"assets",
		"investments",
		"withdrawals",
		"expenses",
		"recurring_expenses",
	)
	s.NoError(err)
}

func (s *AnalyticsSuite) createBusiness(ctx context.Context, workspaceID string) {
	repo := database.NewRepository[business.Business](testEnv.Database)
	biz := &business.Business{
		WorkspaceID:   workspaceID,
		Descriptor:    "main",
		Name:          "Test Business",
		CountryCode:   "AE",
		Currency:      "aed",
		VatRate:       decimal.Zero,
		SafetyBuffer:  decimal.Zero,
		EstablishedAt: time.Now().UTC(),
	}
	s.NoError(repo.CreateOne(ctx, biz))
}

func (s *AnalyticsSuite) analyticsPath(businessDescriptor string, path string) string {
	return "/v1/businesses/" + businessDescriptor + "/analytics" + path
}

func (s *AnalyticsSuite) TestDashboard_Unauthorized() {
	resp, err := s.helper.Client.Get(s.analyticsPath("main", "/dashboard"))
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *AnalyticsSuite) TestDashboard_Success() {
	ctx := context.Background()
	_, workspace, token, err := s.helper.CreateTestUser(ctx, "user@example.com", "ValidPassword123!", "John", "Doe", role.RoleUser)
	s.NoError(err)

	s.createBusiness(ctx, workspace.ID)

	resp, err := s.helper.Client.AuthenticatedRequest(http.MethodGet, s.analyticsPath("main", "/dashboard"), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Contains(result, "businessID")
	s.Contains(result, "revenueLast30Days")
	s.Contains(result, "grossProfitLast30Days")
	s.Contains(result, "openOrdersCount")
	s.Contains(result, "lowStockItemsCount")
	s.Contains(result, "allTimeRevenue")
	s.Contains(result, "safeToDrawAmount")
	s.Contains(result, "salesPerformanceLast30Days")
	s.Contains(result, "liveOrderFunnel")
	s.Contains(result, "topSellingProducts")
	s.Contains(result, "newCustomersTimeSeries")
}

func (s *AnalyticsSuite) TestSalesAnalytics_InvalidDate() {
	ctx := context.Background()
	_, workspace, token, err := s.helper.CreateTestUser(ctx, "user2@example.com", "ValidPassword123!", "Jane", "Doe", role.RoleUser)
	s.NoError(err)

	s.createBusiness(ctx, workspace.ID)

	resp, err := s.helper.Client.AuthenticatedRequest(http.MethodGet, s.analyticsPath("main", "/sales?from=not-a-date"), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *AnalyticsSuite) TestSalesAnalytics_FromAfterTo() {
	ctx := context.Background()
	_, workspace, token, err := s.helper.CreateTestUser(ctx, "user3@example.com", "ValidPassword123!", "Sam", "Doe", role.RoleUser)
	s.NoError(err)

	s.createBusiness(ctx, workspace.ID)

	resp, err := s.helper.Client.AuthenticatedRequest(http.MethodGet, s.analyticsPath("main", "/sales?from=2025-12-10&to=2025-12-01"), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func TestAnalyticsSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(AnalyticsSuite))
}
