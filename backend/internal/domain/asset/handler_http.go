package asset

import (
	"net/http"
	"os"
	"strings"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/request"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpHandler struct {
	svc *Service
}

func NewHttpHandler(svc *Service) *HttpHandler {
	return &HttpHandler{svc: svc}
}

// CreateLogoUpload godoc
// @Summary      Initiate business logo upload
// @Description  Returns a direct upload URL for a business logo (image). Requires manage business permission.
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        body body CreateUploadRequest true "Upload request"
// @Success      200 {object} CreateUploadResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/assets/uploads/logo [post]
func (h *HttpHandler) CreateLogoUpload(c *gin.Context) {
	h.createUpload(c, PurposeBusinessLogo)
}

// CreateProductPhotoUpload godoc
// @Summary      Initiate product photo upload
// @Description  Returns a direct upload URL for a product photo (image). Requires manage inventory permission.
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        body body CreateUploadRequest true "Upload request"
// @Success      200 {object} CreateUploadResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/assets/uploads/product-photo [post]
func (h *HttpHandler) CreateProductPhotoUpload(c *gin.Context) {
	h.createUpload(c, PurposeProductPhoto)
}

// CreateVariantPhotoUpload godoc
// @Summary      Initiate variant photo upload
// @Description  Returns a direct upload URL for a variant photo (image). Requires manage inventory permission.
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        body body CreateUploadRequest true "Upload request"
// @Success      200 {object} CreateUploadResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/assets/uploads/variant-photo [post]
func (h *HttpHandler) CreateVariantPhotoUpload(c *gin.Context) {
	h.createUpload(c, PurposeVariantPhoto)
}

func (h *HttpHandler) createUpload(c *gin.Context, purpose Purpose) {
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
	var req CreateUploadRequest
	if err := request.ValidBody(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	res, err := h.svc.CreateUpload(c.Request.Context(), actor, biz, purpose, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, res)
}

// PutLogoContent godoc
// @Summary      Upload logo content (local provider only)
// @Description  Uploads file bytes to Kyora when storage.provider=local.
// @Tags         assets
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        assetId path string true "Asset ID"
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/assets/uploads/{assetId}/content/business_logo [put]
func (h *HttpHandler) PutLogoContent(c *gin.Context) {
	h.putContent(c, PurposeBusinessLogo)
}

// PutProductPhotoContent godoc
// @Summary      Upload product photo content (local provider only)
// @Tags         assets
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        assetId path string true "Asset ID"
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/assets/uploads/{assetId}/content/product_photo [put]
func (h *HttpHandler) PutProductPhotoContent(c *gin.Context) {
	h.putContent(c, PurposeProductPhoto)
}

// PutVariantPhotoContent godoc
// @Summary      Upload variant photo content (local provider only)
// @Tags         assets
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        assetId path string true "Asset ID"
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/assets/uploads/{assetId}/content/variant_photo [put]
func (h *HttpHandler) PutVariantPhotoContent(c *gin.Context) {
	h.putContent(c, PurposeVariantPhoto)
}

func (h *HttpHandler) putContent(c *gin.Context, expectedPurpose Purpose) {
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
	assetID := strings.TrimSpace(c.Param("assetId"))
	contentType := strings.TrimSpace(c.GetHeader("Content-Type"))
	if _, err := h.svc.StoreLocalContent(c.Request.Context(), actor, biz, assetID, expectedPurpose, contentType, c.Request.Body); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

// CompleteLogoUpload godoc
// @Summary      Complete business logo upload
// @Tags         assets
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        assetId path string true "Asset ID"
// @Success      200 {object} CompleteUploadResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/assets/uploads/{assetId}/complete/business_logo [post]
func (h *HttpHandler) CompleteLogoUpload(c *gin.Context) {
	h.complete(c, PurposeBusinessLogo)
}

// CompleteProductPhotoUpload godoc
// @Summary      Complete product photo upload
// @Tags         assets
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        assetId path string true "Asset ID"
// @Success      200 {object} CompleteUploadResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/assets/uploads/{assetId}/complete/product_photo [post]
func (h *HttpHandler) CompleteProductPhotoUpload(c *gin.Context) {
	h.complete(c, PurposeProductPhoto)
}

// CompleteVariantPhotoUpload godoc
// @Summary      Complete variant photo upload
// @Tags         assets
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        assetId path string true "Asset ID"
// @Success      200 {object} CompleteUploadResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/assets/uploads/{assetId}/complete/variant_photo [post]
func (h *HttpHandler) CompleteVariantPhotoUpload(c *gin.Context) {
	h.complete(c, PurposeVariantPhoto)
}

func (h *HttpHandler) complete(c *gin.Context, expectedPurpose Purpose) {
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
	assetID := strings.TrimSpace(c.Param("assetId"))
	res, err := h.svc.CompleteUpload(c.Request.Context(), actor, biz, assetID, expectedPurpose)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, res)
}

// GetPublicAsset godoc
// @Summary      Get public asset
// @Description  Serves public assets. For local provider, streams bytes. For S3, redirects to the public URL.
// @Tags         assets
// @Produce      octet-stream
// @Param        assetId path string true "Asset ID"
// @Success      302
// @Success      200
// @Failure      404 {object} problem.Problem
// @Router       /v1/public/assets/{assetId} [get]
func (h *HttpHandler) GetPublicAsset(c *gin.Context) {
	assetID := strings.TrimSpace(c.Param("assetId"))
	a, err := h.svc.GetPublicAsset(c.Request.Context(), assetID)
	if err != nil {
		response.Error(c, err)
		return
	}

	provider := strings.ToLower(strings.TrimSpace(viper.GetString(config.StorageProvider)))
	if provider != "local" {
		if a.PublicURL != "" {
			c.Redirect(http.StatusFound, a.PublicURL)
			return
		}
		c.Status(http.StatusNotFound)
		return
	}

	if strings.TrimSpace(a.LocalFilePath) == "" {
		c.Status(http.StatusNotFound)
		return
	}
	if _, err := os.Stat(a.LocalFilePath); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.Header("Content-Type", a.ContentType)
	c.File(a.LocalFilePath)
}
