package asset

import (
	"context"
	"fmt"
	"io"
	"math"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/blob"
	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/spf13/viper"
)

// Service provides business logic for asset management.
type Service struct {
	storage *Storage
	atomic  *database.AtomicProcess
	blob    blob.Provider

	localDir            string
	multipartPartSizeMB int
	maxUploadBytes      int64
}

// NewService creates a new asset service instance.
func NewService(storage *Storage, atomic *database.AtomicProcess, provider blob.Provider) *Service {
	localDir := viper.GetString(config.StorageLocalPath)
	if localDir == "" {
		localDir = "./tmp/assets"
	}

	multipartPartSizeMB := viper.GetInt(config.StorageMultipartPartSize)
	if multipartPartSizeMB <= 0 {
		multipartPartSizeMB = 10
	}

	maxUploadBytes := viper.GetInt64(config.UploadsMaxBytes)
	if maxUploadBytes <= 0 {
		maxUploadBytes = 5 * 1024 * 1024 // 5MB default
	}

	return &Service{
		storage:             storage,
		atomic:              atomic,
		blob:                provider,
		localDir:            localDir,
		multipartPartSizeMB: multipartPartSizeMB,
		maxUploadBytes:      maxUploadBytes,
	}
}

// GenerateUploadURLs generates pre-signed URLs for direct file uploads.
// Supports both S3 multipart uploads and local simple uploads.
func (s *Service) GenerateUploadURLs(ctx context.Context, actor *account.User, biz *business.Business, req *GenerateUploadURLsRequest) (*GenerateUploadURLsResponse, error) {
	if actor == nil {
		return nil, problem.Unauthorized("unauthorized")
	}
	if biz == nil {
		return nil, problem.BadRequest("business is required")
	}
	if len(req.Files) == 0 {
		return nil, problem.BadRequest("at least one file is required").With("field", "files")
	}
	if len(req.Files) > 50 {
		return nil, problem.BadRequest("maximum 50 files per request").With("field", "files")
	}

	// Validate content types (images only for now)
	for i, file := range req.Files {
		if !isImageContentType(file.ContentType) {
			return nil, problem.BadRequest("only image files are supported").
				With("field", fmt.Sprintf("files[%d].contentType", i)).
				With("contentType", file.ContentType)
		}
	}

	uploads := make([]UploadDescriptor, 0, len(req.Files))

	for _, file := range req.Files {
		descriptor, err := s.generateSingleUpload(ctx, actor, biz, file)
		if err != nil {
			return nil, err
		}
		uploads = append(uploads, *descriptor)
	}

	return &GenerateUploadURLsResponse{Uploads: uploads}, nil
}

// generateSingleUpload generates upload instructions for a single file.
func (s *Service) generateSingleUpload(ctx context.Context, actor *account.User, biz *business.Business, file FileUploadRequest) (*UploadDescriptor, error) {
	// Create asset record immediately
	asset := &Asset{
		WorkspaceID:     biz.WorkspaceID,
		BusinessID:      biz.ID,
		CreatedByUserID: actor.ID,
		ContentType:     file.ContentType,
		SizeBytes:       file.SizeBytes,
		ObjectKey:       s.buildObjectKey(biz.ID, "", sanitizeFilename(file.FileName)),
	}

	// Generate assetId before creating upload
	if err := asset.BeforeCreate(nil); err != nil {
		return nil, err
	}

	// Update object key with actual assetId
	asset.ObjectKey = s.buildObjectKey(biz.ID, asset.ID, sanitizeFilename(file.FileName))

	provider := viper.GetString(config.StorageProvider)
	isLocal := provider == "local" || provider == ""

	if isLocal {
		// Local simple upload
		return s.generateLocalUpload(ctx, asset, file)
	}

	// S3 multipart upload
	return s.generateS3MultipartUpload(ctx, asset, file)
}

