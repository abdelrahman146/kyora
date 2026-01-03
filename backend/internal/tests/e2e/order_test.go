package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

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
		"businesses", "shipping_zones", "users", "workspaces", "subscriptions")
	s.NoError(err)
}

func (s *OrderSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "orders", "order_items", "order_notes",
		"customers", "customer_addresses", "products", "variants", "categories",
		"businesses", "shipping_zones", "users", "workspaces", "subscriptions")
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

func (s *OrderSuite) TestCreateOrder_WithShippingZone_ComputesShippingFee() {
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

	zone, err := s.orderHelper.CreateTestShippingZone(ctx, biz.ID, "EG Local", []string{"EG"}, decimal.NewFromInt(25), decimal.NewFromInt(500))
	s.NoError(err)

	// subtotal 400 => shipping fee 25
	payload := map[string]interface{}{
		"customerId":        cust.ID,
		"shippingAddressId": addr.ID,
		"channel":           "instagram",
		"shippingZoneId":    zone.ID,
		"items": []map[string]interface{}{
			{
				"variantId": variant.ID,
				"quantity":  2,
				"unitPrice": 200,
				"unitCost":  100,
			},
		},
	}

	// Respect create-order rate limit (min 1s interval per actor+business).
	// Other tests may have created an order immediately before this one.
	time.Sleep(1100 * time.Millisecond)
	resp, err := s.orderHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/orders", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &created))
	s.Equal(zone.ID, created["shippingZoneId"])
	s.Equal("25", fmt.Sprint(created["shippingFee"]))

	// subtotal 600 => shipping fee 0 (free shipping)
	payload2 := map[string]interface{}{
		"customerId":        cust.ID,
		"shippingAddressId": addr.ID,
		"channel":           "instagram",
		"shippingZoneId":    zone.ID,
		"items": []map[string]interface{}{
			{
				"variantId": variant.ID,
				"quantity":  3,
				"unitPrice": 200,
				"unitCost":  100,
			},
		},
	}
	// Ensure we don't hit the same rate limit for back-to-back order creation.
	time.Sleep(1100 * time.Millisecond)
	resp2, err := s.orderHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/orders", payload2, token)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusCreated, resp2.StatusCode)

	var created2 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &created2))
	s.Equal(zone.ID, created2["shippingZoneId"])
	s.Equal("0", fmt.Sprint(created2["shippingFee"]))
}

