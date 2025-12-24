package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type OrderSuite struct {
	suite.Suite
	accountHelper *AccountTestHelper
	orderHelper   *OrderTestHelper
}

func (s *OrderSuite) SetupSuite() {
	s.accountHelper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
	s.orderHelper = NewOrderTestHelper(testEnv.Database, e2eBaseURL)
}

func (s *OrderSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "orders", "order_items", "order_notes",
		"customers", "customer_addresses", "products", "variants", "categories",
		"businesses", "users", "workspaces", "subscriptions")
	s.NoError(err)
}

func (s *OrderSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "orders", "order_items", "order_notes",
		"customers", "customer_addresses", "products", "variants", "categories",
		"businesses", "users", "workspaces", "subscriptions")
	s.NoError(err)
}

func (s *OrderSuite) TestOrderLifecycle() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.orderHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	cust, addr, err := s.orderHelper.CreateTestCustomer(ctx, biz.ID, "customer@example.com", "Test Customer")
	s.NoError(err)

	cat, err := s.orderHelper.CreateTestCategory(ctx, biz.ID, "Electronics", "electronics")
	s.NoError(err)

	_, variant, err := s.orderHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Test Product",
		decimal.NewFromFloat(100), decimal.NewFromFloat(200), 10)
	s.NoError(err)

	// Create order
	payload := map[string]interface{}{
		"customerId":        cust.ID,
		"shippingAddressId": addr.ID,
		"channel":           "instagram",
		"items": []map[string]interface{}{
			{
				"variantId": variant.ID,
				"quantity":  2,
				"unitPrice": 200,
				"unitCost":  100,
			},
		},
	}

	resp, err := s.orderHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/orders", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &created))
	orderID := created["id"].(string)
	s.NotEmpty(orderID)
	s.Equal(cust.ID, created["customerId"])
	s.Equal("instagram", created["channel"])
	s.Equal(string(order.OrderStatusPending), created["status"])

	// Verify stock reduction
	updatedVariant, err := s.orderHelper.GetVariant(ctx, variant.ID)
	s.NoError(err)
	s.Equal(8, updatedVariant.StockQuantity)

	// Get order
	getResp, err := s.orderHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/test-biz/orders/%s", orderID), nil, token)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusOK, getResp.StatusCode)

	var fetched map[string]interface{}
	s.NoError(testutils.DecodeJSON(getResp, &fetched))
	s.Equal(orderID, fetched["id"])

	// Update order
	updatePayload := map[string]interface{}{
		"channel": "facebook",
	}
	updateResp, err := s.orderHelper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/businesses/test-biz/orders/%s", orderID), updatePayload, token)
	s.NoError(err)
	defer updateResp.Body.Close()
	s.Equal(http.StatusOK, updateResp.StatusCode)

	var updated map[string]interface{}
	s.NoError(testutils.DecodeJSON(updateResp, &updated))
	s.Equal("facebook", updated["channel"])

	// Add note
	notePayload := map[string]interface{}{
		"content": "Customer requested gift wrapping",
	}
	noteResp, err := s.orderHelper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/test-biz/orders/%s/notes", orderID), notePayload, token)
	s.NoError(err)
	defer noteResp.Body.Close()
	s.Equal(http.StatusCreated, noteResp.StatusCode)

	var note map[string]interface{}
	s.NoError(testutils.DecodeJSON(noteResp, &note))
	s.NotEmpty(note["id"])
	s.Equal("Customer requested gift wrapping", note["content"])

	// Delete order (must be pending to delete)
	deleteResp, err := s.orderHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/test-biz/orders/%s", orderID), nil, token)
	s.NoError(err)
	defer deleteResp.Body.Close()
	s.Equal(http.StatusNoContent, deleteResp.StatusCode)

	// Verify deletion
	count, err := s.orderHelper.CountOrders(ctx, biz.ID)
	s.NoError(err)
	s.Equal(int64(0), count)

	// Verify inventory was restocked
	finalVariant, err := s.orderHelper.GetVariant(ctx, variant.ID)
	s.NoError(err)
	s.Equal(10, finalVariant.StockQuantity)
}

func (s *OrderSuite) TestCreateOrderValidation() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	_, err = s.orderHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	tests := []struct {
		name    string
		payload map[string]interface{}
		status  int
	}{
		{"missing customerId", map[string]interface{}{"channel": "instagram", "shippingAddressId": "addr_test", "items": []map[string]interface{}{{"variantId": "var_123", "quantity": 1, "unitPrice": 200, "unitCost": 100}}}, http.StatusBadRequest},
		{"missing items", map[string]interface{}{"customerId": "cus_123", "channel": "instagram", "shippingAddressId": "addr_test"}, http.StatusBadRequest},
		{"empty items", map[string]interface{}{"customerId": "cus_123", "channel": "instagram", "shippingAddressId": "addr_test", "items": []map[string]interface{}{}}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.orderHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/orders", tt.payload, token)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(tt.status, resp.StatusCode)
		})
	}
}

func (s *OrderSuite) TestInsufficientStock() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.orderHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	cust, addr, err := s.orderHelper.CreateTestCustomer(ctx, biz.ID, "customer@example.com", "Test Customer")
	s.NoError(err)

	cat, err := s.orderHelper.CreateTestCategory(ctx, biz.ID, "Electronics", "electronics")
	s.NoError(err)

	_, variant, err := s.orderHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Limited Stock",
		decimal.NewFromFloat(100), decimal.NewFromFloat(200), 3)
	s.NoError(err)

	payload := map[string]interface{}{
		"customerId":        cust.ID,
		"shippingAddressId": addr.ID,
		"channel":           "instagram",
		"items": []map[string]interface{}{
			{
				"variantId": variant.ID,
				"quantity":  5,
				"unitPrice": 200,
				"unitCost":  100,
			},
		},
	}

	resp, err := s.orderHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/orders", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusConflict, resp.StatusCode)

	// Verify stock unchanged
	updatedVariant, err := s.orderHelper.GetVariant(ctx, variant.ID)
	s.NoError(err)
	s.Equal(3, updatedVariant.StockQuantity)
}

func TestOrderSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(OrderSuite))
}
