package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type InventoryProductsSearchFilterSuite struct {
	suite.Suite
	accountHelper   *AccountTestHelper
	inventoryHelper *InventoryTestHelper
}

func (s *InventoryProductsSearchFilterSuite) SetupSuite() {
	s.accountHelper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
	s.inventoryHelper = NewInventoryTestHelper(testEnv.Database, e2eBaseURL)
}

func (s *InventoryProductsSearchFilterSuite) resetDB() {
	s.NoError(testutils.TruncateTables(testEnv.Database, inventoryTables...))
}

func (s *InventoryProductsSearchFilterSuite) SetupTest() {
	s.resetDB()
}

func (s *InventoryProductsSearchFilterSuite) TearDownTest() {
	s.resetDB()
}

func (s *InventoryProductsSearchFilterSuite) TestListProducts_SearchByProductName() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Electronics", "electronics")
	s.NoError(err)

	// Create test products
	prod1, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Wireless Headphones", "Premium audio quality")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod1.ID, "Black", "WH-BLACK-001", "USD", decimal.NewFromInt(50), decimal.NewFromInt(100), 20, 5)
	s.NoError(err)

	prod2, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Gaming Mouse", "RGB gaming mouse")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod2.ID, "Red", "GM-RED-001", "USD", decimal.NewFromInt(30), decimal.NewFromInt(60), 15, 3)
	s.NoError(err)

	// Search by product name
	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/inventory/products?search=Headphones", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	items := result["items"].([]interface{})
	s.True(len(items) >= 1, "Should find at least one product")

	found := false
	for _, it := range items {
		m := it.(map[string]interface{})
		if m["name"] == "Wireless Headphones" {
			found = true
			break
		}
	}
	s.True(found, "Should find Wireless Headphones product")
}

func (s *InventoryProductsSearchFilterSuite) TestListProducts_SearchByVariantSKU() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Accessories", "accessories")
	s.NoError(err)

	// Create test products with distinct SKUs
	prod1, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Phone Case", "Protective case")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod1.ID, "Clear", "PC-CLEAR-XYZ123", "USD", decimal.NewFromInt(5), decimal.NewFromInt(15), 50, 10)
	s.NoError(err)

	prod2, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Screen Protector", "Tempered glass")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod2.ID, "Standard", "SP-STD-ABC456", "USD", decimal.NewFromInt(3), decimal.NewFromInt(10), 100, 20)
	s.NoError(err)

	// Search by SKU
	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/inventory/products?search=XYZ123", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	items := result["items"].([]interface{})
	s.True(len(items) >= 1, "Should find at least one product")

	found := false
	for _, it := range items {
		m := it.(map[string]interface{})
		if m["name"] == "Phone Case" {
			found = true
			break
		}
	}
	s.True(found, "Should find Phone Case product through SKU search")
}

func (s *InventoryProductsSearchFilterSuite) TestListProducts_FilterByCategoryID() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	// Create categories
	cat1, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Books", "books")
	s.NoError(err)
	cat2, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Toys", "toys")
	s.NoError(err)

	// Create products in different categories
	prod1, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat1.ID, "Novel", "Fiction book")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod1.ID, "Hardcover", "BK-HC-001", "USD", decimal.NewFromInt(10), decimal.NewFromInt(20), 30, 5)
	s.NoError(err)

	prod2, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat2.ID, "Action Figure", "Collectible toy")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod2.ID, "Blue", "AF-BLU-001", "USD", decimal.NewFromInt(15), decimal.NewFromInt(30), 20, 5)
	s.NoError(err)

	// Filter by category ID
	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/test-biz/inventory/products?categoryId=%s", cat1.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	items := result["items"].([]interface{})
	s.Len(items, 1, "Should find exactly one product in Books category")
	s.Equal("Novel", items[0].(map[string]interface{})["name"])
}

func (s *InventoryProductsSearchFilterSuite) TestListProducts_FilterByStockStatus_LowStock() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Test", "test")
	s.NoError(err)

	// Create products with different stock levels
	prod1, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "In Stock Product", "")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod1.ID, "Variant1", "IS-V1-001", "USD", decimal.NewFromInt(10), decimal.NewFromInt(20), 20, 5)
	s.NoError(err)

	prod2, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Low Stock Product", "")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod2.ID, "Variant2", "LS-V2-001", "USD", decimal.NewFromInt(10), decimal.NewFromInt(20), 3, 5)
	s.NoError(err)

	prod3, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Out of Stock Product", "")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod3.ID, "Variant3", "OS-V3-001", "USD", decimal.NewFromInt(10), decimal.NewFromInt(20), 0, 5)
	s.NoError(err)

	// Filter by low_stock
	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/inventory/products?stockStatus=low_stock", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	items := result["items"].([]interface{})
	s.True(len(items) >= 1, "Should find at least one low-stock product")

	found := false
	for _, it := range items {
		m := it.(map[string]interface{})
		if m["name"] == "Low Stock Product" {
			found = true
		}
	}
	s.True(found, "Should find the Low Stock Product")
}

func (s *InventoryProductsSearchFilterSuite) TestListProducts_CombinedSearchAndFilters() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	cat1, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Electronics", "electronics")
	s.NoError(err)
	cat2, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Books", "books")
	s.NoError(err)

	// Electronics + Low Stock + contains "Wireless"
	prod1, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat1.ID, "Wireless Keyboard", "")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod1.ID, "Black", "WK-BLK-001", "USD", decimal.NewFromInt(20), decimal.NewFromInt(40), 3, 5)
	s.NoError(err)

	// Electronics + In Stock + contains "Wireless"
	prod2, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat1.ID, "Wireless Mouse", "")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod2.ID, "White", "WM-WHT-001", "USD", decimal.NewFromInt(15), decimal.NewFromInt(30), 20, 5)
	s.NoError(err)

	// Books + Low Stock + contains "Wireless" (edge case)
	prod3, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat2.ID, "Wireless Communications Book", "")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod3.ID, "Paperback", "WCB-PB-001", "USD", decimal.NewFromInt(25), decimal.NewFromInt(50), 2, 5)
	s.NoError(err)

	// Search "Wireless" + Filter by Electronics category + Filter by low_stock
	params := url.Values{}
	params.Add("search", "Wireless")
	params.Add("categoryId", cat1.ID)
	params.Add("stockStatus", "low_stock")

	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/test-biz/inventory/products?%s", params.Encode()), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	items := result["items"].([]interface{})
	s.Len(items, 1, "Should find exactly one product matching all criteria")
	s.Equal("Wireless Keyboard", items[0].(map[string]interface{})["name"])
}

func TestInventoryProductsSearchFilterSuite(t *testing.T) {
	suite.Run(t, new(InventoryProductsSearchFilterSuite))
}
