package asset

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/request"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

// HttpHandler handles HTTP requests for asset management.
type HttpHandler struct {
	svc *Service
}

// NewHttpHandler creates a new HTTP handler for assets.
func NewHttpHandler(svc *Service) *HttpHandler {
	return &HttpHandler{svc: svc}
}

// GenerateUploadURLs godoc
// @Summary      Generate pre-signed upload URLs
// @Description  Generates pre-signed URLs for direct file uploads. Supports multiple files, S3 multipart uploads (with resumable chunking), and local simple uploads. Returns upload descriptors with method, URLs, headers, uploadId (S3), partUrls (S3 multipart), and assetId for each file.
// @Description
// @Description  **S3 Multipart Upload Flow (for large files):**
// @Description  1. Call this endpoint with file details (name, size, contentType)
// @Description  2. Response includes uploadId, partSize, totalParts, and partUrls array
// @Description  3. Client chunks file into parts (each partSize bytes)
// @Description  4. Client uploads each chunk to corresponding partUrl (PUT request)
// @Description  5. Client collects ETag from each upload response header
// @Description  6. Client calls CompleteMultipartUpload endpoint with parts array (partNumber + ETag)
// @Description  7. Asset is ready and can be referenced in product/business APIs
// @Description
// @Description  **Local Upload Flow (for development):**
// @Description  1. Call this endpoint with file details
// @Description  2. Response includes simple POST URL
// @Description  3. Client uploads entire file to URL with Content-Type header
// @Description  4. Asset is ready immediately and can be referenced
// @Description
// @Description  **Using Uploaded Assets:**
// @Description  - After upload completes, use the returned assetId and publicUrl in product photos or business logo
// @Description  - Pass AssetReference object: `{"url": "<publicUrl>", "assetId": "<assetId>", "metadata": {"altText": "...", "caption": "..."}}`
// @Description  - assetId is optional but enables automatic garbage collection
// @Description
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor (ID or slug)"
// @Param        body body GenerateUploadURLsRequest true "Upload request with files array"
// @Success      200 {object} GenerateUploadURLsResponse "Upload descriptors for each file"
// @Failure      400 {object} problem.Problem "Invalid request (bad file details, unsupported content type, too many files)"
// @Failure      401 {object} problem.Problem "Unauthorized (missing or invalid JWT)"
// @Failure      403 {object} problem.Problem "Forbidden (insufficient permissions)"
// @Failure      500 {object} problem.Problem "Internal server error (blob storage failure)"
// @Router       /v1/businesses/{businessDescriptor}/assets/uploads [post]
// @Security     BearerAuth
func (h *HttpHandler) GenerateUploadURLs(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := business.BusinessFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req GenerateUploadURLsRequest
	if err := request.ValidBody(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	res, err := h.svc.GenerateUploadURLs(c.Request.Context(), actor, biz, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusOK, res)
}

// CompleteMultipartUpload godoc
// @Summary      Complete a multipart upload
// @Description  Finalizes a multipart upload after all parts have been uploaded. Client must provide all part numbers with their ETags collected from upload responses. This endpoint assembles the parts on S3 and marks the asset as ready. Only applicable for S3 multipart uploads (not local uploads).
// @Description
// @Description  **Required Data:**
// @Description  - parts: array of objects with partNumber (int) and etag (string)
// @Description  - ETags are returned in the ETag response header from each part upload
// @Description  - Parts must be in order (1, 2, 3, ...)
// @Description
// @Description  **Example Request:**
// @Description  ```json
// @Description  {
// @Description    "parts": [
// @Description      {"partNumber": 1, "etag": "abc123..."},
// @Description      {"partNumber": 2, "etag": "def456..."},
// @Description      {"partNumber": 3, "etag": "ghi789..."}
// @Description    ]
// @Description  }
// @Description  ```
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor (ID or slug)"
// @Param        assetId path string true "Asset ID (returned from GenerateUploadURLs)"
// @Param        body body CompleteMultipartUploadRequest true "Parts with ETags"
// @Success      200 {object} map[string]string "Success message"
// @Failure      400 {object} problem.Problem "Invalid request (missing parts, invalid asset)"
// @Failure      401 {object} problem.Problem "Unauthorized"
// @Failure      404 {object} problem.Problem "Asset not found"
// @Failure      500 {object} problem.Problem "S3 completion failed"
// @Router       /v1/businesses/{businessDescriptor}/assets/uploads/{assetId}/complete [post]
// @Security     BearerAuth
func (h *HttpHandler) CompleteMultipartUpload(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := business.BusinessFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	assetID := c.Param("assetId")
	if assetID == "" {
		response.Error(c, problem.BadRequest("assetId is required"))
		return
	}

	var req CompleteMultipartUploadRequest
	if err := request.ValidBody(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	if err := h.svc.CompleteMultipartUpload(c.Request.Context(), actor, biz, assetID, &req); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusOK, map[string]string{
		"message": "upload completed successfully",
		"assetId": assetID,
	})
}

// GetPublicAsset godoc
// @Summary      Serve public asset
// @Description  Serves a public asset file. For local provider, streams the file from disk with cache headers. For S3, redirects to the public URL. No authentication required - all assets are public by design.
// @Description
// @Description  **Cache Headers:**
// @Description  - Cache-Control: public, max-age=3600 (1 hour)
// @Description  - ETag: MD5 hash of file content (for local provider)
// @Description  - Last-Modified: file modification time
// @Description
// @Description  **Response:**
// @Description  - For local: streams file content with appropriate Content-Type
// @Description  - For S3: 307 Temporary Redirect to S3 public URL
// @Tags         assets
// @Produce      application/octet-stream
// @Param        assetId path string true "Asset ID"
// @Success      200 {file} file "Asset file content (local provider)"
// @Success      307 {string} string "Redirect to S3 public URL (S3 provider)"
// @Failure      404 {object} problem.Problem "Asset not found"
// @Failure      500 {object} problem.Problem "Internal server error"
// @Router       /v1/public/assets/{assetId} [get]
func (h *HttpHandler) GetPublicAsset(c *gin.Context) {
	assetID := c.Param("assetId")
	if assetID == "" {
		response.Error(c, problem.BadRequest("assetId is required"))
		return
	}

	asset, err := h.svc.GetPublicAsset(c.Request.Context(), assetID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// For local files, serve from disk
	if asset.LocalFilePath != "" {
		h.serveLocalFile(c, asset)
		return
	}

	// For S3, redirect to public URL
	if asset.PublicURL != "" {
		c.Redirect(http.StatusTemporaryRedirect, asset.PublicURL)
		return
	}

	response.Error(c, problem.NotFound("asset not accessible"))
}

// serveLocalFile serves a local file with cache headers.
func (h *HttpHandler) serveLocalFile(c *gin.Context, asset *Asset) {
	file, err := os.Open(asset.LocalFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			response.Error(c, problem.NotFound("asset file not found"))
			return
		}
		response.Error(c, problem.InternalError())
		return
	}
	defer file.Close()

	// Get file info for Last-Modified and ETag
	fileInfo, err := file.Stat()
	if err != nil {
		response.Error(c, problem.InternalError())
		return
	}

	// Calculate ETag (MD5 hash of content)
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		response.Error(c, problem.InternalError())
		return
	}
	etag := fmt.Sprintf(`"%s"`, hex.EncodeToString(hash.Sum(nil)))

	// Reset file pointer
	if _, err := file.Seek(0, 0); err != nil {
		response.Error(c, problem.InternalError())
		return
	}

	// Set cache headers
	c.Header("Cache-Control", "public, max-age=3600") // 1 hour
	c.Header("ETag", etag)
	c.Header("Last-Modified", fileInfo.ModTime().UTC().Format(http.TimeFormat))

	// Check If-None-Match (ETag)
	if match := c.GetHeader("If-None-Match"); match == etag {
		c.Status(http.StatusNotModified)
		return
	}

	// Check If-Modified-Since
	if modifiedSince := c.GetHeader("If-Modified-Since"); modifiedSince != "" {
		if t, err := time.Parse(http.TimeFormat, modifiedSince); err == nil {
			if !fileInfo.ModTime().After(t) {
				c.Status(http.StatusNotModified)
				return
			}
		}
	}

	// Serve file
	c.Header("Content-Type", asset.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", asset.SizeBytes))
	c.Status(http.StatusOK)
	io.Copy(c.Writer, file)
}

// UploadLocalContent godoc
// @Summary      Internal: Upload content for local provider
// @Description  Internal endpoint for uploading file content when using local storage provider. Clients should not call this directly - use the URL returned from GenerateUploadURLs instead.
// @Tags         assets
// @Accept       application/octet-stream
// @Produce      json
// @Param        assetId path string true "Asset ID"
// @Param        body body string true "File content (raw bytes)"
// @Success      200 {object} map[string]string "Success message"
// @Failure      400 {object} problem.Problem "Invalid request"
// @Failure      404 {object} problem.Problem "Asset not found"
// @Failure      500 {object} problem.Problem "Failed to store content"
// @Router       /v1/assets/internal/upload/{assetId} [post]
func (h *HttpHandler) UploadLocalContent(c *gin.Context) {
	assetID := c.Param("assetId")
	if assetID == "" {
		response.Error(c, problem.BadRequest("assetId is required"))
		return
	}

	contentType := c.GetHeader("Content-Type")
	if contentType == "" {
		response.Error(c, problem.BadRequest("Content-Type header is required"))
		return
	}

	if err := h.svc.StoreLocalContent(c.Request.Context(), assetID, contentType, c.Request.Body); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusOK, map[string]string{
		"message": "upload completed successfully",
		"assetId": assetID,
	})
}
