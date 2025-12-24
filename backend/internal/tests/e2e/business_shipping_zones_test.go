package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// BusinessShippingZonesSuite tests /v1/businesses/:businessDescriptor/shipping-zones endpoints.
type BusinessShippingZonesSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *BusinessShippingZonesSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *BusinessShippingZonesSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "shipping_zones", "subscriptions", "plans")
	s.NoError(err)
}

func (s *BusinessShippingZonesSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "shipping_zones", "subscriptions", "plans")
	s.NoError(err)
}

func (s *BusinessShippingZonesSuite) createBusiness(ctx context.Context, token string, descriptor string) {
	payload := map[string]interface{}{
		"name":        "Test Business " + descriptor,
		"descriptor":  descriptor,
		"countryCode": "eg",
		"currency":    "usd",
	}
	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)
}

func (s *BusinessShippingZonesSuite) listZones(token string, businessDescriptor string) ([]map[string]interface{}, int) {
	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/"+businessDescriptor+"/shipping-zones", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	status := resp.StatusCode

	if status != http.StatusOK {
		return nil, status
	}
	var zones []map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &zones))
	return zones, status
}

func (s *BusinessShippingZonesSuite) findZoneByName(zones []map[string]interface{}, name string) (map[string]interface{}, bool) {
	for _, z := range zones {
		if z["name"] == name {
			return z, true
		}
	}
	return nil, false
}

func (s *BusinessShippingZonesSuite) TestListShippingZones_IncludesDefaultZone() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))

	s.createBusiness(ctx, token, "test-business")

	zones, status := s.listZones(token, "test-business")
	s.Equal(http.StatusOK, status)
	s.GreaterOrEqual(len(zones), 1)

	z, ok := s.findZoneByName(zones, "EG")
	s.True(ok, "should include default zone named business country")
	s.Equal("USD", z["currency"])
	s.Equal("0", z["shippingCost"])
	s.Equal("0", z["freeShippingThreshold"])

	countriesAny, exists := z["countries"]
	s.True(exists)
	countries, ok := countriesAny.([]interface{})
	s.True(ok)
	s.Len(countries, 1)
	s.Equal("EG", countries[0])
}

func (s *BusinessShippingZonesSuite) TestGetShippingZone_BOPLA_PreventCrossBusinessAccess() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))

	s.createBusiness(ctx, token, "biz-1")
	s.createBusiness(ctx, token, "biz-2")

	createZonePayload := map[string]interface{}{
		"name":                  "Alexandria",
		"countries":             []string{"EG"},
		"shippingCost":          "25",
		"freeShippingThreshold": "500",
	}
	createResp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/biz-1/shipping-zones", createZonePayload, token)
	s.NoError(err)
	defer createResp.Body.Close()
	s.Equal(http.StatusCreated, createResp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(createResp, &created))
	zoneID, _ := created["id"].(string)
	s.NotEmpty(zoneID)

	getResp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/businesses/biz-2/shipping-zones/"+zoneID, nil, token)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusNotFound, getResp.StatusCode)
}

func (s *BusinessShippingZonesSuite) TestListShippingZones_CrossWorkspaceIsolation() {
	ctx := context.Background()
	_, ws1, token1, err := s.helper.CreateTestUser(ctx, "admin1@example.com", "ValidPassword123!", "Admin", "One", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws1.ID))

	_, _, token2, err := s.helper.CreateTestUser(ctx, "admin2@example.com", "ValidPassword123!", "Admin", "Two", role.RoleAdmin)
	s.NoError(err)

	s.createBusiness(ctx, token1, "biz")

	_, status := s.listZones(token2, "biz")
	s.Equal(http.StatusNotFound, status)
}

func (s *BusinessShippingZonesSuite) TestCreateShippingZone_PersistsAndLists() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))

	s.createBusiness(ctx, token, "biz")

	createZonePayload := map[string]interface{}{
		"name":                  "GCC",
		"countries":             []string{"SA", "AE"},
		"shippingCost":          "10",
		"freeShippingThreshold": "100",
	}
	createResp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses/biz/shipping-zones", createZonePayload, token)
	s.NoError(err)
	defer createResp.Body.Close()
	s.Equal(http.StatusCreated, createResp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(createResp, &created))
	s.Equal("GCC", created["name"])
	s.Equal("USD", created["currency"], "currency should be derived from business")

	// Verify it persisted (default zone + created zone).
	bizRepo := database.NewRepository[business.Business](testEnv.Database)
	biz, err := bizRepo.FindOne(ctx, bizRepo.ScopeEquals(business.BusinessSchema.Descriptor, "biz"))
	s.NoError(err)

	zoneRepo := database.NewRepository[business.ShippingZone](testEnv.Database)
	count, err := zoneRepo.Count(ctx, zoneRepo.ScopeBusinessID(biz.ID))
	s.NoError(err)
	s.Equal(int64(2), count)

	zones, status := s.listZones(token, "biz")
	s.Equal(http.StatusOK, status)
	_, ok := s.findZoneByName(zones, "GCC")
	s.True(ok)
}

func TestBusinessShippingZonesSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BusinessShippingZonesSuite))
}
