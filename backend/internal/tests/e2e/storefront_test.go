package e2e_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type StorefrontSuite struct {
	suite.Suite
	client *testutils.HTTPClient
	db     *database.Database
}

func (s *StorefrontSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient(e2eBaseURL)
	s.db = testEnv.Database
}

func (s *StorefrontSuite) SetupTest() {
	s.NoError(testutils.TruncateTables(s.db,
		"storefront_requests",
		"order_notes",
		"order_items",
		"orders",
		"customer_addresses",
		"customers",
		"variants",
		"products",
		"categories",
		"businesses",
		"workspaces",
		"users",
	))
}

func (s *StorefrontSuite) TearDownTest() {
	s.NoError(testutils.TruncateTables(s.db,
		"storefront_requests",
		"order_notes",
		"order_items",
		"orders",
		"customer_addresses",
		"customers",
		"variants",
		"products",
		"categories",
		"businesses",
		"workspaces",
		"users",
	))
}

func (s *StorefrontSuite) createWorkspace(ctx context.Context) *account.Workspace {
	repo := database.NewRepository[account.Workspace](s.db)
	ws := &account.Workspace{OwnerID: "usr_owner"}
	s.NoError(repo.CreateOne(ctx, ws))
	return ws
}

func (s *StorefrontSuite) createBusiness(ctx context.Context, wsID, descriptor string, enabled bool) *business.Business {
	repo := database.NewRepository[business.Business](s.db)
	biz := &business.Business{
		WorkspaceID:       wsID,
		Descriptor:        descriptor,
		Name:              "Test Business",
		Brand:             "Brand",
		CountryCode:       "EG",
		Currency:          "USD",
		StorefrontEnabled: enabled,
		StorefrontTheme:   business.StorefrontTheme{},
		SupportEmail:      "support@example.com",
		PhoneNumber:       "+201234567890",
		WhatsappNumber:    "+201234567890",
		WebsiteURL:        "https://example.com",
		InstagramURL:      "https://instagram.com/example",
		FacebookURL:       "https://facebook.com/example",
		TikTokURL:         "https://tiktok.com/@example",
		XURL:              "https://x.com/example",
		SnapchatURL:       "https://snapchat.com/add/example",
		VatRate:           decimal.RequireFromString("0.14"),
		SafetyBuffer:      decimal.RequireFromString("100"),
	}
	s.NoError(repo.CreateOne(ctx, biz))
	return biz
}

func (s *StorefrontSuite) createCategory(ctx context.Context, businessID string) *inventory.Category {
	repo := database.NewRepository[inventory.Category](s.db)
	cat := &inventory.Category{BusinessID: businessID, Name: "Cat", Descriptor: "cat"}
	s.NoError(repo.CreateOne(ctx, cat))
	return cat
}

func (s *StorefrontSuite) createProduct(ctx context.Context, businessID, categoryID string) *inventory.Product {
	repo := database.NewRepository[inventory.Product](s.db)
	p := &inventory.Product{BusinessID: businessID, CategoryID: categoryID, Name: "Product", Description: "Desc"}
	s.NoError(repo.CreateOne(ctx, p))
	return p
}

func (s *StorefrontSuite) createVariant(ctx context.Context, businessID, productID string, stock int) *inventory.Variant {
	repo := database.NewRepository[inventory.Variant](s.db)
	v := &inventory.Variant{
		BusinessID:         businessID,
		ProductID:          productID,
		Code:               "DEFAULT",
		Name:               "Product - DEFAULT",
		SKU:                "SKU-1",
		CostPrice:          decimal.RequireFromString("10"),
		SalePrice:          decimal.RequireFromString("25"),
		Currency:           "USD",
		StockQuantity:      stock,
		StockQuantityAlert: 1,
	}
	s.NoError(repo.CreateOne(ctx, v))
	return v
}

