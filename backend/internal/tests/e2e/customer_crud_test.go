package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

type CustomerCRUDSuite struct {
	suite.Suite
	accountHelper  *AccountTestHelper
	customerHelper *CustomerTestHelper
}

func (s *CustomerCRUDSuite) SetupSuite() {
	s.accountHelper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
	s.customerHelper = NewCustomerTestHelper(testEnv.Database, "http://localhost:18080")
}

func (s *CustomerCRUDSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "customers", "customer_addresses", "customer_notes", "subscriptions")
	s.NoError(err)
}

func (s *CustomerCRUDSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "customers", "customer_addresses", "customer_notes", "subscriptions")
	s.NoError(err)
}

func (s *CustomerCRUDSuite) TestCreateCustomer_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	payload := map[string]interface{}{
		"name":        "John Doe",
		"email":       "john@example.com",
		"countryCode": "eg",
		"gender":      "male",
		"phoneNumber": "1234567890",
		"phoneCode":   "+20",
	}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("POST", "/v1/customers", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.NotEmpty(result["id"])
	s.Equal(biz.ID, result["businessId"])
	s.Equal("John Doe", result["name"])
	s.Equal("john@example.com", result["email"])
	s.Equal("EG", result["countryCode"])
	s.Equal("male", result["gender"])
	s.Equal("1234567890", result["phoneNumber"])
	s.Equal("+20", result["phoneCode"])
	s.Contains(result, "joinedAt")

	// Verify DB state
	dbCustomer, err := s.customerHelper.GetCustomer(ctx, result["id"].(string))
	s.NoError(err)
	s.Equal("John Doe", dbCustomer.Name)
	s.Equal("john@example.com", dbCustomer.Email.String)
}

func (s *CustomerCRUDSuite) TestCreateCustomer_ValidationErrors() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	tests := []struct {
		name    string
		payload map[string]interface{}
	}{
		{"missing name", map[string]interface{}{"email": "test@example.com", "countryCode": "eg"}},
		{"missing email", map[string]interface{}{"name": "Test", "countryCode": "eg"}},
		{"missing countryCode", map[string]interface{}{"name": "Test", "email": "test@example.com"}},
		{"invalid email", map[string]interface{}{"name": "Test", "email": "invalid-email", "countryCode": "eg"}},
		{"invalid countryCode length", map[string]interface{}{"name": "Test", "email": "test@example.com", "countryCode": "e"}},
		{"invalid gender", map[string]interface{}{"name": "Test", "email": "test@example.com", "countryCode": "eg", "gender": "invalid"}},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.customerHelper.Client.AuthenticatedRequest("POST", "/v1/customers", tt.payload, token)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func (s *CustomerCRUDSuite) TestCreateCustomer_DuplicateEmail() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	// Create first customer
	_, err = s.customerHelper.CreateTestCustomer(ctx, biz.ID, "duplicate@example.com", "First Customer")
	s.NoError(err)

	// Try to create second customer with same email in same business
	payload := map[string]interface{}{
		"name":        "Second Customer",
		"email":       "duplicate@example.com",
		"countryCode": "eg",
	}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("POST", "/v1/customers", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusConflict, resp.StatusCode)
}

func (s *CustomerCRUDSuite) TestCreateCustomer_SameEmailDifferentBusiness() {
	ctx := context.Background()
	_, ws, _, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	// Create two businesses
	biz1, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "biz-1")
	s.NoError(err)
	biz2, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "biz-2")
	s.NoError(err)

	// Create customer in first business
	_, err = s.customerHelper.CreateTestCustomer(ctx, biz1.ID, "shared@example.com", "Customer 1")
	s.NoError(err)

	// Should succeed - email uniqueness is scoped per business
	customer2, err := s.customerHelper.CreateTestCustomer(ctx, biz2.ID, "shared@example.com", "Customer 2")
	s.NoError(err)
	s.Equal(biz2.ID, customer2.BusinessID)
}

