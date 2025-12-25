package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

type BusinessAssetsSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *BusinessAssetsSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *BusinessAssetsSuite) resetDB() {
	s.NoError(testutils.TruncateTables(testEnv.Database,
		"users",
		"workspaces",
		"businesses",
		"shipping_zones",
		"uploaded_assets",
		"subscriptions",
		"plans",
	))
}

func (s *BusinessAssetsSuite) SetupTest() {
	s.resetDB()
}

func (s *BusinessAssetsSuite) TearDownTest() {
	s.resetDB()
}

func (s *BusinessAssetsSuite) uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s_%s@example.com", prefix, strings.ToLower(id.Base62(8)))
}

func (s *BusinessAssetsSuite) TestCreateBusiness_WithUploadedLogo() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, s.uniqueEmail("admin"), "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))

	bizDescriptor := fmt.Sprintf("logo-%s", strings.ToLower(id.Base62(6)))
	content := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0x00, 0x00, 0x00, 0x00}
	_, logoURL := UploadTestAsset(ctx, s.T(), s.helper.Client, AssetUploadParams{
		BusinessDescriptor: bizDescriptor,
		Token:              token,
		CreatePath:         "assets/uploads/logo",
		CompletePurpose:    "business_logo",
		ContentType:        "image/png",
		Content:            content,
	})

	payload := map[string]interface{}{
		"name":        "Logo Biz",
		"descriptor":  bizDescriptor,
		"countryCode": "eg",
		"currency":    "usd",
		"logoUrl":     logoURL,
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	biz := result["business"].(map[string]interface{})
	s.Equal(logoURL, biz["logoUrl"])
	s.Equal(bizDescriptor, biz["descriptor"])

	getResp, err := s.helper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/%s", bizDescriptor), nil, token)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusOK, getResp.StatusCode)

	var fetched map[string]interface{}
	s.NoError(testutils.DecodeJSON(getResp, &fetched))
	fetchedBiz := fetched["business"].(map[string]interface{})
	s.Equal(logoURL, fetchedBiz["logoUrl"])
}

func TestBusinessAssetsSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BusinessAssetsSuite))
}
