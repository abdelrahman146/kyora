package e2e_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

type AssetUploadsSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *AssetUploadsSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *AssetUploadsSuite) uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s_%s@example.com", prefix, id.Base62(10))
}

func (s *AssetUploadsSuite) createBusiness(ctx context.Context, token string) (descriptor string) {
	descriptor = fmt.Sprintf("biz-%s", strings.ToLower(id.Base62(10)))
	payload := map[string]interface{}{
		"name":        "Test Business",
		"descriptor":  descriptor,
		"countryCode": "eg",
		"currency":    "usd",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/businesses", payload, token)
	s.NoError(err)
	if s.NotNil(resp) && resp.StatusCode != http.StatusCreated {
		body, _ := testutils.ReadBody(resp)
		s.FailNow("unexpected create business response", "status=%d body=%s", resp.StatusCode, body)
	}
	s.Equal(http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Len(result, 1)
	s.Contains(result, "business")
	biz := result["business"].(map[string]interface{})
	s.Equal(descriptor, biz["descriptor"])
	return descriptor
}

func (s *AssetUploadsSuite) TestLogoUpload_LocalProvider_HappyPath_PublicGET() {
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
	s.Contains(created, "assetId")
	s.Contains(created, "upload")
	s.Contains(created, "publicUrl")
	assetID := created["assetId"].(string)
	upload := created["upload"].(map[string]interface{})
	s.Equal("PUT", upload["method"])
	uploadURL := upload["url"].(string)
	s.Contains(uploadURL, fmt.Sprintf("/v1/businesses/%s/assets/uploads/%s/content/business_logo", bizDescriptor, assetID))

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
	s.Equal(assetID, completed["assetId"])
	s.Equal("ready", completed["status"])
	publicURL := completed["publicUrl"].(string)
	s.Contains(publicURL, "/v1/public/assets/")

	getResp, err := s.helper.Client.Get(fmt.Sprintf("/v1/public/assets/%s", assetID))
	s.NoError(err)
	s.Equal(http.StatusOK, getResp.StatusCode)
	s.Equal("image/png", getResp.Header.Get("Content-Type"))
	body, err := io.ReadAll(getResp.Body)
	s.NoError(err)
	getResp.Body.Close()
	s.Equal(content, body)

	// Integration with existing business flow: update the business logo URL.
	patchPayload := map[string]interface{}{"logoUrl": publicURL}
	patchResp, err := s.helper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/businesses/%s", bizDescriptor), patchPayload, token)
	s.NoError(err)
	defer patchResp.Body.Close()
	s.Equal(http.StatusOK, patchResp.StatusCode)

	getBizResp, err := s.helper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/%s", bizDescriptor), nil, token)
	s.NoError(err)
	s.Equal(http.StatusOK, getBizResp.StatusCode)
	var bizResult map[string]interface{}
	s.NoError(testutils.DecodeJSON(getBizResp, &bizResult))
	s.Contains(bizResult, "business")
	biz := bizResult["business"].(map[string]interface{})
	s.Equal(publicURL, biz["logoUrl"])
}

func (s *AssetUploadsSuite) TestCreateUpload_IdempotencyReplay_SameResponse() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, s.uniqueEmail("admin"), "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))
	bizDescriptor := s.createBusiness(ctx, token)

	idemKey := fmt.Sprintf("idem_%s", id.Base62(10))
	createReq := map[string]interface{}{
		"idempotencyKey": idemKey,
		"fileName":       "logo.png",
		"contentType":    "image/png",
		"sizeBytes":      int64(12),
	}

	r1, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/assets/uploads/logo", bizDescriptor), createReq, token)
	s.NoError(err)
	s.Equal(http.StatusOK, r1.StatusCode)
	var res1 map[string]interface{}
	s.NoError(testutils.DecodeJSON(r1, &res1))
	id1 := res1["assetId"].(string)

	// Asset upload create has a minimum spacing rate-limit (to blunt abuse).
	// Sleep slightly over that threshold so we can deterministically test replay.
	time.Sleep(300 * time.Millisecond)

	r2, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/assets/uploads/logo", bizDescriptor), createReq, token)
	s.NoError(err)
	s.Equal(http.StatusOK, r2.StatusCode)
	var res2 map[string]interface{}
	s.NoError(testutils.DecodeJSON(r2, &res2))
	id2 := res2["assetId"].(string)

	s.Equal(id1, id2)
}

func (s *AssetUploadsSuite) TestCreateUpload_IdempotencyConflict_WhenRequestDiffers() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, s.uniqueEmail("admin"), "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))
	bizDescriptor := s.createBusiness(ctx, token)

	idemKey := fmt.Sprintf("idem_%s", id.Base62(10))
	base := map[string]interface{}{
		"idempotencyKey": idemKey,
		"fileName":       "logo.png",
		"contentType":    "image/png",
		"sizeBytes":      int64(12),
	}
	conflict := map[string]interface{}{
		"idempotencyKey": idemKey,
		"fileName":       "logo.png",
		"contentType":    "image/png",
		"sizeBytes":      int64(13),
	}

	r1, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/assets/uploads/logo", bizDescriptor), base, token)
	s.NoError(err)
	r1.Body.Close()
	s.Equal(http.StatusOK, r1.StatusCode)

	r2, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/assets/uploads/logo", bizDescriptor), conflict, token)
	s.NoError(err)
	defer r2.Body.Close()
	s.Equal(http.StatusConflict, r2.StatusCode)
}

