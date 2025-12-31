# Kyora Asset Upload Guide

**Last Updated:** December 31, 2025  
**Target Audience:** AI Agents, Backend Developers, Frontend Developers

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [AssetReference Type](#assetreference-type)
4. [Upload Flows](#upload-flows)
5. [API Endpoints](#api-endpoints)
6. [Frontend Implementation Examples](#frontend-implementation-examples)
7. [Backend Usage Patterns](#backend-usage-patterns)
8. [Migration from Legacy System](#migration-from-legacy-system)
9. [Configuration](#configuration)
10. [Troubleshooting](#troubleshooting)

---

## Overview

Kyora's asset management system is designed for **simplicity and flexibility**. Unlike traditional systems that track asset lifecycle states, Kyora uses a **pre-signed URL approach** where:

- **No complex state management** (no "pending", "ready", "orphan" states)
- **Assets are generated immediately** with unique IDs
- **Clients handle the upload** using pre-signed URLs
- **Automatic garbage collection** cleans up unused assets
- **Flexible storage** supports both S3-compatible storage and local filesystem

### Key Principles

1. **Simplicity First**: No tracking of upload states, no idempotency keys, no purpose enums
2. **Client Responsibility**: Frontend handles chunking, retries, and progress tracking
3. **AssetReference Pattern**: Other domains store structured asset metadata (URL, assetId, metadata)
4. **Provider Agnostic**: Works with S3, MinIO, local filesystem, or any S3-compatible storage

---

## Architecture

```
┌─────────────────┐
│   Frontend      │
│  (React/Web)    │
└────────┬────────┘
         │ 1. Request pre-signed URLs
         ▼
┌─────────────────────────────────┐
│  POST /businesses/:desc/assets/  │
│         uploads                  │
│  ─────────────────────────────  │
│  GenerateUploadURLs Handler      │
└────────┬────────────────────────┘
         │ 2. Generate assetId + URLs
         ▼
┌─────────────────────────────────┐
│    Asset Service                 │
│  • S3: Create multipart upload   │
│  • Local: Return POST endpoint   │
└────────┬────────────────────────┘
         │ 3. Return upload info
         ▼
┌─────────────────┐
│   Frontend      │
│  • S3: PUT to   │
│    part URLs    │
│  • Local: POST  │
│    to endpoint  │
└────────┬────────┘
         │ 4. Upload complete
         ▼
┌─────────────────────────────────┐
│  POST /businesses/:desc/assets/  │
│    uploads/:assetId/complete     │
│  (S3 multipart only)             │
└─────────────────────────────────┘
```

### Storage Providers

#### S3-Compatible (AWS S3, MinIO, DigitalOcean Spaces, etc.)
- **Multipart uploads** for files > 10MB (configurable)
- **Pre-signed PUT URLs** for each part (max 10,000 parts)
- **Client-side chunking** and parallel uploads
- **ETag collection** for part completion
- **Server-side assembly** via CompleteMultipartUpload

#### Local Filesystem
- **Direct POST endpoint** (`/v1/assets/internal/upload/:assetId`)
- **Single-request upload** (suitable for development)
- **File storage** in configurable directory (default: `./tmp/assets`)
- **Simple and fast** for local testing

---

## AssetReference Type

The `AssetReference` type is the cornerstone of asset management in Kyora. It's defined in `internal/platform/types/asset/asset.go` to avoid circular dependencies.

### Structure

```go
// AssetReference represents a reference to an asset
type AssetReference struct {
    URL      string         `json:"url" binding:"required,max=2048"`
    AssetID  *string        `json:"assetId,omitempty"`
    Metadata *AssetMetadata `json:"metadata,omitempty"`
}

// AssetMetadata provides optional semantic information
type AssetMetadata struct {
    AltText string `json:"altText,omitempty"`  // Accessibility text
    Caption string `json:"caption,omitempty"`  // Display caption
    Width   *int   `json:"width,omitempty"`    // Image width in pixels
    Height  *int   `json:"height,omitempty"`   // Image height in pixels
}
```

### Key Features

- **JSONB Storage**: Stored as JSONB in PostgreSQL for flexibility
- **URL Required**: Always includes the public URL for immediate display
- **Optional AssetID**: Links to Asset table for garbage collection
- **Rich Metadata**: Support for accessibility and display information
- **Null-Safe**: Can be `nil` or `{}` when no asset is present

### Usage in Other Domains

#### Business Logo (Optional Single Asset)
```go
type Business struct {
    // ...
    Logo *asset.AssetReference `gorm:"column:logo;type:jsonb" json:"logo,omitempty"`
}
```

#### Inventory Photos (Required List)
```go
type Product struct {
    // ...
    Photos AssetReferenceList `gorm:"column:photos;type:jsonb;not null;default:'[]'" json:"photos"`
}

// AssetReferenceList is a JSONB-backed array
type AssetReferenceList []asset.AssetReference
```

---

## Upload Flows

### Flow 1: S3 Multipart Upload (Production)

**Use Case:** Files larger than part size (default 10MB), production environments

#### Step 1: Request Upload URLs

**Request:**
```http
POST /v1/businesses/my-store/assets/uploads
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "files": [
    {
      "fileName": "product-photo.jpg",
      "contentType": "image/jpeg",
      "sizeBytes": 52428800  // 50MB
    }
  ]
}
```

**Response:**
```json
{
  "uploads": [
    {
      "assetId": "ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1",
      "fileName": "product-photo.jpg",
      "contentType": "image/jpeg",
      "sizeBytes": 52428800,
      "publicUrl": "https://cdn.kyora.app/assets/my-store/ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1/product-photo.jpg",
      "uploadType": "multipart",
      "uploadId": "abc123xyz...",
      "partSize": 10485760,  // 10MB
      "totalParts": 5,
      "partUrls": [
        "https://s3.amazonaws.com/bucket/path?partNumber=1&uploadId=...",
        "https://s3.amazonaws.com/bucket/path?partNumber=2&uploadId=...",
        "https://s3.amazonaws.com/bucket/path?partNumber=3&uploadId=...",
        "https://s3.amazonaws.com/bucket/path?partNumber=4&uploadId=...",
        "https://s3.amazonaws.com/bucket/path?partNumber=5&uploadId=..."
      ],
      "expiresAt": "2025-12-31T12:30:00Z"
    }
  ]
}
```

#### Step 2: Upload Parts to S3

The client splits the file into chunks and uploads each part:

```javascript
const uploadPart = async (url, chunk, partNumber) => {
  const response = await fetch(url, {
    method: 'PUT',
    body: chunk,
    headers: {
      'Content-Type': 'application/octet-stream'
    }
  });
  
  // Extract ETag from response headers (critical!)
  const etag = response.headers.get('ETag').replace(/"/g, '');
  return { partNumber, etag };
};

// Upload all parts (potentially in parallel)
const completedParts = [];
for (let i = 0; i < upload.totalParts; i++) {
  const start = i * upload.partSize;
  const end = Math.min(start + upload.partSize, file.size);
  const chunk = file.slice(start, end);
  
  const part = await uploadPart(upload.partUrls[i], chunk, i + 1);
  completedParts.push(part);
}
```

#### Step 3: Complete Multipart Upload

**Request:**
```http
POST /v1/businesses/my-store/assets/uploads/ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1/complete
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "uploadId": "abc123xyz...",
  "parts": [
    { "partNumber": 1, "etag": "d41d8cd98f00b204e9800998ecf8427e" },
    { "partNumber": 2, "etag": "e4d909c290d0fb1ca068ffaddf22cbd0" },
    { "partNumber": 3, "etag": "a1b2c3d4e5f6789012345678901234ab" },
    { "partNumber": 4, "etag": "b2c3d4e5f6789012345678901234abcd" },
    { "partNumber": 5, "etag": "c3d4e5f6789012345678901234abcdef" }
  ]
}
```

**Response:**
```json
{
  "assetId": "ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1",
  "publicUrl": "https://cdn.kyora.app/assets/my-store/ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1/product-photo.jpg",
  "status": "completed"
}
```

#### Step 4: Use AssetReference

Now use the asset in your domain:

```json
{
  "name": "Premium T-Shirt",
  "photos": [
    {
      "url": "https://cdn.kyora.app/assets/my-store/ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1/product-photo.jpg",
      "assetId": "ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1",
      "metadata": {
        "altText": "Premium white t-shirt front view",
        "width": 2000,
        "height": 2000
      }
    }
  ]
}
```

---

### Flow 2: Local Simple Upload (Development)

**Use Case:** Small files, local development, testing

#### Step 1: Request Upload URLs

**Request:**
```http
POST /v1/businesses/my-store/assets/uploads
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "files": [
    {
      "fileName": "logo.png",
      "contentType": "image/png",
      "sizeBytes": 51200  // 50KB
    }
  ]
}
```

**Response:**
```json
{
  "uploads": [
    {
      "assetId": "ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1",
      "fileName": "logo.png",
      "contentType": "image/png",
      "sizeBytes": 51200,
      "publicUrl": "http://localhost:8080/v1/public/assets/ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1",
      "uploadType": "simple",
      "uploadUrl": "http://localhost:8080/v1/assets/internal/upload/ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1",
      "expiresAt": "2025-12-31T12:30:00Z"
    }
  ]
}
```

#### Step 2: Upload File Directly

```javascript
const formData = new FormData();
formData.append('file', file);

const response = await fetch(upload.uploadUrl, {
  method: 'POST',
  body: formData
});

const result = await response.json();
// { "assetId": "ast_...", "publicUrl": "...", "status": "completed" }
```

#### Step 3: Use AssetReference

```json
{
  "name": "My Store",
  "logo": {
    "url": "http://localhost:8080/v1/public/assets/ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1",
    "assetId": "ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1",
    "metadata": {
      "altText": "My Store Logo"
    }
  }
}
```

---

## API Endpoints

### POST /v1/businesses/:businessDescriptor/assets/uploads

**Purpose:** Generate pre-signed upload URLs for one or more files

**Authentication:** Required (JWT Bearer token)

**Business Scoping:** Yes (via `:businessDescriptor` param)

**Request Body:**
```typescript
{
  files: Array<{
    fileName: string;      // Max 255 chars
    contentType: string;   // MIME type, max 128 chars
    sizeBytes: number;     // File size in bytes, must be > 0
  }>  // Min 1, max 50 files per request
}
```

**Response:**
```typescript
{
  uploads: Array<{
    assetId: string;           // Generated KSUID with "ast_" prefix
    fileName: string;
    contentType: string;
    sizeBytes: number;
    publicUrl: string;         // CDN URL for serving the asset
    uploadType: "multipart" | "simple";
    
    // For multipart uploads (S3):
    uploadId?: string;         // S3 multipart upload ID
    partSize?: number;         // Bytes per part (default 10MB)
    totalParts?: number;       // Number of parts
    partUrls?: string[];       // Pre-signed PUT URLs for each part
    
    // For simple uploads (local):
    uploadUrl?: string;        // Direct POST endpoint
    
    expiresAt: string;         // ISO 8601 timestamp
  }>
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body, files array empty/too large
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: User doesn't have access to the business
- `500 Internal Server Error`: Storage provider error

---

### POST /v1/businesses/:businessDescriptor/assets/uploads/:assetId/complete

**Purpose:** Complete a multipart upload (S3 only)

**Authentication:** Required (JWT Bearer token)

**Business Scoping:** Yes (via `:businessDescriptor` param)

**Request Body:**
```typescript
{
  uploadId: string;          // S3 multipart upload ID
  parts: Array<{
    partNumber: number;      // 1-based part number
    etag: string;            // ETag returned by S3 for this part
  }>  // Must match totalParts from generation
}
```

**Response:**
```typescript
{
  assetId: string;
  publicUrl: string;
  status: "completed"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid assetId, missing uploadId, invalid parts array
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: User doesn't have access to the business
- `404 Not Found`: Asset not found
- `500 Internal Server Error`: S3 completion failed

---

### GET /v1/public/assets/:assetId

**Purpose:** Serve public asset (local provider only)

**Authentication:** Not required (public endpoint)

**Business Scoping:** No

**Response:**
- File content with appropriate `Content-Type`, `Cache-Control`, `ETag`, and `Last-Modified` headers
- Supports `If-None-Match` for efficient caching (304 Not Modified)

**Error Responses:**
- `404 Not Found`: Asset not found or not accessible
- `500 Internal Server Error`: File read error

---

### POST /v1/assets/internal/upload/:assetId (Internal Only)

**Purpose:** Direct file upload endpoint for local provider

**Authentication:** Not required (internal endpoint, should be protected by network rules)

**Business Scoping:** No

**Request:** Multipart form data with `file` field

**Response:**
```typescript
{
  assetId: string;
  publicUrl: string;
  status: "completed"
}
```

---

## Frontend Implementation Examples

### React Hook: useAssetUpload

```typescript
import { useState } from 'react';
import { useAuth } from './useAuth';

interface UploadProgress {
  fileName: string;
  progress: number; // 0-100
  status: 'pending' | 'uploading' | 'completed' | 'error';
  assetId?: string;
  publicUrl?: string;
  error?: string;
}

export const useAssetUpload = (businessDescriptor: string) => {
  const { token } = useAuth();
  const [uploads, setUploads] = useState<Map<string, UploadProgress>>(new Map());

  const updateUpload = (fileName: string, update: Partial<UploadProgress>) => {
    setUploads(prev => {
      const next = new Map(prev);
      next.set(fileName, { ...prev.get(fileName)!, ...update });
      return next;
    });
  };

  const uploadFiles = async (files: File[]): Promise<AssetReference[]> => {
    // Step 1: Request upload URLs
    const response = await fetch(
      `/v1/businesses/${businessDescriptor}/assets/uploads`,
      {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          files: files.map(f => ({
            fileName: f.name,
            contentType: f.type,
            sizeBytes: f.size
          }))
        })
      }
    );

    if (!response.ok) {
      throw new Error(`Upload request failed: ${response.statusText}`);
    }

    const { uploads: uploadDescriptors } = await response.json();

    // Step 2: Upload each file
    const results = await Promise.all(
      files.map((file, index) => {
        const descriptor = uploadDescriptors[index];
        
        setUploads(prev => new Map(prev).set(file.name, {
          fileName: file.name,
          progress: 0,
          status: 'pending'
        }));

        return descriptor.uploadType === 'multipart'
          ? uploadMultipart(file, descriptor)
          : uploadSimple(file, descriptor);
      })
    );

    return results;
  };

  const uploadMultipart = async (file: File, descriptor: any): Promise<AssetReference> => {
    updateUpload(file.name, { status: 'uploading' });

    const completedParts: { partNumber: number; etag: string }[] = [];

    // Upload parts with progress tracking
    for (let i = 0; i < descriptor.totalParts; i++) {
      const start = i * descriptor.partSize;
      const end = Math.min(start + descriptor.partSize, file.size);
      const chunk = file.slice(start, end);

      const response = await fetch(descriptor.partUrls[i], {
        method: 'PUT',
        body: chunk,
        headers: { 'Content-Type': 'application/octet-stream' }
      });

      if (!response.ok) {
        updateUpload(file.name, {
          status: 'error',
          error: `Part ${i + 1} upload failed`
        });
        throw new Error(`Part ${i + 1} upload failed`);
      }

      const etag = response.headers.get('ETag')?.replace(/"/g, '') || '';
      completedParts.push({ partNumber: i + 1, etag });

      // Update progress
      const progress = Math.round(((i + 1) / descriptor.totalParts) * 100);
      updateUpload(file.name, { progress });
    }

    // Step 3: Complete multipart upload
    const completeResponse = await fetch(
      `/v1/businesses/${businessDescriptor}/assets/uploads/${descriptor.assetId}/complete`,
      {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          uploadId: descriptor.uploadId,
          parts: completedParts
        })
      }
    );

    if (!completeResponse.ok) {
      updateUpload(file.name, {
        status: 'error',
        error: 'Failed to complete upload'
      });
      throw new Error('Failed to complete multipart upload');
    }

    updateUpload(file.name, {
      status: 'completed',
      progress: 100,
      assetId: descriptor.assetId,
      publicUrl: descriptor.publicUrl
    });

    return {
      url: descriptor.publicUrl,
      assetId: descriptor.assetId,
      metadata: {
        altText: file.name
      }
    };
  };

  const uploadSimple = async (file: File, descriptor: any): Promise<AssetReference> => {
    updateUpload(file.name, { status: 'uploading', progress: 50 });

    const formData = new FormData();
    formData.append('file', file);

    const response = await fetch(descriptor.uploadUrl, {
      method: 'POST',
      body: formData
    });

    if (!response.ok) {
      updateUpload(file.name, {
        status: 'error',
        error: 'Upload failed'
      });
      throw new Error('Simple upload failed');
    }

    updateUpload(file.name, {
      status: 'completed',
      progress: 100,
      assetId: descriptor.assetId,
      publicUrl: descriptor.publicUrl
    });

    return {
      url: descriptor.publicUrl,
      assetId: descriptor.assetId,
      metadata: {
        altText: file.name
      }
    };
  };

  return { uploadFiles, uploads: Array.from(uploads.values()) };
};
```

### Usage Example

```typescript
const ProductPhotoUploader = () => {
  const { uploadFiles, uploads } = useAssetUpload('my-store');
  const [photos, setPhotos] = useState<AssetReference[]>([]);

  const handleFileSelect = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(event.target.files || []);
    
    try {
      const assetRefs = await uploadFiles(files);
      setPhotos(prev => [...prev, ...assetRefs]);
    } catch (error) {
      console.error('Upload failed:', error);
    }
  };

  return (
    <div>
      <input type="file" multiple accept="image/*" onChange={handleFileSelect} />
      
      {uploads.map(upload => (
        <div key={upload.fileName}>
          <span>{upload.fileName}</span>
          <progress value={upload.progress} max="100" />
          <span>{upload.status}</span>
        </div>
      ))}

      <div>
        {photos.map(photo => (
          <img key={photo.assetId} src={photo.url} alt={photo.metadata?.altText} />
        ))}
      </div>
    </div>
  );
};
```

---

## Backend Usage Patterns

### Creating Products with Photos

```go
func (s *Service) CreateProduct(
    ctx context.Context,
    actor *account.User,
    biz *business.Business,
    req *CreateProductRequest,
) (*Product, error) {
    // Validate photos limit
    if len(req.Photos) > s.maxPhotosPerProduct {
        return nil, problem.BadRequest("too many photos").
            With("max", s.maxPhotosPerProduct).
            With("provided", len(req.Photos))
    }

    // Convert to AssetReferenceList for JSONB storage
    photos := AssetReferenceList(req.Photos)

    product := &Product{
        BusinessID:  biz.ID,
        Name:        req.Name,
        Description: req.Description,
        Photos:      photos, // Stored as JSONB
        CategoryID:  req.CategoryID,
    }

    if err := s.storage.products.CreateOne(ctx, product); err != nil {
        return nil, err
    }

    return product, nil
}
```

### Updating Business Logo

```go
func (s *Service) UpdateBusiness(
    ctx context.Context,
    actor *account.User,
    biz *business.Business,
    input *UpdateBusinessInput,
) (*business.Business, error) {
    // ... other updates ...

    if input.Logo != nil {
        // Logo is *asset.AssetReference (can be nil)
        biz.Logo = input.Logo
    }

    if err := s.storage.UpdateOne(ctx, biz); err != nil {
        return nil, err
    }

    return biz, nil
}
```

### Querying Assets for Garbage Collection

```go
// Find assets that are not referenced in any business or inventory
func (s *Storage) ListUnusedAssets(
    ctx context.Context,
    minAge time.Duration,
    limit int,
) ([]*Asset, error) {
    cutoff := time.Now().Add(-minAge)

    var assets []*Asset
    err := s.asset.FindMany(
        ctx,
        &assets,
        s.asset.WithLimit(limit),
        s.asset.ScopeWhere("created_at < ?", cutoff),
        // Use JSON extraction to check if assetId is not referenced
        s.asset.ScopeWhere(`
            NOT EXISTS (
                SELECT 1 FROM businesses
                WHERE logo->>'assetId' = uploaded_assets.id
            )
            AND NOT EXISTS (
                SELECT 1 FROM products
                WHERE EXISTS (
                    SELECT 1 FROM jsonb_array_elements(photos) AS photo
                    WHERE photo->>'assetId' = uploaded_assets.id
                )
            )
            AND NOT EXISTS (
                SELECT 1 FROM variants
                WHERE EXISTS (
                    SELECT 1 FROM jsonb_array_elements(photos) AS photo
                    WHERE photo->>'assetId' = uploaded_assets.id
                )
            )
        `),
    )

    return assets, err
}
```

---

## Migration from Legacy System

### Key Differences

| Aspect | Legacy System | New System |
|--------|--------------|------------|
| **State Management** | Complex (Pending → Ready → Orphan) | None (immediate generation) |
| **Purpose Field** | Required enum (Logo, ProductPhoto, etc.) | Removed (context from domain) |
| **Idempotency** | Tracked via idempotency keys | Client responsibility |
| **Endpoints** | 9 purpose-specific endpoints | 2 universal endpoints |
| **Storage Format** | Plain URL strings | Structured AssetReference |
| **Metadata** | Not supported | Rich metadata support |

### Migration Steps

1. **Update Frontend Code**
   - Replace purpose-specific upload calls with universal `GenerateUploadURLs`
   - Add ETag collection for S3 multipart uploads
   - Update asset references to use AssetReference structure

2. **Update Backend Models**
   - Change `LogoURL string` to `Logo *asset.AssetReference`
   - Change `Photos []string` to `Photos AssetReferenceList`
   - Remove purpose parameters from service methods

3. **Data Migration** (if needed)
   ```sql
   -- Migrate business logos
   UPDATE businesses
   SET logo = jsonb_build_object(
       'url', logo_url,
       'assetId', NULL,
       'metadata', NULL
   )
   WHERE logo_url IS NOT NULL AND logo_url != '';

   -- Migrate product photos
   UPDATE products
   SET photos = (
       SELECT jsonb_agg(
           jsonb_build_object(
               'url', photo_url,
               'assetId', NULL,
               'metadata', NULL
           )
       )
       FROM unnest(string_to_array(photos_text, ',')) AS photo_url
   );
   ```

4. **Remove Legacy Code**
   - Delete old handlers and routes
   - Remove Purpose, Status, Visibility enums
   - Clean up idempotency tracking

---

## Configuration

### Environment Variables

```yaml
# Storage configuration
storage:
  provider: "s3"  # or "local"
  local_path: "./tmp/assets"  # Local storage directory
  multipart_part_size_mb: 10  # S3 part size in MB (5-5000)

# S3-compatible storage (when provider=s3)
s3:
  endpoint: "https://s3.amazonaws.com"
  region: "us-east-1"
  bucket: "kyora-assets"
  access_key_id: "${AWS_ACCESS_KEY_ID}"
  secret_access_key: "${AWS_SECRET_ACCESS_KEY}"
  cdn_url: "https://cdn.kyora.app"  # Optional CDN prefix

# Asset limits
inventory:
  max_photos_per_product: 10  # Maximum photos per product/variant
```

### Configurable Limits

- **Max files per request**: 50 (hard-coded)
- **Max file name length**: 255 characters
- **Max content type length**: 128 characters
- **Max URL length**: 2048 characters
- **Part size**: 5MB - 5GB (configurable, default 10MB)
- **Max parts**: 10,000 (S3 limit)
- **Photos per product**: Configurable (default 10)

---

## Troubleshooting

### Issue: "Part upload failed with 403 Forbidden"

**Cause:** Pre-signed URL expired (default 1 hour)

**Solution:** Request new upload URLs and retry

### Issue: "Multipart completion failed: InvalidPart"

**Cause:** ETag mismatch or missing parts

**Solution:** 
- Ensure all ETags are collected correctly (without quotes)
- Verify all parts were uploaded successfully
- Check part numbers are 1-based and sequential

### Issue: "Asset not found after upload"

**Cause:** Asset record not created or database issue

**Solution:**
- Check asset service logs for creation errors
- Verify database connection
- Ensure proper business scoping

### Issue: "Local uploads fail with permission error"

**Cause:** Storage directory not writable

**Solution:**
- Create directory: `mkdir -p ./tmp/assets`
- Set permissions: `chmod 755 ./tmp/assets`
- Verify `storage.local_path` configuration

### Issue: "S3 uploads work but CDN serves 404"

**Cause:** CDN not configured or cache not invalidated

**Solution:**
- Verify `s3.cdn_url` configuration
- Check CDN origin settings point to S3 bucket
- Invalidate CDN cache if needed
- Use S3 direct URL temporarily

### Issue: "Photos not displaying in frontend"

**Cause:** CORS not configured on S3 bucket

**Solution:**
```json
{
  "CORSRules": [
    {
      "AllowedOrigins": ["*"],
      "AllowedMethods": ["GET", "PUT"],
      "AllowedHeaders": ["*"],
      "ExposeHeaders": ["ETag"],
      "MaxAgeSeconds": 3600
    }
  ]
}
```

---

## Best Practices

### For Frontend Developers

1. **Always collect ETags** for S3 multipart uploads
2. **Implement retry logic** for failed part uploads
3. **Show progress indicators** using part completion
4. **Validate file types** before uploading
5. **Use AssetReference format** consistently
6. **Add meaningful alt text** for accessibility
7. **Handle upload errors** gracefully with user feedback

### For Backend Developers

1. **Always use AssetReference** instead of plain strings
2. **Validate photo limits** before creating records
3. **Use database-agnostic queries** for GC
4. **Cache public assets** with proper headers
5. **Log upload failures** with context
6. **Monitor orphan asset growth** regularly
7. **Test with both S3 and local providers**

### For DevOps

1. **Configure CDN** for production asset serving
2. **Set up S3 lifecycle policies** for old multipart uploads
3. **Monitor storage costs** and usage patterns
4. **Enable S3 versioning** for disaster recovery
5. **Set up CloudFront** or similar CDN
6. **Configure bucket CORS** properly
7. **Implement asset backup strategy**

---

## Support & Resources

- **Backend Code**: `internal/domain/asset/`
- **Asset Types**: `internal/platform/types/asset/`
- **API Routes**: `internal/server/routes.go`
- **Configuration**: `.kyora.yaml.example`
- **Seed Examples**: `cmd/seed.go`

For questions or issues, consult the Kyora development team or check the inline code documentation.

---

**Document Version:** 1.0  
**Last Reviewed:** December 31, 2025
