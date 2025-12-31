# Kyora Asset Upload Guide

**Last Updated:** January 2025  
**Target Audience:** AI Agents, Backend Developers, Frontend Developers

---

## Table of Contents

1. [Overview](#overview)
2. [Key Features](#key-features)
3. [Architecture](#architecture)
4. [File Type Support](#file-type-support)
5. [Thumbnail Generation](#thumbnail-generation)
6. [CDN Integration](#cdn-integration)
7. [AssetReference Type](#assetreference-type)
8. [Upload Flows](#upload-flows)
9. [API Endpoints](#api-endpoints)
10. [Frontend Implementation](#frontend-implementation)
11. [Backend Usage Patterns](#backend-usage-patterns)
12. [Configuration](#configuration)
13. [Migration Guide](#migration-guide)
14. [Troubleshooting](#troubleshooting)

---

## Overview

Kyora's asset management system provides a **production-ready, flexible solution** for handling media uploads with support for multiple file types, automatic thumbnail generation, and CDN delivery. The system is designed around these core principles:

- **Simplicity**: No complex state management (no "pending", "ready", "orphan" states)
- **Immediate Generation**: Assets are created with unique IDs instantly
- **Client-Side Control**: Frontend handles chunking, retries, and progress tracking
- **Provider Agnostic**: Works with S3, MinIO, CloudFront, Cloudflare, BunnyCDN, or local filesystem
- **Rich Metadata**: Support for accessibility, dimensions, and thumbnails
- **Automatic Cleanup**: Garbage collection removes unused assets

### Design Philosophy

1. **Pre-Signed URL Approach**: No server-side upload handling, clients upload directly to storage
2. **AssetReference Pattern**: Domains store structured asset metadata (URL, CDN URL, thumbnails, metadata)
3. **File Type Validation**: Configurable extensions and size limits per category
4. **Thumbnail Support**: Automatic thumbnail generation for images and videos
5. **CDN-First**: Optimized for edge caching with immutable URLs

---

## Key Features

### 📁 Multi-Format Support
- **Images**: JPEG, PNG, WebP, GIF, HEIC, HEIF
- **Videos**: MP4, MOV, AVI, MKV, WebM
- **Audio**: MP3, WAV, OGG, M4A, AAC
- **Documents**: PDF, DOC, DOCX, TXT, RTF, ODT
- **Compressed**: ZIP, TAR, GZ, RAR, 7Z
- **Configurable**: Extensions and size limits per category

### 🖼️ Automatic Thumbnails
- **Smart Generation**: Automatically creates thumbnails for images and videos
- **Client-Side Processing**: Frontend generates thumbnails before upload
- **JPEG Format**: Optimized JPEG format with configurable quality (default 80%)
- **Size Control**: Maximum dimension of 512px (configurable)
- **Tight Coupling**: Thumbnails linked via `ThumbnailAssetID` for garbage collection

### 🚀 CDN Optimization
- **Edge Delivery**: Full CDN support (CloudFront, Cloudflare, Fastly, BunnyCDN)
- **Immutable Caching**: 1-year cache headers with immutable flag
- **Multiple URLs**: Both CDN (fast) and original URLs (fallback)
- **Thumbnail CDN**: Separate CDN URLs for thumbnails
- **Graceful Degradation**: Falls back to storage URLs if CDN not configured

### 🔒 Security & Validation
- **File Type Validation**: Server-side extension and MIME type validation
- **Size Limits**: Per-category size restrictions (configurable)
- **Business Scoping**: All assets scoped to specific business
- **Authentication Required**: JWT Bearer token for all uploads
- **Pre-Signed URLs**: Time-limited upload URLs (24 hours for S3)

---

## Architecture

## Architecture

```
┌─────────────────┐
│   Frontend      │
│  (React/Web)    │
└────────┬────────┘
         │ 1. Request pre-signed URLs + thumbnails
         ▼
┌──────────────────────────────────────┐
│  POST /businesses/:desc/assets/      │
│         uploads                       │
│  ────────────────────────────────────│
│  • Validate file types & sizes        │
│  • Generate asset IDs                │
│  • Determine thumbnail requirement   │
│  • Generate CDN URLs                 │
└────────┬─────────────────────────────┘
         │ 2. Return upload descriptors
         │    (main file + thumbnail if needed)
         ▼
┌─────────────────────────────┐
│   Frontend                   │
│  • Generate thumbnail (if    │
│    image/video)              │
│  • S3: PUT parts to URLs     │
│  • Local: POST to endpoint   │
│  • Upload thumbnail          │
└────────┬────────────────────┘
         │ 3. Upload complete
         ▼
┌──────────────────────────────────────┐
│  POST /businesses/:desc/assets/      │
│    uploads/:assetId/complete         │
│  (S3 multipart only)                 │
│  ────────────────────────────────────│
│  • Assemble S3 parts                 │
│  • Complete thumbnail (if exists)    │
└──────────────────────────────────────┘
         │
         ▼
┌──────────────────────────────────────┐
│  Asset Record in Database            │
│  ────────────────────────────────────│
│  • assetId: ast_...                  │
│  • objectKey: path in storage        │
│  • publicUrl: storage URL            │
│  • cdnUrl: CDN URL (if configured)   │
│  • fileCategory: image/video/etc     │
│  • thumbnailAssetId: (if generated)  │
│  • thumbnailObjectKey: thumb path    │
│  • thumbnailPublicUrl: storage URL   │
│  • thumbnailCdnUrl: CDN URL          │
└──────────────────────────────────────┘
```

### Components

#### Frontend Layer
- Validates file selection
- Requests upload URLs from API
- Generates thumbnails for images/videos (canvas-based compression)
- Handles chunked uploads for large files
- Tracks upload progress
- Manages retries and error handling

#### API Layer (`internal/domain/asset`)
- **Handler**: HTTP endpoints (generate URLs, complete uploads, serve assets)
- **Service**: Business logic (validation, URL generation, thumbnail orchestration, CDN integration)
- **Validator**: File type and size validation (`file_types.go`)
- **Storage**: Database operations (GORM/PostgreSQL)
- **Blob Provider**: S3 multipart or local filesystem

#### Storage Providers

##### S3-Compatible (Production)
- AWS S3, MinIO, DigitalOcean Spaces, Backblaze B2
- Multipart uploads for files > part size (default 10MB)
- Pre-signed PUT URLs for each part (max 10,000 parts)
- Client-side chunking and parallel uploads
- ETag collection for part completion
- Server-side assembly via CompleteMultipartUpload

##### Local Filesystem (Development)
- Direct POST endpoint (`/v1/assets/internal/upload/:assetId`)
- Single-request upload
- File storage in configurable directory (default: `./tmp/assets`)
- Simple and fast for local testing
- No multipart complexity

---

## File Type Support

### Supported File Categories

The system supports five file categories, each with configurable extensions and size limits:

| Category | Default Extensions | Default Max Size | Thumbnail Support |
|----------|-------------------|------------------|-------------------|
| **Image** | jpg, jpeg, png, webp, gif, heic, heif | 10 MB | ✅ Yes |
| **Video** | mp4, mov, avi, mkv, webm | 100 MB | ✅ Yes |
| **Audio** | mp3, wav, ogg, m4a, aac | 20 MB | ❌ No |
| **Document** | pdf, doc, docx, txt, rtf, odt | 10 MB | ❌ No |
| **Compressed** | zip, tar, gz, rar, 7z | 50 MB | ❌ No |

### File Type Validation

Validation happens in two stages:

1. **Extension Check**: File extension must be in the allowed list for its category
2. **Size Check**: File size must not exceed the category's maximum size

**Example Configuration** (`.kyora.yaml`):

```yaml
uploads:
  allowed_extensions:
    image: ["jpg", "jpeg", "png", "webp", "gif", "heic", "heif"]
    video: ["mp4", "mov", "avi", "mkv", "webm"]
    audio: ["mp3", "wav", "ogg", "m4a", "aac"]
    document: ["pdf", "doc", "docx", "txt", "rtf", "odt"]
    compressed: ["zip", "tar", "gz", "rar", "7z"]
  max_size_bytes:
    image: 10485760      # 10 MB
    video: 104857600     # 100 MB
    audio: 20971520      # 20 MB
    document: 10485760   # 10 MB
    compressed: 52428800 # 50 MB
    other: 5242880       # 5 MB (fallback)
```

### File Type Validator Implementation

The `FileTypeValidator` (in `internal/domain/asset/file_types.go`) provides:

```go
type FileTypeValidator struct {
    allowedExtensions map[FileCategory][]string
    maxSizeBytes      map[FileCategory]int64
}

// Check if file is allowed
func (v *FileTypeValidator) IsAllowed(fileName string) bool

// Get file category (image, video, audio, document, compressed)
func (v *FileTypeValidator) GetCategory(fileName string) FileCategory

// Get max size for file
func (v *FileTypeValidator) GetMaxSize(fileName string) int64

// Check if file needs thumbnail
func (v *FileTypeValidator) NeedsThumbnail(fileName string) bool
```

### Frontend File Type Validation

**Before Upload**:

```typescript
const ACCEPTED_IMAGES = '.jpg,.jpeg,.png,.webp,.gif,.heic,.heif';
const ACCEPTED_VIDEOS = '.mp4,.mov,.avi,.mkv,.webm';
const ACCEPTED_DOCUMENTS = '.pdf,.doc,.docx,.txt,.rtf,.odt';

<input
  type="file"
  accept={ACCEPTED_IMAGES}
  multiple
  onChange={handleFileSelect}
/>
```

**During Selection**:

```typescript
const MAX_IMAGE_SIZE = 10 * 1024 * 1024; // 10MB
const MAX_VIDEO_SIZE = 100 * 1024 * 1024; // 100MB

const validateFile = (file: File): string | null => {
  const ext = file.name.split('.').pop()?.toLowerCase();
  
  // Check image
  if (['jpg', 'jpeg', 'png', 'webp', 'gif', 'heic', 'heif'].includes(ext!)) {
    if (file.size > MAX_IMAGE_SIZE) {
      return `Image too large: ${(file.size / 1024 / 1024).toFixed(1)}MB (max 10MB)`;
    }
    return null;
  }
  
  // Check video
  if (['mp4', 'mov', 'avi', 'mkv', 'webm'].includes(ext!)) {
    if (file.size > MAX_VIDEO_SIZE) {
      return `Video too large: ${(file.size / 1024 / 1024).toFixed(1)}MB (max 100MB)`;
    }
    return null;
  }
  
  return `Unsupported file type: ${ext}`;
};
```

---

## Thumbnail Generation

### Overview

The system automatically generates thumbnails for images and videos to improve loading performance and provide preview capabilities. Thumbnails are:

- **Client-generated**: Frontend creates thumbnails before upload
- **JPEG format**: Optimized format with configurable quality
- **512px max dimension**: Maintains aspect ratio (configurable)
- **Tightly coupled**: Linked via `ThumbnailAssetID` for lifecycle management
- **Separate uploads**: Thumbnails get their own asset IDs and upload URLs

### When Thumbnails Are Generated

| File Category | Thumbnail Support | Use Case |
|--------------|-------------------|----------|
| **Image** | ✅ Yes | Product photos, profile pictures, gallery images |
| **Video** | ✅ Yes | Video previews, video player posters |
| **Audio** | ❌ No | Not applicable |
| **Document** | ❌ No | Use document icon instead |
| **Compressed** | ❌ No | Use archive icon instead |

### API Response with Thumbnail

When you upload an image or video, the API response includes **two** upload descriptors:

```json
{
  "uploads": [
    {
      "assetId": "ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1",
      "fileName": "product-photo.jpg",
      "contentType": "image/jpeg",
      "sizeBytes": 5242880,
      "publicUrl": "https://s3.amazonaws.com/bucket/my-store/ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1/product-photo.jpg",
      "cdnUrl": "https://cdn.kyora.app/assets/my-store/ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1/product-photo.jpg",
      "uploadType": "simple",
      "uploadUrl": "http://localhost:8080/v1/assets/internal/upload/ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1",
      "isThumbnail": false,
      "thumbnail": {
        "assetId": "ast_3BrLMzZ4o9NqX0C6WsDgHrO5Ud2",
        "fileName": "product-photo_thumb.jpg",
        "contentType": "image/jpeg",
        "publicUrl": "https://s3.amazonaws.com/bucket/my-store/ast_3BrLMzZ4o9NqX0C6WsDgHrO5Ud2/product-photo_thumb.jpg",
        "cdnUrl": "https://cdn.kyora.app/assets/my-store/ast_3BrLMzZ4o9NqX0C6WsDgHrO5Ud2/product-photo_thumb.jpg",
        "uploadUrl": "http://localhost:8080/v1/assets/internal/upload/ast_3BrLMzZ4o9NqX0C6WsDgHrO5Ud2"
      },
      "expiresAt": "2025-12-31T12:30:00Z"
    }
  ]
}
```

### Frontend Thumbnail Generation

#### React Hook: `useThumbnailGenerator`

```typescript
import { useCallback } from 'react';

interface ThumbnailConfig {
  maxDimension?: number; // Default 512
  quality?: number; // Default 0.8 (80%)
  format?: 'image/jpeg' | 'image/png'; // Default jpeg
}

export const useThumbnailGenerator = (config: ThumbnailConfig = {}) => {
  const maxDimension = config.maxDimension || 512;
  const quality = config.quality || 0.8;
  const format = config.format || 'image/jpeg';

  const generateThumbnail = useCallback(async (file: File): Promise<Blob> => {
    return new Promise((resolve, reject) => {
      const img = new Image();
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');

      img.onload = () => {
        // Calculate dimensions maintaining aspect ratio
        let { width, height } = img;
        if (width > height) {
          if (width > maxDimension) {
            height = (height * maxDimension) / width;
            width = maxDimension;
          }
        } else {
          if (height > maxDimension) {
            width = (width * maxDimension) / height;
            height = maxDimension;
          }
        }

        canvas.width = width;
        canvas.height = height;
        ctx?.drawImage(img, 0, 0, width, height);

        canvas.toBlob(
          (blob) => {
            if (blob) {
              resolve(blob);
            } else {
              reject(new Error('Failed to generate thumbnail'));
            }
          },
          format,
          quality
        );
      };

      img.onerror = () => reject(new Error('Failed to load image'));
      img.src = URL.createObjectURL(file);
    });
  }, [maxDimension, quality, format]);

  const generateVideoThumbnail = useCallback(async (file: File, seekTo: number = 1.0): Promise<Blob> => {
    return new Promise((resolve, reject) => {
      const video = document.createElement('video');
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');

      video.onloadedmetadata = () => {
        video.currentTime = Math.min(seekTo, video.duration);
      };

      video.onseeked = () => {
        // Calculate dimensions
        let { videoWidth: width, videoHeight: height } = video;
        if (width > height) {
          if (width > maxDimension) {
            height = (height * maxDimension) / width;
            width = maxDimension;
          }
        } else {
          if (height > maxDimension) {
            width = (width * maxDimension) / height;
            height = maxDimension;
          }
        }

        canvas.width = width;
        canvas.height = height;
        ctx?.drawImage(video, 0, 0, width, height);

        canvas.toBlob(
          (blob) => {
            if (blob) {
              resolve(blob);
            } else {
              reject(new Error('Failed to generate video thumbnail'));
            }
          },
          format,
          quality
        );
      };

      video.onerror = () => reject(new Error('Failed to load video'));
      video.src = URL.createObjectURL(file);
      video.load();
    });
  }, [maxDimension, quality, format]);

  return { generateThumbnail, generateVideoThumbnail };
};
```

#### Complete Upload Flow with Thumbnails

```typescript
const { uploadFiles, uploads } = useAssetUpload('my-store');
const { generateThumbnail, generateVideoThumbnail } = useThumbnailGenerator();

const handleFileSelect = async (event: React.ChangeEvent<HTMLInputElement>) => {
  const files = Array.from(event.target.files || []);

  try {
    // Step 1: Request upload URLs (API will include thumbnail descriptors)
    const response = await fetch(`/v1/businesses/my-store/assets/uploads`, {
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
    });

    const { uploads: descriptors } = await response.json();

    // Step 2: Upload each file + thumbnail
    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      const descriptor = descriptors[i];

      // Upload main file
      await uploadSimple(file, descriptor);

      // Upload thumbnail if provided
      if (descriptor.thumbnail) {
        let thumbnailBlob: Blob;

        if (file.type.startsWith('image/')) {
          thumbnailBlob = await generateThumbnail(file);
        } else if (file.type.startsWith('video/')) {
          thumbnailBlob = await generateVideoThumbnail(file);
        } else {
          continue; // No thumbnail needed
        }

        // Convert blob to file
        const thumbnailFile = new File([thumbnailBlob], descriptor.thumbnail.fileName, {
          type: 'image/jpeg'
        });

        // Upload thumbnail
        const formData = new FormData();
        formData.append('file', thumbnailFile);

        await fetch(descriptor.thumbnail.uploadUrl, {
          method: 'POST',
          body: formData
        });
      }
    }
  } catch (error) {
    console.error('Upload failed:', error);
  }
};
```

### Backend Thumbnail Handling

When generating upload URLs, the service checks if the file needs a thumbnail:

```go
// In internal/domain/asset/service.go

func (s *Service) GenerateUploadURLs(ctx context.Context, req *GenerateUploadURLsRequest) (*GenerateUploadURLsResponse, error) {
    descriptors := make([]*UploadDescriptor, len(req.Files))
    
    for i, fileReq := range req.Files {
        // Validate file type and size
        if !s.validator.IsAllowed(fileReq.FileName) {
            return nil, problem.BadRequest("file type not allowed")
        }
        
        category := s.validator.GetCategory(fileReq.FileName)
        maxSize := s.validator.GetMaxSize(fileReq.FileName)
        
        if fileReq.SizeBytes > maxSize {
            return nil, problem.BadRequest("file too large")
        }
        
        // Generate main upload descriptor
        descriptor := s.generateUpload(ctx, req, fileReq, category)
        
        // Generate thumbnail if needed
        if s.validator.NeedsThumbnail(fileReq.FileName) {
            thumbnail := s.generateThumbnailUpload(ctx, req, descriptor, fileReq.FileName)
            descriptor.Thumbnail = thumbnail
        }
        
        descriptors[i] = descriptor
    }
    
    return &GenerateUploadURLsResponse{Uploads: descriptors}, nil
}
```

### AssetReference with Thumbnails

When using assets in domain models, the AssetReference includes all thumbnail URLs:

```typescript
interface AssetReference {
  url: string;                    // CDN URL (primary, fast)
  originalUrl?: string;           // Storage URL (fallback)
  thumbnailUrl?: string;          // CDN thumbnail URL
  thumbnailOriginalUrl?: string;  // Storage thumbnail URL
  assetId?: string;              // For garbage collection
  metadata?: {
    altText?: string;
    caption?: string;
    width?: number;
    height?: number;
  };
}
```

**Display Pattern**:

```tsx
const ProductImage = ({ photo }: { photo: AssetReference }) => {
  return (
    <img
      src={photo.thumbnailUrl || photo.url}  // Use thumbnail for preview
      srcSet={`${photo.thumbnailUrl} 512w, ${photo.url} 2000w`}
      sizes="(max-width: 768px) 512px, 2000px"
      alt={photo.metadata?.altText || 'Product photo'}
      loading="lazy"
      onClick={() => openLightbox(photo.url)}  // Full size on click
    />
  );
};
```

---

## CDN Integration

### Overview

Kyora supports **generic HTTP CDN integration** for fast, edge-cached asset delivery. The system works with any CDN that supports HTTP origin pulls (CloudFront, Cloudflare, Fastly, BunnyCDN, etc.).

### CDN URL Structure

Each asset gets **four URLs**:

1. **`url`** (CDN URL) - Primary, fast, edge-cached
2. **`originalUrl`** (Storage URL) - Fallback, direct from S3/local
3. **`thumbnailUrl`** (CDN thumbnail) - Fast thumbnail delivery
4. **`thumbnailOriginalUrl`** (Storage thumbnail) - Fallback thumbnail

**Example**:

```json
{
  "url": "https://cdn.kyora.app/assets/my-store/ast_ABC123/photo.jpg",
  "originalUrl": "https://s3.amazonaws.com/kyora-assets/my-store/ast_ABC123/photo.jpg",
  "thumbnailUrl": "https://cdn.kyora.app/assets/my-store/ast_DEF456/photo_thumb.jpg",
  "thumbnailOriginalUrl": "https://s3.amazonaws.com/kyora-assets/my-store/ast_DEF456/photo_thumb.jpg"
}
```

### Configuration

Add CDN base URL to your config (`.kyora.yaml`):

```yaml
storage:
  cdn_base_url: "https://cdn.kyora.app"  # Your CDN domain

s3:
  endpoint: "https://s3.amazonaws.com"
  region: "us-east-1"
  bucket: "kyora-assets"
```

**Without CDN** (graceful degradation):

```yaml
storage:
  cdn_base_url: ""  # Empty = use storage URLs directly

# Asset URLs will be storage URLs:
# url: "https://s3.amazonaws.com/kyora-assets/..."
```

### CDN Setup Examples

#### AWS CloudFront

1. **Create CloudFront Distribution**:
   - Origin: Your S3 bucket (e.g., `kyora-assets.s3.amazonaws.com`)
   - Origin Protocol Policy: HTTPS Only
   - Viewer Protocol Policy: Redirect HTTP to HTTPS
   - Allowed HTTP Methods: GET, HEAD, OPTIONS
   - Cache Policy: CachingOptimized (recommended)

2. **CORS Configuration** (S3 bucket):

```json
{
  "CORSRules": [
    {
      "AllowedOrigins": ["*"],
      "AllowedMethods": ["GET", "HEAD"],
      "AllowedHeaders": ["*"],
      "ExposeHeaders": ["ETag", "Content-Length"],
      "MaxAgeSeconds": 3600
    }
  ]
}
```

3. **Update Kyora Config**:

```yaml
storage:
  cdn_base_url: "https://d1234567890abc.cloudfront.net"  # CloudFront domain
```

#### Cloudflare

1. **Add CNAME Record**:
   - Type: CNAME
   - Name: `cdn` (or any subdomain)
   - Target: Your S3 bucket endpoint
   - Proxy status: Proxied (orange cloud)

2. **Page Rules**:
   - URL: `cdn.yourdomain.com/assets/*`
   - Cache Level: Cache Everything
   - Edge Cache TTL: 1 year

3. **Update Kyora Config**:

```yaml
storage:
  cdn_base_url: "https://cdn.yourdomain.com"
```

#### BunnyCDN

1. **Create Pull Zone**:
   - Origin URL: Your S3 bucket HTTPS URL
   - Pull Zone Name: Choose unique name
   - Enable CORS headers

2. **Update Kyora Config**:

```yaml
storage:
  cdn_base_url: "https://yourpullzone.b-cdn.net"
```

### Cache Headers

Kyora serves assets with **immutable cache headers** for maximum CDN efficiency:

```http
HTTP/1.1 200 OK
Content-Type: image/jpeg
Cache-Control: public, max-age=31536000, immutable
CDN-Cache-Control: max-age=31536000
ETag: "d41d8cd98f00b204e9800998ecf8427e"
Last-Modified: Wed, 31 Dec 2025 10:00:00 GMT
```

**Cache Strategy**:

- **1-year cache**: `max-age=31536000` (365 days)
- **Immutable**: Browser never revalidates cached assets
- **CDN-specific**: `CDN-Cache-Control` for edge caching
- **ETag support**: Efficient cache validation if needed
- **Content-addressable**: Asset IDs in URLs ensure uniqueness

### CDN Implementation Details

The system generates CDN URLs by replacing the storage domain with the CDN domain:

```go
// internal/domain/asset/cdn.go

func GenerateCDNURL(storageURL, cdnBaseURL string) string {
    if cdnBaseURL == "" {
        return storageURL  // No CDN configured, use storage URL
    }
    
    parsedStorage, err := url.Parse(storageURL)
    if err != nil {
        return storageURL
    }
    
    parsedCDN, err := url.Parse(cdnBaseURL)
    if err != nil {
        return storageURL
    }
    
    // Replace storage host with CDN host, keep path
    parsedStorage.Scheme = parsedCDN.Scheme
    parsedStorage.Host = parsedCDN.Host
    
    return parsedStorage.String()
}
```

**Example Transformation**:

```
Storage URL:  https://s3.amazonaws.com/kyora-assets/my-store/ast_ABC/photo.jpg
CDN Base URL: https://cdn.kyora.app
Result:       https://cdn.kyora.app/kyora-assets/my-store/ast_ABC/photo.jpg
```

### Frontend CDN Usage

Always prefer CDN URLs, with fallback to original URLs:

```typescript
const AssetImage = ({ asset }: { asset: AssetReference }) => {
  const [src, setSrc] = useState(asset.url);  // Start with CDN URL
  const [thumbnailSrc, setThumbnailSrc] = useState(asset.thumbnailUrl);

  return (
    <img
      src={thumbnailSrc || src}
      srcSet={`${thumbnailSrc} 512w, ${src} 2000w`}
      sizes="(max-width: 768px) 512px, 2000px"
      alt={asset.metadata?.altText}
      onError={() => {
        // Fallback to storage URL if CDN fails
        if (src === asset.url && asset.originalUrl) {
          setSrc(asset.originalUrl);
        }
        if (thumbnailSrc === asset.thumbnailUrl && asset.thumbnailOriginalUrl) {
          setThumbnailSrc(asset.thumbnailOriginalUrl);
        }
      }}
    />
  );
};
```

### CDN Performance Benefits

| Metric | Without CDN | With CDN | Improvement |
|--------|-------------|----------|-------------|
| **TTFB** (Time to First Byte) | 200-500ms | 10-50ms | **10x faster** |
| **Global Latency** | High (single region) | Low (edge cache) | **Global** |
| **Bandwidth Costs** | S3 egress fees | CDN costs (cheaper) | **50-70% savings** |
| **Origin Load** | Every request | Only cache misses | **99% reduction** |
| **Scalability** | Limited by origin | Unlimited (edge) | **Infinite** |

---

## AssetReference Type

## AssetReference Type

The `AssetReference` type is the cornerstone of asset management. It's defined in `internal/platform/types/asset/asset.go` to avoid circular dependencies.

### Complete Structure

```go
// AssetReference represents a reference to an asset with CDN and thumbnail support
type AssetReference struct {
    // Primary URLs (CDN-optimized)
    URL              string         `json:"url" binding:"required,max=2048"`
    OriginalURL      string         `json:"originalUrl,omitempty" binding:"max=2048"`
    
    // Thumbnail URLs (for images/videos)
    ThumbnailURL     string         `json:"thumbnailUrl,omitempty" binding:"max=2048"`
    ThumbnailOriginalURL string     `json:"thumbnailOriginalUrl,omitempty" binding:"max=2048"`
    
    // Asset tracking
    AssetID          *string        `json:"assetId,omitempty"`
    
    // Rich metadata
    Metadata         *AssetMetadata `json:"metadata,omitempty"`
}

// AssetMetadata provides optional semantic information
type AssetMetadata struct {
    AltText string `json:"altText,omitempty"`  // Accessibility text
    Caption string `json:"caption,omitempty"`  // Display caption
    Width   *int   `json:"width,omitempty"`    // Image width in pixels
    Height  *int   `json:"height,omitempty"`   // Image height in pixels
}
```

### Field Descriptions

| Field | Required | Description | Example |
|-------|----------|-------------|---------|
| `url` | **Yes** | CDN URL (primary, fast) | `https://cdn.kyora.app/assets/...` |
| `originalUrl` | No | Storage URL (fallback) | `https://s3.amazonaws.com/bucket/...` |
| `thumbnailUrl` | No | CDN thumbnail URL | `https://cdn.kyora.app/assets/...thumb.jpg` |
| `thumbnailOriginalUrl` | No | Storage thumbnail URL | `https://s3.amazonaws.com/bucket/...thumb.jpg` |
| `assetId` | No | Asset ID for GC | `ast_2ZqJKxY3n8LmW9B5VrDfGpN4Tc1` |
| `metadata` | No | Accessibility & display data | `{ altText: "Product photo" }` |

### Storage Format

AssetReference is stored as **JSONB in PostgreSQL**, enabling flexible queries and updates:

```sql
-- Business logo (single optional asset)
CREATE TABLE businesses (
    ...
    logo JSONB,
    ...
);

-- Product photos (array of assets)
CREATE TABLE products (
    ...
    photos JSONB NOT NULL DEFAULT '[]'::jsonb,
    ...
);
```

### Usage Examples

#### Single Optional Asset (Business Logo)

```go
type Business struct {
    Logo *asset.AssetReference `gorm:"column:logo;type:jsonb" json:"logo,omitempty"`
}

// Create business with logo
business := &Business{
    Name: "My Store",
    Logo: &asset.AssetReference{
        URL:         "https://cdn.kyora.app/assets/my-store/ast_ABC123/logo.png",
        OriginalURL: "https://s3.amazonaws.com/bucket/my-store/ast_ABC123/logo.png",
        AssetID:     ptr("ast_ABC123"),
        Metadata: &asset.AssetMetadata{
            AltText: "My Store Logo",
        },
    },
}

// Remove logo
business.Logo = nil
```

#### Array of Required Assets (Product Photos)

```go
type Product struct {
    Photos inventory.AssetReferenceList `gorm:"column:photos;type:jsonb;not null;default:'[]'" json:"photos"`
}

// Create product with photos
product := &Product{
    Name: "Premium T-Shirt",
    Photos: inventory.AssetReferenceList{
        {
            URL:              "https://cdn.kyora.app/assets/my-store/ast_ABC123/photo1.jpg",
            OriginalURL:      "https://s3.amazonaws.com/bucket/my-store/ast_ABC123/photo1.jpg",
            ThumbnailURL:     "https://cdn.kyora.app/assets/my-store/ast_DEF456/photo1_thumb.jpg",
            ThumbnailOriginalURL: "https://s3.amazonaws.com/bucket/my-store/ast_DEF456/photo1_thumb.jpg",
            AssetID:          ptr("ast_ABC123"),
            Metadata: &asset.AssetMetadata{
                AltText: "Premium T-Shirt Front View",
                Width:   ptr(2000),
                Height:  ptr(2000),
            },
        },
        {
            URL:              "https://cdn.kyora.app/assets/my-store/ast_GHI789/photo2.jpg",
            OriginalURL:      "https://s3.amazonaws.com/bucket/my-store/ast_GHI789/photo2.jpg",
            ThumbnailURL:     "https://cdn.kyora.app/assets/my-store/ast_JKL012/photo2_thumb.jpg",
            ThumbnailOriginalURL: "https://s3.amazonaws.com/bucket/my-store/ast_JKL012/photo2_thumb.jpg",
            AssetID:          ptr("ast_GHI789"),
            Metadata: &asset.AssetMetadata{
                AltText: "Premium T-Shirt Back View",
                Width:   ptr(2000),
                Height:  ptr(2000),
            },
        },
    },
}
```

---

## Upload Flows

### Flow 1: S3 Multipart Upload with Thumbnail (Production)

**Use Case**: Large images/videos in production with CDN delivery

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
