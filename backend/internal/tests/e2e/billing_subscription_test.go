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
	err := testutils.TruncateTables(testEnv.Database,
		"plans",
		"subscriptions",
		"billing_invoice_records",
		"stripe_events",
		"users",
		"workspaces",
	)
	s.NoError(err)
}

func (s *BillingSubscriptionSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database,
		"plans",
		"subscriptions",
		"billing_invoice_records",
		"stripe_events",
		"users",
		"workspaces",
	)
	s.NoError(err)
}

func (s *BillingSubscriptionSuite) TestSubscription_CreateGetCancel() {
	ctx := s.T().Context()
	descriptor := s.helper.UniqueSlug("starter")
	_, err := s.helper.CreatePlan(ctx, descriptor, decimal.Zero, billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)

	_, ws, adminToken, err := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("admin"), role.RoleAdmin)
	s.NoError(err)

	respCreate, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": descriptor}, adminToken)
	s.NoError(err)
	defer respCreate.Body.Close()
	s.Require().Equal(http.StatusOK, respCreate.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(respCreate, &created))
	s.Require().Contains(created, "id")
	s.Require().Contains(created, "workspaceId")
	s.Require().Contains(created, "planId")
	s.Require().Contains(created, "stripeSubId")
	s.Require().Contains(created, "status")
	if gotWs, ok := created["workspaceId"].(string); ok {
		s.Equal(ws.ID, gotWs)
	}

	respGet, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/subscription", nil, adminToken)
	s.NoError(err)
	defer respGet.Body.Close()
	s.Require().Equal(http.StatusOK, respGet.StatusCode)

	var got map[string]interface{}
	s.NoError(testutils.DecodeJSON(respGet, &got))
	s.Require().Contains(got, "id")
	s.Require().Contains(got, "workspaceId")
	s.Require().Contains(got, "planId")
	s.Require().Contains(got, "stripeSubId")
	s.Require().Contains(got, "status")

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
	descriptor := s.helper.UniqueSlug("starter")
	_, err := s.helper.CreatePlan(ctx, descriptor, decimal.Zero, billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)

	_, _, memberToken, err := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("member"), role.RoleUser)
	s.NoError(err)

	resp, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": descriptor}, memberToken)
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