func (s *AssetUploadsSuite) TestPutContent_ValidatesSizeAndContentType() {
	ctx := context.Background()
	_, ws, token, err := s.helper.CreateTestUser(ctx, s.uniqueEmail("admin"), "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws.ID))
	bizDescriptor := s.createBusiness(ctx, token)

	content := []byte("1234567890")
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
	assetID := created["assetId"].(string)
	upload := created["upload"].(map[string]interface{})
	uploadURL := upload["url"].(string)

	s.Run("size mismatch", func() {
		bad := append([]byte{}, content...)
		bad = append(bad, 'x')
		putHeaders := map[string]string{"Content-Type": "image/png"}
		putResp, err := s.helper.Client.AuthenticatedRequestRaw("PUT", uploadURL, bad, putHeaders, token)
		s.NoError(err)
		defer putResp.Body.Close()
		s.Equal(http.StatusBadRequest, putResp.StatusCode)
	})

	s.Run("content type mismatch", func() {
		putHeaders := map[string]string{"Content-Type": "image/jpeg"}
		putResp, err := s.helper.Client.AuthenticatedRequestRaw("PUT", uploadURL, content, putHeaders, token)
		s.NoError(err)
		defer putResp.Body.Close()
		s.Equal(http.StatusConflict, putResp.StatusCode)
	})

	putHeaders := map[string]string{"Content-Type": "image/png"}
	okPut, err := s.helper.Client.AuthenticatedRequestRaw("PUT", uploadURL, content, putHeaders, token)
	s.NoError(err)
	defer okPut.Body.Close()
	s.Equal(http.StatusNoContent, okPut.StatusCode)

	completeResp, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/assets/uploads/%s/complete/business_logo", bizDescriptor, assetID), nil, token)
	s.NoError(err)
	defer completeResp.Body.Close()
	s.Equal(http.StatusOK, completeResp.StatusCode)
}

func (s *AssetUploadsSuite) TestCompleteUpload_FailsWhenNotUploadedYet() {
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
	assetID := created["assetId"].(string)

	completeResp, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/assets/uploads/%s/complete/business_logo", bizDescriptor, assetID), nil, token)
	s.NoError(err)
	defer completeResp.Body.Close()
	s.Equal(http.StatusConflict, completeResp.StatusCode)
}

func (s *AssetUploadsSuite) TestCrossWorkspaceIsolation_ReturnsNotFound() {
	ctx := context.Background()
	_, ws1, token1, err := s.helper.CreateTestUser(ctx, s.uniqueEmail("admin1"), "ValidPassword123!", "Admin", "One", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.helper.CreateTestSubscription(ctx, ws1.ID))
	bizDescriptor := s.createBusiness(ctx, token1)

	_, _, token2, err := s.helper.CreateTestUser(ctx, s.uniqueEmail("admin2"), "ValidPassword123!", "Admin", "Two", role.RoleAdmin)
	s.NoError(err)

	createReq := map[string]interface{}{
		"idempotencyKey": fmt.Sprintf("idem_%s", id.Base62(10)),
		"fileName":       "logo.png",
		"contentType":    "image/png",
		"sizeBytes":      int64(12),
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/assets/uploads/logo", bizDescriptor), createReq, token2)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func TestAssetUploadsSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(AssetUploadsSuite))
}
