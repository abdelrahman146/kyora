package blob

import (
	"context"
	"time"
)

// Provider defines the minimum interface Kyora needs to support direct-to-blob uploads
// (via presigned URLs) and basic object verification.
//
// This abstraction is intentionally small to keep it easy to swap implementations
// (DigitalOcean Spaces, AWS S3, GCS, Azure Blob, etc.).
type Provider interface {
	// PresignPut returns a URL (and required headers) that allows uploading the object directly.
	PresignPut(ctx context.Context, in PresignPutInput) (*PresignPutOutput, error)

	// Head returns basic metadata for an existing object.
	Head(ctx context.Context, key string) (*ObjectInfo, error)

	// Delete removes an object. Must be idempotent: deleting a missing object should return nil.
	Delete(ctx context.Context, key string) error

	// PublicURL returns a public URL for the given key, if the provider is configured for public access.
	// If the provider cannot produce a public URL, ok will be false.
	PublicURL(key string) (url string, ok bool)
}

type PresignPutInput struct {
	Key         string
	ContentType string
	// SizeBytes is optional; some providers can enforce it, others cannot.
	SizeBytes int64
	ExpiresIn time.Duration
}

type PresignPutOutput struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	// ExpiresAt is the server-side expected expiry.
	ExpiresAt time.Time `json:"expiresAt"`
}

type ObjectInfo struct {
	Key          string    `json:"key"`
	SizeBytes    int64     `json:"sizeBytes"`
	ContentType  string    `json:"contentType"`
	ETag         string    `json:"etag"`
	LastModified time.Time `json:"lastModified"`
}
