package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type InventoryProductsSuite struct {
	suite.Suite
	accountHelper   *AccountTestHelper
	inventoryHelper *InventoryTestHelper
}

func (s *InventoryProductsSuite) SetupSuite() {
	s.accountHelper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
	s.inventoryHelper = NewInventoryTestHelper(testEnv.Database, "http://localhost:18080")
}

func (s *InventoryProductsSuite) resetDB() {
	s.NoError(testutils.TruncateTables(testEnv.Database, inventoryTables...))
}

func (s *InventoryProductsSuite) SetupTest() {
	s.resetDB()
}

func (s *InventoryProductsSuite) TearDownTest() {
	s.resetDB()
}

func (s *InventoryProductsSuite) TestCreateProduct_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "cat")
	s.NoError(err)

	payload := map[string]interface{}{"name": "Product 1", "description": "Desc", "categoryId": cat.ID}
	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/inventory/products", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.NotEmpty(result["id"])
	s.Equal(biz.ID, result["businessId"])
	s.Equal("Product 1", result["name"])
	s.Equal("Desc", result["description"])
	s.Equal(cat.ID, result["categoryId"])
}

func (s *InventoryProductsSuite) TestCreateProduct_CategoryNotFound() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	_, err = s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	payload := map[string]interface{}{"name": "Product", "description": "Desc", "categoryId": "cat_does_not_exist"}
	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/inventory/products", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *InventoryProductsSuite) TestListProducts_PaginationAndOrdering() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "cat")
	s.NoError(err)

	_, err = s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "A Product", "")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "B Product", "")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "C Product", "")
	s.NoError(err)

	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/inventory/products?page=1&pageSize=2&orderBy=name", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var page1 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &page1))
	s.Equal(float64(1), page1["page"])
	s.Equal(float64(2), page1["pageSize"])
	s.Equal(float64(3), page1["totalCount"])
	s.Equal(true, page1["hasMore"])
	items := page1["items"].([]interface{})
	s.Len(items, 2)
	s.Equal("A Product", items[0].(map[string]interface{})["name"])

	resp2, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/inventory/products?page=2&pageSize=2&orderBy=name", nil, token)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusOK, resp2.StatusCode)

	var page2 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &page2))
	s.Equal(false, page2["hasMore"])
	items2 := page2["items"].([]interface{})
	s.Len(items2, 1)
	s.Equal("C Product", items2[0].(map[string]interface{})["name"])
}

func (s *InventoryProductsSuite) TestGetProduct_IncludesVariants() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "cat")
	s.NoError(err)
	prod, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Product", "")
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod.ID, "RED", "SKU-RED", "USD", decimal.NewFromInt(1), decimal.NewFromInt(2), 5, 1)
	s.NoError(err)
	_, err = s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod.ID, "BLU", "SKU-BLU", "USD", decimal.NewFromInt(1), decimal.NewFromInt(2), 6, 1)
	s.NoError(err)

	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/test-biz/inventory/products/%s", prod.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Equal(prod.ID, result["id"])
	variants := result["variants"].([]interface{})
	s.Len(variants, 2)
	for _, it := range variants {
		v := it.(map[string]interface{})
		s.NotEmpty(v["id"])
		s.Equal(prod.ID, v["productId"])
		s.Contains(v, "sku")
		s.Contains(v, "stockQuantity")
	}
}

func (s *InventoryProductsSuite) TestUpdateProduct_RenamesVariants() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "cat")
	s.NoError(err)
	prod, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Old", "")
	s.NoError(err)
	v1, err := s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod.ID, "RED", "SKU-RED", "USD", decimal.NewFromInt(1), decimal.NewFromInt(2), 5, 1)
	s.NoError(err)

	payload := map[string]interface{}{"name": "New"}
	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/businesses/test-biz/inventory/products/%s", prod.ID), payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	variants := result["variants"].([]interface{})
	s.Len(variants, 1)
	updatedV := variants[0].(map[string]interface{})
	s.Equal("New - RED", updatedV["name"])

	getVar, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/test-biz/inventory/variants/%s", v1.ID), nil, token)
	s.NoError(err)
	defer getVar.Body.Close()
	s.Equal(http.StatusOK, getVar.StatusCode)

	var vRes map[string]interface{}
	s.NoError(testutils.DecodeJSON(getVar, &vRes))
	s.Equal("New - RED", vRes["name"])
}

func (s *InventoryProductsSuite) TestDeleteProduct_CascadesVariants() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "cat")
	s.NoError(err)
	prod, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Product", "")
	s.NoError(err)
	v1, err := s.inventoryHelper.CreateTestVariant(ctx, biz.ID, prod.ID, "RED", "SKU-RED", "USD", decimal.NewFromInt(1), decimal.NewFromInt(2), 5, 1)
	s.NoError(err)

	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/test-biz/inventory/products/%s", prod.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNoContent, resp.StatusCode)

	getVar, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/test-biz/inventory/variants/%s", v1.ID), nil, token)
	s.NoError(err)
	defer getVar.Body.Close()
	s.Equal(http.StatusNotFound, getVar.StatusCode)

	count, err := s.inventoryHelper.CountVariants(ctx, biz.ID)
	s.NoError(err)
	s.Equal(int64(0), count)
}

func (s *InventoryProductsSuite) TestCreateProductWithVariants_Success_SkuAutoGenerated() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "cat")
	s.NoError(err)

	payload := map[string]interface{}{
		"product": map[string]interface{}{"name": "P", "description": "D", "categoryId": cat.ID},
		"variants": []interface{}{
			map[string]interface{}{
				"code":               "RED",
				"sku":                "",
				"costPrice":          "10",
				"salePrice":          "15",
				"stockQuantity":      2,
				"stockQuantityAlert": 5,
			},
			map[string]interface{}{
				"code":               "BLU",
				"sku":                "SKU-BLU",
				"costPrice":          "1",
				"salePrice":          "2",
				"stockQuantity":      0,
				"stockQuantityAlert": 0,
			},
		},
	}

	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/inventory/products/with-variants", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.NotEmpty(result["id"])
	s.Equal("P", result["name"])
	variants := result["variants"].([]interface{})
	s.Len(variants, 2)

	first := variants[0].(map[string]interface{})
	second := variants[1].(map[string]interface{})
	if first["code"].(string) != "RED" {
		first, second = second, first
	}
	s.Equal("RED", first["code"])
	s.Equal("P - RED", first["name"])
	s.NotEmpty(first["sku"])
	s.Equal("USD", first["currency"])

	s.Equal("BLU", second["code"])
	s.Equal("P - BLU", second["name"])
	s.Equal("SKU-BLU", second["sku"])
}

func (s *InventoryProductsSuite) TestViewAllowed_ManageForbidden_ForUserRole() {
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
	memberToken, err := auth.NewJwtToken(users[1].ID, ws.ID)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "cat")
	s.NoError(err)
	prod, err := s.inventoryHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "P", "")
	s.NoError(err)

	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/inventory/products", nil, memberToken)
	s.NoError(err)
	resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	resp2, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/test-biz/inventory/products/%s", prod.ID), nil, memberToken)
	s.NoError(err)
	resp2.Body.Close()
	s.Equal(http.StatusOK, resp2.StatusCode)

	createPayload := map[string]interface{}{"name": "X", "description": "", "categoryId": cat.ID}
	resp3, err := s.inventoryHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/inventory/products", createPayload, memberToken)
	s.NoError(err)
	resp3.Body.Close()
	s.Equal(http.StatusForbidden, resp3.StatusCode)
}

func TestInventoryProductsSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(InventoryProductsSuite))
}
