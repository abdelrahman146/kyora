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

type BillingPaymentMethodsSuite struct {
	suite.Suite
	helper *BillingTestHelper
}

func (s *BillingPaymentMethodsSuite) SetupSuite() {
	s.helper = NewBillingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *BillingPaymentMethodsSuite) TestSetupIntent_Authz() {
	ctx := s.T().Context()
	_, _, adminToken := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("admin"), role.RoleAdmin)
	_, _, memberToken := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("member"), role.RoleUser)

	respUnauthed, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/payment-methods/setup-intent", nil, "")
	s.NoError(err)
	defer respUnauthed.Body.Close()
	s.Equal(http.StatusUnauthorized, respUnauthed.StatusCode)

	respForbidden, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/payment-methods/setup-intent", nil, memberToken)
	s.NoError(err)
	defer respForbidden.Body.Close()
	s.Equal(http.StatusForbidden, respForbidden.StatusCode)

	respOK, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/payment-methods/setup-intent", nil, adminToken)
	s.NoError(err)
	defer respOK.Body.Close()
	s.Equal(http.StatusOK, respOK.StatusCode)

	var body map[string]interface{}
	s.NoError(testutils.DecodeJSON(respOK, &body))
	s.Len(body, 1, "response should have exactly 1 field")
	s.Contains(body, "clientSecret")
	secret, ok := body["clientSecret"].(string)
	s.True(ok)
	s.NotEmpty(secret)
}

func (s *BillingPaymentMethodsSuite) TestAttachPaymentMethod_Validation_AndDetailsReflectDefault() {
	ctx := s.T().Context()
	descriptor := s.helper.UniqueSlug("paid")
	_, err := s.helper.CreatePlan(ctx, descriptor, decimal.NewFromInt(10), billing.PlanLimit{MaxOrdersPerMonth: 1000, MaxTeamMembers: 10, MaxBusinesses: 5})
	s.NoError(err)

	_, _, adminToken := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("admin"), role.RoleAdmin)

	// Create subscription so we have a customer.
	respSub, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/subscription", map[string]interface{}{"planDescriptor": descriptor}, adminToken)
	s.NoError(err)
	defer respSub.Body.Close()
	s.Require().Equal(http.StatusOK, respSub.StatusCode)

	// Validation: missing paymentMethodId
	respBad, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/payment-methods/attach", map[string]interface{}{}, adminToken)
	s.NoError(err)
	defer respBad.Body.Close()
	s.Equal(http.StatusBadRequest, respBad.StatusCode)

	// Invalid PM id should be 400
	respInvalid, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/payment-methods/attach", map[string]interface{}{"paymentMethodId": "pm_does_not_exist"}, adminToken)
	s.NoError(err)
	defer respInvalid.Body.Close()
	s.Equal(http.StatusBadRequest, respInvalid.StatusCode)

	pmID, err := s.helper.CreateStripeCardPaymentMethod()
	s.NoError(err)

	respAttach, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/payment-methods/attach", map[string]interface{}{"paymentMethodId": pmID}, adminToken)
	s.NoError(err)
	defer respAttach.Body.Close()
	s.Equal(http.StatusOK, respAttach.StatusCode)

	respDetails, err := s.helper.Client().AuthenticatedRequest("GET", "/v1/billing/subscription/details", nil, adminToken)
	s.NoError(err)
	defer respDetails.Body.Close()
	s.Equal(http.StatusOK, respDetails.StatusCode)

	var details map[string]interface{}
	s.NoError(testutils.DecodeJSON(respDetails, &details))
	s.Contains(details, "subscription")
	s.Contains(details, "plan")
	s.Contains(details, "paymentMethod")

	pm, ok := details["paymentMethod"].(map[string]interface{})
	s.True(ok)
	s.Contains(pm, "id")
	s.Contains(pm, "brand")
	s.Contains(pm, "last4")
	s.Contains(pm, "expMonth")
	s.Contains(pm, "expYear")

	// Ensure the attached PM is reflected.
	if gotID, ok := pm["id"].(string); ok {
		s.Equal(pmID, gotID)
	}
}

func TestBillingPaymentMethodsSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BillingPaymentMethodsSuite))
}
