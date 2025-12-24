package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/accounting"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

// AnalyticsSuite tests analytics endpoints
type AnalyticsSuite struct {
	suite.Suite
	analyticsHelper *AnalyticsTestHelper
}

func (s *AnalyticsSuite) SetupSuite() {
	s.analyticsHelper = NewAnalyticsTestHelper(testEnv.Database, e2eBaseURL)
}

func (s *AnalyticsSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database,
		"users", "workspaces", "businesses", "subscriptions",
		"orders", "order_items", "customers", "customer_addresses",
		"products", "variants", "categories",
		"expenses", "investments", "withdrawals", "assets")
}

func (s *AnalyticsSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database,
		"users", "workspaces", "businesses", "subscriptions",
		"orders", "order_items", "customers", "customer_addresses",
		"products", "variants", "categories",
		"expenses", "investments", "withdrawals", "assets")
}

func (s *AnalyticsSuite) TestDashboard_WithData() {
	ctx := context.Background()

	_, ws, token, err := testutils.CreateAuthenticatedUser(ctx, testEnv.Database, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(testutils.CreateTestSubscription(ctx, testEnv.Database, ws.ID))

	biz, err := s.analyticsHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	// Create basic data
	cat, err := s.analyticsHelper.CreateTestCategory(ctx, biz.ID, "Electronics", "electronics")
	s.NoError(err)
	product, err := s.analyticsHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Test Product")
	s.NoError(err)
	variant, err := s.analyticsHelper.CreateTestVariant(ctx, biz.ID, product.ID, "Variant 1",
		decimal.NewFromInt(50), decimal.NewFromInt(100), 3) // Low stock
	s.NoError(err)

	cust, err := s.analyticsHelper.CreateTestCustomer(ctx, biz.ID, "customer@example.com", "Customer")
	s.NoError(err)
	addr, err := s.analyticsHelper.CreateTestAddress(ctx, cust.ID)
	s.NoError(err)

	// Create an open order
	orderedAt := time.Date(2025, 1, 20, 10, 0, 0, 0, time.UTC)
	_, err = s.analyticsHelper.CreateTestOrder(ctx, biz.ID, cust.ID, addr.ID, "instagram", order.OrderStatusPending,
		[]OrderItemData{{VariantID: variant.ID, Quantity: 1, UnitPrice: decimal.NewFromInt(100), UnitCost: decimal.NewFromInt(50)}},
		orderedAt)
	s.NoError(err)

	resp, err := s.analyticsHelper.Client.AuthenticatedRequest("GET",
		fmt.Sprintf("/v1/businesses/%s/analytics/dashboard", biz.Descriptor), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Equal(biz.ID, result["businessID"])
	s.Contains(result, "revenueLast30Days")
	s.Contains(result, "grossProfitLast30Days")
	s.Contains(result, "openOrdersCount")
	s.Equal(float64(1), result["openOrdersCount"])
	s.Contains(result, "lowStockItemsCount")
	s.Equal(float64(1), result["lowStockItemsCount"])
	s.Contains(result, "allTimeRevenue")
	s.Contains(result, "safeToDrawAmount")
	s.Contains(result, "salesPerformanceLast30Days")
	s.Contains(result, "liveOrderFunnel")
	s.Contains(result, "topSellingProducts")
	s.Contains(result, "newCustomersTimeSeries")
}

func (s *AnalyticsSuite) TestInventoryAnalytics() {
	ctx := context.Background()

	_, ws, token, err := testutils.CreateAuthenticatedUser(ctx, testEnv.Database, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(testutils.CreateTestSubscription(ctx, testEnv.Database, ws.ID))

	biz, err := s.analyticsHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	cat, err := s.analyticsHelper.CreateTestCategory(ctx, biz.ID, "Electronics", "electronics")
	s.NoError(err)
	product, err := s.analyticsHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Product")
	s.NoError(err)
	_, err = s.analyticsHelper.CreateTestVariant(ctx, biz.ID, product.ID, "Variant",
		decimal.NewFromInt(100), decimal.NewFromInt(200), 10)
	s.NoError(err)

	resp, err := s.analyticsHelper.Client.AuthenticatedRequest("GET",
		fmt.Sprintf("/v1/businesses/%s/analytics/inventory", biz.Descriptor), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Equal(biz.ID, result["businessID"])
	s.Contains(result, "totalInventoryValue")
	s.Equal("1000", result["totalInventoryValue"]) // 100*10
	s.Contains(result, "totalInStock")
	s.Equal(float64(10), result["totalInStock"])
}

func (s *AnalyticsSuite) TestSalesAnalytics() {
	ctx := context.Background()

	_, ws, token, err := testutils.CreateAuthenticatedUser(ctx, testEnv.Database, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(testutils.CreateTestSubscription(ctx, testEnv.Database, ws.ID))

	biz, err := s.analyticsHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	cat, err := s.analyticsHelper.CreateTestCategory(ctx, biz.ID, "Electronics", "electronics")
	s.NoError(err)
	product, err := s.analyticsHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Test Product")
	s.NoError(err)
	variant, err := s.analyticsHelper.CreateTestVariant(ctx, biz.ID, product.ID, "Variant 1",
		decimal.NewFromInt(50), decimal.NewFromInt(100), 100)
	s.NoError(err)

	cust, err := s.analyticsHelper.CreateTestCustomer(ctx, biz.ID, "customer@example.com", "Customer")
	s.NoError(err)
	addr, err := s.analyticsHelper.CreateTestAddress(ctx, cust.ID)
	s.NoError(err)

	orderedAt1 := time.Date(2025, 1, 10, 10, 0, 0, 0, time.UTC)
	_, err = s.analyticsHelper.CreateTestOrder(ctx, biz.ID, cust.ID, addr.ID, "instagram", order.OrderStatusPending,
		[]OrderItemData{{VariantID: variant.ID, Quantity: 1, UnitPrice: decimal.NewFromInt(100), UnitCost: decimal.NewFromInt(50)}},
		orderedAt1)
	s.NoError(err)

	orderedAt2 := time.Date(2025, 1, 20, 10, 0, 0, 0, time.UTC)
	_, err = s.analyticsHelper.CreateTestOrder(ctx, biz.ID, cust.ID, addr.ID, "whatsapp", order.OrderStatusFulfilled,
		[]OrderItemData{{VariantID: variant.ID, Quantity: 2, UnitPrice: decimal.NewFromInt(100), UnitCost: decimal.NewFromInt(50)}},
		orderedAt2)
	s.NoError(err)

	resp, err := s.analyticsHelper.Client.AuthenticatedRequest("GET",
		fmt.Sprintf("/v1/businesses/%s/analytics/sales?from=2025-01-01&to=2025-01-31", biz.Descriptor), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Equal(biz.ID, result["businessID"])
	s.Equal("300", result["totalRevenue"]) // 100 + 200
	s.Equal(float64(2), result["totalOrders"])
	s.Equal(float64(3), result["itemsSold"]) // 1 + 2
	s.Equal("150", result["averageOrderValue"])
	s.Contains(result, "grossProfit")
	s.Contains(result, "numberOfSalesOverTime")
	s.Contains(result, "revenueOverTime")
	s.Contains(result, "topSellingProducts")
	s.Contains(result, "orderStatusBreakdown")
	s.Contains(result, "salesByCountry")
	s.Contains(result, "salesByChannel")
}

func (s *AnalyticsSuite) TestCustomerAnalytics() {
	ctx := context.Background()

	_, ws, token, err := testutils.CreateAuthenticatedUser(ctx, testEnv.Database, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(testutils.CreateTestSubscription(ctx, testEnv.Database, ws.ID))

	biz, err := s.analyticsHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	cat, err := s.analyticsHelper.CreateTestCategory(ctx, biz.ID, "Electronics", "electronics")
	s.NoError(err)
	product, err := s.analyticsHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Test Product")
	s.NoError(err)
	variant, err := s.analyticsHelper.CreateTestVariant(ctx, biz.ID, product.ID, "Variant 1",
		decimal.NewFromInt(50), decimal.NewFromInt(150), 100)
	s.NoError(err)

	// One existing customer (created before range) and one new customer (created within range)
	createdOld := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	custOld, err := s.analyticsHelper.CreateTestCustomerAt(ctx, biz.ID, "old@example.com", "Old Customer", createdOld)
	s.NoError(err)
	addrOld, err := s.analyticsHelper.CreateTestAddress(ctx, custOld.ID)
	s.NoError(err)

	createdNew := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)
	custNew, err := s.analyticsHelper.CreateTestCustomerAt(ctx, biz.ID, "new@example.com", "New Customer", createdNew)
	s.NoError(err)
	addrNew, err := s.analyticsHelper.CreateTestAddress(ctx, custNew.ID)
	s.NoError(err)

	// CustomerOld: 2 orders (150 + 150)
	_, err = s.analyticsHelper.CreateTestOrder(ctx, biz.ID, custOld.ID, addrOld.ID, "instagram", order.OrderStatusFulfilled,
		[]OrderItemData{{VariantID: variant.ID, Quantity: 1, UnitPrice: decimal.NewFromInt(150), UnitCost: decimal.NewFromInt(50)}},
		time.Date(2025, 1, 10, 10, 0, 0, 0, time.UTC))
	s.NoError(err)
	_, err = s.analyticsHelper.CreateTestOrder(ctx, biz.ID, custOld.ID, addrOld.ID, "instagram", order.OrderStatusFulfilled,
		[]OrderItemData{{VariantID: variant.ID, Quantity: 1, UnitPrice: decimal.NewFromInt(150), UnitCost: decimal.NewFromInt(50)}},
		time.Date(2025, 1, 11, 10, 0, 0, 0, time.UTC))
	s.NoError(err)

	// CustomerNew: 2 orders (50 + 50)
	_, err = s.analyticsHelper.CreateTestOrder(ctx, biz.ID, custNew.ID, addrNew.ID, "whatsapp", order.OrderStatusPending,
		[]OrderItemData{{VariantID: variant.ID, Quantity: 1, UnitPrice: decimal.NewFromInt(50), UnitCost: decimal.NewFromInt(50)}},
		time.Date(2025, 1, 20, 10, 0, 0, 0, time.UTC))
	s.NoError(err)
	_, err = s.analyticsHelper.CreateTestOrder(ctx, biz.ID, custNew.ID, addrNew.ID, "whatsapp", order.OrderStatusPending,
		[]OrderItemData{{VariantID: variant.ID, Quantity: 1, UnitPrice: decimal.NewFromInt(50), UnitCost: decimal.NewFromInt(50)}},
		time.Date(2025, 1, 21, 10, 0, 0, 0, time.UTC))
	s.NoError(err)

	// Marketing spend to compute CAC (newCustomers=1)
	_, err = s.analyticsHelper.CreateTestExpense(ctx, biz.ID, accounting.ExpenseCategoryMarketing,
		decimal.NewFromInt(100), time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC))
	s.NoError(err)

	resp, err := s.analyticsHelper.Client.AuthenticatedRequest("GET",
		fmt.Sprintf("/v1/businesses/%s/analytics/customers?from=2025-01-01&to=2025-01-31", biz.Descriptor), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Equal(biz.ID, result["businessID"])
	s.Equal(float64(1), result["newCustomers"])
	s.Equal(float64(2), result["returningCustomers"])
	s.Equal("1", result["repeatCustomerRate"])
	s.Equal("200", result["averageRevenuePerCustomer"]) // (300+100)/2
	s.Equal("100", result["customerAcquisitionCost"])
	s.Equal("200", result["customerLifetimeValue"])          // avgOrderValue(100) * avgPurchaseFrequency(2)
	s.Equal("2", result["averageCustomerPurchaseFrequency"]) // 4 orders / 2 unique purchasers
	s.Contains(result, "newCustomersOverTime")
	s.Contains(result, "topCustomersByRevenue")

	// Top customers list should be ordered by revenue
	topCustomers, ok := result["topCustomersByRevenue"].([]interface{})
	s.True(ok)
	s.GreaterOrEqual(len(topCustomers), 2)
	first := topCustomers[0].(map[string]interface{})
	s.Equal(custOld.ID, first["id"])
}

func (s *AnalyticsSuite) TestFinancialPosition() {
	ctx := context.Background()

	user, ws, token, err := testutils.CreateAuthenticatedUser(ctx, testEnv.Database, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(testutils.CreateTestSubscription(ctx, testEnv.Database, ws.ID))

	biz, err := s.analyticsHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	_, err = s.analyticsHelper.CreateTestInvestment(ctx, biz.ID, user.ID, decimal.NewFromInt(1000),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	s.NoError(err)

	resp, err := s.analyticsHelper.Client.AuthenticatedRequest("GET",
		fmt.Sprintf("/v1/businesses/%s/analytics/reports/financial-position", biz.Descriptor), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Equal(biz.ID, result["businessID"])
	s.Contains(result, "totalAssets")
	s.Contains(result, "ownerInvestment")
	s.Equal("1000", result["ownerInvestment"])
}

func (s *AnalyticsSuite) TestProfitAndLoss() {
	ctx := context.Background()

	_, ws, token, err := testutils.CreateAuthenticatedUser(ctx, testEnv.Database, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(testutils.CreateTestSubscription(ctx, testEnv.Database, ws.ID))

	biz, err := s.analyticsHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	occurredOn := time.Date(2025, 1, 22, 0, 0, 0, 0, time.UTC)
	_, err = s.analyticsHelper.CreateTestExpense(ctx, biz.ID, accounting.ExpenseCategoryShipping,
		decimal.NewFromInt(50), occurredOn)
	s.NoError(err)

	resp, err := s.analyticsHelper.Client.AuthenticatedRequest("GET",
		fmt.Sprintf("/v1/businesses/%s/analytics/reports/profit-and-loss", biz.Descriptor), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Equal(biz.ID, result["businessID"])
	s.Contains(result, "revenue")
	s.Contains(result, "totalExpenses")
	s.Equal("50", result["totalExpenses"])
	s.Contains(result, "netProfit")
}

func (s *AnalyticsSuite) TestCashFlow() {
	ctx := context.Background()

	user, ws, token, err := testutils.CreateAuthenticatedUser(ctx, testEnv.Database, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(testutils.CreateTestSubscription(ctx, testEnv.Database, ws.ID))

	biz, err := s.analyticsHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	investedAt := time.Date(2025, 1, 21, 10, 0, 0, 0, time.UTC)
	_, err = s.analyticsHelper.CreateTestInvestment(ctx, biz.ID, user.ID, decimal.NewFromInt(1000), investedAt)
	s.NoError(err)

	resp, err := s.analyticsHelper.Client.AuthenticatedRequest("GET",
		fmt.Sprintf("/v1/businesses/%s/analytics/reports/cash-flow", biz.Descriptor), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Equal(biz.ID, result["businessID"])
	s.Contains(result, "cashFromOwner")
	s.Equal("1000", result["cashFromOwner"])
	s.Contains(result, "totalCashIn")
	s.Contains(result, "netCashFlow")
}

func (s *AnalyticsSuite) TestAnalytics_Unauthorized() {
	ctx := context.Background()

	_, ws, _, err := testutils.CreateAuthenticatedUser(ctx, testEnv.Database, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)

	biz, err := s.analyticsHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	resp, err := s.analyticsHelper.Client.Get(fmt.Sprintf("/v1/businesses/%s/analytics/dashboard", biz.Descriptor))
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func TestAnalyticsSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(AnalyticsSuite))
}
