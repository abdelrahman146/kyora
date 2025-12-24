package e2e_test

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
)

// OrderTestHelper provides reusable helpers for order E2E tests
type OrderTestHelper struct {
	db     *database.Database
	Client *testutils.HTTPClient
}

func NewOrderTestHelper(db *database.Database, baseURL string) *OrderTestHelper {
	return &OrderTestHelper{
		db:     db,
		Client: testutils.NewHTTPClient(baseURL),
	}
}

// CreateTestBusiness creates a business for testing
func (h *OrderTestHelper) CreateTestBusiness(ctx context.Context, workspaceID, descriptor string) (*business.Business, error) {
	bizRepo := database.NewRepository[business.Business](h.db)
	biz := &business.Business{
		WorkspaceID:  workspaceID,
		Descriptor:   descriptor,
		Name:         "Test Business",
		CountryCode:  "EG",
		Currency:     "USD",
		VatRate:      decimal.NewFromFloat(0.14),
		SafetyBuffer: decimal.NewFromFloat(100),
	}
	if err := bizRepo.CreateOne(ctx, biz); err != nil {
		return nil, err
	}
	return biz, nil
}

// CreateTestCustomer creates a customer with an address for testing
func (h *OrderTestHelper) CreateTestCustomer(ctx context.Context, businessID, email, name string) (*customer.Customer, *customer.CustomerAddress, error) {
	customerRepo := database.NewRepository[customer.Customer](h.db)
	addressRepo := database.NewRepository[customer.CustomerAddress](h.db)

	cust := &customer.Customer{
		BusinessID:  businessID,
		Email:       transformer.ToNullableString(email),
		Name:        name,
		CountryCode: "EG",
		Gender:      customer.GenderMale,
	}
	if err := customerRepo.CreateOne(ctx, cust); err != nil {
		return nil, nil, err
	}

	addr := &customer.CustomerAddress{
		CustomerID:  cust.ID,
		CountryCode: "EG",
		State:       "Cairo",
		City:        "Cairo",
		Street:      transformer.ToNullableString("123 Test Street"),
		PhoneCode:   "+20",
		PhoneNumber: "1234567890",
	}
	if err := addressRepo.CreateOne(ctx, addr); err != nil {
		return nil, nil, err
	}

	return cust, addr, nil
}

// CreateTestCategory creates a product category for testing
func (h *OrderTestHelper) CreateTestCategory(ctx context.Context, businessID, name, descriptor string) (*inventory.Category, error) {
	catRepo := database.NewRepository[inventory.Category](h.db)
	cat := &inventory.Category{
		BusinessID: businessID,
		Name:       name,
		Descriptor: descriptor,
	}
	if err := catRepo.CreateOne(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

// CreateTestProduct creates a product with a variant for testing
func (h *OrderTestHelper) CreateTestProduct(ctx context.Context, businessID, categoryID, name string, cost, price decimal.Decimal, stock int) (*inventory.Product, *inventory.Variant, error) {
	productRepo := database.NewRepository[inventory.Product](h.db)
	variantRepo := database.NewRepository[inventory.Variant](h.db)

	product := &inventory.Product{
		BusinessID:  businessID,
		CategoryID:  categoryID,
		Name:        name,
		Description: name + " description",
	}
	if err := productRepo.CreateOne(ctx, product); err != nil {
		return nil, nil, err
	}

	variant := &inventory.Variant{
		ProductID:     product.ID,
		BusinessID:    businessID,
		Name:          name + " Default",
		CostPrice:     cost,
		SalePrice:     price,
		StockQuantity: stock,
		Currency:      "USD",
	}
	if err := variantRepo.CreateOne(ctx, variant); err != nil {
		return nil, nil, err
	}

	return product, variant, nil
}

// CreateTestShippingZone creates a shipping zone for a business.
func (h *OrderTestHelper) CreateTestShippingZone(ctx context.Context, businessID, name string, countries []string, shippingCost, freeShippingThreshold decimal.Decimal) (*business.ShippingZone, error) {
	zoneRepo := database.NewRepository[business.ShippingZone](h.db)
	z := &business.ShippingZone{
		BusinessID:            businessID,
		Name:                  name,
		Countries:             business.CountryCodeList(countries),
		Currency:              "USD",
		ShippingCost:          shippingCost,
		FreeShippingThreshold: freeShippingThreshold,
	}
	if err := zoneRepo.CreateOne(ctx, z); err != nil {
		return nil, err
	}
	return z, nil
}

// GetOrder retrieves an order by ID
func (h *OrderTestHelper) GetOrder(ctx context.Context, orderID string) (*order.Order, error) {
	orderRepo := database.NewRepository[order.Order](h.db)
	return orderRepo.FindByID(ctx, orderID)
}

// GetVariant retrieves a variant by ID
func (h *OrderTestHelper) GetVariant(ctx context.Context, variantID string) (*inventory.Variant, error) {
	variantRepo := database.NewRepository[inventory.Variant](h.db)
	return variantRepo.FindByID(ctx, variantID)
}

// CountOrders counts orders for a business
func (h *OrderTestHelper) CountOrders(ctx context.Context, businessID string) (int64, error) {
	orderRepo := database.NewRepository[order.Order](h.db)
	return orderRepo.Count(ctx, orderRepo.ScopeBusinessID(businessID))
}

// CountOrderNotes counts notes for an order
func (h *OrderTestHelper) CountOrderNotes(ctx context.Context, orderID string) (int64, error) {
	noteRepo := database.NewRepository[order.OrderNote](h.db)
	return noteRepo.Count(ctx, noteRepo.ScopeEquals(order.OrderNoteSchema.OrderID, orderID))
}
