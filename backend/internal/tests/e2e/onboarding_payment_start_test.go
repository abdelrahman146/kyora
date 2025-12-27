package e2e_test

import (
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// OnboardingPaymentStartSuite tests POST /v1/onboarding/payment/start endpoint.
type OnboardingPaymentStartSuite struct {
	suite.Suite
	client *testutils.HTTPClient
	helper *OnboardingTestHelper
}

func (s *OnboardingPaymentStartSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient(e2eBaseURL)
	s.helper = NewOnboardingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *OnboardingPaymentStartSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses")
}

func (s *OnboardingPaymentStartSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses")
}

func (s *OnboardingPaymentStartSuite) TestPaymentStart_ValidationErrors() {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name:           "missing payload",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing sessionToken",
			payload: map[string]interface{}{
				"successUrl": "https://example.com/s",
				"cancelUrl":  "https://example.com/c",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing successUrl",
			payload: map[string]interface{}{
				"sessionToken": "tok",
				"cancelUrl":    "https://example.com/c",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing cancelUrl",
			payload: map[string]interface{}{
				"sessionToken": "tok",
				"successUrl":   "https://example.com/s",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid url",
			payload: map[string]interface{}{
				"sessionToken": "tok",
				"successUrl":   "not-a-url",
				"cancelUrl":    "https://example.com/c",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.client.Post("/v1/onboarding/payment/start", tt.payload)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(tt.expectedStatus, resp.StatusCode)

			var result map[string]interface{}
			s.NoError(testutils.DecodeJSON(resp, &result))
			s.Contains(result, "type")
			s.Contains(result, "title")
			s.Contains(result, "status")
			s.Equal(float64(tt.expectedStatus), result["status"])
		})
	}
}

func (s *OnboardingPaymentStartSuite) TestPaymentStart_SessionNotFound() {
	resp, err := s.client.Post("/v1/onboarding/payment/start", map[string]interface{}{
		"sessionToken": "sess_does_not_exist",
		"successUrl":   "https://example.com/success",
		"cancelUrl":    "https://example.com/cancel",
	})
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *OnboardingPaymentStartSuite) TestPaymentStart_FreePlan_ReturnsEmptyCheckoutURL() {
	token, err := s.helper.CreateOnboardingSession("freepay@example.com", "starter")
	s.NoError(err)
	s.NotEmpty(token)

	resp, err := s.client.Post("/v1/onboarding/payment/start", map[string]interface{}{
		"sessionToken": token,
		"successUrl":   "https://example.com/success",
		"cancelUrl":    "https://example.com/cancel",
	})
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Len(result, 1, "response should have exactly 1 field")
	s.Contains(result, "checkoutUrl")
	checkout, ok := result["checkoutUrl"].(string)
	s.True(ok)
	s.Equal("", checkout)
}

func TestOnboardingPaymentStartSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(OnboardingPaymentStartSuite))
}

/* package e2e_test

import (
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// OnboardingPaymentStartSuite tests POST /v1/onboarding/payment/start endpoint.
type OnboardingPaymentStartSuite struct {
	suite.Suite
	client *testutils.HTTPClient
	helper *OnboardingTestHelper
}

func (s *OnboardingPaymentStartSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient(e2eBaseURL)
	s.helper = NewOnboardingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *OnboardingPaymentStartSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses")
}

func (s *OnboardingPaymentStartSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses")
}

func (s *OnboardingPaymentStartSuite) TestPaymentStart_ValidationErrors() {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name:           "missing payload",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing sessionToken",
			payload: map[string]interface{}{
				"successUrl": "https://example.com/s",
				"cancelUrl":  "https://example.com/c",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing successUrl",
			payload: map[string]interface{}{
				"sessionToken": "tok",
				"cancelUrl":    "https://example.com/c",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing cancelUrl",
			payload: map[string]interface{}{
				"sessionToken": "tok",
				"successUrl":   "https://example.com/s",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid url",
			payload: map[string]interface{}{
				"sessionToken": "tok",
				"successUrl":   "not-a-url",
				"cancelUrl":    "https://example.com/c",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.client.Post("/v1/onboarding/payment/start", tt.payload)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(tt.expectedStatus, resp.StatusCode)

			var result map[string]interface{}
			s.NoError(testutils.DecodeJSON(resp, &result))
			s.Contains(result, "type")
			s.Contains(result, "title")
			s.Contains(result, "status")
			s.Equal(float64(tt.expectedStatus), result["status"])
		})
	}
}

func (s *OnboardingPaymentStartSuite) TestPaymentStart_SessionNotFound() {
	resp, err := s.client.Post("/v1/onboarding/payment/start", map[string]interface{}{
		"sessionToken": "sess_does_not_exist",
		"successUrl":   "https://example.com/success",
		"cancelUrl":    "https://example.com/cancel",
	})
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *OnboardingPaymentStartSuite) TestPaymentStart_FreePlan_ReturnsEmptyCheckoutURL() {
	token, err := s.helper.CreateOnboardingSession("freepay@example.com", "starter")
	s.NoError(err)
	s.NotEmpty(token)

	resp, err := s.client.Post("/v1/onboarding/payment/start", map[string]interface{}{
		"sessionToken": token,
		"successUrl":   "https://example.com/success",
		"cancelUrl":    "https://example.com/cancel",
	})
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Len(result, 1, "response should have exactly 1 field")
	s.Contains(result, "checkoutUrl")
	checkout, ok := result["checkoutUrl"].(string)
	s.True(ok)
	s.Equal("", checkout)
}

func TestOnboardingPaymentStartSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(OnboardingPaymentStartSuite))
}

*/
