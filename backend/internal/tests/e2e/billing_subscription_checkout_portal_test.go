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

type BillingSubscriptionCheckoutPortalSuite struct {
	suite.Suite
	helper *BillingTestHelper
}

func (s *BillingSubscriptionCheckoutPortalSuite) SetupSuite() {
	s.helper = NewBillingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *BillingSubscriptionCheckoutPortalSuite) SetupTest() {
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

func (s *BillingSubscriptionCheckoutPortalSuite) TearDownTest() {
	s.SetupTest()
}

func (s *BillingSubscriptionCheckoutPortalSuite) TestCheckoutSession_AuthzAndFreePlan() {
	ctx := s.T().Context()
	_, err := s.helper.CreatePlan(ctx, "free", decimal.Zero, billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)

	_, _, adminToken := s.helper.CreateTestUser(ctx, "admin@example.com", role.RoleAdmin)
	_, _, memberToken := s.helper.CreateTestUser(ctx, "member@example.com", role.RoleUser)

	payload := map[string]interface{}{
		"planDescriptor": "free",
		"successUrl":     "https://example.com/success",
		"cancelUrl":      "https://example.com/cancel",
	}

	respUnauthed, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/checkout/session", payload, "")
	s.NoError(err)
	defer respUnauthed.Body.Close()
	s.Require().Equal(http.StatusUnauthorized, respUnauthed.StatusCode)

	respForbidden, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/checkout/session", payload, memberToken)
	s.NoError(err)
	defer respForbidden.Body.Close()
	s.Require().Equal(http.StatusForbidden, respForbidden.StatusCode)

	resp, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/checkout/session", payload, adminToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Require().Len(result, 2, "response should have exactly 2 fields")
	s.Require().Contains(result, "url")
	s.Require().Contains(result, "checkoutUrl")
	url, ok := result["url"].(string)
	s.Require().True(ok)
	checkoutURL, ok := result["checkoutUrl"].(string)
	s.Require().True(ok)
	s.Equal("", url)
	s.Equal("", checkoutURL)

	// Free plan should create a subscription without redirect
	respSub, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/subscription", nil, adminToken)
	s.NoError(err)
	defer respSub.Body.Close()
	s.Require().Equal(http.StatusOK, respSub.StatusCode)
}

func (s *BillingSubscriptionCheckoutPortalSuite) TestCheckoutSession_PaidPlan_ReturnsURL() {
	ctx := s.T().Context()
	_, err := s.helper.CreatePlan(ctx, "paid", decimal.NewFromInt(10), billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)

	_, _, adminToken := s.helper.CreateTestUser(ctx, "admin@example.com", role.RoleAdmin)

	payload := map[string]interface{}{
		"planDescriptor": "paid",
		"successUrl":     "https://example.com/success",
		"cancelUrl":      "https://example.com/cancel",
	}

	resp, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/checkout/session", payload, adminToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Require().Len(result, 2, "response should have exactly 2 fields")
	s.Require().Contains(result, "url")
	s.Require().Contains(result, "checkoutUrl")
	url, ok := result["url"].(string)
	s.Require().True(ok)
	checkoutURL, ok := result["checkoutUrl"].(string)
	s.Require().True(ok)
	s.NotEmpty(url)
	s.Equal(url, checkoutURL)
}

func (s *BillingSubscriptionCheckoutPortalSuite) TestBillingPortalSession_ReturnsURL() {
	ctx := s.T().Context()
	_, _, adminToken := s.helper.CreateTestUser(ctx, "admin@example.com", role.RoleAdmin)

	payload := map[string]interface{}{
		"returnUrl": "https://example.com/return",
	}

	resp, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/portal/session", payload, adminToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Require().Len(result, 2, "response should have exactly 2 fields")
	s.Require().Contains(result, "url")
	s.Require().Contains(result, "portalUrl")
	url, ok := result["url"].(string)
	s.Require().True(ok)
	portalURL, ok := result["portalUrl"].(string)
	s.Require().True(ok)
	s.NotEmpty(url)
	s.Equal(url, portalURL)
}

func TestBillingSubscriptionCheckoutPortalSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BillingSubscriptionCheckoutPortalSuite))
}
