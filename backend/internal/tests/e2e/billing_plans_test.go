package e2e_test

import (
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type BillingPlansSuite struct {
	suite.Suite
	helper *BillingTestHelper
}

func (s *BillingPlansSuite) SetupSuite() {
	s.helper = NewBillingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *BillingPlansSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database,
		"stripe_events",
		"subscriptions",
		"plans",
		"users",
		"workspaces",
	)
}

func (s *BillingPlansSuite) TearDownTest() {
	s.SetupTest()
}

func (s *BillingPlansSuite) TestPlans_ListAndGet() {
	ctx := s.T().Context()
	_, err := s.helper.CreatePlan(ctx, "starter", decimal.Zero, billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)

	resp, err := s.helper.Client().Get("/v1/billing/plans")
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var plans []map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &plans))
	s.NotEmpty(plans)

	resp2, err := s.helper.Client().Get("/v1/billing/plans/starter")
	s.NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusOK, resp2.StatusCode)

	var plan map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &plan))
	s.Contains(plan, "descriptor")
	s.Equal("starter", plan["descriptor"])
}

func (s *BillingPlansSuite) TestPlans_Get_NotFound() {
	resp, err := s.helper.Client().Get("/v1/billing/plans/does-not-exist")
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func TestBillingPlansSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BillingPlansSuite))
}
