package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/abdelrahman146/kyora/internal/domain/asset"
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

type AssetDeleteSuite struct {
	suite.Suite
	helper       *AccountTestHelper
	assetStorage *asset.Storage
}

func (s *AssetDeleteSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
	cacheClient := cache.NewConnection([]string{testEnv.CacheAddr})
	s.assetStorage = asset.NewStorage(testEnv.Database, cacheClient)
}

func (s *AssetDeleteSuite) SetupTest() {
	s.NoError(testutils.TruncateTables(testEnv.Database,
		"users",
		"workspaces",
		"sessions",
		"subscriptions",
		"plans",
		"businesses",
		"products",
		"variants",
		"categories",
		"uploaded_assets",
	))
}

func (s *AssetDeleteSuite) TearDownTest() {
	s.NoError(testutils.TruncateTables(testEnv.Database,
		"users",
		"workspaces",
		"sessions",
		"subscriptions",
		"plans",
		"businesses",
		"products",
		"variants",
		"categories",
		"uploaded_assets",
	))
}

func (s *AssetDeleteSuite) uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s_%s@example.com", prefix, id.Base62(10))
}

func (s *AssetDeleteSuite) createBusiness(ctx context.Context, token string) string {
	descriptor := fmt.Sprintf("biz-%s", strings.ToLower(id.Base62(10)))
	payload := map[string]interface{}{
		"name":        "Test Business",
		"descriptor":  descriptor,
		"countryCode": "eg",
		"currency":    "usd",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)
	return descriptor
}

func (s *AssetDeleteSuite) getBusinessID(ctx context.Context, descriptor, token string) string {
	resp, err := s.helper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/%s", descriptor), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var bizResult map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &bizResult))
	bizObj := bizResult["business"].(map[string]interface{})
	return bizObj["id"].(string)
}

func (s *AssetDeleteSuite) TestDelete_UnreferencedReadyAsset() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, s.uniqueEmail("admin"), "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))

	descriptor := s.createBusiness(ctx, token)

	content := []byte{0x89, 'P', 'N', 'G'}
	assetID, _ := UploadTestAsset(ctx, s.T(), s.helper.Client, AssetUploadParams{
		BusinessDescriptor: descriptor,
		Token:              token,
		CreatePath:         "assets/uploads/product-photo",
		CompletePurpose:    "product_photo",
		ContentType:        "image/png",
		Content:            content,
	})

	businessID := s.getBusinessID(ctx, descriptor, token)
	a, err := s.assetStorage.GetByID(ctx, businessID, assetID)
	s.NoError(err)
	s.NotNil(a)
	localPath := a.LocalFilePath
	if strings.TrimSpace(localPath) != "" {
		_, statErr := os.Stat(localPath)
		s.NoError(statErr)
	}

	delResp, err := s.helper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/%s/assets/%s", descriptor, assetID), nil, token)
	s.NoError(err)
	delResp.Body.Close()
	s.Equal(http.StatusNoContent, delResp.StatusCode)

	_, err = s.assetStorage.FindByID(ctx, assetID)
	s.Error(err)
	s.True(database.IsRecordNotFound(err))

	if strings.TrimSpace(localPath) != "" {
		_, statErr := os.Stat(localPath)
		s.True(os.IsNotExist(statErr))
	}

	getResp, err := s.helper.Client.Get(fmt.Sprintf("/v1/public/assets/%s", assetID))
	s.NoError(err)
	getResp.Body.Close()
	s.Equal(http.StatusNotFound, getResp.StatusCode)
}

func (s *AssetDeleteSuite) TestDelete_ReferencedAssetRejected() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, s.uniqueEmail("admin"), "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))

	descriptor := s.createBusiness(ctx, token)

	content := []byte{0x89, 'P', 'N', 'G'}
	assetID, publicURL := UploadTestAsset(ctx, s.T(), s.helper.Client, AssetUploadParams{
		BusinessDescriptor: descriptor,
		Token:              token,
		CreatePath:         "assets/uploads/logo",
		CompletePurpose:    "business_logo",
		ContentType:        "image/png",
		Content:            content,
	})

	patchPayload := map[string]interface{}{"logoUrl": publicURL}
	patchResp, err := s.helper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/businesses/%s", descriptor), patchPayload, token)
	s.NoError(err)
	patchResp.Body.Close()
	s.Equal(http.StatusOK, patchResp.StatusCode)

	delResp, err := s.helper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/%s/assets/%s", descriptor, assetID), nil, token)
	s.NoError(err)
	delResp.Body.Close()
	s.Equal(http.StatusConflict, delResp.StatusCode)

	getResp, err := s.helper.Client.Get(fmt.Sprintf("/v1/public/assets/%s", assetID))
	s.NoError(err)
	getResp.Body.Close()
	s.Equal(http.StatusOK, getResp.StatusCode)
}

func (s *AssetDeleteSuite) TestDelete_PermissionDeniedForViewer() {
	ctx := context.Background()

	workspace, users, err := testutils.CreateWorkspaceWithUsers(ctx, testEnv.Database, s.uniqueEmail("owner"), "ValidPassword123!", []struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
		Role      role.Role
	}{
		{Email: s.uniqueEmail("viewer"), Password: "ValidPassword123!", FirstName: "View", LastName: "User", Role: role.RoleUser},
	})
	s.NoError(err)
	s.Len(users, 2)
	owner := users[0]
	viewer := users[1]

	s.NoError(testutils.CreateTestSubscription(ctx, testEnv.Database, workspace.ID))

	ownerToken, err := testutils.LoginAndGetToken(s.helper.Client, owner.Email, "ValidPassword123!")
	s.NoError(err)
	viewerToken, err := testutils.LoginAndGetToken(s.helper.Client, viewer.Email, "ValidPassword123!")
	s.NoError(err)

	descriptor := s.createBusiness(ctx, ownerToken)

	content := []byte{0x89, 'P', 'N', 'G'}
	assetID, _ := UploadTestAsset(ctx, s.T(), s.helper.Client, AssetUploadParams{
		BusinessDescriptor: descriptor,
		Token:              ownerToken,
		CreatePath:         "assets/uploads/variant-photo",
		CompletePurpose:    "variant_photo",
		ContentType:        "image/png",
		Content:            content,
	})

	delResp, err := s.helper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/%s/assets/%s", descriptor, assetID), nil, viewerToken)
	s.NoError(err)
	delResp.Body.Close()
	s.Equal(http.StatusForbidden, delResp.StatusCode)
}

func TestAssetDeleteSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(AssetDeleteSuite))
}
