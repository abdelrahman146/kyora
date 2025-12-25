package e2e_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/asset"
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

type AssetGCSuite struct {
	suite.Suite
	helper       *AccountTestHelper
	assetStorage *asset.Storage
	assetService *asset.Service
}

func (s *AssetGCSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
	cacheClient := cache.NewConnection([]string{testEnv.CacheAddr})
	s.assetStorage = asset.NewStorage(testEnv.Database, cacheClient)
	s.assetService = asset.NewService(s.assetStorage, database.NewAtomicProcess(testEnv.Database), nil)
}

func (s *AssetGCSuite) SetupTest() {
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

func (s *AssetGCSuite) TearDownTest() {
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

func (s *AssetGCSuite) uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s_%s@example.com", prefix, id.Base62(10))
}

func (s *AssetGCSuite) createBusiness(ctx context.Context, token string) (descriptor string) {
	descriptor = fmt.Sprintf("biz-%s", strings.ToLower(id.Base62(10)))
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

func (s *AssetGCSuite) getBusinessID(ctx context.Context, bizDescriptor, token string) string {
	resp, err := s.helper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/%s", bizDescriptor), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var bizResult map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &bizResult))
	bizObj := bizResult["business"].(map[string]interface{})
	return bizObj["id"].(string)
}

func (s *AssetGCSuite) TestGC_DeletesExpiredPendingUpload() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, s.uniqueEmail("admin"), "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))

	bizDescriptor := s.createBusiness(ctx, token)

	createReq := map[string]interface{}{
		"idempotencyKey": fmt.Sprintf("idem_%s", id.Base62(10)),
		"fileName":       "logo.png",
		"contentType":    "image/png",
		"sizeBytes":      int64(12),
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/assets/uploads/logo", bizDescriptor), createReq, token)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &created))
	resp.Body.Close()

	assetID := created["assetId"].(string)
	businessID := s.getBusinessID(ctx, bizDescriptor, token)

	a, err := s.assetStorage.GetByID(ctx, businessID, assetID)
	s.NoError(err)
	s.NotNil(a)
	past := time.Now().UTC().Add(-2 * time.Hour)
	a.UploadExpiresAt = &past
	s.NoError(s.assetStorage.Update(ctx, a))

	gcNow := time.Now().UTC()
	gcRes, err := s.assetService.GarbageCollect(ctx, asset.GarbageCollectOptions{Now: gcNow, PendingLimit: 100, OrphanLimit: 1, OrphanMinAge: time.Hour})
	s.NoError(err)
	s.GreaterOrEqual(gcRes.DeletedAssets, 1)

	_, err = s.assetStorage.FindByID(ctx, assetID)
	s.Error(err)
	s.True(database.IsRecordNotFound(err))
}

func (s *AssetGCSuite) TestGC_DeletesReadyOrphanAndRemovesLocalFile() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, s.uniqueEmail("admin"), "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))
	bizDescriptor := s.createBusiness(ctx, token)

	content := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0x00, 0x00, 0x00, 0x00}
	createReq := map[string]interface{}{
		"idempotencyKey": fmt.Sprintf("idem_%s", id.Base62(10)),
		"fileName":       "logo.png",
		"contentType":    "image/png",
		"sizeBytes":      int64(len(content)),
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/assets/uploads/logo", bizDescriptor), createReq, token)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &created))
	resp.Body.Close()
	assetID := created["assetId"].(string)
	upload := created["upload"].(map[string]interface{})
	uploadURL := upload["url"].(string)

	putHeaders := map[string]string{"Content-Type": "image/png"}
	putResp, err := s.helper.Client.AuthenticatedRequestRaw("PUT", uploadURL, content, putHeaders, token)
	s.NoError(err)
	putResp.Body.Close()
	s.Equal(http.StatusNoContent, putResp.StatusCode)

	completeResp, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/assets/uploads/%s/complete/business_logo", bizDescriptor, assetID), nil, token)
	s.NoError(err)
	s.Equal(http.StatusOK, completeResp.StatusCode)
	var completed map[string]interface{}
	s.NoError(testutils.DecodeJSON(completeResp, &completed))
	completeResp.Body.Close()

	businessID := s.getBusinessID(ctx, bizDescriptor, token)
	a, err := s.assetStorage.GetByID(ctx, businessID, assetID)
	s.NoError(err)
	s.NotNil(a)
	s.NotEmpty(a.LocalFilePath)
	_, statErr := os.Stat(a.LocalFilePath)
	s.NoError(statErr)

	// Confirm public GET works before GC.
	getResp, err := s.helper.Client.Get(fmt.Sprintf("/v1/public/assets/%s", assetID))
	s.NoError(err)
	s.Equal(http.StatusOK, getResp.StatusCode)
	_, _ = io.ReadAll(getResp.Body)
	getResp.Body.Close()

	gcNow := time.Now().UTC().Add(2 * time.Minute)
	gcRes, err := s.assetService.GarbageCollect(ctx, asset.GarbageCollectOptions{Now: gcNow, PendingLimit: 1, OrphanLimit: 100, OrphanMinAge: 0})
	s.NoError(err)
	s.GreaterOrEqual(gcRes.DeletedAssets, 1)

	_, err = s.assetStorage.FindByID(ctx, assetID)
	s.Error(err)
	s.True(database.IsRecordNotFound(err))

	_, statErr = os.Stat(a.LocalFilePath)
	s.True(os.IsNotExist(statErr))

	getAfter, err := s.helper.Client.Get(fmt.Sprintf("/v1/public/assets/%s", assetID))
	s.NoError(err)
	getAfter.Body.Close()
	s.Equal(http.StatusNotFound, getAfter.StatusCode)
}

func (s *AssetGCSuite) TestGC_DoesNotDeleteReferencedReadyAsset() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, s.uniqueEmail("admin"), "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))
	bizDescriptor := s.createBusiness(ctx, token)

	content := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0x00, 0x00, 0x00, 0x00}
	createReq := map[string]interface{}{
		"idempotencyKey": fmt.Sprintf("idem_%s", id.Base62(10)),
		"fileName":       "logo.png",
		"contentType":    "image/png",
		"sizeBytes":      int64(len(content)),
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/assets/uploads/logo", bizDescriptor), createReq, token)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &created))
	resp.Body.Close()
	assetID := created["assetId"].(string)
	upload := created["upload"].(map[string]interface{})
	uploadURL := upload["url"].(string)

	putHeaders := map[string]string{"Content-Type": "image/png"}
	putResp, err := s.helper.Client.AuthenticatedRequestRaw("PUT", uploadURL, content, putHeaders, token)
	s.NoError(err)
	putResp.Body.Close()
	s.Equal(http.StatusNoContent, putResp.StatusCode)

	completeResp, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/assets/uploads/%s/complete/business_logo", bizDescriptor, assetID), nil, token)
	s.NoError(err)
	s.Equal(http.StatusOK, completeResp.StatusCode)
	var completed map[string]interface{}
	s.NoError(testutils.DecodeJSON(completeResp, &completed))
	completeResp.Body.Close()
	publicURL := completed["publicUrl"].(string)

	patchPayload := map[string]interface{}{"logoUrl": publicURL}
	patchResp, err := s.helper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/businesses/%s", bizDescriptor), patchPayload, token)
	s.NoError(err)
	patchResp.Body.Close()
	s.Equal(http.StatusOK, patchResp.StatusCode)

	gcNow := time.Now().UTC().Add(2 * time.Minute)
	gcRes, err := s.assetService.GarbageCollect(ctx, asset.GarbageCollectOptions{Now: gcNow, PendingLimit: 1, OrphanLimit: 100, OrphanMinAge: 0})
	s.NoError(err)
	s.Equal(0, gcRes.Errors)

	getAfter, err := s.helper.Client.Get(fmt.Sprintf("/v1/public/assets/%s", assetID))
	s.NoError(err)
	getAfter.Body.Close()
	s.Equal(http.StatusOK, getAfter.StatusCode)
}

func TestAssetGCSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(AssetGCSuite))
}