func (s *StorefrontSuite) TestGetCatalog_Success() {
	ctx := context.Background()
	ws := s.createWorkspace(ctx)
	biz := s.createBusiness(ctx, ws.ID, "biz", true)
	cat := s.createCategory(ctx, biz.ID)
	prod := s.createProduct(ctx, biz.ID, cat.ID)
	variant := s.createVariant(ctx, biz.ID, prod.ID, 100)

	resp, err := s.client.Get("/v1/storefront/" + biz.StorefrontPublicID + "/catalog")
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Len(result, 3)
	s.Contains(result, "business")
	s.Contains(result, "categories")
	s.Contains(result, "products")

	b := result["business"].(map[string]interface{})
	s.Equal(biz.ID, b["id"])
	s.Equal(biz.StorefrontPublicID, b["storefrontPublicId"])
	s.Equal(true, b["storefrontEnabled"])
	s.Equal("EG", b["countryCode"])
	s.Equal("USD", b["currency"])

	cats := result["categories"].([]interface{})
	s.Len(cats, 1)
	c0 := cats[0].(map[string]interface{})
	s.Equal(cat.ID, c0["id"])

	prods := result["products"].([]interface{})
	s.Len(prods, 1)
	p0 := prods[0].(map[string]interface{})
	s.Equal(prod.ID, p0["id"])
	vars := p0["variants"].([]interface{})
	s.Len(vars, 1)
	v0 := vars[0].(map[string]interface{})
	s.Equal(variant.ID, v0["id"])
	s.Equal("25", v0["salePrice"])
	s.Equal("USD", v0["currency"])
}

func (s *StorefrontSuite) TestGetCatalog_DisabledStorefront() {
	ctx := context.Background()
	ws := s.createWorkspace(ctx)
	biz := s.createBusiness(ctx, ws.ID, "biz", false)

	resp, err := s.client.Get("/v1/storefront/" + biz.StorefrontPublicID + "/catalog")
	s.NoError(err)
	defer resp.Body.Close()
	// Disabled storefronts are intentionally hidden by the business lookup (treated as not found).
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *StorefrontSuite) TestCreateOrder_IdempotentAndCreatesNote() {
	ctx := context.Background()
	ws := s.createWorkspace(ctx)
	biz := s.createBusiness(ctx, ws.ID, "biz", true)
	cat := s.createCategory(ctx, biz.ID)
	prod := s.createProduct(ctx, biz.ID, cat.ID)
	variant := s.createVariant(ctx, biz.ID, prod.ID, 1000)

	payload := map[string]interface{}{
		"customer": map[string]interface{}{
			"email":             "buyer@example.com",
			"name":              "Buyer",
			"phoneNumber":       "+201111111111",
			"instagramUsername": "buyer_ig",
		},
		"shippingAddress": map[string]interface{}{
			"countryCode": "EG",
			"state":       "Cairo",
			"city":        "Cairo",
			"street":      "Test St",
			"zipCode":     "12345",
			"phoneCode":   "+20",
			"phoneNumber": "1111111111",
		},
		"items": []map[string]interface{}{
			{"variantId": variant.ID, "quantity": 2, "specialRequest": "No onions"},
		},
	}
	body, sErr := json.Marshal(payload)
	s.NoError(sErr)

	headers := map[string]string{
		"Content-Type":    "application/json",
		"Idempotency-Key": "idem-1",
	}

	resp, err := s.client.PostRaw("/v1/storefront/"+biz.StorefrontPublicID+"/orders", body, headers)
	s.NoError(err)
	s.Equal(http.StatusCreated, resp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &created))
	s.Len(created, 6)
	s.Contains(created, "orderId")
	s.Contains(created, "orderNumber")
	s.Equal("pending", created["status"])
	s.Equal("pending", created["paymentStatus"])
	s.Equal("USD", created["currency"])
	s.Equal("57", created["total"], "total should be server-side computed (2 * 25 + 14% VAT)")
	orderID := created["orderId"].(string)

	// Replay same request with same idempotency key => same order id, no new order.
	replay, err := s.client.PostRaw("/v1/storefront/"+biz.StorefrontPublicID+"/orders", body, headers)
	s.NoError(err)
	s.Equal(http.StatusCreated, replay.StatusCode)

	var replayed map[string]interface{}
	s.NoError(testutils.DecodeJSON(replay, &replayed))
	s.Equal(orderID, replayed["orderId"])

	orderRepo := database.NewRepository[order.Order](s.db)
	count, err := orderRepo.Count(ctx, orderRepo.ScopeBusinessID(biz.ID))
	s.NoError(err)
	s.Equal(int64(1), count)

	// Customer is upserted by email in this business.
	custRepo := database.NewRepository[customer.Customer](s.db)
	custEmail := strings.ToLower("buyer@example.com")
	cust, err := custRepo.FindOne(ctx,
		custRepo.ScopeBusinessID(biz.ID),
		custRepo.ScopeEquals(customer.CustomerSchema.Email, custEmail),
	)
	s.NoError(err)
	s.Equal("Buyer", cust.Name)
	s.Equal(transformer.ToNullableString(custEmail), cust.Email)

	// Note created and includes the special request.
	noteRepo := database.NewRepository[order.OrderNote](s.db)
	notes, err := noteRepo.FindMany(ctx, noteRepo.ScopeEquals(order.OrderNoteSchema.OrderID, orderID))
	s.NoError(err)
	s.Len(notes, 1)
	s.Contains(notes[0].Content, "Special requests:")
	s.Contains(notes[0].Content, "No onions")
}

