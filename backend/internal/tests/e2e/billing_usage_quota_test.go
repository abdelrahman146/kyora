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
	testutils.TruncateTables(testEnv.Database,
		"stripe_events",
		"order_items",
		"orders",
		"customer_addresses",
		"customers",
		"businesses",
		"subscriptions",
		"plans",
		"users",
		"workspaces",
	)
}

func (s *BillingUsageQuotaSuite) TearDownTest() {
	s.SetupTest()
}

func (s *BillingUsageQuotaSuite) TestUsageAndQuota_HappyPath() {
	ctx := s.T().Context()
	limits := billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5}
	_, err := s.helper.CreatePlan(ctx, "starter", decimal.Zero, limits)
	s.NoError(err)

	_, ws, token := s.helper.CreateTestUser(ctx, "admin@example.com", role.RoleAdmin)
	_, err = s.helper.AddWorkspaceUser(ctx, ws.ID, "member@example.com", role.RoleUser)
	s.NoError(err)

	respSub, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": "starter"}, token)
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
	s.Equal(float64(1), usage["ordersPerMonth"])
	s.Equal(float64(2), usage["teamMembers"])
	s.Equal(float64(1), usage["businesses"])

	respQuota, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/usage/quota?type=team_members", nil, token)
	s.NoError(err)
	defer respQuota.Body.Close()
	s.Equal(http.StatusOK, respQuota.StatusCode)

	var quota map[string]interface{}
	s.NoError(testutils.DecodeJSON(respQuota, &quota))
	s.Equal("team_members", quota["type"])
	s.Equal(float64(2), quota["used"])
	s.Equal(float64(limits.MaxTeamMembers), quota["limit"])
}

func (s *BillingUsageQuotaSuite) TestUsageQuota_ValidationErrors() {
	ctx := s.T().Context()
	_, err := s.helper.CreatePlan(ctx, "starter", decimal.Zero, billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)
	_, _, token := s.helper.CreateTestUser(ctx, "admin@example.com", role.RoleAdmin)

	// needs subscription to avoid 500 on quota read
	respSub, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": "starter"}, token)
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
