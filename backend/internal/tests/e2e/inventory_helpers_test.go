package e2e_test

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
)

var inventoryTables = []string{
	"users",
	"workspaces",
	"businesses",
	"categories",
	"products",
	"variants",
	"uploaded_assets",
	"subscriptions",
	"plans",
}

// InventoryTestHelper provides reusable helpers for inventory E2E tests.
type InventoryTestHelper struct {
	db     *database.Database
	Client *testutils.HTTPClient
}

func NewInventoryTestHelper(db *database.Database, baseURL string) *InventoryTestHelper {
	return &InventoryTestHelper{
		db:     db,
		Client: testutils.NewHTTPClient(baseURL),
	}
}

func (h *InventoryTestHelper) CreateTestBusiness(ctx context.Context, workspaceID, descriptor string) (*business.Business, error) {
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

func (h *InventoryTestHelper) CreateTestCategory(ctx context.Context, businessID, name, descriptor string) (*inventory.Category, error) {
	repo := database.NewRepository[inventory.Category](h.db)
	cat := &inventory.Category{BusinessID: businessID, Name: name, Descriptor: descriptor}
	if err := repo.CreateOne(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (h *InventoryTestHelper) CreateTestProduct(ctx context.Context, businessID, categoryID, name, description string) (*inventory.Product, error) {
	repo := database.NewRepository[inventory.Product](h.db)
	prod := &inventory.Product{BusinessID: businessID, Name: name, Description: description, CategoryID: categoryID}
	if err := repo.CreateOne(ctx, prod); err != nil {
		return nil, err
	}
	return prod, nil
}

func (h *InventoryTestHelper) CreateTestVariant(
	ctx context.Context,
	businessID, productID, code, sku, currency string,
	costPrice, salePrice decimal.Decimal,
	stockQuantity, stockQuantityAlert int,
) (*inventory.Variant, error) {
	repo := database.NewRepository[inventory.Variant](h.db)
	variant := &inventory.Variant{
		BusinessID:         businessID,
		ProductID:          productID,
		Code:               code,
		Name:               "variant",
		SKU:                sku,
		CostPrice:          costPrice,
		SalePrice:          salePrice,
		Currency:           currency,
		StockQuantity:      stockQuantity,
		StockQuantityAlert: stockQuantityAlert,
	}
	if err := repo.CreateOne(ctx, variant); err != nil {
		return nil, err
	}
	return variant, nil
}

func (h *InventoryTestHelper) GetCategory(ctx context.Context, id string) (*inventory.Category, error) {
	repo := database.NewRepository[inventory.Category](h.db)
	return repo.FindByID(ctx, id)
}

func (h *InventoryTestHelper) GetProduct(ctx context.Context, id string) (*inventory.Product, error) {
	repo := database.NewRepository[inventory.Product](h.db)
	return repo.FindByID(ctx, id)
}

func (h *InventoryTestHelper) GetVariant(ctx context.Context, id string) (*inventory.Variant, error) {
	repo := database.NewRepository[inventory.Variant](h.db)
	return repo.FindByID(ctx, id)
}

func (h *InventoryTestHelper) CountProducts(ctx context.Context, businessID string) (int64, error) {
	repo := database.NewRepository[inventory.Product](h.db)
	return repo.Count(ctx, repo.ScopeBusinessID(businessID))
}

func (h *InventoryTestHelper) CountVariants(ctx context.Context, businessID string) (int64, error) {
	repo := database.NewRepository[inventory.Variant](h.db)
	return repo.Count(ctx, repo.ScopeBusinessID(businessID))
}
