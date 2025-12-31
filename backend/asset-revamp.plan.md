Plan: Revamp Asset Management to Simple Direct Upload Flow

Complete overhaul of asset management from complex purpose-driven tracking to a streamlined pre-signed URL approach with resumable multipart uploads, automatic assetId generation, and garbage collection—removing all legacy code and updating inventory/business domains.

Steps

1. Design New Asset API with Pre-Signed URL Generation
Create service.go methods:
GenerateUploadURLs(ctx, actor, workspace, request) - single endpoint supporting multiple files, single file, chunked uploads, and resumable multipart uploads
Always generate and return assetId for every file in request
For S3/cloud: generate S3 multipart presigned URLs with uploadId, part numbers, and configurable part size for client-side chunking (resumable)
For local dev: return simple POST endpoint URLs that write entire file to configurable local directory (no chunking—keep it simple)
Return upload descriptors (method, URL, headers, uploadId for S3, partSize) + assetId for each file
Remove purpose enum, status enum, visibility enum, idempotency entirely; treat all uploads uniformly

2. Simplify Asset Model and Storage with Typed Metadata
Update model.go and storage.go:
Remove purpose, visibility, status, upload_expires_at, completed_at, request_hash, idempotency_key fields completely
Keep id, workspace_id, business_id, created_by_user_id, object_key, public_url, content_type, size_bytes, local_file_path
Add S3 multipart fields: upload_id (nullable), completed_parts (JSONB array of ETags), is_multipart (boolean), total_parts (int)
Define AssetMetadata struct: {AltText string, Caption string, Width *int, Height *int} (all optional, stored as JSONB)
Define AssetReference struct: {URL string, AssetID *string, Metadata *AssetMetadata}
Add storage methods: TrackUploadPart, GetUploadProgress, MarkUploadComplete, ListUnusedAssets (simplified orphan detection without Postgres JSONB operators)
Remove all purpose-based queries, idempotency queries, Postgres-specific reference detection queries, and IsReferenced method

3. Update Configuration for Asset Management
Modify config.go and .kyora.yaml.example:
Add storage.local_path (default: ./tmp/assets) - local file storage directory
Add storage.multipart_part_size_mb (default: 10) - S3 multipart part size in megabytes
Add inventory.max_photos_per_product (default: 10) - max photos per product/variant
Remove deprecated asset-related config keys from example file

4. Update Inventory Domain to Accept Typed Asset Metadata
Modify model.go and service.go:
Change Photos field from PhotoURLList ([]string) to []asset.AssetReference using strongly-typed struct
Update CreateProduct, UpdateProduct, CreateVariant, UpdateVariant to accept new AssetReference format with AssetMetadata
Validate photo count against inventory.max_photos_per_product config (default: 10)
No validation of assetId existence or URL consistency (optional field for GC tracking only, allow external URLs)
Delete PhotoURLList type and all related helper methods completely

5. Update Business Domain to Accept Typed Asset Metadata
Modify model.go and service.go:
Change LogoURL from string to *asset.AssetReference (nullable) using strongly-typed struct from asset domain
Update CreateBusiness, UpdateBusiness to accept new AssetReference format with optional AssetMetadata
No validation of assetId existence or URL consistency (keep loose for flexibility with external URLs)
Remove trimming/validation logic; store as provided
Delete all logo-specific asset-handling code and related helpers

6. Update Blob Platform for Multipart and Configuration
Enhance backend/internal/platform/blob/s3.go and backend/internal/platform/blob/provider.go:
Add CreateMultipartUpload(ctx, key, contentType) (uploadId, error) to Provider interface
Add PresignMultipartPart(ctx, key, uploadId, partNumber, expiry) (url, error) to generate presigned URLs for individual parts
Add CompleteMultipartUpload(ctx, key, uploadId, parts) error to finalize multipart upload with ETags
Add AbortMultipartUpload(ctx, key, uploadId) error for cleanup (optional - leave abandoned uploads for S3 lifecycle policies)
Implement in S3 provider using AWS SDK v2 multipart operations
Local provider returns errors for multipart methods (not supported for simplicity)


7. Consolidate Asset HTTP Endpoints with Cache Headers
Rewrite handler_http.go and update routes.go:
Replace all 9 purpose-specific endpoints with 2 universal endpoints:
POST /businesses/:desc/assets/uploads - generate pre-signed URLs (multiple files, S3 multipart for cloud, simple POST for local)
GET /v1/public/assets/:assetId - serve public assets with cache headers (Cache-Control: max-age=3600, ETag based on file hash)
Remove CompleteUpload, StoreLocalContent, DeleteAsset, and all purpose-specific handlers entirely
For local provider: add simple internal handler that accepts entire file POST and writes to configured storage.local_path directory (no chunking)
Use standard workspace/business middleware; delete backend/internal/domain/asset/middleware_http.go completely
Remove all asset-related routes from routes.go except the 2 new ones