func (s *StorefrontSuite) TestCreateOrder_IdempotencyHashMismatch() {
	ctx := context.Background()
	ws := s.createWorkspace(ctx)
	biz := s.createBusiness(ctx, ws.ID, "biz", true)
	cat := s.createCategory(ctx, biz.ID)
	prod := s.createProduct(ctx, biz.ID, cat.ID)
	variant := s.createVariant(ctx, biz.ID, prod.ID, 1000)

	payload1 := map[string]interface{}{
		"customer": map[string]interface{}{"email": "buyer@example.com", "name": "Buyer"},
		"shippingAddress": map[string]interface{}{
			"countryCode": "EG",
			"state":       "Cairo",
			"city":        "Cairo",
			"phoneCode":   "+20",
			"phoneNumber": "1111111111",
		},
		"items": []map[string]interface{}{{"variantId": variant.ID, "quantity": 1}},
	}
	payload2 := map[string]interface{}{
		"customer": map[string]interface{}{"email": "buyer@example.com", "name": "Buyer"},
		"shippingAddress": map[string]interface{}{
			"countryCode": "EG",
			"state":       "Cairo",
			"city":        "Cairo",
			"phoneCode":   "+20",
			"phoneNumber": "1111111111",
		},
		"items": []map[string]interface{}{{"variantId": variant.ID, "quantity": 2}},
	}
	b1, _ := json.Marshal(payload1)
	b2, _ := json.Marshal(payload2)
	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "idem-x"}

	resp1, err := s.client.PostRaw("/v1/storefront/"+biz.StorefrontPublicID+"/orders", b1, headers)
	s.NoError(err)
	s.Equal(http.StatusCreated, resp1.StatusCode)

	resp2, err := s.client.PostRaw("/v1/storefront/"+biz.StorefrontPublicID+"/orders", b2, headers)
	s.NoError(err)
	s.Equal(http.StatusConflict, resp2.StatusCode)

	orderRepo := database.NewRepository[order.Order](s.db)
	count, err := orderRepo.Count(ctx, orderRepo.ScopeBusinessID(biz.ID))
	s.NoError(err)
	s.Equal(int64(1), count)
}

