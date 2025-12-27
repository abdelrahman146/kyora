package e2e_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
) // OnboardingStartSuite tests POST /v1/onboarding/start endpoint
type OnboardingStartSuite struct {
	suite.Suite
	client *testutils.HTTPClient
	helper *OnboardingTestHelper
}

func (s *OnboardingStartSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient(e2eBaseURL)
	s.helper = NewOnboardingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *OnboardingStartSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses")
}

func (s *OnboardingStartSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses")
}

func (s *OnboardingStartSuite) TestStart_ValidFreePlan() {
	s.helper.EnsureTestPlan("starter", "Starter Plan", 0.0)

	tests := []struct {
		name           string
		email          string
		planDescriptor string
		expectedStatus int
		checkResponse  func(body map[string]interface{})
	}{
		{
			name:           "new user with free plan",
			email:          "newuser@example.com",
			planDescriptor: "starter",
			expectedStatus: http.StatusOK,
			checkResponse: func(body map[string]interface{}) {
				s.NotEmpty(body["sessionToken"], "should return session token")
				s.Equal("plan_selected", body["stage"], "stage should be plan_selected")
				s.Equal(false, body["isPaid"], "should not be paid plan")
				// Verify response has exactly the expected fields
				s.Len(body, 3, "response should have exactly 3 fields")
				s.Contains(body, "sessionToken")
				s.Contains(body, "stage")
				s.Contains(body, "isPaid")
			},
		},
		{
			name:           "different email with free plan",
			email:          "another@example.com",
			planDescriptor: "starter",
			expectedStatus: http.StatusOK,
			checkResponse: func(body map[string]interface{}) {
				s.NotEmpty(body["sessionToken"])
				s.Equal("plan_selected", body["stage"])
				s.Equal(false, body["isPaid"])
				s.Len(body, 3, "response should have exactly 3 fields")
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload := map[string]interface{}{
				"email":          tt.email,
				"planDescriptor": tt.planDescriptor,
			}
			resp, err := s.client.Post("/v1/onboarding/start", payload)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(tt.expectedStatus, resp.StatusCode)

			var result map[string]interface{}
			s.NoError(testutils.DecodeJSON(resp, &result))
			if tt.checkResponse != nil {
				tt.checkResponse(result)
			}
		})
	}
}

func (s *OnboardingStartSuite) TestStart_ValidPaidPlan() {
	s.helper.EnsureTestPlan("professional", "Professional Plan", 99.0)

	payload := map[string]interface{}{
		"email":          "paiduser@example.com",
		"planDescriptor": "professional",
	}
	resp, err := s.client.Post("/v1/onboarding/start", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.NotEmpty(result["sessionToken"])
	s.Equal("plan_selected", result["stage"])
	s.Equal(true, result["isPaid"], "should be paid plan")
	// Verify all expected fields are present
	s.Len(result, 3, "response should have exactly 3 fields")
	s.Contains(result, "sessionToken")
	s.Contains(result, "stage")
	s.Contains(result, "isPaid")
}

func (s *OnboardingStartSuite) TestStart_ResumeExistingSession() {
	s.helper.EnsureTestPlan("starter", "Starter Plan", 0.0)

	// Create first session
	payload := map[string]interface{}{
		"email":          "resume@example.com",
		"planDescriptor": "starter",
	}
	resp1, err := s.client.Post("/v1/onboarding/start", payload)
	s.NoError(err)
	defer resp1.Body.Close()
	s.Equal(http.StatusOK, resp1.StatusCode)

	var result1 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp1, &result1))
	token1 := result1["sessionToken"].(string)

	// Create another session with same email - should resume
	resp2, err := s.client.Post("/v1/onboarding/start", payload)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusOK, resp2.StatusCode)

	var result2 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &result2))
	token2 := result2["sessionToken"].(string)

	// Tokens should be the same (resumed session)
	s.Equal(token1, token2, "should resume existing session")
}

func (s *OnboardingStartSuite) TestStart_EmailAlreadyRegistered() {
	// Create a registered user first
	err := s.helper.CreateTestUser("existing@example.com", "Password123!", "John", "Doe")
	s.NoError(err, "should create test user")
	s.helper.EnsureTestPlan("starter", "Starter Plan", 0.0)

	payload := map[string]interface{}{
		"email":          "existing@example.com",
		"planDescriptor": "starter",
	}
	resp, err := s.client.Post("/v1/onboarding/start", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusConflict, resp.StatusCode, "should reject already registered email")
}

func (s *OnboardingStartSuite) TestStart_InvalidEmail() {
	s.helper.EnsureTestPlan("starter", "Starter Plan", 0.0)

	tests := []struct {
		name  string
		email string
	}{
		{"invalid format", "notanemail"},
		{"missing @", "user.com"},
		{"missing domain", "user@"},
		{"empty", ""},
		{"spaces", "user @example.com"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload := map[string]interface{}{
				"email":          tt.email,
				"planDescriptor": "starter",
			}
			resp, err := s.client.Post("/v1/onboarding/start", payload)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func (s *OnboardingStartSuite) TestStart_InvalidPlan() {
	tests := []struct {
		name           string
		planDescriptor string
	}{
		{"nonexistent plan", "nonexistent"},
		{"empty plan", ""},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload := map[string]interface{}{
				"email":          "test@example.com",
				"planDescriptor": tt.planDescriptor,
			}
			resp, err := s.client.Post("/v1/onboarding/start", payload)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func (s *OnboardingStartSuite) TestStart_MissingFields() {
	tests := []struct {
		name    string
		payload map[string]interface{}
	}{
		{"missing email", map[string]interface{}{"planDescriptor": "starter"}},
		{"missing plan", map[string]interface{}{"email": "test@example.com"}},
		{"empty payload", map[string]interface{}{}},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.client.Post("/v1/onboarding/start", tt.payload)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func (s *OnboardingStartSuite) TestStart_SQLInjectionAttempts() {
	s.helper.EnsureTestPlan("starter", "Starter Plan", 0.0)

	sqlInjections := []string{
		"'; DROP TABLE users; --",
		"admin'--",
		"' OR '1'='1",
		"'; DELETE FROM users WHERE 1=1; --",
	}

	for i, injection := range sqlInjections {
		s.Run(fmt.Sprintf("injection %d", i), func() {
			payload := map[string]interface{}{
				"email":          injection,
				"planDescriptor": "starter",
			}
			resp, err := s.client.Post("/v1/onboarding/start", payload)
			s.NoError(err)
			defer resp.Body.Close()
			// Should reject as invalid email or handle safely
			s.True(resp.StatusCode >= 400, "should reject SQL injection")
		})
	}
}

func (s *OnboardingStartSuite) TestStart_XSSAttempts() {
	s.helper.EnsureTestPlan("starter", "Starter Plan", 0.0)

	xssPayloads := []string{
		"<script>alert('xss')</script>@example.com",
		"user<img src=x>@example.com",
	}

	for _, xss := range xssPayloads {
		payload := map[string]interface{}{
			"email":          xss,
			"planDescriptor": "starter",
		}
		resp, err := s.client.Post("/v1/onboarding/start", payload)
		s.NoError(err)
		defer resp.Body.Close()
		// Should reject as invalid email or sanitize
		s.True(resp.StatusCode >= 400, "should reject XSS attempt")
	}
}

func TestOnboardingStartSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(OnboardingStartSuite))
}
