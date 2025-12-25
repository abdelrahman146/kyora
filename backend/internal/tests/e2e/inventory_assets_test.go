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

type InventoryAssetsSuite struct {
	suite.Suite
	accountHelper   *AccountTestHelper
	inventoryHelper *InventoryTestHelper
}

func (s *InventoryAssetsSuite) SetupSuite() {
	s.accountHelper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
	s.inventoryHelper = NewInventoryTestHelper(testEnv.Database, e2eBaseURL)
}

func (s *InventoryAssetsSuite) resetDB() {
	s.NoError(testutils.TruncateTables(testEnv.Database, inventoryTables...))
}

func (s *InventoryAssetsSuite) SetupTest() {
	s.resetDB()
}

func (s *InventoryAssetsSuite) TearDownTest() {
	s.resetDB()
}

func (s *InventoryAssetsSuite) uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s_%s@example.com", prefix, strings.ToLower(id.Base62(8)))
}

func (s *InventoryAssetsSuite) TestCreateProduct_WithUploadedPhoto() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, s.uniqueEmail("admin"), "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	bizDescriptor := fmt.Sprintf("prod-photo-%s", strings.ToLower(id.Base62(6)))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, bizDescriptor)
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "cat")
	s.NoError(err)

	content := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0x00, 0x00, 0x00, 0x00}
	_, publicURL := UploadTestAsset(ctx, s.T(), s.inventoryHelper.Client, AssetUploadParams{
		BusinessDescriptor: bizDescriptor,
		Token:              token,
		CreatePath:         "assets/uploads/product-photo",
		CompletePurpose:    "product_photo",
		ContentType:        "image/png",
		Content:            content,
	})

	payload := map[string]interface{}{
		"name":        "Photo Product",
		"description": "Desc",
		"categoryId":  cat.ID,
		"photos":      []string{publicURL},
	}
	resp, err := s.inventoryHelper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/inventory/products", bizDescriptor), payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var created map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &created))
	s.Equal([]interface{}{publicURL}, created["photos"])

	getResp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/%s/inventory/products/%s", bizDescriptor, created["id"]), nil, token)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusOK, getResp.StatusCode)

	var fetched map[string]interface{}
	s.NoError(testutils.DecodeJSON(getResp, &fetched))
	s.Equal([]interface{}{publicURL}, fetched["photos"])
}

func (s *InventoryAssetsSuite) TestCreateAndUpdateVariant_WithUploadedPhotos() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, s.uniqueEmail("admin"), "ValidPassword123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	bizDescriptor := fmt.Sprintf("variant-photo-%s", strings.ToLower(id.Base62(6)))
	biz, err := s.inventoryHelper.CreateTestBusiness(ctx, ws.ID, bizDescriptor)
	s.NoError(err)
	cat, err := s.inventoryHelper.CreateTestCategory(ctx, biz.ID, "Cat", "cat")
	s.NoError(err)

	prodPayload := map[string]interface{}{"name": "Base Product", "description": "", "categoryId": cat.ID}
	prodResp, err := s.inventoryHelper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/inventory/products", bizDescriptor), prodPayload, token)
	s.NoError(err)
	defer prodResp.Body.Close()
	s.Equal(http.StatusCreated, prodResp.StatusCode)
	var prod map[string]interface{}
	s.NoError(testutils.DecodeJSON(prodResp, &prod))
	prodID := prod["id"].(string)

	initialPhoto := []byte{'p', 'h', 'o', 't', 'o', '1'}
	_, initialURL := UploadTestAsset(ctx, s.T(), s.inventoryHelper.Client, AssetUploadParams{
		BusinessDescriptor: bizDescriptor,
		Token:              token,
		CreatePath:         "assets/uploads/variant-photo",
		CompletePurpose:    "variant_photo",
		ContentType:        "image/png",
		Content:            initialPhoto,
	})

	variantPayload := map[string]interface{}{
		"productId":          prodID,
		"code":               "RED",
		"sku":                "",
		"photos":             []string{initialURL},
		"costPrice":          "10",
		"salePrice":          "15",
		"stockQuantity":      2,
		"stockQuantityAlert": 1,
	}
	varResp, err := s.inventoryHelper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/%s/inventory/variants", bizDescriptor), variantPayload, token)
	s.NoError(err)
	defer varResp.Body.Close()
	s.Equal(http.StatusCreated, varResp.StatusCode)

	var createdVariant map[string]interface{}
	s.NoError(testutils.DecodeJSON(varResp, &createdVariant))
	variantID := createdVariant["id"].(string)
	s.Equal([]interface{}{initialURL}, createdVariant["photos"])

	updatedPhoto := []byte{'p', 'h', 'o', 't', 'o', '2'}
	_, updatedURL := UploadTestAsset(ctx, s.T(), s.inventoryHelper.Client, AssetUploadParams{
		BusinessDescriptor: bizDescriptor,
		Token:              token,
		CreatePath:         "assets/uploads/variant-photo",
		CompletePurpose:    "variant_photo",
		ContentType:        "image/png",
		Content:            updatedPhoto,
	})

	patchPayload := map[string]interface{}{"photos": []string{updatedURL}}
	patchResp, err := s.inventoryHelper.Client.AuthenticatedRequest("PATCH", fmt.Sprintf("/v1/businesses/%s/inventory/variants/%s", bizDescriptor, variantID), patchPayload, token)
	s.NoError(err)
	defer patchResp.Body.Close()
	s.Equal(http.StatusOK, patchResp.StatusCode)

	var patched map[string]interface{}
	s.NoError(testutils.DecodeJSON(patchResp, &patched))
	s.Equal([]interface{}{updatedURL}, patched["photos"])

	getResp, err := s.inventoryHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/%s/inventory/variants/%s", bizDescriptor, variantID), nil, token)
	s.NoError(err)
	defer getResp.Body.Close()
	s.Equal(http.StatusOK, getResp.StatusCode)

	var fetched map[string]interface{}
	s.NoError(testutils.DecodeJSON(getResp, &fetched))
	s.Equal([]interface{}{updatedURL}, fetched["photos"])
}

func TestInventoryAssetsSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(InventoryAssetsSuite))
}
