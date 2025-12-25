package e2e_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/require"
)

// AssetUploadParams describes a complete local upload flow for E2E tests.
type AssetUploadParams struct {
	BusinessDescriptor string
	Token              string
	CreatePath         string // e.g. assets/uploads/product-photo
	CompletePurpose    string // e.g. product_photo
	ContentType        string
	Content            []byte
	IdempotencyKey     string
}

// UploadTestAsset performs create -> PUT content -> complete and returns asset ID and public URL.
func UploadTestAsset(ctx context.Context, t testing.TB, client *testutils.HTTPClient, params AssetUploadParams) (string, string) {
	t.Helper()

	size := int64(len(params.Content))
	idem := params.IdempotencyKey
	if idem == "" {
		s := id.Base62(10)
		i := id.Base62(10)
		idem = fmt.Sprintf("idem_%s_%s", s, i)
	}
	ct := params.ContentType
	if ct == "" {
		ct = "image/png"
	}

	createReq := map[string]interface{}{
		"idempotencyKey": idem,
		"fileName":       "upload.png",
		"contentType":    ct,
		"sizeBytes":      size,
	}

	createPath := fmt.Sprintf("/v1/businesses/%s/%s", params.BusinessDescriptor, params.CreatePath)
	createResp, err := client.AuthenticatedRequest("POST", createPath, createReq, params.Token)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, createResp.StatusCode)

	var created map[string]interface{}
	require.NoError(t, testutils.DecodeJSON(createResp, &created))
	createResp.Body.Close()

	upload := created["upload"].(map[string]interface{})
	uploadURL := upload["url"].(string)
	headers := map[string]string{"Content-Type": ct}
	if rawHeaders, ok := upload["headers"].(map[string]interface{}); ok {
		for k, v := range rawHeaders {
			if vs, ok := v.(string); ok {
				headers[k] = vs
			}
		}
	}

	putResp, err := client.AuthenticatedRequestRaw("PUT", uploadURL, params.Content, headers, params.Token)
	require.NoError(t, err)
	putResp.Body.Close()
	require.Equal(t, http.StatusNoContent, putResp.StatusCode)

	completePath := fmt.Sprintf("/v1/businesses/%s/assets/uploads/%s/complete/%s", params.BusinessDescriptor, created["assetId"], params.CompletePurpose)
	completeResp, err := client.AuthenticatedRequest("POST", completePath, nil, params.Token)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, completeResp.StatusCode)

	var completed map[string]interface{}
	require.NoError(t, testutils.DecodeJSON(completeResp, &completed))
	completeResp.Body.Close()

	assetID := completed["assetId"].(string)
	publicURL := completed["publicUrl"].(string)

	getResp, err := client.Get(fmt.Sprintf("/v1/public/assets/%s", assetID))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, getResp.StatusCode)
	body, err := io.ReadAll(getResp.Body)
	require.NoError(t, err)
	getResp.Body.Close()
	require.Equal(t, params.Content, body)

	return assetID, publicURL
}
