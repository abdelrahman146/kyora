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

func (s *BillingPlansSuite) TestPlans_ListAndGet() {
	ctx := s.T().Context()
	descriptor := s.helper.UniqueSlug("starter")
	_, err := s.helper.CreatePlan(ctx, descriptor, decimal.Zero, billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)

	resp, err := s.helper.Client().Get("/v1/billing/plans")
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var plans []map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &plans))
	s.NotEmpty(plans)

	// Assert our plan is present and shaped correctly.
	found := false
	for _, p := range plans {
		if p["descriptor"] == descriptor {
			found = true
			s.Contains(p, "id")
			s.Contains(p, "descriptor")
			s.Contains(p, "name")
			s.Contains(p, "description")
			s.Contains(p, "stripePlanId")
			s.Contains(p, "price")
			s.Contains(p, "currency")
			s.Contains(p, "billingCycle")
			s.Contains(p, "features")
			s.Contains(p, "limits")
			break
		}
	}
	s.True(found, "expected plan to be present in list")

	resp2, err := s.helper.Client().Get("/v1/billing/plans/" + descriptor)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusOK, resp2.StatusCode)

	var plan map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &plan))
	s.Contains(plan, "id")
	s.Contains(plan, "descriptor")
	s.Equal(descriptor, plan["descriptor"])
	s.Contains(plan, "stripePlanId")
	s.Contains(plan, "price")
	s.Contains(plan, "currency")
	s.Contains(plan, "billingCycle")
	s.Contains(plan, "features")
	s.Contains(plan, "limits")
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