func (s *OrderSuite) TestCreateOrder_WithShippingZone_RejectsCountryNotInZone() {
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

	// Zone does not include EG, but address is EG.
	zone, err := s.orderHelper.CreateTestShippingZone(ctx, biz.ID, "US Only", []string{"US"}, decimal.NewFromInt(25), decimal.NewFromInt(500))
	s.NoError(err)

	payload := map[string]interface{}{
		"customerId":        cust.ID,
		"shippingAddressId": addr.ID,
		"channel":           "instagram",
		"shippingZoneId":    zone.ID,
		"items": []map[string]interface{}{
			{
				"variantId": variant.ID,
				"quantity":  1,
				"unitPrice": 200,
				"unitCost":  100,
			},
		},
	}

	resp, err := s.orderHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/orders", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *OrderSuite) TestUpdateOrder_SetShippingZone_RecomputesShippingFee() {
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

	zone, err := s.orderHelper.CreateTestShippingZone(ctx, biz.ID, "EG Local", []string{"EG"}, decimal.NewFromInt(25), decimal.NewFromInt(500))
	s.NoError(err)

	createPayload := map[string]interface{}{
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
	resp, err := s.orderHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/orders", createPayload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &created))
	orderID := created["id"].(string)
	s.NotEmpty(orderID)

	updatePayload := map[string]interface{}{
		"shippingZoneId": zone.ID,
	}
	updateResp, err := s.orderHelper.Client.AuthenticatedRequest("PATCH", "/v1/businesses/test-biz/orders/"+orderID, updatePayload, token)
	s.NoError(err)
	defer updateResp.Body.Close()
	s.Equal(http.StatusOK, updateResp.StatusCode)

	var updated map[string]interface{}
	s.NoError(testutils.DecodeJSON(updateResp, &updated))
	s.Equal(zone.ID, updated["shippingZoneId"])
	s.Equal("25", fmt.Sprint(updated["shippingFee"]))
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

func (s *OrderSuite) TestListOrders_Search_ByCustomerAndOrderNumber() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.orderHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	cust, addr, err := s.orderHelper.CreateTestCustomer(ctx, biz.ID, "customer@example.com", "Search Customer")
	s.NoError(err)

	cat, err := s.orderHelper.CreateTestCategory(ctx, biz.ID, "Electronics", "electronics")
	s.NoError(err)

	_, variant, err := s.orderHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Test Product", decimal.NewFromFloat(100), decimal.NewFromFloat(200), 10)
	s.NoError(err)

	// Create order
	payload := map[string]interface{}{
		"customerId":        cust.ID,
		"shippingAddressId": addr.ID,
		"channel":           "instagram",
		"items": []map[string]interface{}{
			{
				"variantId": variant.ID,
				"quantity":  1,
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
	orderNumber, ok := created["orderNumber"].(string)
	s.True(ok)
	s.NotEmpty(orderNumber)

	// Search by customer name (join)
	listResp, err := s.orderHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/orders?search=Search%20Customer", nil, token)
	s.NoError(err)
	defer listResp.Body.Close()
	s.Equal(http.StatusOK, listResp.StatusCode)

	var listResult map[string]interface{}
	s.NoError(testutils.DecodeJSON(listResp, &listResult))
	s.Equal(float64(1), listResult["totalCount"])
	items := listResult["items"].([]interface{})
	s.Len(items, 1)
	s.Equal(orderID, items[0].(map[string]interface{})["id"])

	// Search by order number
	listResp2, err := s.orderHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/orders?search="+orderNumber, nil, token)
	s.NoError(err)
	defer listResp2.Body.Close()
	s.Equal(http.StatusOK, listResp2.StatusCode)

	var listResult2 map[string]interface{}
	s.NoError(testutils.DecodeJSON(listResp2, &listResult2))
	s.Equal(float64(1), listResult2["totalCount"])
}

func (s *OrderSuite) TestListOrders_Filter_ByPlatform() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.orderHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	cust, addr, err := s.orderHelper.CreateTestCustomer(ctx, biz.ID, "customer@example.com", "Platform Customer")
	s.NoError(err)

	cat, err := s.orderHelper.CreateTestCategory(ctx, biz.ID, "Electronics", "electronics")
	s.NoError(err)

	_, variant, err := s.orderHelper.CreateTestProduct(ctx, biz.ID, cat.ID, "Test Product", decimal.NewFromFloat(100), decimal.NewFromFloat(200), 50)
	s.NoError(err)

	create := func(channel string) (string, error) {
		payload := map[string]interface{}{
			"customerId":        cust.ID,
			"shippingAddressId": addr.ID,
			"channel":           channel,
			"items": []map[string]interface{}{
				{
					"variantId": variant.ID,
					"quantity":  1,
					"unitPrice": 200,
					"unitCost":  100,
				},
			},
		}

		resp, err := s.orderHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/orders", payload, token)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			return "", fmt.Errorf("expected 201, got %d", resp.StatusCode)
		}
		var created map[string]interface{}
		if err := testutils.DecodeJSON(resp, &created); err != nil {
			return "", err
		}
		id, _ := created["id"].(string)
		return id, nil
	}

	instagramID, err := create("instagram")
	s.NoError(err)
	time.Sleep(1100 * time.Millisecond)
	whatsappID, err := create("whatsapp")
	s.NoError(err)

	// Filter by a single platform
	listResp, err := s.orderHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/orders?socialPlatforms=instagram", nil, token)
	s.NoError(err)
	defer listResp.Body.Close()
	s.Equal(http.StatusOK, listResp.StatusCode)

	var listResult map[string]interface{}
	s.NoError(testutils.DecodeJSON(listResp, &listResult))
	s.Equal(float64(1), listResult["totalCount"])
	items := listResult["items"].([]interface{})
	s.Len(items, 1)
	s.Equal(instagramID, items[0].(map[string]interface{})["id"])

	// Filter by multiple platforms (repeatable query param)
	listResp2, err := s.orderHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/orders?socialPlatforms=instagram&socialPlatforms=whatsapp", nil, token)
	s.NoError(err)
	defer listResp2.Body.Close()
	s.Equal(http.StatusOK, listResp2.StatusCode)

	var listResult2 map[string]interface{}
	s.NoError(testutils.DecodeJSON(listResp2, &listResult2))
	s.Equal(float64(2), listResult2["totalCount"])
	items2 := listResult2["items"].([]interface{})
	s.Len(items2, 2)

	got := map[string]bool{}
	for _, it := range items2 {
		id, _ := it.(map[string]interface{})["id"].(string)
		got[id] = true
	}
	s.True(got[instagramID])
	s.True(got[whatsappID])
}

func (s *OrderSuite) TestListOrders_Search_TooLong() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	_, err = s.orderHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	long := strings.Repeat("a", 400)
	resp, err := s.orderHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/orders?search="+long, nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *OrderSuite) TestListOrders_Search_CrossWorkspaceIsolation() {
	ctx := context.Background()

	// Workspace 1
	_, ws1, token1, err := s.accountHelper.CreateTestUser(ctx, "admin1@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws1.ID))
	biz1, err := s.orderHelper.CreateTestBusiness(ctx, ws1.ID, "biz-1")
	s.NoError(err)
	c1, a1, err := s.orderHelper.CreateTestCustomer(ctx, biz1.ID, "one@example.com", "Shared Name")
	s.NoError(err)
	cat1, err := s.orderHelper.CreateTestCategory(ctx, biz1.ID, "Cat", "cat")
	s.NoError(err)
	_, v1, err := s.orderHelper.CreateTestProduct(ctx, biz1.ID, cat1.ID, "P1", decimal.NewFromInt(1), decimal.NewFromInt(2), 10)
	s.NoError(err)

	resp1, err := s.orderHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/biz-1/orders", map[string]interface{}{
		"customerId":        c1.ID,
		"shippingAddressId": a1.ID,
		"channel":           "instagram",
		"items": []map[string]interface{}{{
			"variantId": v1.ID,
			"quantity":  1,
			"unitPrice": 2,
			"unitCost":  1,
		}},
	}, token1)
	s.NoError(err)
	resp1.Body.Close()
	s.Equal(http.StatusCreated, resp1.StatusCode)

	// Workspace 2
	_, ws2, token2, err := s.accountHelper.CreateTestUser(ctx, "admin2@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws2.ID))
	biz2, err := s.orderHelper.CreateTestBusiness(ctx, ws2.ID, "biz-2")
	s.NoError(err)
	c2, a2, err := s.orderHelper.CreateTestCustomer(ctx, biz2.ID, "two@example.com", "Shared Name")
	s.NoError(err)
	cat2, err := s.orderHelper.CreateTestCategory(ctx, biz2.ID, "Cat", "cat")
	s.NoError(err)
	_, v2, err := s.orderHelper.CreateTestProduct(ctx, biz2.ID, cat2.ID, "P2", decimal.NewFromInt(1), decimal.NewFromInt(2), 10)
	s.NoError(err)

	resp2, err := s.orderHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/biz-2/orders", map[string]interface{}{
		"customerId":        c2.ID,
		"shippingAddressId": a2.ID,
		"channel":           "instagram",
		"items": []map[string]interface{}{{
			"variantId": v2.ID,
			"quantity":  1,
			"unitPrice": 2,
			"unitCost":  1,
		}},
	}, token2)
	s.NoError(err)
	resp2.Body.Close()
	s.Equal(http.StatusCreated, resp2.StatusCode)

	// Search within each business must not leak.
	list1, err := s.orderHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/biz-1/orders?search=Shared%20Name", nil, token1)
	s.NoError(err)
	defer list1.Body.Close()
	var r1 map[string]interface{}
	s.NoError(testutils.DecodeJSON(list1, &r1))
	s.Equal(float64(1), r1["totalCount"])

	list2, err := s.orderHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/biz-2/orders?search=Shared%20Name", nil, token2)
	s.NoError(err)
	defer list2.Body.Close()
	var r2 map[string]interface{}
	s.NoError(testutils.DecodeJSON(list2, &r2))
	s.Equal(float64(1), r2["totalCount"])
}

func TestOrderSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(OrderSuite))
}
