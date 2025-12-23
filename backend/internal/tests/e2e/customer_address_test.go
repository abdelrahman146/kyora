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

type CustomerAddressSuite struct {
	suite.Suite
	accountHelper  *AccountTestHelper
	customerHelper *CustomerTestHelper
}

func (s *CustomerAddressSuite) SetupSuite() {
	s.accountHelper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
	s.customerHelper = NewCustomerTestHelper(testEnv.Database, "http://localhost:18080")
}

func (s *CustomerAddressSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "customers", "customer_addresses", "subscriptions")
	s.NoError(err)
}

func (s *CustomerAddressSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "customers", "customer_addresses", "subscriptions")
	s.NoError(err)
}

func (s *CustomerAddressSuite) TestCreateAddress_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	payload := map[string]interface{}{
		"countryCode": "eg",
		"state":       "Cairo",
		"city":        "Cairo",
		"street":      "123 Test St",
		"phoneCode":   "+20",
		"phone":       "1234567890",
		"zipCode":     "12345",
	}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/addresses", customer.ID), payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.NotEmpty(result["id"])
	s.Equal(customer.ID, result["customerId"])
	s.Equal("EG", result["countryCode"])
	s.Equal("Cairo", result["state"])
	s.Equal("Cairo", result["city"])
	s.Equal("123 Test St", result["street"])
	s.Equal("+20", result["phoneCode"])
	s.Equal("1234567890", result["phoneNumber"])
	s.Equal("12345", result["zipCode"])

	// Verify DB state
	count, err := s.customerHelper.CountAddresses(ctx, customer.ID)
	s.NoError(err)
	s.Equal(int64(1), count)
}

func (s *CustomerAddressSuite) TestCreateAddress_ValidationErrors() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	tests := []struct {
		name    string
		payload map[string]interface{}
	}{
		{"missing countryCode", map[string]interface{}{"state": "Cairo", "city": "Cairo", "phoneCode": "+20", "phone": "1234567890"}},
		{"missing state", map[string]interface{}{"countryCode": "eg", "city": "Cairo", "phoneCode": "+20", "phone": "1234567890"}},
		{"missing city", map[string]interface{}{"countryCode": "eg", "state": "Cairo", "phoneCode": "+20", "phone": "1234567890"}},
		{"missing phoneCode", map[string]interface{}{"countryCode": "eg", "state": "Cairo", "city": "Cairo", "phone": "1234567890"}},
		{"missing phone", map[string]interface{}{"countryCode": "eg", "state": "Cairo", "city": "Cairo", "phoneCode": "+20"}},
		{"invalid countryCode length", map[string]interface{}{"countryCode": "e", "state": "Cairo", "city": "Cairo", "phoneCode": "+20", "phone": "1234567890"}},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.customerHelper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/addresses", customer.ID), tt.payload, token)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func (s *CustomerAddressSuite) TestCreateAddress_CustomerNotFound() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	payload := map[string]interface{}{
		"countryCode": "eg",
		"state":       "Cairo",
		"city":        "Cairo",
		"phoneCode":   "+20",
		"phone":       "1234567890",
	}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/customers/cus_nonexistent/addresses", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *CustomerAddressSuite) TestListAddresses_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	// Create multiple addresses
	for i := 0; i < 3; i++ {
		_, err := s.customerHelper.CreateTestAddress(ctx, customer.ID)
		s.NoError(err)
	}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/addresses", customer.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result []interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Len(result, 3)
}

func (s *CustomerAddressSuite) TestListAddresses_CustomerNotFound() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	resp, err := s.customerHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/customers/cus_nonexistent/addresses", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *CustomerAddressSuite) TestUpdateAddress_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	address, err := s.customerHelper.CreateTestAddress(ctx, customer.ID)
	s.NoError(err)

	payload := map[string]interface{}{
		"city":   "Giza",
		"street": "456 Updated St",
	}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/addresses/%s", customer.ID, address.ID), payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Equal("Giza", result["city"])
	s.Equal("456 Updated St", result["street"])

	// Verify DB state
	dbAddress, err := s.customerHelper.GetAddress(ctx, address.ID)
	s.NoError(err)
	s.Equal("Giza", dbAddress.City)
	s.Equal("456 Updated St", dbAddress.Street.String)
}

