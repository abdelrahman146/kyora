package e2e_test

import (
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

type HTTPBodyLimitSuite struct {
	suite.Suite
	client *testutils.HTTPClient
}

func (s *HTTPBodyLimitSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient("http://localhost:18080")
}

func (s *HTTPBodyLimitSuite) TestRejectsOverLimitBody() {
	if testServer == nil {
		s.T().Skip("Test server not initialized")
	}

	// Default configured max is 1 MiB (see config.HTTPMaxBodyBytes default).
	overLimitBody := make([]byte, 1024*1024+1)

	resp, err := s.client.PostRaw("/healthz", overLimitBody, map[string]string{
		"Content-Type": "application/json",
	})
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusRequestEntityTooLarge, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Equal(float64(http.StatusRequestEntityTooLarge), result["status"])
	s.Equal("Payload Too Large", result["title"])
}

func TestHTTPBodyLimitSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(HTTPBodyLimitSuite))
}
