package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

// OnboardingSuite is a testify suite for end-to-end onboarding tests
type OnboardingSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
}

// SetupSuite runs once before all tests in the suite
func (s *OnboardingSuite) SetupSuite() {
	s.baseURL = "http://localhost:18080"
	s.httpClient = &http.Client{}
}

// TearDownSuite runs once after all tests in the suite
func (s *OnboardingSuite) TearDownSuite() {
	// Cleanup if needed per suite
}

// SetupTest runs before each test
func (s *OnboardingSuite) SetupTest() {
	// Per-test setup if needed
}

// TearDownTest runs after each test
func (s *OnboardingSuite) TearDownTest() {
	// Per-test cleanup if needed
}

// TestHealthCheck verifies the server is responding
func (s *OnboardingSuite) TestHealthCheck() {
	resp, err := s.httpClient.Get(s.baseURL + "/healthz")
	s.NoError(err, "health check request should succeed")
	s.Equal(http.StatusOK, resp.StatusCode, "health check should return 200")
	defer resp.Body.Close()
}

// TestLivenessCheck verifies the liveness endpoint
func (s *OnboardingSuite) TestLivenessCheck() {
	resp, err := s.httpClient.Get(s.baseURL + "/livez")
	s.NoError(err, "liveness check request should succeed")
	s.Equal(http.StatusOK, resp.StatusCode, "liveness check should return 200")
	defer resp.Body.Close()
}

// Helper method to make JSON POST requests
func (s *OnboardingSuite) postJSON(path string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.baseURL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return s.httpClient.Do(req)
}

// Helper method to make authenticated requests
func (s *OnboardingSuite) authenticatedRequest(method, path string, body []byte, token string) (*http.Response, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}
	req, err := http.NewRequest(method, s.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.AddCookie(&http.Cookie{Name: "jwt", Value: token})
	}
	return s.httpClient.Do(req)
}

// Example test - you can expand this based on your onboarding API
func (s *OnboardingSuite) TestOnboardingFlow() {
	s.T().Skip("Implement onboarding flow tests based on your API endpoints")

	// Example structure:
	// 1. Register new user
	// 2. Verify email (mock)
	// 3. Create workspace
	// 4. Verify workspace created

	// payload := map[string]interface{}{
	// 	"email":    "test@example.com",
	// 	"password": "TestPass123!",
	// 	"name":     "Test User",
	// }
	// resp, err := s.postJSON("/api/auth/register", payload)
	// s.NoError(err)
	// s.Equal(http.StatusCreated, resp.StatusCode)
	// defer resp.Body.Close()
}

// TestOnboardingSuite runs the OnboardingSuite
func TestOnboardingSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized - skipping e2e tests")
	}
	suite.Run(t, new(OnboardingSuite))
}