8. Add Comprehensive API Documentation with OpenAPI/Swagger
Update handler_http.go with detailed godoc comments:
Add OpenAPI/Swagger annotations for POST /businesses/:desc/assets/uploads showing request/response schemas, multipart flow, S3 vs local differences
Document multipart upload flow: request structure (files array with name/size/contentType), response structure (uploadDescriptors with method/URL/headers/uploadId/partSize/assetId)
Include examples for single file, multiple files, and large file chunking scenarios
Document client-side responsibilities: chunking files, uploading parts, collecting ETags, assembling completion request
Add inline comments explaining S3 multipart mechanics (5MB min part size, max 10,000 parts, ETag collection)


9. Create Asset Upload Guide for AI Agents
Create backend/docs/ASSET_UPLOAD_GUIDE.md:
Explain the complete upload flow for both S3 multipart and local simple uploads
Provide step-by-step instructions for AI agents/clients: request pre-signed URLs → chunk files (if S3) → upload parts → collect ETags → reference assets in product/business APIs
Include code examples for chunking files, calculating part numbers, uploading with presigned URLs, handling errors/retries
Document how to use AssetReference in product photos and business logo fields
Explain garbage collection behavior and how assetId enables automatic cleanup
Target technical audience (developers, AI coding agents) with clear, actionable guidance

10. Revamp Garbage Collection for Simplified Detection
Rewrite backend/internal/domain/asset/service_gc.go and update assets_gc.go:
Remove "pending expired" phase (no status field anymore)
GC logic: find all assets where id is NOT present in any business.logo.assetId OR products.photos[].assetId OR variants.photos[].assetId
Use generic SQL/GORM queries with JSON extraction functions instead of Postgres-specific JSONB operators for database portability
Delete blob objects, local files, and DB records for unused assets older than configurable threshold (default: 24 hours)
Keep dry-run support and summary reporting
Update CLI command to remove irrelevant flags (--pending-limit, --orphan-min-age) and add --min-age (default: 24h)

11. Update All Test Cases for Asset Domain
Rewrite or remove test files related to assets:
Delete test files testing old endpoints: any E2E tests in e2e that test asset upload/complete/delete flows with purpose-based endpoints
Create new E2E tests for POST /businesses/:desc/assets/uploads covering: single file, multiple files, S3 multipart response structure, local provider response structure, error cases
Create new E2E tests for GET /v1/public/assets/:assetId validating local provider file serving, cache headers, ETag handling, 404 for non-existent assets
Update inventory E2E tests in e2e to use new AssetReference format (with URL, AssetID, Metadata) instead of plain URL strings
Update business E2E tests to use new AssetReference format for logo (nullable)
Create GC command tests for new orphan detection logic and dry-run behavior
Update test helpers in testutils to support creating test assets with new schema


12. Remove All Legacy Asset Code
Delete deprecated files and code:
Remove error definitions in errors.go for removed operations (upload not complete, fingerprint mismatch, asset in use, rate limit exceeded, etc.)
Remove rate limiting logic from service layer (no longer needed without idempotency)
Remove all purpose-specific constants, enums, and status enums from model.go
Clean up storage.go removing FindByBusinessAndIdempotencyKey, ListExpiredPending, ListReadyOrphans methods
Delete middleware file entirely: backend/internal/domain/asset/middleware_http.go
Remove all legacy routes and handlers from routes.go


Further Considerations
1. S3 Multipart Completion API: Should we provide a backend endpoint to complete multipart uploads (receiving ETags from client), or require clients to call S3's CompleteMultipartUpload API directly? Option A: Backend endpoint (simpler for clients, we track completion in DB). Option B: Direct S3 call (one less API call, less backend load). Recommendation: Option A - add POST /businesses/:desc/assets/uploads/:assetId/complete endpoint that accepts ETags array and calls S3 CompleteMultipartUpload, then updates DB with final URL.

2. Local Provider Upload Size Limits: Should local uploads enforce max file size limits to prevent disk exhaustion? Recommendation: Yes, use same uploads.max_bytes config (default 5MB) for local provider, return 413 Payload Too Large if exceeded.

3. Asset Cleanup Age Threshold: 24 hours for orphan cleanup might be too aggressive if clients upload assets hours before creating products. Should we increase default? Recommendation: Change default to 48 hours (2 days) to provide comfortable buffer for legitimate multi-step workflows.