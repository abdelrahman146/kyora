package e2e_test

import (
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// OnboardingStartSuite tests the POST /api/onboarding/start endpoint
type OnboardingStartSuite struct {
	suite.Suite
	client *testutils.HTTPClient
}

func (s *OnboardingStartSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient("http://localhost:18080")
}

func (s *OnboardingStartSuite) SetupTest() {
	// Clear relevant tables before each test
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions")
}

func (s *OnboardingStartSuite) TearDownTest() {
	// Clean up after each test
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions")
}

func (s *OnboardingStartSuite) TestOnboardingStart_ValidEmail() {
	// Table-driven test for multiple scenarios
	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "valid email",
			payload: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid email format",
			payload: map[string]interface{}{
				"email": "invalid-email",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing email",
			payload: map[string]interface{}{
				"name": "Test User",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty payload",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.client.Post("/api/onboarding/start", tt.payload)
			s.NoError(err, "request should not error")
			s.Equal(tt.expectedStatus, resp.StatusCode, "status code should match")
			defer resp.Body.Close()

			if tt.expectedStatus == http.StatusOK {
				var result map[string]interface{}
				err = testutils.DecodeJSON(resp, &result)
				s.NoError(err, "should decode response")
				s.NotEmpty(result, "response should not be empty")
			}
		})
	}
}

func TestOnboardingStartSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(OnboardingStartSuite))
}