func (s *StorefrontSuite) TestCreateOrder_PreventsCrossTenantVariant() {
	ctx := context.Background()
	ws := s.createWorkspace(ctx)
	bizA := s.createBusiness(ctx, ws.ID, "biz-a", true)
	bizB := s.createBusiness(ctx, ws.ID, "biz-b", true)

	catA := s.createCategory(ctx, bizA.ID)
	prodA := s.createProduct(ctx, bizA.ID, catA.ID)
	s.createVariant(ctx, bizA.ID, prodA.ID, 1000)

	catB := s.createCategory(ctx, bizB.ID)
	prodB := s.createProduct(ctx, bizB.ID, catB.ID)
	variantB := s.createVariant(ctx, bizB.ID, prodB.ID, 1000)

	payload := map[string]interface{}{
		"customer": map[string]interface{}{"email": "buyer@example.com", "name": "Buyer"},
		"shippingAddress": map[string]interface{}{
			"countryCode": "EG",
			"state":       "Cairo",
			"city":        "Cairo",
			"phoneCode":   "+20",
			"phoneNumber": "1111111111",
		},
		"items": []map[string]interface{}{{"variantId": variantB.ID, "quantity": 1}},
	}
	body, _ := json.Marshal(payload)
	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "idem-ct"}

	resp, err := s.client.PostRaw("/v1/storefront/"+bizA.StorefrontPublicID+"/orders", body, headers)
	s.NoError(err)
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *StorefrontSuite) TestCreateOrder_RejectsUnknownFields_MassAssignmentProtection() {
	ctx := context.Background()
	ws := s.createWorkspace(ctx)
	biz := s.createBusiness(ctx, ws.ID, "biz", true)
	cat := s.createCategory(ctx, biz.ID)
	prod := s.createProduct(ctx, biz.ID, cat.ID)
	variant := s.createVariant(ctx, biz.ID, prod.ID, 1000)

	payload := map[string]interface{}{
		"customer": map[string]interface{}{"email": "buyer@example.com", "name": "Buyer"},
		"shippingAddress": map[string]interface{}{
			"countryCode": "EG",
			"state":       "Cairo",
			"city":        "Cairo",
			"phoneCode":   "+20",
			"phoneNumber": "1111111111",
		},
		"items":  []map[string]interface{}{{"variantId": variant.ID, "quantity": 1}},
		"status": "paid", // not allowed
	}
	body, _ := json.Marshal(payload)
	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "idem-ma"}

	resp, err := s.client.PostRaw("/v1/storefront/"+biz.StorefrontPublicID+"/orders", body, headers)
	s.NoError(err)
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *StorefrontSuite) TestCreateOrder_RateLimited() {
	ctx := context.Background()
	ws := s.createWorkspace(ctx)
	biz := s.createBusiness(ctx, ws.ID, "biz", true)
	cat := s.createCategory(ctx, biz.ID)
	prod := s.createProduct(ctx, biz.ID, cat.ID)
	variant := s.createVariant(ctx, biz.ID, prod.ID, 1000)

	payload := map[string]interface{}{
		"customer": map[string]interface{}{"email": "buyer@example.com", "name": "Buyer"},
		"shippingAddress": map[string]interface{}{
			"countryCode": "EG",
			"state":       "Cairo",
			"city":        "Cairo",
			"phoneCode":   "+20",
			"phoneNumber": "1111111111",
		},
		"items": []map[string]interface{}{{"variantId": variant.ID, "quantity": 1}},
	}
	body, _ := json.Marshal(payload)

	path := "/v1/storefront/" + biz.StorefrontPublicID + "/orders"
	{
		headers := map[string]string{
			"Content-Type":    "application/json",
			"Idempotency-Key": "rl-0",
		}
		resp, err := s.client.PostRaw(path, body, headers)
		s.NoError(err)
		s.Equal(http.StatusCreated, resp.StatusCode)
		resp.Body.Close()
	}

	// Storefront rate limiting includes a min-interval (1s) to prevent bursts.
	// A second request immediately after a successful one should be rate-limited.
	{
		headers := map[string]string{
			"Content-Type":    "application/json",
			"Idempotency-Key": "rl-1",
		}
		resp, err := s.client.PostRaw(path, body, headers)
		s.NoError(err)
		s.Equal(http.StatusTooManyRequests, resp.StatusCode)
		resp.Body.Close()
	}
}

func TestStorefrontSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(StorefrontSuite))
}
