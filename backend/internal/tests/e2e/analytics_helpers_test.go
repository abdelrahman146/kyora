package e2e_test

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/accounting"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
)

// AnalyticsTestHelper provides reusable helpers for analytics E2E tests
type AnalyticsTestHelper struct {
	db     *database.Database
	Client *testutils.HTTPClient
}

func NewAnalyticsTestHelper(db *database.Database, baseURL string) *AnalyticsTestHelper {
	return &AnalyticsTestHelper{
		db:     db,
		Client: testutils.NewHTTPClient(baseURL),
	}
}

func (h *AnalyticsTestHelper) CreateTestBusiness(ctx context.Context, workspaceID, descriptor string) (*business.Business, error) {
	bizRepo := database.NewRepository[business.Business](h.db)
	biz := &business.Business{
		WorkspaceID:   workspaceID,
		Descriptor:    descriptor,
		Name:          "Test Business",
		CountryCode:   "EG",
		Currency:      "USD",
		VatRate:       decimal.NewFromFloat(0.14),
		SafetyBuffer:  decimal.NewFromFloat(100),
		EstablishedAt: time.Now().UTC(),
	}
	if err := bizRepo.CreateOne(ctx, biz); err != nil {
		return nil, err
	}
	return biz, nil
}

func (h *AnalyticsTestHelper) CreateTestCategory(ctx context.Context, businessID, name, descriptor string) (*inventory.Category, error) {
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

func (h *AnalyticsTestHelper) CreateTestProduct(ctx context.Context, businessID, categoryID, name string) (*inventory.Product, error) {
	productRepo := database.NewRepository[inventory.Product](h.db)
	product := &inventory.Product{
		BusinessID:  businessID,
		CategoryID:  categoryID,
		Name:        name,
		Description: name + " description",
	}
	if err := productRepo.CreateOne(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (h *AnalyticsTestHelper) CreateTestVariant(ctx context.Context, businessID, productID, name string, cost, price decimal.Decimal, stock int) (*inventory.Variant, error) {
	variantRepo := database.NewRepository[inventory.Variant](h.db)
	variant := &inventory.Variant{
		ProductID:          productID,
		BusinessID:         businessID,
		Name:               name,
		CostPrice:          cost,
		SalePrice:          price,
		StockQuantity:      stock,
		StockQuantityAlert: 5,
		Currency:           "USD",
	}
	if err := variantRepo.CreateOne(ctx, variant); err != nil {
		return nil, err
	}
	return variant, nil
}

func (h *AnalyticsTestHelper) CreateTestCustomer(ctx context.Context, businessID, email, name string) (*customer.Customer, error) {
	customerRepo := database.NewRepository[customer.Customer](h.db)
	cust := &customer.Customer{
		BusinessID:  businessID,
		Email:       transformer.ToNullableString(email),
		Name:        name,
		CountryCode: "EG",
		Gender:      customer.GenderMale,
	}
	if err := customerRepo.CreateOne(ctx, cust); err != nil {
		return nil, err
	}
	return cust, nil
}

func (h *AnalyticsTestHelper) CreateTestCustomerAt(ctx context.Context, businessID, email, name string, createdAt time.Time) (*customer.Customer, error) {
	customerRepo := database.NewRepository[customer.Customer](h.db)
	cust := &customer.Customer{
		BusinessID:  businessID,
		Email:       transformer.ToNullableString(email),
		Name:        name,
		CountryCode: "EG",
		Gender:      customer.GenderMale,
	}
	cust.CreatedAt = createdAt
	if err := customerRepo.CreateOne(ctx, cust); err != nil {
		return nil, err
	}
	return cust, nil
}

func (h *AnalyticsTestHelper) CreateTestAddress(ctx context.Context, customerID string) (*customer.CustomerAddress, error) {
	addressRepo := database.NewRepository[customer.CustomerAddress](h.db)
	addr := &customer.CustomerAddress{
		CustomerID:  customerID,
		CountryCode: "EG",
		State:       "Cairo",
		City:        "Cairo",
		Street:      transformer.ToNullableString("123 Test St"),
		PhoneCode:   "+20",
		PhoneNumber: "1234567890",
	}
	if err := addressRepo.CreateOne(ctx, addr); err != nil {
		return nil, err
	}
	return addr, nil
}

type OrderItemData struct {
	VariantID string
	Quantity  int
	UnitPrice decimal.Decimal
	UnitCost  decimal.Decimal
}

func (h *AnalyticsTestHelper) CreateTestOrder(ctx context.Context, businessID, customerID, shippingAddressID string, channel string, status order.OrderStatus, items []OrderItemData, createdAt time.Time) (*order.Order, error) {
	orderRepo := database.NewRepository[order.Order](h.db)
	orderItemRepo := database.NewRepository[order.OrderItem](h.db)

	// Calculate order totals
	var subtotal, cogs decimal.Decimal
	for _, item := range items {
		itemTotal := item.UnitPrice.Mul(decimal.NewFromInt(int64(item.Quantity)))
		itemCogs := item.UnitCost.Mul(decimal.NewFromInt(int64(item.Quantity)))
		subtotal = subtotal.Add(itemTotal)
		cogs = cogs.Add(itemCogs)
	}

	ord := &order.Order{
		BusinessID:        businessID,
		OrderNumber:       fmt.Sprintf("ord-%d", createdAt.UnixNano()),
		CustomerID:        customerID,
		ShippingAddressID: shippingAddressID,
		Channel:           channel,
		Status:            status,
		Subtotal:          subtotal,
		Total:             subtotal,
		COGS:              cogs,
		Currency:          "USD",
	}
	ord.CreatedAt = createdAt
	ord.OrderedAt = createdAt
	if err := orderRepo.CreateOne(ctx, ord); err != nil {
		return nil, err
	}

	// Create order items
	variantRepo := database.NewRepository[inventory.Variant](h.db)
	for _, itemData := range items {
		// Get variant to retrieve product ID
		variant, err := variantRepo.FindByID(ctx, itemData.VariantID)
		if err != nil {
			return nil, err
		}

		item := &order.OrderItem{
			OrderID:   ord.ID,
			ProductID: variant.ProductID,
			VariantID: itemData.VariantID,
			Quantity:  itemData.Quantity,
			UnitPrice: itemData.UnitPrice,
			UnitCost:  itemData.UnitCost,
			Total:     itemData.UnitPrice.Mul(decimal.NewFromInt(int64(itemData.Quantity))),
			TotalCost: itemData.UnitCost.Mul(decimal.NewFromInt(int64(itemData.Quantity))),
			Currency:  "USD",
		}
		if err := orderItemRepo.CreateOne(ctx, item); err != nil {
			return nil, err
		}
	}

	return ord, nil
}

func (h *AnalyticsTestHelper) CreateTestExpense(ctx context.Context, businessID string, category accounting.ExpenseCategory, amount decimal.Decimal, occurredOn time.Time) (*accounting.Expense, error) {
	expenseRepo := database.NewRepository[accounting.Expense](h.db)
	expense := &accounting.Expense{
		BusinessID: businessID,
		Category:   category,
		Type:       accounting.ExpenseTypeOneTime,
		Amount:     amount,
		OccurredOn: occurredOn,
		Note:       sql.NullString{String: "Test expense", Valid: true},
	}
	if err := expenseRepo.CreateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (h *AnalyticsTestHelper) CreateTestInvestment(ctx context.Context, businessID string, investorID string, amount decimal.Decimal, investedAt time.Time) (*accounting.Investment, error) {
	investmentRepo := database.NewRepository[accounting.Investment](h.db)
	investment := &accounting.Investment{
		BusinessID: businessID,
		InvestorID: investorID,
		Amount:     amount,
		InvestedAt: investedAt,
		Note:       "Test investment",
	}
	if err := investmentRepo.CreateOne(ctx, investment); err != nil {
		return nil, err
	}
	return investment, nil
}

func (h *AnalyticsTestHelper) CreateTestWithdrawal(ctx context.Context, businessID string, withdrawerID string, amount decimal.Decimal, withdrawnAt time.Time) (*accounting.Withdrawal, error) {
	withdrawalRepo := database.NewRepository[accounting.Withdrawal](h.db)
	withdrawal := &accounting.Withdrawal{
		BusinessID:   businessID,
		WithdrawerID: withdrawerID,
		Amount:       amount,
		WithdrawnAt:  withdrawnAt,
		Note:         "Test withdrawal",
	}
	if err := withdrawalRepo.CreateOne(ctx, withdrawal); err != nil {
		return nil, err
	}
	return withdrawal, nil
}

func (h *AnalyticsTestHelper) CreateTestAsset(ctx context.Context, businessID string, assetType accounting.AssetType, name string, value decimal.Decimal, purchasedAt time.Time) (*accounting.Asset, error) {
	assetRepo := database.NewRepository[accounting.Asset](h.db)
	asset := &accounting.Asset{
		BusinessID:  businessID,
		Name:        name,
		Type:        assetType,
		Value:       value,
		PurchasedAt: purchasedAt,
	}
	if err := assetRepo.CreateOne(ctx, asset); err != nil {
		return nil, err
	}
	return asset, nil
}