func (s *CustomerCRUDSuite) TestCreateCustomer_UnauthorizedUser() {
	ctx := context.Background()
	_, ws, _, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	_, _, userToken, err := s.accountHelper.CreateTestUser(ctx, "user@example.com", "Password123!", "User", "User", role.RoleUser)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	payload := map[string]interface{}{
		"name":        "John Doe",
		"email":       "john@example.com",
		"countryCode": "eg",
	}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("POST", "/v1/customers", payload, userToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusForbidden, resp.StatusCode)
}

func (s *CustomerCRUDSuite) TestGetCustomer_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	resp, err := s.customerHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/customers/%s", customer.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Equal(customer.ID, result["id"])
	s.Equal("John Doe", result["name"])
	s.Equal("john@example.com", result["email"])
}

func (s *CustomerCRUDSuite) TestGetCustomer_NotFound() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	resp, err := s.customerHelper.Client.AuthenticatedRequest("GET", "/v1/customers/cus_nonexistent", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *CustomerCRUDSuite) TestGetCustomer_CrossWorkspaceIsolation() {
	ctx := context.Background()

	// Workspace 1
	_, ws1, _, err := s.accountHelper.CreateTestUser(ctx, "admin1@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws1.ID))
	biz1, err := s.customerHelper.CreateTestBusiness(ctx, ws1.ID, "biz-1")
	s.NoError(err)
	customer1, err := s.customerHelper.CreateTestCustomer(ctx, biz1.ID, "customer1@example.com", "Customer 1")
	s.NoError(err)

	// Workspace 2
	_, ws2, token2, err := s.accountHelper.CreateTestUser(ctx, "admin2@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws2.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws2.ID, "biz-2")
	s.NoError(err)

	// Workspace 2 should not access Workspace 1's customer
	resp, err := s.customerHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/customers/%s", customer1.ID), nil, token2)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *CustomerCRUDSuite) TestListCustomers_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	// Create multiple customers
	for i := 0; i < 5; i++ {
		_, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, fmt.Sprintf("customer%d@example.com", i), fmt.Sprintf("Customer %d", i))
		s.NoError(err)
	}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("GET", "/v1/customers", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result, "items")
	s.Contains(result, "totalCount")
	s.Contains(result, "page")
	s.Contains(result, "pageSize")

	items := result["items"].([]interface{})
	s.Len(items, 5)
	s.Equal(float64(5), result["totalCount"])
}

func (s *CustomerCRUDSuite) TestListCustomers_Pagination() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	// Create 15 customers
	for i := 0; i < 15; i++ {
		_, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, fmt.Sprintf("customer%d@example.com", i), fmt.Sprintf("Customer %d", i))
		s.NoError(err)
	}

	// Page 1
	resp, err := s.customerHelper.Client.AuthenticatedRequest("GET", "/v1/customers?page=1&pageSize=10", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var page1 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &page1))
	items1 := page1["items"].([]interface{})
	s.Len(items1, 10)
	s.Equal(float64(15), page1["totalCount"])
	s.Equal(true, page1["hasMore"])

	// Page 2
	resp, err = s.customerHelper.Client.AuthenticatedRequest("GET", "/v1/customers?page=2&pageSize=10", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var page2 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &page2))
	items2 := page2["items"].([]interface{})
	s.Len(items2, 5)
	s.Equal(false, page2["hasMore"])
}

func (s *CustomerCRUDSuite) TestListCustomers_CrossWorkspaceIsolation() {
	ctx := context.Background()

	// Workspace 1
	_, ws1, token1, err := s.accountHelper.CreateTestUser(ctx, "admin1@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws1.ID))
	biz1, err := s.customerHelper.CreateTestBusiness(ctx, ws1.ID, "biz-1")
	s.NoError(err)
	_, err = s.customerHelper.CreateTestCustomer(ctx, biz1.ID, "ws1customer@example.com", "WS1 Customer")
	s.NoError(err)

	// Workspace 2
	_, ws2, token2, err := s.accountHelper.CreateTestUser(ctx, "admin2@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws2.ID))
	biz2, err := s.customerHelper.CreateTestBusiness(ctx, ws2.ID, "biz-2")
	s.NoError(err)
	_, err = s.customerHelper.CreateTestCustomer(ctx, biz2.ID, "ws2customer@example.com", "WS2 Customer")
	s.NoError(err)

	// Each workspace should only see its own customers
	resp1, err := s.customerHelper.Client.AuthenticatedRequest("GET", "/v1/customers", nil, token1)
	s.NoError(err)
	defer resp1.Body.Close()
	var result1 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp1, &result1))
	s.Equal(float64(1), result1["totalCount"])

	resp2, err := s.customerHelper.Client.AuthenticatedRequest("GET", "/v1/customers", nil, token2)
	s.NoError(err)
	defer resp2.Body.Close()
	var result2 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &result2))
	s.Equal(float64(1), result2["totalCount"])
}

