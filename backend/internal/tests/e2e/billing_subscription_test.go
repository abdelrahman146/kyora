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

type BillingSubscriptionSuite struct {
	suite.Suite
	helper *BillingTestHelper
}

func (s *BillingSubscriptionSuite) SetupSuite() {
	s.helper = NewBillingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *BillingSubscriptionSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database,
		"stripe_events",
		"subscriptions",
		"plans",
		"users",
		"workspaces",
	)
}

func (s *BillingSubscriptionSuite) TearDownTest() {
	s.SetupTest()
}

func (s *BillingSubscriptionSuite) TestSubscription_CreateGetCancel() {
	ctx := s.T().Context()
	_, err := s.helper.CreatePlan(ctx, "starter", decimal.Zero, billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)

	_, _, adminToken := s.helper.CreateTestUser(ctx, "admin@example.com", role.RoleAdmin)

	respCreate, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": "starter"}, adminToken)
	s.NoError(err)
	defer respCreate.Body.Close()
	s.Require().Equal(http.StatusOK, respCreate.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(respCreate, &created))
	s.Require().Contains(created, "workspaceId")
	s.Require().Contains(created, "planId")
	s.Require().Contains(created, "stripeSubId")
	s.Require().Contains(created, "status")

	respGet, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/subscription", nil, adminToken)
	s.NoError(err)
	defer respGet.Body.Close()
	s.Require().Equal(http.StatusOK, respGet.StatusCode)

	respCancel, err := s.helper.Client().AuthenticatedRequest("DELETE", "/v1/billing/subscription", nil, adminToken)
	s.NoError(err)
	defer respCancel.Body.Close()
	s.Require().Equal(http.StatusNoContent, respCancel.StatusCode)

	respGet2, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/subscription", nil, adminToken)
	s.NoError(err)
	defer respGet2.Body.Close()
	s.Require().Equal(http.StatusOK, respGet2.StatusCode)

	var after map[string]interface{}
	s.NoError(testutils.DecodeJSON(respGet2, &after))
	if status, ok := after["status"].(string); ok {
		s.Equal("canceled", status)
	}
}

func (s *BillingSubscriptionSuite) TestSubscription_Create_ForbiddenForMember() {
	ctx := s.T().Context()
	_, err := s.helper.CreatePlan(ctx, "starter", decimal.Zero, billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)

	_, _, memberToken := s.helper.CreateTestUser(ctx, "member@example.com", role.RoleUser)

	resp, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": "starter"}, memberToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusForbidden, resp.StatusCode)
}

func TestBillingSubscriptionSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BillingSubscriptionSuite))
}