// generateLocalUpload creates a simple local upload endpoint.
func (s *Service) generateLocalUpload(ctx context.Context, asset *Asset, file FileUploadRequest) (*UploadDescriptor, error) {
	// For local, we store files in configured directory
	localPath := filepath.Join(s.localDir, asset.ID)
	asset.LocalFilePath = localPath
	asset.PublicURL = s.buildPublicURL(asset)

	if err := s.storage.Create(ctx, asset); err != nil {
		return nil, err
	}

	baseURL := viper.GetString(config.HTTPBaseURL)
	uploadURL := fmt.Sprintf("%s/v1/assets/internal/upload/%s", baseURL, asset.ID)

	return &UploadDescriptor{
		AssetID:     asset.ID,
		FileName:    file.FileName,
		Method:      "POST",
		URL:         uploadURL,
		Headers:     map[string]string{"Content-Type": file.ContentType},
		PublicURL:   asset.PublicURL,
		ContentType: file.ContentType,
		SizeBytes:   file.SizeBytes,
	}, nil
}

// generateS3MultipartUpload creates a multipart upload with presigned part URLs.
func (s *Service) generateS3MultipartUpload(ctx context.Context, asset *Asset, file FileUploadRequest) (*UploadDescriptor, error) {
	if s.blob == nil {
		return nil, problem.InternalError()
	}

	// Calculate parts
	partSizeBytes := int64(s.multipartPartSizeMB) * 1024 * 1024
	totalParts := int(math.Ceil(float64(file.SizeBytes) / float64(partSizeBytes)))

	if totalParts > 10000 {
		return nil, problem.BadRequest("file too large").
			With("maxParts", 10000).
			With("partSize", partSizeBytes).
			With("maxFileSize", 10000*partSizeBytes)
	}

	// Initiate multipart upload
	uploadID, err := s.blob.CreateMultipartUpload(ctx, asset.ObjectKey, file.ContentType)
	if err != nil {
		return nil, problem.InternalError().WithError(err)
	}

	// Generate presigned URLs for each part
	partURLs := make([]PartURLInfo, totalParts)
	expiresIn := 24 * time.Hour // Long expiry for large uploads

	for i := 1; i <= totalParts; i++ {
		partURL, err := s.blob.PresignMultipartPart(ctx, asset.ObjectKey, uploadID, i, expiresIn)
		if err != nil {
			// Abort multipart upload on failure
			_ = s.blob.AbortMultipartUpload(ctx, asset.ObjectKey, uploadID)
			return nil, problem.InternalError().WithError(err)
		}
		partURLs[i-1] = PartURLInfo{
			PartNumber: i,
			URL:        partURL,
		}
	}

	// Get public URL
	publicURL, ok := s.blob.PublicURL(asset.ObjectKey)
	if !ok {
		publicURL = ""
	}
	asset.PublicURL = publicURL
	asset.UploadID = uploadID
	asset.IsMultipart = true
	asset.TotalParts = totalParts

	if err := s.storage.Create(ctx, asset); err != nil {
		// Abort multipart upload on failure
		_ = s.blob.AbortMultipartUpload(ctx, asset.ObjectKey, uploadID)
		return nil, err
	}

	return &UploadDescriptor{
		AssetID:     asset.ID,
		FileName:    file.FileName,
		Method:      "PUT",
		UploadID:    uploadID,
		PartSize:    partSizeBytes,
		TotalParts:  totalParts,
		PartURLs:    partURLs,
		PublicURL:   publicURL,
		ContentType: file.ContentType,
		SizeBytes:   file.SizeBytes,
	}, nil
}

