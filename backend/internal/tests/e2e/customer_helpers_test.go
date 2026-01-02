package e2e_test

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
)

// CustomerTestHelper provides reusable helpers for customer tests
type CustomerTestHelper struct {
	db     *database.Database
	Client *testutils.HTTPClient
}

// NewCustomerTestHelper creates a new customer test helper
func NewCustomerTestHelper(db *database.Database, baseURL string) *CustomerTestHelper {
	return &CustomerTestHelper{
		db:     db,
		Client: testutils.NewHTTPClient(baseURL),
	}
}

// CreateTestBusiness creates a business for testing
func (h *CustomerTestHelper) CreateTestBusiness(ctx context.Context, workspaceID, descriptor string) (*business.Business, error) {
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

// CreateTestCustomer creates a customer for testing
func (h *CustomerTestHelper) CreateTestCustomer(ctx context.Context, businessID, email, name string) (*customer.Customer, error) {
	customerRepo := database.NewRepository[customer.Customer](h.db)
	cust := &customer.Customer{
		BusinessID:  businessID,
		Name:        name,
		Email:       transformer.ToNullableString(email),
		CountryCode: "EG",
		Gender:      customer.GenderMale,
	}
	if err := customerRepo.CreateOne(ctx, cust); err != nil {
		return nil, err
	}
	return cust, nil
}

// GetCustomer retrieves a customer by ID
func (h *CustomerTestHelper) GetCustomer(ctx context.Context, customerID string) (*customer.Customer, error) {
	customerRepo := database.NewRepository[customer.Customer](h.db)
	return customerRepo.FindByID(ctx, customerID)
}

// CreateTestAddress creates an address for testing
func (h *CustomerTestHelper) CreateTestAddress(ctx context.Context, customerID string) (*customer.CustomerAddress, error) {
	addressRepo := database.NewRepository[customer.CustomerAddress](h.db)
	addr := &customer.CustomerAddress{
		CustomerID:  customerID,
		CountryCode: "EG",
		State:       "Cairo",
		City:        "Cairo",
		Street:      transformer.ToNullableString("123 Test St"),
		PhoneCode:   "+20",
		PhoneNumber: "1234567890",
		ZipCode:     transformer.ToNullableString("12345"),
	}
	if err := addressRepo.CreateOne(ctx, addr); err != nil {
		return nil, err
	}
	return addr, nil
}

// GetAddress retrieves an address by ID
func (h *CustomerTestHelper) GetAddress(ctx context.Context, addressID string) (*customer.CustomerAddress, error) {
	addressRepo := database.NewRepository[customer.CustomerAddress](h.db)
	return addressRepo.FindByID(ctx, addressID)
}

// CreateTestNote creates a note for testing
func (h *CustomerTestHelper) CreateTestNote(ctx context.Context, customerID, content string) (*customer.CustomerNote, error) {
	noteRepo := database.NewRepository[customer.CustomerNote](h.db)
	note := &customer.CustomerNote{
		CustomerID: customerID,
		Content:    content,
	}
	if err := noteRepo.CreateOne(ctx, note); err != nil {
		return nil, err
	}
	return note, nil
}

// GetNote retrieves a note by ID
func (h *CustomerTestHelper) GetNote(ctx context.Context, noteID string) (*customer.CustomerNote, error) {
	noteRepo := database.NewRepository[customer.CustomerNote](h.db)
	return noteRepo.FindByID(ctx, noteID)
}

// CountCustomers counts customers for a business
func (h *CustomerTestHelper) CountCustomers(ctx context.Context, businessID string) (int64, error) {
	customerRepo := database.NewRepository[customer.Customer](h.db)
	return customerRepo.Count(ctx, customerRepo.ScopeBusinessID(businessID))
}

// CountAddresses counts addresses for a customer
func (h *CustomerTestHelper) CountAddresses(ctx context.Context, customerID string) (int64, error) {
	addressRepo := database.NewRepository[customer.CustomerAddress](h.db)
	return addressRepo.Count(ctx, addressRepo.ScopeEquals(customer.CustomerAddressSchema.CustomerID, customerID))
}

// CountNotes counts notes for a customer
func (h *CustomerTestHelper) CountNotes(ctx context.Context, customerID string) (int64, error) {
	noteRepo := database.NewRepository[customer.CustomerNote](h.db)
	return noteRepo.Count(ctx, noteRepo.ScopeEquals(customer.CustomerNoteSchema.CustomerID, customerID))
}

// CreateTestCustomerWithCountry creates a customer with a specific country
func (h *CustomerTestHelper) CreateTestCustomerWithCountry(ctx context.Context, businessID, email, name, countryCode string) (*customer.Customer, error) {
	customerRepo := database.NewRepository[customer.Customer](h.db)
	cust := &customer.Customer{
		BusinessID:  businessID,
		Name:        name,
		Email:       transformer.ToNullableString(email),
		CountryCode: countryCode,
		Gender:      customer.GenderMale,
	}
	if err := customerRepo.CreateOne(ctx, cust); err != nil {
		return nil, err
	}
	return cust, nil
}

// CreateTestCustomerWithSocial creates a customer with a specific social media platform
func (h *CustomerTestHelper) CreateTestCustomerWithSocial(ctx context.Context, businessID, email, name, platform, username string) (*customer.Customer, error) {
	customerRepo := database.NewRepository[customer.Customer](h.db)
	cust := &customer.Customer{
		BusinessID:  businessID,
		Name:        name,
		Email:       transformer.ToNullableString(email),
		CountryCode: "EG",
		Gender:      customer.GenderMale,
	}

	switch platform {
	case "instagram":
		cust.InstagramUsername = transformer.ToNullableString(username)
	case "tiktok":
		cust.TikTokUsername = transformer.ToNullableString(username)
	case "facebook":
		cust.FacebookUsername = transformer.ToNullableString(username)
	case "x":
		cust.XUsername = transformer.ToNullableString(username)
	case "snapchat":
		cust.SnapchatUsername = transformer.ToNullableString(username)
	case "whatsapp":
		cust.WhatsappNumber = transformer.ToNullableString(username)
	}

	if err := customerRepo.CreateOne(ctx, cust); err != nil {
		return nil, err
	}
	return cust, nil
}

// SetCustomerSocial updates a customer's social media platform
func (h *CustomerTestHelper) SetCustomerSocial(ctx context.Context, customerID, platform, username string) error {
	customerRepo := database.NewRepository[customer.Customer](h.db)
	cust, err := customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return err
	}

	switch platform {
	case "instagram":
		cust.InstagramUsername = transformer.ToNullableString(username)
	case "tiktok":
		cust.TikTokUsername = transformer.ToNullableString(username)
	case "facebook":
		cust.FacebookUsername = transformer.ToNullableString(username)
	case "x":
		cust.XUsername = transformer.ToNullableString(username)
	case "snapchat":
		cust.SnapchatUsername = transformer.ToNullableString(username)
	case "whatsapp":
		cust.WhatsappNumber = transformer.ToNullableString(username)
	}

	return customerRepo.UpdateOne(ctx, cust)
}
