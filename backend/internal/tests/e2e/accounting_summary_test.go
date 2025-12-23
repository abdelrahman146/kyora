package e2e_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/nullable"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

// SummarySuite tests /v1/businesses/:businessDescriptor/accounting/summary endpoint.
type SummarySuite struct {
	suite.Suite
	helper *AccountingTestHelper
}

func (s *SummarySuite) SetupSuite() {
	s.helper = NewAccountingTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *SummarySuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "orders", "order_items", "order_notes", "investments", "withdrawals", "expenses")
}

func (s *SummarySuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "orders", "order_items", "order_notes", "investments", "withdrawals", "expenses")
}

func (s *SummarySuite) TestSummary_EmptyWorkspace() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/summary", nil, ws.AdminToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &body))
	s.Contains(body, "totalAssetValue")
	s.Contains(body, "totalInvestments")
	s.Contains(body, "totalWithdrawals")
	s.Contains(body, "totalExpenses")
	s.Contains(body, "safeToDrawAmount")
	s.Contains(body, "currency")
}

func (s *SummarySuite) TestSummary_ComputedFields_WithInvestmentsWithdrawalsExpenses() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	now := time.Now().UTC()
	// Orders require a valid customer_id (empty string violates FK constraints), so create a customer first.
	customerRepo := database.NewRepository[customer.Customer](testEnv.Database)
	cus := &customer.Customer{
		BusinessID:  ws.Business.ID,
		Name:        "Test Customer",
		CountryCode: "US",
		Gender:      customer.GenderOther,
		Email:       nullable.NewString("customer@example.com"),
		JoinedAt:    now,
	}
	s.NoError(customerRepo.CreateOne(ctx, cus))

	addressRepo := database.NewRepository[customer.CustomerAddress](testEnv.Database)
	addr := &customer.CustomerAddress{
		CustomerID:  cus.ID,
		CountryCode: "US",
		State:       "CA",
		City:        "San Francisco",
		PhoneCode:   "+1",
		PhoneNumber: "5551234",
		Street:      nullable.NewString("1 Market St"),
	}
	s.NoError(addressRepo.CreateOne(ctx, addr))

	// SafeToDrawAmount is based on order revenue/COGS (not investments), so seed at least one order.
	orderRepo := database.NewRepository[order.Order](testEnv.Database)
	ord := &order.Order{
		OrderNumber:       "ord-1",
		BusinessID:        ws.Business.ID,
		CustomerID:        cus.ID,
		ShippingAddressID: addr.ID,
		Channel:           "instagram",
		Currency:          ws.Business.Currency,
		OrderedAt:         now,
		Total:             decimal.RequireFromString("1000"),
		COGS:              decimal.RequireFromString("200"),
		Status:            order.OrderStatusFulfilled,
		PaymentStatus:     order.OrderPaymentStatusPaid,
		PaymentMethod:     order.OrderPaymentMethodBankTransfer,
	}
	s.NoError(orderRepo.CreateOne(ctx, ord))

	inv := map[string]interface{}{"investorId": ws.Admin.ID, "amount": "1000.00", "investedAt": now}
	exp := map[string]interface{}{"category": "software", "type": "one_time", "amount": "100.00", "occurredOn": now}
	wd := map[string]interface{}{"withdrawerId": ws.Admin.ID, "amount": "50.00", "withdrawnAt": now}

	resp1, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/investments", inv, ws.AdminToken)
	s.NoError(err)
	resp1.Body.Close()
	resp2, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/expenses", exp, ws.AdminToken)
	s.NoError(err)
	resp2.Body.Close()
	resp3, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/withdrawals", wd, ws.AdminToken)
	s.NoError(err)
	resp3.Body.Close()

	summaryResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/summary", nil, ws.AdminToken)
	s.NoError(err)
	defer summaryResp.Body.Close()
	s.Equal(http.StatusOK, summaryResp.StatusCode)

	var body map[string]interface{}
	s.NoError(testutils.DecodeJSON(summaryResp, &body))

	s.Equal("1000", body["totalInvestments"])
	s.Equal("50", body["totalWithdrawals"])
	s.Equal("100", body["totalExpenses"])
	// If SafetyBuffer is not set, it defaults to last-30-days expenses
	// 1000 - 200 - 100 - 50 - 100 = 550
	s.Equal("550", body["safeToDrawAmount"])
}

