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

type BillingSubscriptionAdvancedSuite struct {
	suite.Suite
	helper *BillingTestHelper
}

func (s *BillingSubscriptionAdvancedSuite) SetupSuite() {
	s.helper = NewBillingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *BillingSubscriptionAdvancedSuite) SetupTest() {
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

func (s *BillingSubscriptionAdvancedSuite) TearDownTest() {
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

func (s *BillingSubscriptionAdvancedSuite) TestEstimateProration_AndScheduleChange() {
	ctx := s.T().Context()

	starter := s.helper.UniqueSlug("starter")
	pro := s.helper.UniqueSlug("pro")
	_, err := s.helper.CreatePlan(ctx, starter, decimal.NewFromInt(10), billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)
	_, err = s.helper.CreatePlan(ctx, pro, decimal.NewFromInt(20), billing.PlanLimit{MaxOrdersPerMonth: 2000, MaxTeamMembers: 20, MaxBusinesses: 10})
	s.NoError(err)

	_, _, adminToken := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("admin"), role.RoleAdmin)

	respSub, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": starter}, adminToken)
	s.NoError(err)
	defer respSub.Body.Close()
	s.Require().Equal(http.StatusOK, respSub.StatusCode)

	respEstimate, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription/estimate-proration", map[string]interface{}{"newPlanDescriptor": pro}, adminToken)
	s.NoError(err)
	defer respEstimate.Body.Close()
	s.Equal(http.StatusOK, respEstimate.StatusCode)
	var est map[string]interface{}
	s.NoError(testutils.DecodeJSON(respEstimate, &est))
	s.Len(est, 1)
	s.Contains(est, "amount")

	// Validation: missing newPlanDescriptor
	respBad, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription/estimate-proration", map[string]interface{}{}, adminToken)
	s.NoError(err)
	defer respBad.Body.Close()
	s.Equal(http.StatusBadRequest, respBad.StatusCode)

	// Schedule change: invalid effective date should be 400 (API contract)
	respBadDate, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription/schedule-change", map[string]interface{}{
		"planDescriptor": pro,
		"effectiveDate":  "not-a-date",
	}, adminToken)
	s.NoError(err)
	defer respBadDate.Body.Close()
	s.Equal(http.StatusBadRequest, respBadDate.StatusCode)

	respSched, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription/schedule-change", map[string]interface{}{
		"planDescriptor": pro,
		"effectiveDate":  "2030-01-01",
		"prorationMode":  "none",
	}, adminToken)
	s.NoError(err)
	defer respSched.Body.Close()
	s.Equal(http.StatusOK, respSched.StatusCode)
	var schedule map[string]interface{}
	s.NoError(testutils.DecodeJSON(respSched, &schedule))
	s.Contains(schedule, "id")
}

func (s *BillingSubscriptionAdvancedSuite) TestTrialEndpoints_AreConsistent() {
	ctx := s.T().Context()
	descriptor := s.helper.UniqueSlug("starter")
	_, err := s.helper.CreatePlan(ctx, descriptor, decimal.Zero, billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)

	_, _, adminToken := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("admin"), role.RoleAdmin)
	respSub, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": descriptor}, adminToken)
	s.NoError(err)
	defer respSub.Body.Close()
	s.Require().Equal(http.StatusOK, respSub.StatusCode)

	respTrial, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/subscription/trial", nil, adminToken)
	s.NoError(err)
	defer respTrial.Body.Close()
	s.Equal(http.StatusOK, respTrial.StatusCode)

	var trial map[string]interface{}
	s.NoError(testutils.DecodeJSON(respTrial, &trial))
	s.Contains(trial, "isInTrial")
	s.Contains(trial, "trialEnd")
	s.Contains(trial, "daysRemaining")

	// Extend trial should be a 400 when not trialing (API contract)
	respExtend, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription/trial/extend", map[string]interface{}{"additionalDays": 5}, adminToken)
	s.NoError(err)
	defer respExtend.Body.Close()
	s.Equal(http.StatusBadRequest, respExtend.StatusCode)
}

func (s *BillingSubscriptionAdvancedSuite) TestGracePeriod_Validation() {
	ctx := s.T().Context()
	descriptor := s.helper.UniqueSlug("starter")
	_, err := s.helper.CreatePlan(ctx, descriptor, decimal.Zero, billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)

	_, _, adminToken := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("admin"), role.RoleAdmin)
	respSub, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": descriptor}, adminToken)
	s.NoError(err)
	defer respSub.Body.Close()
	s.Require().Equal(http.StatusOK, respSub.StatusCode)

	// Validation: graceDays required/min/max enforced
	respBad, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription/grace-period", map[string]interface{}{}, adminToken)
	s.NoError(err)
	defer respBad.Body.Close()
	s.Equal(http.StatusBadRequest, respBad.StatusCode)

	respBad2, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription/grace-period", map[string]interface{}{"graceDays": 0}, adminToken)
	s.NoError(err)
	defer respBad2.Body.Close()
	s.Equal(http.StatusBadRequest, respBad2.StatusCode)

	// Not past-due -> should be 400 (API contract)
	respNotPastDue, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription/grace-period", map[string]interface{}{"graceDays": 5}, adminToken)
	s.NoError(err)
	defer respNotPastDue.Body.Close()
	s.Equal(http.StatusBadRequest, respNotPastDue.StatusCode)
}

func (s *BillingSubscriptionAdvancedSuite) TestResumeSubscription_AfterCancel() {
	ctx := s.T().Context()
	descriptor := s.helper.UniqueSlug("starter")
	_, err := s.helper.CreatePlan(ctx, descriptor, decimal.Zero, billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)

	_, _, adminToken := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("admin"), role.RoleAdmin)

	respSub, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": descriptor}, adminToken)
	s.NoError(err)
	defer respSub.Body.Close()
	s.Require().Equal(http.StatusOK, respSub.StatusCode)

	respCancel, err := s.helper.Client().AuthenticatedRequest("DELETE", "/v1/billing/subscription", nil, adminToken)
	s.NoError(err)
	defer respCancel.Body.Close()
	s.Require().Equal(http.StatusNoContent, respCancel.StatusCode)

	respResume, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription/resume", map[string]interface{}{}, adminToken)
	s.NoError(err)
	defer respResume.Body.Close()
	s.Equal(http.StatusOK, respResume.StatusCode)
	var sub map[string]interface{}
	s.NoError(testutils.DecodeJSON(respResume, &sub))
	s.Contains(sub, "id")
	s.Contains(sub, "status")
}

func TestBillingSubscriptionAdvancedSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BillingSubscriptionAdvancedSuite))
}