func (s *CustomerAddressSuite) TestUpdateAddress_AddressNotFound() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	payload := map[string]interface{}{"city": "Giza"}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/addresses/addr_nonexistent", customer.ID), payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *CustomerAddressSuite) TestUpdateAddress_WrongCustomer() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer1, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "customer1@example.com", "Customer 1")
	s.NoError(err)

	customer2, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "customer2@example.com", "Customer 2")
	s.NoError(err)

	address1, err := s.customerHelper.CreateTestAddress(ctx, customer1.ID)
	s.NoError(err)

	payload := map[string]interface{}{"city": "Hacked"}

	// Try to update customer1's address through customer2's endpoint
	resp, err := s.customerHelper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/addresses/%s", customer2.ID, address1.ID), payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)

	// Verify address unchanged
	dbAddress, err := s.customerHelper.GetAddress(ctx, address1.ID)
	s.NoError(err)
	s.Equal("Cairo", dbAddress.City)
}

func (s *CustomerAddressSuite) TestDeleteAddress_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	address, err := s.customerHelper.CreateTestAddress(ctx, customer.ID)
	s.NoError(err)

	resp, err := s.customerHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/addresses/%s", customer.ID, address.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify address is soft deleted
	count, err := s.customerHelper.CountAddresses(ctx, customer.ID)
	s.NoError(err)
	s.Equal(int64(0), count)
}

func (s *CustomerAddressSuite) TestDeleteAddress_AddressNotFound() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	resp, err := s.customerHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/addresses/addr_nonexistent", customer.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *CustomerAddressSuite) TestDeleteAddress_WrongCustomer() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer1, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "customer1@example.com", "Customer 1")
	s.NoError(err)

	customer2, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "customer2@example.com", "Customer 2")
	s.NoError(err)

	address1, err := s.customerHelper.CreateTestAddress(ctx, customer1.ID)
	s.NoError(err)

	// Try to delete customer1's address through customer2's endpoint
	resp, err := s.customerHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/addresses/%s", customer2.ID, address1.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)

	// Verify address still exists
	count, err := s.customerHelper.CountAddresses(ctx, customer1.ID)
	s.NoError(err)
	s.Equal(int64(1), count)
}

func (s *CustomerAddressSuite) TestAddressOperations_CrossWorkspaceIsolation() {
	ctx := context.Background()

	// Workspace 1
	_, ws1, token1, err := s.accountHelper.CreateTestUser(ctx, "admin1@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws1.ID))
	biz1, err := s.customerHelper.CreateTestBusiness(ctx, ws1.ID, "biz-1")
	s.NoError(err)
	customer1, err := s.customerHelper.CreateTestCustomer(ctx, biz1.ID, "customer1@example.com", "Customer 1")
	s.NoError(err)
	address1, err := s.customerHelper.CreateTestAddress(ctx, customer1.ID)
	s.NoError(err)

	// Workspace 2
	_, ws2, token2, err := s.accountHelper.CreateTestUser(ctx, "admin2@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws2.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws2.ID, "biz-2")
	s.NoError(err)

	// WS2 should not list WS1's customer addresses
	resp, err := s.customerHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/biz-1/customers/%s/addresses", customer1.ID), nil, token2)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)

	// WS2 should not update WS1's address
	payload := map[string]interface{}{"city": "Hacked"}
	resp, err = s.customerHelper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/businesses/biz-1/customers/%s/addresses/%s", customer1.ID, address1.ID), payload, token2)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)

	// WS2 should not delete WS1's address
	resp, err = s.customerHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/biz-1/customers/%s/addresses/%s", customer1.ID, address1.ID), nil, token2)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)

	// Verify WS1 address unchanged and exists
	count, err := s.customerHelper.CountAddresses(ctx, customer1.ID)
	s.NoError(err)
	s.Equal(int64(1), count)

	// WS1 can still access their own address
	resp, err = s.customerHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/biz-1/customers/%s/addresses", customer1.ID), nil, token1)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)
}

func TestCustomerAddressSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(CustomerAddressSuite))
}
