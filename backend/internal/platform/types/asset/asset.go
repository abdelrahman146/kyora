package asset

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
)

// AssetMetadata represents optional metadata for an asset.
// Stored as JSONB in database but strongly typed in Go.
type AssetMetadata struct {
	AltText string `json:"altText,omitempty"`
	Caption string `json:"caption,omitempty"`
	Width   *int   `json:"width,omitempty"`
	Height  *int   `json:"height,omitempty"`
}

// AssetReference represents a reference to an asset, used in other domains (inventory, business).
// This allows flexible asset referencing with optional metadata.
// URL and ThumbnailURL are CDN URLs (primary access). OriginalURL is the storage provider URL (fallback).
type AssetReference struct {
	URL                  string         `json:"url" binding:"required,max=2048"`
	OriginalURL          *string        `json:"originalUrl,omitempty" binding:"omitempty,max=2048"`
	ThumbnailURL         *string        `json:"thumbnailUrl,omitempty" binding:"omitempty,max=2048"`
	ThumbnailOriginalURL *string        `json:"thumbnailOriginalUrl,omitempty" binding:"omitempty,max=2048"`
	AssetID              *string        `json:"assetId,omitempty" binding:"omitempty"`
	Metadata             *AssetMetadata `json:"metadata,omitempty"`
}

// Value implements driver.Valuer for JSONB storage.
func (a AssetReference) Value() (driver.Value, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan implements sql.Scanner for JSONB retrieval.
func (a *AssetReference) Scan(value any) error {
	if a == nil {
		return problem.InternalError().WithError(errors.New("AssetReference scan into nil receiver"))
	}
	if value == nil {
		*a = AssetReference{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return problem.InternalError().WithError(errors.New("unexpected scan type for AssetReference"))
	}
}
