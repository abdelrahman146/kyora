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

func TestAssetUploadSuite(t *testing.T) {
	suite.Run(t, new(AssetUploadSuite))
}