// CompleteMultipartUpload completes a multipart upload after all parts are uploaded.
func (s *Service) CompleteMultipartUpload(ctx context.Context, actor *account.User, biz *business.Business, assetID string, req *CompleteMultipartUploadRequest) error {
	if actor == nil {
		return problem.Unauthorized("unauthorized")
	}
	if biz == nil {
		return problem.BadRequest("business is required")
	}

	asset, err := s.storage.GetByID(ctx, biz.ID, assetID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			return problem.NotFound("asset not found").With("assetId", assetID)
		}
		return err
	}

	if !asset.IsMultipart {
		return problem.BadRequest("asset is not a multipart upload").With("assetId", assetID)
	}

	if asset.UploadID == "" {
		return problem.BadRequest("asset has no upload ID").With("assetId", assetID)
	}

	if len(req.Parts) == 0 {
		return problem.BadRequest("at least one part is required").With("field", "parts")
	}

	// Convert to blob.CompletedPart
	completedParts := make([]blob.CompletedPart, len(req.Parts))
	for i, part := range req.Parts {
		completedParts[i] = blob.CompletedPart{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
		}
	}

	// Complete on S3
	if err := s.blob.CompleteMultipartUpload(ctx, asset.ObjectKey, asset.UploadID, completedParts); err != nil {
		return problem.InternalError().WithError(err)
	}

	// Mark upload as complete in DB
	if err := s.storage.MarkUploadComplete(ctx, assetID); err != nil {
		return err
	}

	return nil
}

// StoreLocalContent stores uploaded content for local provider.
// This is an internal handler called by the local upload endpoint.
func (s *Service) StoreLocalContent(ctx context.Context, assetID string, contentType string, r io.Reader) error {
	asset, err := s.storage.FindByID(ctx, assetID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			return problem.NotFound("asset not found").With("assetId", assetID)
		}
		return err
	}

	if asset.LocalFilePath == "" {
		return problem.BadRequest("asset is not configured for local upload").With("assetId", assetID)
	}

	// Verify content type matches
	if asset.ContentType != contentType {
		return problem.BadRequest("content type mismatch").
			With("expected", asset.ContentType).
			With("actual", contentType)
	}

	// Ensure directory exists
	dir := filepath.Dir(asset.LocalFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return problem.InternalError().WithError(err)
	}

	// Write file
	f, err := os.Create(asset.LocalFilePath)
	if err != nil {
		return problem.InternalError().WithError(err)
	}
	defer f.Close()

	written, err := io.Copy(f, r)
	if err != nil {
		return problem.InternalError().WithError(err)
	}

	// Verify size matches
	if written != asset.SizeBytes {
		// Clean up partial file
		_ = os.Remove(asset.LocalFilePath)
		return problem.BadRequest("size mismatch").
			With("expected", asset.SizeBytes).
			With("actual", written)
	}

	return nil
}

// GetPublicAsset returns an asset for public serving.
// No authentication required - assets are public by design.
func (s *Service) GetPublicAsset(ctx context.Context, assetID string) (*Asset, error) {
	asset, err := s.storage.FindByID(ctx, assetID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			return nil, problem.NotFound("asset not found").With("assetId", assetID)
		}
		return nil, err
	}

	return asset, nil
}

// buildObjectKey creates the blob storage key for an asset.
func (s *Service) buildObjectKey(businessID, assetID, fileName string) string {
	if assetID == "" {
		assetID = "temp"
	}
	sanitized := sanitizeFilename(fileName)
	return fmt.Sprintf("assets/%s/%s/%s", businessID, assetID, sanitized)
}

// buildPublicURL constructs the public URL for an asset.
func (s *Service) buildPublicURL(asset *Asset) string {
	baseURL := viper.GetString(config.HTTPBaseURL)
	return fmt.Sprintf("%s/v1/public/assets/%s", baseURL, asset.ID)
}

// isImageContentType checks if the content type is a supported image format.
func isImageContentType(ct string) bool {
	ct = strings.ToLower(strings.TrimSpace(ct))
	switch ct {
	case "image/jpeg", "image/jpg", "image/png", "image/webp", "image/gif":
		return true
	default:
		return false
	}
}

// sanitizeFilename removes dangerous characters from filenames.
func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "file"
	}
	// Remove path traversal and exotic chars
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "..", "_")
	name = strings.ReplaceAll(name, " ", "_")
	if len(name) > 80 {
		ext := filepath.Ext(name)
		base := name[:len(name)-len(ext)]
		if len(base) > 75 {
			base = base[:75]
		}
		name = base + ext
	}
	return name
}

// getMimeType attempts to detect MIME type from file extension.
func getMimeType(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return "application/octet-stream"
	}
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream"
	}
	return mimeType
}