func (s *CustomerCRUDSuite) TestUpdateCustomer_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	payload := map[string]interface{}{
		"name":   "John Updated",
		"gender": "other",
	}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/customers/%s", customer.ID), payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Equal("John Updated", result["name"])
	s.Equal("other", result["gender"])
	s.Equal("john@example.com", result["email"])

	// Verify DB state
	dbCustomer, err := s.customerHelper.GetCustomer(ctx, customer.ID)
	s.NoError(err)
	s.Equal("John Updated", dbCustomer.Name)
}

func (s *CustomerCRUDSuite) TestUpdateCustomer_NotFound() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	payload := map[string]interface{}{"name": "Updated"}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("PATCH", "/v1/customers/cus_nonexistent", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *CustomerCRUDSuite) TestUpdateCustomer_CrossWorkspaceIsolation() {
	ctx := context.Background()

	// Workspace 1
	_, ws1, _, err := s.accountHelper.CreateTestUser(ctx, "admin1@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws1.ID))
	biz1, err := s.customerHelper.CreateTestBusiness(ctx, ws1.ID, "biz-1")
	s.NoError(err)
	customer1, err := s.customerHelper.CreateTestCustomer(ctx, biz1.ID, "customer1@example.com", "Customer 1")
	s.NoError(err)

	// Workspace 2
	_, ws2, token2, err := s.accountHelper.CreateTestUser(ctx, "admin2@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws2.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws2.ID, "biz-2")
	s.NoError(err)

	payload := map[string]interface{}{"name": "Hacked"}

	// Workspace 2 should not update Workspace 1's customer
	resp, err := s.customerHelper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/customers/%s", customer1.ID), payload, token2)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)

	// Verify original customer unchanged
	dbCustomer, err := s.customerHelper.GetCustomer(ctx, customer1.ID)
	s.NoError(err)
	s.Equal("Customer 1", dbCustomer.Name)
}

func (s *CustomerCRUDSuite) TestDeleteCustomer_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	resp, err := s.customerHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/customers/%s", customer.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify customer is soft deleted
	count, err := s.customerHelper.CountCustomers(ctx, biz.ID)
	s.NoError(err)
	s.Equal(int64(0), count)
}

func (s *CustomerCRUDSuite) TestDeleteCustomer_NotFound() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	resp, err := s.customerHelper.Client.AuthenticatedRequest("DELETE", "/v1/customers/cus_nonexistent", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *CustomerCRUDSuite) TestDeleteCustomer_CrossWorkspaceIsolation() {
	ctx := context.Background()

	// Workspace 1
	_, ws1, _, err := s.accountHelper.CreateTestUser(ctx, "admin1@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws1.ID))
	biz1, err := s.customerHelper.CreateTestBusiness(ctx, ws1.ID, "biz-1")
	s.NoError(err)
	customer1, err := s.customerHelper.CreateTestCustomer(ctx, biz1.ID, "customer1@example.com", "Customer 1")
	s.NoError(err)

	// Workspace 2
	_, ws2, token2, err := s.accountHelper.CreateTestUser(ctx, "admin2@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws2.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws2.ID, "biz-2")
	s.NoError(err)

	// Workspace 2 should not delete Workspace 1's customer
	resp, err := s.customerHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/customers/%s", customer1.ID), nil, token2)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)

	// Verify customer still exists
	count, err := s.customerHelper.CountCustomers(ctx, biz1.ID)
	s.NoError(err)
	s.Equal(int64(1), count)
}

func TestCustomerCRUDSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(CustomerCRUDSuite))
}
