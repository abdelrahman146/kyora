package asset

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
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

type Purpose string

const (
	PurposeBusinessLogo Purpose = "business_logo"
	PurposeProductPhoto Purpose = "product_photo"
	PurposeVariantPhoto Purpose = "variant_photo"
)

type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusReady   Status = "ready"
)

type Asset struct {
	gorm.Model
	ID              string             `gorm:"column:id;primaryKey;type:text" json:"id"`
	WorkspaceID     string             `gorm:"column:workspace_id;type:text;not null;index" json:"workspaceId"`
	Workspace       *account.Workspace `gorm:"foreignKey:WorkspaceID;references:ID" json:"workspace,omitempty"`
	BusinessID      string             `gorm:"column:business_id;type:text;not null;index;uniqueIndex:uniq_business_idem" json:"businessId"`
	Business        *business.Business `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	CreatedByUserID string             `gorm:"column:created_by_user_id;type:text;not null;index" json:"createdByUserId"`
	Purpose         Purpose            `gorm:"column:purpose;type:text;not null;index" json:"purpose"`
	Visibility      Visibility         `gorm:"column:visibility;type:text;not null;default:'public';index" json:"visibility"`
	Status          Status             `gorm:"column:status;type:text;not null;default:'pending';index" json:"status"`
	ObjectKey       string             `gorm:"column:object_key;type:text;not null;uniqueIndex" json:"objectKey"`
	PublicURL       string             `gorm:"column:public_url;type:text" json:"publicUrl"`
	ContentType     string             `gorm:"column:content_type;type:text;not null" json:"contentType"`
	SizeBytes       int64              `gorm:"column:size_bytes;type:bigint;not null" json:"sizeBytes"`
	IdempotencyKey  string             `gorm:"column:idempotency_key;type:text;uniqueIndex:uniq_business_idem" json:"idempotencyKey"`
	RequestHash     string             `gorm:"column:request_hash;type:text;not null" json:"requestHash"`
	UploadExpiresAt *time.Time         `gorm:"column:upload_expires_at;type:timestamp with time zone" json:"uploadExpiresAt,omitempty"`
	CompletedAt     *time.Time         `gorm:"column:completed_at;type:timestamp with time zone" json:"completedAt,omitempty"`
	LocalFilePath   string             `gorm:"column:local_file_path;type:text" json:"-"`
}

var AssetSchema = struct {
	ID              schema.Field
	WorkspaceID     schema.Field
	BusinessID      schema.Field
	CreatedByUserID schema.Field
	Purpose         schema.Field
	Visibility      schema.Field
	Status          schema.Field
	ObjectKey       schema.Field
	PublicURL       schema.Field
	ContentType     schema.Field
	SizeBytes       schema.Field
	IdempotencyKey  schema.Field
	RequestHash     schema.Field
	UploadExpiresAt schema.Field
	CompletedAt     schema.Field
	CreatedAt       schema.Field
	UpdatedAt       schema.Field
	DeletedAt       schema.Field
}{
	ID:              schema.NewField("id", "id"),
	WorkspaceID:     schema.NewField("workspace_id", "workspaceId"),
	BusinessID:      schema.NewField("business_id", "businessId"),
	CreatedByUserID: schema.NewField("created_by_user_id", "createdByUserId"),
	Purpose:         schema.NewField("purpose", "purpose"),
	Visibility:      schema.NewField("visibility", "visibility"),
	Status:          schema.NewField("status", "status"),
	ObjectKey:       schema.NewField("object_key", "objectKey"),
	PublicURL:       schema.NewField("public_url", "publicUrl"),
	ContentType:     schema.NewField("content_type", "contentType"),
	SizeBytes:       schema.NewField("size_bytes", "sizeBytes"),
	IdempotencyKey:  schema.NewField("idempotency_key", "idempotencyKey"),
	RequestHash:     schema.NewField("request_hash", "requestHash"),
	UploadExpiresAt: schema.NewField("upload_expires_at", "uploadExpiresAt"),
	CompletedAt:     schema.NewField("completed_at", "completedAt"),
	CreatedAt:       schema.NewField("created_at", "createdAt"),
	UpdatedAt:       schema.NewField("updated_at", "updatedAt"),
	DeletedAt:       schema.NewField("deleted_at", "deletedAt"),
}

func (m *Asset) TableName() string { return AssetTable }

func (m *Asset) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(AssetPrefix)
	}
	return nil
}

type CreateUploadRequest struct {
	IdempotencyKey string `json:"idempotencyKey" binding:"omitempty,max=128"`
	FileName       string `json:"fileName" binding:"required,max=255"`
	ContentType    string `json:"contentType" binding:"required,max=128"`
	SizeBytes      int64  `json:"sizeBytes" binding:"required,gt=0"`
}

type CreateUploadResponse struct {
	AssetID   string            `json:"assetId"`
	Upload    *UploadDescriptor `json:"upload"`
	PublicURL string            `json:"publicUrl"`
	Status    Status            `json:"status"`
	ExpiresAt *time.Time        `json:"expiresAt,omitempty"`
}

type UploadDescriptor struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type CompleteUploadResponse struct {
	AssetID   string `json:"assetId"`
	PublicURL string `json:"publicUrl"`
	Status    Status `json:"status"`
}

func NormalizeIdempotencyKey(k string) (string, error) {
	k = strings.TrimSpace(k)
	if k == "" {
		return "", nil
	}
	if len(k) > 128 {
		return "", problem.BadRequest("idempotencyKey is too long").With("field", "idempotencyKey")
	}
	return k, nil
}

func RequestFingerprint(purpose Purpose, visibility Visibility, fileName, contentType string, sizeBytes int64) string {
	h := sha256.Sum256([]byte(string(purpose) + "|" + string(visibility) + "|" + strings.TrimSpace(fileName) + "|" + strings.TrimSpace(contentType) + "|" + fmt.Sprintf("%d", sizeBytes)))
	return hex.EncodeToString(h[:])
}

func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "file"
	}
	// Keep it simple; avoid path traversal and exotic chars.
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "..", "_")
	if len(name) > 80 {
		name = name[:80]
	}
	return name
}
