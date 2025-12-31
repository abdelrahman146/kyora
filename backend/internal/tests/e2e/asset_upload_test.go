package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/asset"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type AssetUploadSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *AssetUploadSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *AssetUploadSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "subscriptions", "plans", "uploaded_assets")
	s.NoError(err)
}

func (s *AssetUploadSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "subscriptions", "plans", "uploaded_assets")
	s.NoError(err)
}

func (s *AssetUploadSuite) createBusiness(ctx context.Context, workspaceID, descriptor string) string {
	bizRepo := database.NewRepository[business.Business](testEnv.Database)
	biz := &business.Business{
		WorkspaceID:  workspaceID,
		Descriptor:   descriptor,
		Name:         "Test Business",
		CountryCode:  "EG",
		Currency:     "USD",
		VatRate:      decimal.NewFromFloat(0.14),
		SafetyBuffer: decimal.NewFromFloat(100),
	}
	s.NoError(bizRepo.CreateOne(ctx, biz))
	return biz.ID
}

func (s *AssetUploadSuite) createSubscription(ctx context.Context, workspaceID string) {
	planRepo := database.NewRepository[billing.Plan](testEnv.Database)
	stripePlanID := "price_test"
	plan := &billing.Plan{
		Descriptor:   "test-plan",
		Name:         "Test Plan",
		Description:  "Test Plan for E2E",
		StripePlanID: &stripePlanID,
		Price:        decimal.NewFromFloat(10.00),
		Currency:     "usd",
		BillingCycle: billing.BillingCycleMonthly,
		Features: billing.PlanFeature{
			CustomerManagement:  true,
			InventoryManagement: true,
			OrderManagement:     true,
		},
		Limits: billing.PlanLimit{
			MaxOrdersPerMonth: 1000,
			MaxTeamMembers:    5,
			MaxBusinesses:     3,
		},
	}
	s.NoError(planRepo.CreateOne(ctx, plan))

	subRepo := database.NewRepository[billing.Subscription](testEnv.Database)
	sub := &billing.Subscription{
		WorkspaceID:      workspaceID,
		PlanID:           plan.ID,
		Status:           billing.SubscriptionStatusActive,
		StripeSubID:      "sub_test",
		CurrentPeriodEnd: time.Now().Add(30 * 24 * time.Hour),
	}
	s.NoError(subRepo.CreateOne(ctx, sub))
}

