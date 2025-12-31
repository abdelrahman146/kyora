package asset

import (
	"encoding/json"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/types/asset"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"gorm.io/gorm"
)

const (
	// AssetTable is the database table for uploaded file assets.
	// NOTE: this is intentionally not "assets" to avoid colliding with the
	// accounting domain's assets table.
	AssetTable  = "uploaded_assets"
	AssetStruct = "Asset"
	AssetPrefix = "ast"
)

// Asset represents an uploaded file in blob storage.
// This model tracks both S3 multipart uploads and simple local file uploads.
type Asset struct {
	gorm.Model
	ID              string             `gorm:"column:id;primaryKey;type:text" json:"id"`
	WorkspaceID     string             `gorm:"column:workspace_id;type:text;not null;index" json:"workspaceId"`
	Workspace       *account.Workspace `gorm:"foreignKey:WorkspaceID;references:ID" json:"workspace,omitempty"`
	BusinessID      string             `gorm:"column:business_id;type:text;not null;index" json:"businessId"`
	Business        *business.Business `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	CreatedByUserID string             `gorm:"column:created_by_user_id;type:text;not null;index" json:"createdByUserId"`
	ObjectKey       string             `gorm:"column:object_key;type:text;not null;uniqueIndex" json:"objectKey"`
	PublicURL       string             `gorm:"column:public_url;type:text" json:"publicUrl"`
	CDNURL          string             `gorm:"column:cdn_url;type:text" json:"cdnUrl"`
	ContentType     string             `gorm:"column:content_type;type:text;not null" json:"contentType"`
	FileCategory    string             `gorm:"column:file_category;type:text;not null" json:"fileCategory"`
	SizeBytes       int64              `gorm:"column:size_bytes;type:bigint;not null" json:"sizeBytes"`
	LocalFilePath   string             `gorm:"column:local_file_path;type:text" json:"-"`

	// Thumbnail support (tight coupling)
	ThumbnailAssetID   *string `gorm:"column:thumbnail_asset_id;type:text;index" json:"thumbnailAssetId,omitempty"`
	ThumbnailObjectKey string  `gorm:"column:thumbnail_object_key;type:text" json:"thumbnailObjectKey,omitempty"`
	ThumbnailPublicURL string  `gorm:"column:thumbnail_public_url;type:text" json:"thumbnailPublicUrl,omitempty"`
	ThumbnailCDNURL    string  `gorm:"column:thumbnail_cdn_url;type:text" json:"thumbnailCdnUrl,omitempty"`

	// Multipart upload tracking (S3 only)
	UploadID       string          `gorm:"column:upload_id;type:text" json:"uploadId,omitempty"`
	IsMultipart    bool            `gorm:"column:is_multipart;type:boolean;default:false" json:"isMultipart"`
	TotalParts     int             `gorm:"column:total_parts;type:integer;default:0" json:"totalParts,omitempty"`
	CompletedParts json.RawMessage `gorm:"column:completed_parts;type:jsonb;default:'[]'" json:"completedParts,omitempty"`
}

// AssetSchema defines the field mappings for ordering and searching.
var AssetSchema = struct {
	ID                 schema.Field
	WorkspaceID        schema.Field
	BusinessID         schema.Field
	CreatedByUserID    schema.Field
	ObjectKey          schema.Field
	PublicURL          schema.Field
	CDNURL             schema.Field
	ContentType        schema.Field
	FileCategory       schema.Field
	SizeBytes          schema.Field
	ThumbnailAssetID   schema.Field
	ThumbnailObjectKey schema.Field
	ThumbnailPublicURL schema.Field
	ThumbnailCDNURL    schema.Field
	UploadID           schema.Field
	IsMultipart        schema.Field
	TotalParts         schema.Field
	CreatedAt          schema.Field
	UpdatedAt          schema.Field
	DeletedAt          schema.Field
}{
	ID:                 schema.NewField("id", "id"),
	WorkspaceID:        schema.NewField("workspace_id", "workspaceId"),
	BusinessID:         schema.NewField("business_id", "businessId"),
	CreatedByUserID:    schema.NewField("created_by_user_id", "createdByUserId"),
	ObjectKey:          schema.NewField("object_key", "objectKey"),
	PublicURL:          schema.NewField("public_url", "publicUrl"),
	CDNURL:             schema.NewField("cdn_url", "cdnUrl"),
	ContentType:        schema.NewField("content_type", "contentType"),
	FileCategory:       schema.NewField("file_category", "fileCategory"),
	SizeBytes:          schema.NewField("size_bytes", "sizeBytes"),
	ThumbnailAssetID:   schema.NewField("thumbnail_asset_id", "thumbnailAssetId"),
	ThumbnailObjectKey: schema.NewField("thumbnail_object_key", "thumbnailObjectKey"),
	ThumbnailPublicURL: schema.NewField("thumbnail_public_url", "thumbnailPublicUrl"),
	ThumbnailCDNURL:    schema.NewField("thumbnail_cdn_url", "thumbnailCdnUrl"),
	UploadID:           schema.NewField("upload_id", "uploadId"),
	IsMultipart:        schema.NewField("is_multipart", "isMultipart"),
	TotalParts:         schema.NewField("total_parts", "totalParts"),
	CreatedAt:          schema.NewField("created_at", "createdAt"),
	UpdatedAt:          schema.NewField("updated_at", "updatedAt"),
	DeletedAt:          schema.NewField("deleted_at", "deletedAt"),
}

func (m *Asset) TableName() string { return AssetTable }

func (m *Asset) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(AssetPrefix)
	}
	return nil
}

// AssetMetadata and AssetReference are defined in internal/platform/types/asset
// to avoid circular dependencies with other domains.
type AssetMetadata = asset.AssetMetadata
type AssetReference = asset.AssetReference

// PartInfo represents a completed part of a multipart upload.
type PartInfo struct {
	PartNumber int    `json:"partNumber"`
	ETag       string `json:"etag"`
}

// GenerateUploadURLsRequest is the request to generate pre-signed URLs for file uploads.
type GenerateUploadURLsRequest struct {
	Files []FileUploadRequest `json:"files" binding:"required,min=1,max=50,dive"`
}

// FileUploadRequest represents a single file to be uploaded.
type FileUploadRequest struct {
	FileName    string `json:"fileName" binding:"required,max=255"`
	ContentType string `json:"contentType" binding:"required,max=128"`
	SizeBytes   int64  `json:"sizeBytes" binding:"required,gt=0"`
}

// GenerateUploadURLsResponse is the response containing upload descriptors for each file.
type GenerateUploadURLsResponse struct {
	Uploads []UploadDescriptor `json:"uploads"`
}

// UploadDescriptor contains all information needed to upload a file.
type UploadDescriptor struct {
	AssetID     string            `json:"assetId"`
	FileName    string            `json:"fileName"`
	Method      string            `json:"method"`
	URL         string            `json:"url,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	UploadID    string            `json:"uploadId,omitempty"`
	PartSize    int64             `json:"partSize,omitempty"`
	TotalParts  int               `json:"totalParts,omitempty"`
	PartURLs    []PartURLInfo     `json:"partUrls,omitempty"`
	PublicURL   string            `json:"publicUrl"`
	CDNURL      string            `json:"cdnUrl"`
	ContentType string            `json:"contentType"`
	SizeBytes   int64             `json:"sizeBytes"`

	// Thumbnail support
	IsThumbnail bool                 `json:"isThumbnail,omitempty"`
	Thumbnail   *ThumbnailDescriptor `json:"thumbnail,omitempty"`
}

// ThumbnailDescriptor contains upload info for the thumbnail.
type ThumbnailDescriptor struct {
	AssetID     string            `json:"assetId"`
	FileName    string            `json:"fileName"`
	Method      string            `json:"method"`
	URL         string            `json:"url,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	UploadID    string            `json:"uploadId,omitempty"`
	PartSize    int64             `json:"partSize,omitempty"`
	TotalParts  int               `json:"totalParts,omitempty"`
	PartURLs    []PartURLInfo     `json:"partUrls,omitempty"`
	PublicURL   string            `json:"publicUrl"`
	CDNURL      string            `json:"cdnUrl"`
	ContentType string            `json:"contentType"`
}

// PartURLInfo contains the pre-signed URL for a specific part.
type PartURLInfo struct {
	PartNumber int    `json:"partNumber"`
	URL        string `json:"url"`
}

// CompleteMultipartUploadRequest is the request to complete a multipart upload.
type CompleteMultipartUploadRequest struct {
	Parts []PartInfo `json:"parts" binding:"required,min=1,dive"`
}
