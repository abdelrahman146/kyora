package e2e_test

import (
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type BillingUsageQuotaSuite struct {
	suite.Suite
	helper *BillingTestHelper
}

func (s *BillingUsageQuotaSuite) SetupSuite() {
	s.helper = NewBillingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *BillingUsageQuotaSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database,
		"order_items",
		"orders",
		"customer_addresses",
		"customers",
		"businesses",
		"plans",
		"subscriptions",
		"billing_invoice_records",
		"stripe_events",
		"users",
		"workspaces",
	)
	s.NoError(err)
}

func (s *BillingUsageQuotaSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database,
		"order_items",
		"orders",
		"customer_addresses",
		"customers",
		"businesses",
		"plans",
		"subscriptions",
		"billing_invoice_records",
		"stripe_events",
		"users",
		"workspaces",
	)
	s.NoError(err)
}

func (s *BillingUsageQuotaSuite) TestUsageAndQuota_HappyPath() {
	ctx := s.T().Context()
	limits := billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5}
	descriptor := s.helper.UniqueSlug("starter")
	_, err := s.helper.CreatePlan(ctx, descriptor, decimal.Zero, limits)
	s.NoError(err)

	_, ws, token, err := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("admin"), role.RoleAdmin)
	s.NoError(err)
	_, err = s.helper.AddWorkspaceUser(ctx, ws.ID, s.helper.UniqueEmail("member"), role.RoleUser)
	s.NoError(err)

	respSub, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": descriptor}, token)
	s.NoError(err)
	defer respSub.Body.Close()
	s.Equal(http.StatusOK, respSub.StatusCode)

	biz, err := s.helper.CreateBusiness(ctx, ws.ID)
	s.NoError(err)
	cust, addr, err := s.helper.CreateCustomerWithAddress(ctx, biz.ID)
	s.NoError(err)
	_, err = s.helper.CreateOrder(ctx, biz.ID, cust.ID, addr.ID)
	s.NoError(err)

	respUsage, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/usage", nil, token)
	s.NoError(err)
	defer respUsage.Body.Close()
	s.Equal(http.StatusOK, respUsage.StatusCode)

	var usage map[string]interface{}
	s.NoError(testutils.DecodeJSON(respUsage, &usage))
	s.Contains(usage, "ordersPerMonth")
	s.Contains(usage, "teamMembers")
	s.Contains(usage, "businesses")
	s.Equal(float64(1), usage["ordersPerMonth"])
	s.Equal(float64(2), usage["teamMembers"])
	s.Equal(float64(1), usage["businesses"])

	respQuota, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/usage/quota?type=team_members", nil, token)
	s.NoError(err)
	defer respQuota.Body.Close()
	s.Equal(http.StatusOK, respQuota.StatusCode)

	var quota map[string]interface{}
	s.NoError(testutils.DecodeJSON(respQuota, &quota))
	s.Contains(quota, "type")
	s.Contains(quota, "used")
	s.Contains(quota, "limit")
	s.Equal("team_members", quota["type"])
	s.Equal(float64(2), quota["used"])
	s.Equal(float64(limits.MaxTeamMembers), quota["limit"])
}

func (s *BillingUsageQuotaSuite) TestUsageQuota_ValidationErrors() {
	ctx := s.T().Context()
	descriptor := s.helper.UniqueSlug("starter")
	_, err := s.helper.CreatePlan(ctx, descriptor, decimal.Zero, billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)
	_, _, token, err := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("admin"), role.RoleAdmin)
	s.NoError(err)

	// needs subscription to avoid 500 on quota read
	respSub, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": descriptor}, token)
	s.NoError(err)
	defer respSub.Body.Close()
	s.Equal(http.StatusOK, respSub.StatusCode)

	respMissing, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/usage/quota", nil, token)
	s.NoError(err)
	defer respMissing.Body.Close()
	s.Equal(http.StatusBadRequest, respMissing.StatusCode)

	respInvalid, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/usage/quota?type=wat", nil, token)
	s.NoError(err)
	defer respInvalid.Body.Close()
	s.Equal(http.StatusBadRequest, respInvalid.StatusCode)
}

func TestBillingUsageQuotaSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BillingUsageQuotaSuite))
}
