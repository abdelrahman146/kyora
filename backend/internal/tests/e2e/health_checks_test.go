package e2e_test

import (
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// HealthCheckSuite tests the GET /healthz endpoint
type HealthCheckSuite struct {
	suite.Suite
	client *testutils.HTTPClient
}

func (s *HealthCheckSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient("http://localhost:18080")
}

func (s *HealthCheckSuite) SetupTest() {
	// No database cleanup needed for health checks
}

func (s *HealthCheckSuite) TearDownTest() {
	// No cleanup needed
}

func (s *HealthCheckSuite) TestHealthCheck_Success() {
	resp, err := s.client.Get("/healthz")
	s.NoError(err, "health check request should succeed")
	s.Equal(http.StatusOK, resp.StatusCode, "health check should return 200")
	defer resp.Body.Close()

	body, err := testutils.ReadBody(resp)
	s.NoError(err, "should read response body")
	s.Equal("ok", body, "health check should return 'ok'")
}

func TestHealthCheckSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(HealthCheckSuite))
}

// LivenessCheckSuite tests the GET /livez endpoint
type LivenessCheckSuite struct {
	suite.Suite
	client *testutils.HTTPClient
}

func (s *LivenessCheckSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient("http://localhost:18080")
}

func (s *LivenessCheckSuite) SetupTest() {
	// No database cleanup needed for liveness checks
}

func (s *LivenessCheckSuite) TearDownTest() {
	// No cleanup needed
}

func (s *LivenessCheckSuite) TestLivenessCheck_Success() {
	resp, err := s.client.Get("/livez")
	s.NoError(err, "liveness check request should succeed")
	s.Equal(http.StatusOK, resp.StatusCode, "liveness check should return 200")
	defer resp.Body.Close()

	body, err := testutils.ReadBody(resp)
	s.NoError(err, "should read response body")
	s.Equal("ok", body, "liveness check should return 'ok'")
}

func TestLivenessCheckSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(LivenessCheckSuite))
}