func (s *SummarySuite) TestSummary_DateRange_AppliesToSafeToDrawAndTotals() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	within := time.Date(2025, 12, 10, 12, 0, 0, 0, time.UTC)
	outside := time.Date(2025, 10, 10, 12, 0, 0, 0, time.UTC)

	customerRepo := database.NewRepository[customer.Customer](testEnv.Database)
	cus := &customer.Customer{
		BusinessID:  ws.Business.ID,
		Name:        "Test Customer",
		CountryCode: "US",
		Gender:      customer.GenderOther,
		Email:       nullable.NewString("customer-range@example.com"),
		JoinedAt:    within,
	}
	s.NoError(customerRepo.CreateOne(ctx, cus))

	addressRepo := database.NewRepository[customer.CustomerAddress](testEnv.Database)
	addr := &customer.CustomerAddress{
		CustomerID:  cus.ID,
		CountryCode: "US",
		State:       "CA",
		City:        "San Francisco",
		PhoneCode:   "+1",
		PhoneNumber: "5551234",
		Street:      nullable.NewString("1 Market St"),
	}
	s.NoError(addressRepo.CreateOne(ctx, addr))

	orderRepo := database.NewRepository[order.Order](testEnv.Database)
	ordWithin := &order.Order{
		OrderNumber:       "ord-within",
		BusinessID:        ws.Business.ID,
		CustomerID:        cus.ID,
		ShippingAddressID: addr.ID,
		Channel:           "instagram",
		Currency:          ws.Business.Currency,
		OrderedAt:         within,
		Total:             decimal.RequireFromString("1000"),
		COGS:              decimal.RequireFromString("200"),
		Status:            order.OrderStatusFulfilled,
		PaymentStatus:     order.OrderPaymentStatusPaid,
		PaymentMethod:     order.OrderPaymentMethodBankTransfer,
	}
	ordOutside := &order.Order{
		OrderNumber:       "ord-outside",
		BusinessID:        ws.Business.ID,
		CustomerID:        cus.ID,
		ShippingAddressID: addr.ID,
		Channel:           "instagram",
		Currency:          ws.Business.Currency,
		OrderedAt:         outside,
		Total:             decimal.RequireFromString("999"),
		COGS:              decimal.RequireFromString("111"),
		Status:            order.OrderStatusFulfilled,
		PaymentStatus:     order.OrderPaymentStatusPaid,
		PaymentMethod:     order.OrderPaymentMethodBankTransfer,
	}
	s.NoError(orderRepo.CreateOne(ctx, ordWithin))
	s.NoError(orderRepo.CreateOne(ctx, ordOutside))

	invWithin := map[string]interface{}{"investorId": ws.Admin.ID, "amount": "1000.00", "investedAt": within}
	expWithin := map[string]interface{}{"category": "software", "type": "one_time", "amount": "100.00", "occurredOn": within}
	wdWithin := map[string]interface{}{"withdrawerId": ws.Admin.ID, "amount": "50.00", "withdrawnAt": within}

	expOutside := map[string]interface{}{"category": "software", "type": "one_time", "amount": "999.00", "occurredOn": outside}
	wdOutside := map[string]interface{}{"withdrawerId": ws.Admin.ID, "amount": "999.00", "withdrawnAt": outside}

	resp1, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/investments", invWithin, ws.AdminToken)
	s.NoError(err)
	resp1.Body.Close()
	resp2, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/expenses", expWithin, ws.AdminToken)
	s.NoError(err)
	resp2.Body.Close()
	resp3, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/withdrawals", wdWithin, ws.AdminToken)
	s.NoError(err)
	resp3.Body.Close()
	resp4, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/expenses", expOutside, ws.AdminToken)
	s.NoError(err)
	resp4.Body.Close()
	resp5, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/withdrawals", wdOutside, ws.AdminToken)
	s.NoError(err)
	resp5.Body.Close()

	from := "2025-12-01"
	to := "2025-12-31"
	summaryResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/summary?from="+from+"&to="+to, nil, ws.AdminToken)
	s.NoError(err)
	defer summaryResp.Body.Close()
	s.Equal(http.StatusOK, summaryResp.StatusCode)

	var body map[string]interface{}
	s.NoError(testutils.DecodeJSON(summaryResp, &body))

	s.Equal("1000", body["totalInvestments"])
	s.Equal("50", body["totalWithdrawals"])
	s.Equal("100", body["totalExpenses"])
	// SafetyBuffer defaults to last-30-days expenses anchored to `to`
	// (1000 - 200) - 100 - 50 - 100 = 550
	s.Equal("550", body["safeToDrawAmount"])
	s.Equal(from, body["from"])
	s.Equal(to, body["to"])
}

func (s *SummarySuite) TestSummary_Permissions_MemberCanView() {
	ctx := context.Background()
	ws, err := s.helper.CreateWorkspaceWithAdminAndMemberAndBusiness(ctx)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+ws.Business.Descriptor+"/accounting/summary", nil, ws.MemberToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)
}

func TestSummarySuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(SummarySuite))
}