func (s *AssetUploadSuite) TestGenerateUploadURLs_Success() {
	ctx := context.Background()

	// Create test user, workspace, business
	_, ws, token, err := s.helper.CreateTestUser(ctx, "test@example.com", "Password123!", "Test", "User", role.RoleAdmin)
	s.NoError(err)
	s.createSubscription(ctx, ws.ID)
	bizID := s.createBusiness(ctx, ws.ID, "test-biz")

	// Request upload URLs
	reqBody := asset.GenerateUploadURLsRequest{
		Files: []asset.FileUploadRequest{
			{FileName: "logo.png", ContentType: "image/png", SizeBytes: 50_000},
		},
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/businesses/test-biz/assets/uploads", e2eBaseURL), bytes.NewReader(body))
	s.NoError(err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var result asset.GenerateUploadURLsResponse
	s.NoError(json.NewDecoder(resp.Body).Decode(&result))

	s.Len(result.Uploads, 1)
	upload := result.Uploads[0]

	// Verify assetId is returned
	s.NotEmpty(upload.AssetID)

	// Verify asset record was created
	assetRepo := database.NewRepository[asset.Asset](testEnv.Database)
	dbAsset, err := assetRepo.FindByID(ctx, upload.AssetID)
	s.NoError(err)
	s.Equal(ws.ID, dbAsset.WorkspaceID)
	s.Equal(bizID, dbAsset.BusinessID)
	s.Equal("image/png", dbAsset.ContentType)
	s.Equal(int64(50_000), dbAsset.SizeBytes)
}

func (s *AssetUploadSuite) TestGenerateUploadURLs_MultipleFiles() {
	ctx := context.Background()

	_, ws, token, err := s.helper.CreateTestUser(ctx, "test2@example.com", "Password123!", "Test", "User", role.RoleAdmin)
	s.NoError(err)
	s.createSubscription(ctx, ws.ID)
	s.createBusiness(ctx, ws.ID, "test-biz")

	// Request upload URLs for multiple files
	reqBody := asset.GenerateUploadURLsRequest{
		Files: []asset.FileUploadRequest{
			{FileName: "photo1.jpg", ContentType: "image/jpeg", SizeBytes: 100_000},
			{FileName: "photo2.jpg", ContentType: "image/jpeg", SizeBytes: 150_000},
			{FileName: "photo3.png", ContentType: "image/png", SizeBytes: 80_000},
		},
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/businesses/test-biz/assets/uploads", e2eBaseURL), bytes.NewReader(body))
	s.NoError(err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var result asset.GenerateUploadURLsResponse
	s.NoError(json.NewDecoder(resp.Body).Decode(&result))

	// Should have descriptors for all 3 files
	s.Len(result.Uploads, 3)

	// Each should have unique assetId
	assetIDs := make(map[string]bool)
	for _, upload := range result.Uploads {
		s.NotEmpty(upload.AssetID)
		s.NotContains(assetIDs, upload.AssetID)
		assetIDs[upload.AssetID] = true
	}
}

func (s *AssetUploadSuite) TestGenerateUploadURLs_ValidationErrors() {
	ctx := context.Background()

	_, ws, token, err := s.helper.CreateTestUser(ctx, "test3@example.com", "Password123!", "Test", "User", role.RoleAdmin)
	s.NoError(err)
	s.createSubscription(ctx, ws.ID)
	s.createBusiness(ctx, ws.ID, "test-biz")

	testCases := []struct {
		name           string
		request        asset.GenerateUploadURLsRequest
		expectedStatus int
	}{
		{
			name:           "empty files array",
			request:        asset.GenerateUploadURLsRequest{Files: []asset.FileUploadRequest{}},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing fileName",
			request: asset.GenerateUploadURLsRequest{
				Files: []asset.FileUploadRequest{{ContentType: "image/png", SizeBytes: 1000}},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing content type",
			request: asset.GenerateUploadURLsRequest{
				Files: []asset.FileUploadRequest{{FileName: "test.png", SizeBytes: 1000}},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "zero size",
			request: asset.GenerateUploadURLsRequest{
				Files: []asset.FileUploadRequest{{FileName: "test.png", ContentType: "image/png", SizeBytes: 0}},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			body, _ := json.Marshal(tc.request)
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/businesses/test-biz/assets/uploads", e2eBaseURL), bytes.NewReader(body))
			s.NoError(err)
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tc.expectedStatus, resp.StatusCode)
		})
	}
}

func (s *AssetUploadSuite) TestGenerateUploadURLs_WithoutAuthentication() {
	reqBody := asset.GenerateUploadURLsRequest{
		Files: []asset.FileUploadRequest{
			{FileName: "test.png", ContentType: "image/png", SizeBytes: 1000},
		},
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/businesses/test-biz/assets/uploads", e2eBaseURL), bytes.NewReader(body))
	s.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *AssetUploadSuite) TestGenerateUploadURLs_FileTypeValidation() {
	ctx := context.Background()

	_, ws, token, err := s.helper.CreateTestUser(ctx, "filetest@example.com", "Password123!", "Test", "User", role.RoleAdmin)
	s.NoError(err)
	s.createSubscription(ctx, ws.ID)
	s.createBusiness(ctx, ws.ID, "test-biz")

	testCases := []struct {
		name           string
		fileName       string
		contentType    string
		sizeBytes      int64
		expectSuccess  bool
		expectedReason string
	}{
		// Images - should succeed
		{
			name:          "valid jpeg image",
			fileName:      "photo.jpg",
			contentType:   "image/jpeg",
			sizeBytes:     5_000_000, // 5MB
			expectSuccess: true,
		},
		{
			name:          "valid png image",
			fileName:      "logo.png",
			contentType:   "image/png",
			sizeBytes:     2_000_000,
			expectSuccess: true,
		},
		{
			name:          "valid webp image",
			fileName:      "banner.webp",
			contentType:   "image/webp",
			sizeBytes:     3_000_000,
			expectSuccess: true,
		},
		// Videos - should succeed
		{
			name:          "valid mp4 video",
			fileName:      "demo.mp4",
			contentType:   "video/mp4",
			sizeBytes:     50_000_000, // 50MB
			expectSuccess: true,
		},
		{
			name:          "valid mov video",
			fileName:      "tutorial.mov",
			contentType:   "video/quicktime",
			sizeBytes:     60_000_000,
			expectSuccess: true,
		},
		// Documents - should succeed
		{
			name:          "valid pdf document",
			fileName:      "invoice.pdf",
			contentType:   "application/pdf",
			sizeBytes:     5_000_000,
			expectSuccess: true,
		},
		// Unsupported types - should fail
		{
			name:           "unsupported exe file",
			fileName:       "virus.exe",
			contentType:    "application/x-msdownload",
			sizeBytes:      1_000,
			expectSuccess:  false,
			expectedReason: "file type not allowed",
		},
		{
			name:           "unsupported php file",
			fileName:       "malicious.php",
			contentType:    "application/x-httpd-php",
			sizeBytes:      500,
			expectSuccess:  false,
			expectedReason: "file type not allowed",
		},
		// Size limit violations - should fail
		{
			name:           "image too large",
			fileName:       "huge.jpg",
			contentType:    "image/jpeg",
			sizeBytes:      15_000_000, // 15MB > 10MB default limit
			expectSuccess:  false,
			expectedReason: "file too large",
		},
		{
			name:           "video too large",
			fileName:       "movie.mp4",
			contentType:    "video/mp4",
			sizeBytes:      150_000_000, // 150MB > 100MB default limit
			expectSuccess:  false,
			expectedReason: "file too large",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			reqBody := asset.GenerateUploadURLsRequest{
				Files: []asset.FileUploadRequest{
					{FileName: tc.fileName, ContentType: tc.contentType, SizeBytes: tc.sizeBytes},
				},
			}
			body, _ := json.Marshal(reqBody)

			req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/businesses/test-biz/assets/uploads", e2eBaseURL), bytes.NewReader(body))
			s.NoError(err)
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			s.NoError(err)
			defer resp.Body.Close()

			if tc.expectSuccess {
				s.Equal(http.StatusOK, resp.StatusCode, "Expected success for %s", tc.name)
			} else {
				s.Equal(http.StatusBadRequest, resp.StatusCode, "Expected validation error for %s", tc.name)
			}
		})
	}
}

func (s *AssetUploadSuite) TestGenerateUploadURLs_ThumbnailGeneration() {
	ctx := context.Background()

	_, ws, token, err := s.helper.CreateTestUser(ctx, "thumb@example.com", "Password123!", "Test", "User", role.RoleAdmin)
	s.NoError(err)
	s.createSubscription(ctx, ws.ID)
	s.createBusiness(ctx, ws.ID, "test-biz")

	testCases := []struct {
		name             string
		fileName         string
		contentType      string
		sizeBytes        int64
		expectThumbnail  bool
		expectedCategory string
	}{
		{
			name:             "image should have thumbnail",
			fileName:         "product.jpg",
			contentType:      "image/jpeg",
			sizeBytes:        5_000_000,
			expectThumbnail:  true,
			expectedCategory: "image",
		},
		{
			name:             "video should have thumbnail",
			fileName:         "demo.mp4",
			contentType:      "video/mp4",
			sizeBytes:        30_000_000,
			expectThumbnail:  true,
			expectedCategory: "video",
		},
		{
			name:             "pdf should not have thumbnail",
			fileName:         "invoice.pdf",
			contentType:      "application/pdf",
			sizeBytes:        2_000_000,
			expectThumbnail:  false,
			expectedCategory: "document",
		},
		{
			name:             "audio should not have thumbnail",
			fileName:         "podcast.mp3",
			contentType:      "audio/mpeg",
			sizeBytes:        10_000_000,
			expectThumbnail:  false,
			expectedCategory: "audio",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			reqBody := asset.GenerateUploadURLsRequest{
				Files: []asset.FileUploadRequest{
					{FileName: tc.fileName, ContentType: tc.contentType, SizeBytes: tc.sizeBytes},
				},
			}
			body, _ := json.Marshal(reqBody)

			req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/businesses/test-biz/assets/uploads", e2eBaseURL), bytes.NewReader(body))
			s.NoError(err)
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(http.StatusOK, resp.StatusCode)

			var result asset.GenerateUploadURLsResponse
			s.NoError(json.NewDecoder(resp.Body).Decode(&result))

			s.Len(result.Uploads, 1)
			upload := result.Uploads[0]

			if tc.expectThumbnail {
				s.NotNil(upload.Thumbnail, "Expected thumbnail for %s", tc.name)
				s.NotEmpty(upload.Thumbnail.AssetID)
				s.NotEmpty(upload.Thumbnail.URL)
				s.NotEmpty(upload.Thumbnail.PublicURL)
				s.NotEmpty(upload.Thumbnail.CDNURL)
				s.Equal("image/jpeg", upload.Thumbnail.ContentType, "Thumbnails should always be JPEG")

				// Verify thumbnail asset was created in DB
				assetRepo := database.NewRepository[asset.Asset](testEnv.Database)
				thumbAsset, err := assetRepo.FindByID(ctx, upload.Thumbnail.AssetID)
				s.NoError(err)
				s.Equal("image/jpeg", thumbAsset.ContentType)
				s.Equal("image", thumbAsset.FileCategory)
			} else {
				s.Nil(upload.Thumbnail, "Expected no thumbnail for %s", tc.name)
			}

			// Verify main asset has correct category
			assetRepo := database.NewRepository[asset.Asset](testEnv.Database)
			mainAsset, err := assetRepo.FindByID(ctx, upload.AssetID)
			s.NoError(err)
			s.Equal(tc.expectedCategory, mainAsset.FileCategory)
		})
	}
}

func (s *AssetUploadSuite) TestGenerateUploadURLs_CDNUrls() {
	ctx := context.Background()

	_, ws, token, err := s.helper.CreateTestUser(ctx, "cdn@example.com", "Password123!", "Test", "User", role.RoleAdmin)
	s.NoError(err)
	s.createSubscription(ctx, ws.ID)
	s.createBusiness(ctx, ws.ID, "test-biz")

	reqBody := asset.GenerateUploadURLsRequest{
		Files: []asset.FileUploadRequest{
			{FileName: "logo.png", ContentType: "image/png", SizeBytes: 50_000},
		},
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/businesses/test-biz/assets/uploads", e2eBaseURL), bytes.NewReader(body))
	s.NoError(err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var result asset.GenerateUploadURLsResponse
	s.NoError(json.NewDecoder(resp.Body).Decode(&result))

	s.Len(result.Uploads, 1)
	upload := result.Uploads[0]

	// Verify CDN URLs are present (even if CDN not configured, they should equal public URLs)
	s.NotEmpty(upload.PublicURL)
	s.NotEmpty(upload.CDNURL)

	// Verify asset record has CDN URL
	assetRepo := database.NewRepository[asset.Asset](testEnv.Database)
	dbAsset, err := assetRepo.FindByID(ctx, upload.AssetID)
	s.NoError(err)
	s.NotEmpty(dbAsset.PublicURL)
	s.NotEmpty(dbAsset.CDNURL)

	// If thumbnail exists, verify its CDN URLs too
	if upload.Thumbnail != nil {
		s.NotEmpty(upload.Thumbnail.PublicURL)
		s.NotEmpty(upload.Thumbnail.CDNURL)

		thumbAsset, err := assetRepo.FindByID(ctx, upload.Thumbnail.AssetID)
		s.NoError(err)
		s.NotEmpty(thumbAsset.PublicURL)
		s.NotEmpty(thumbAsset.CDNURL)
	}
}

func TestAssetUploadSuite(t *testing.T) {
	suite.Run(t, new(AssetUploadSuite))
}
