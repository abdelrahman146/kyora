---
description: Asset Upload (backend contract + frontend integration flow)
applyTo: "**/*"
---

# Asset Upload System

**SSOT Hierarchy:**

- Parent: copilot-instructions.md
- Related Backend: backend-core.instructions.md (architecture), backend-testing.instructions.md (testing)
- Related Frontend: forms.instructions.md (FileUploadField, ImageUploadField), ky.instructions.md (HTTP client)

**When to Read:**

- Implementing file uploads
- Business logo updates
- Product/variant photo management
- Understanding asset GC flow
- Backend asset API contract

---

## Why this exists

Kyora’s asset system is intentionally **simple and frontend-driven**:

- Backend **never** receives file bytes in production. Frontend uploads directly to blob storage using pre-signed URLs.
- Backend creates an **asset record immediately** (with `ast...` id), then the client uploads bytes, then optionally completes multipart.
- Domains don’t store “asset objects”. They store **`AssetReference`** (JSONB) on the owning entity (business/product/variant).

This doc describes the **canonical contract** and the **only supported flow** that frontend/mobile/web should implement.

## Core invariants (don’t fight these)

- **All assets are public by design.** Treat the returned `url` as public.
- **No pending states.** There is no “draft asset” vs “final asset”. The record exists as soon as you request upload URLs.
- **Asset ownership is business-scoped.** Upload endpoints are under `/v1/businesses/{businessDescriptor}` and are protected by workspace + business middleware.
- **Garbage collection depends on `assetId` being referenced.** If you store only URLs (without `assetId`), the system can’t reliably determine ownership and cleanup.

## System summary (keep this mental model)

- **Pre-signed URL approach:** clients upload directly to storage; backend does not proxy bytes.
- **Client-side control:** frontend handles chunking, retries, progress, and resuming multipart.
- **Provider agnostic:** works with S3-compatible storage in production and local filesystem in development.
- **CDN-first URLs:** backend returns `cdnUrl` and `publicUrl` to render (CDN should be configured in prod, but public/storage URLs still work as fallback).
- **Immutable caching (local provider):** `GET /v1/public/assets/:assetId` sets long-lived immutable cache headers only when streaming from local disk.
- **Server-side validation:** backend validates file type and size per category; uploads are business-scoped and require JWT.

### Supported categories (defaults; configurable)

The backend classifies uploads by `contentType` into these categories:

- **Images** (e.g., jpg/jpeg/png/webp/gif/heic/heif) — default max **10 MB** — thumbnails supported
- **Videos** (e.g., mp4/mov/avi/mkv/webm) — default max **100 MB** — thumbnails supported
- **Audio** (e.g., mp3/wav/ogg/m4a/aac) — default max **20 MB**
- **Documents** (e.g., pdf/doc/docx/txt/rtf/odt) — default max **10 MB**
- **Compressed** (e.g., zip/tar/gz/rar/7z) — default max **50 MB**

These defaults are configurable via:

- `uploads.allowed_extensions` (per category)
- `uploads.max_size_bytes` (per category)

### Thumbnails (images/videos)

- Backend decides if a file needs a thumbnail and may return a nested thumbnail upload descriptor.
- Client generates the thumbnail bytes (image: canvas resize/compress; video: extracted frame) and uploads them like any other file.
- Thumbnails are tightly coupled to the main asset for garbage collection.

## Types you must use

### `AssetReference` (stored on other domains)

Source of truth: `backend/internal/platform/types/asset/asset.go`

Example shape:

```json
{
  "url": "https://cdn.../businesses/.../ast_xxx/filename.jpg",
  "originalUrl": "https://s3...",
  "thumbnailUrl": "https://cdn.../thumb.jpg",
  "thumbnailOriginalUrl": "https://s3...",
  "assetId": "ast_xxx",
  "metadata": {
    "altText": "...",
    "caption": "...",
    "width": 123,
    "height": 456
  }
}
```

Rules:

- `url` is **required**.
- `assetId` is **required in practice** for GC and future migrations. Backend validation allows it to be optional for backward compatibility, but all new code must send it.
  - Portal-web TypeScript types define `assetId` as **required** (non-optional).
  - Backend Go validation uses `binding:"omitempty"` but this is for flexibility only.
- `originalUrl` / `thumbnailOriginalUrl` are optional fallbacks; today the upload API mainly returns public/CDN URLs, so most clients will populate only `url` (+ `thumbnailUrl` if provided by their thumbnail workflow).

### Where `AssetReference` is used today

- Business logo:
  - `Business.Logo` is `*AssetReference` (JSONB)
  - Request fields: `CreateBusinessInput.Logo`, `UpdateBusinessInput.Logo`
- Product photos:
  - `Product.Photos` is `AssetReferenceList` (JSONB list)
  - Request fields: `CreateProductRequest.Photos`, `UpdateProductRequest.Photos`
- Variant photos:
  - `Variant.Photos` is `AssetReferenceList` (JSONB list)
  - Request fields: `CreateVariantRequest.Photos`, `UpdateVariantRequest.Photos`
- Product + variants in one go:
  - `CreateProductWithVariantsRequest` supports `product.photos` and `variants[].photos`

Important limits:

- Product photos: max **10**
- Variant photos: max **10**
- Upload URL generation: max **50 files per request**

## Backend endpoints (contract)

### 1) Generate upload descriptors (always the first call)

`POST /v1/businesses/{businessDescriptor}/assets/uploads`

AuthZ:

- Requires `Authorization: Bearer <JWT>`
- Requires actor permission: `ActionManage` on `ResourceInventory`

Request:

```json
{
  "files": [
    {
      "fileName": "my-photo.jpg",
      "contentType": "image/jpeg",
      "sizeBytes": 123456
    }
  ]
}
```

Response:

- Returns `uploads[]` where each upload has:
  - `assetId`
  - `method` and either:
    - S3 multipart: `partSize`, `totalParts`, `partUrls` (array of `{partNumber,url}`), `uploadId`, `method: "PUT"`
    - Local dev: `url`, `headers`, `method: "POST"`
  - `publicUrl` and `cdnUrl` (treat these as display URLs)
  - Optional `thumbnail` descriptor when the backend decides a thumbnail is needed

Important contract details:

- `partUrls` is an array of objects: `[{ "partNumber": 1, "url": "..." }, ...]` (not a string array).
- Part URLs are currently issued with a long expiry (24h) to support large uploads.
- Thumbnail uploads are:
  - Local provider: a single `POST` to an internal URL (same pattern as main local upload).
  - S3 provider: a single pre-signed `PUT` to a URL + required headers (not multipart).

### 2) Upload file bytes (client-to-storage)

#### S3 multipart (production)

- Chunk the file into `partSize` bytes.
- For each `partUrls[i].url` upload using `PUT`.
- Collect the `ETag` header from each part upload response.
  - Browser note: your bucket CORS must expose `ETag`.

#### Local provider (development)

- Upload the full file in one request:
  - `POST {upload.url}`
  - Set `Content-Type` from `upload.headers`
  - Body is raw bytes

### 3) Complete multipart uploads (S3 only)

`POST /v1/businesses/{businessDescriptor}/assets/uploads/{assetId}/complete`

Request:

```json
{
  "parts": [
    { "partNumber": 1, "etag": "..." },
    { "partNumber": 2, "etag": "..." }
  ]
}
```

Rules:

- `partNumber` is **1-based**.
- `etag` must be sent **without quotes** (strip `"` if your HTTP client returns them).

### 4) Public serving (fallback)

`GET /v1/public/assets/{assetId}`

- No auth.
- Local provider streams bytes with immutable cache headers.
- S3 provider redirects to the public URL.

In most frontends you should **not** call this endpoint directly; just render the stored `url`/`cdnUrl`. This endpoint exists mainly for local dev and as a stable public retrieval mechanism.

## The supported frontend flow (web + mobile)

### A) Business creation + logo

You can create a business with or without a logo. However, **uploading requires an existing business** because upload endpoints are business-scoped.

Recommended sequence:

1. Create the business (no logo):
   - `POST /v1/businesses` (or onboarding flow if you’re inside onboarding)
2. Upload logo:
   - `POST /v1/businesses/{businessDescriptor}/assets/uploads` with the logo file
   - Upload bytes (and call `complete` if multipart)
3. Patch business logo:
   - `PATCH /v1/businesses/{businessDescriptor}` with:

```json
{
  "logo": {
    "url": "<upload.cdnUrl>",
    "assetId": "<upload.assetId>",
    "metadata": { "altText": "Business logo" }
  }
}
```

Why this order:

- Prevents “broken image” states (logo referenced before bytes exist)
- Ensures the asset is reachable and owned before saving the reference

### B) Categories → products → variants + photos

Inventory entities are business-scoped.

Canonical creation options:

- **Option 1 (common):** Create product → create variants
- **Option 2:** Create product with variants atomically via `/inventory/products/with-variants`

In both options, photo flow is identical:

1. Upload product/variant photos first.
2. Build `AssetReference` objects from the upload descriptors.
3. Send them in `photos` when creating/updating the product/variant.

#### Option 1: product then variants

- Create product:
  - `POST /v1/businesses/{businessDescriptor}/inventory/products`
  - Include `photos: AssetReference[]` if you already uploaded them
- Create variant:
  - `POST /v1/businesses/{businessDescriptor}/inventory/variants`
  - Include `photos: AssetReference[]` if you already uploaded them

#### Option 2: product with variants (single request)

- `POST /v1/businesses/{businessDescriptor}/inventory/products/with-variants`

This request supports photos at both levels:

```json
{
  "product": {
    "name": "...",
    "description": "...",
    "categoryId": "cat_...",
    "photos": [{ "url": "...", "assetId": "ast_..." }]
  },
  "variants": [
    {
      "code": "RED",
      "photos": [{ "url": "...", "assetId": "ast_..." }],
      "costPrice": "10",
      "salePrice": "15",
      "stockQuantity": 5,
      "stockQuantityAlert": 1
    }
  ]
}
```

### When to upload vs when to save references

- **Upload first, then persist references**.
- “Usable” here means:
  - the bytes were uploaded successfully, and
  - for S3 multipart uploads: `/assets/uploads/{assetId}/complete` succeeded.

Safer recommended patterns:

- **Best for data integrity (recommended):**
  - Upload files first.
  - Only after upload completion, create/update business/product/variant with `AssetReference`.

- **Best for UX (still safe):**
  - Create the product/variant without photos.
  - Upload in background.
  - `PATCH` the entity’s `photos` only after upload completion.
  - Keep a client-side “pending uploads” record (local DB/storage) so you can resume after app restarts.

Do **not** persist an `AssetReference` to the backend just because you already received `cdnUrl/publicUrl` from `GenerateUploadURLs`. Those URLs are returned before the bytes exist.

## Completion tracking (current reality)

- The `complete` endpoint is required to assemble multipart uploads on S3.
- The backend does not currently persist a "completed" flag in the `uploaded_assets` table.
  - `MarkUploadComplete` is a no-op today (returns `nil` immediately).
  - Comment in code: "For now, just return nil - this will be implemented in the GC rewrite"
- Do not build UX that depends on the server knowing multipart completion state client-side tracking is the source of truth.

## Multipart upload failures and resume

Kyora uses **S3 multipart** for non-local providers. The backend currently returns all required part URLs up-front.

### What you should persist client-side while uploading

To support resume, persist this state locally (SQLite/Room/CoreData/IndexedDB/localStorage—whatever fits the platform):

- `businessDescriptor`
- `assetId`
- `uploadId`
- `partSize`
- `totalParts`
- `partUrls` (array of `{ partNumber, url }`, or enough info to map partNumber → URL)
- `completedParts`: map of `partNumber -> etag`
- (optional) a stable reference to the local file so the app can continue reading bytes

This is **client state only**. Kyora intentionally does not have a server-side “pending upload” concept.

### Retry and resume rules (S3 multipart)

- **Retry individual parts aggressively.** A part upload (`PUT` to a presigned URL) is safe to retry.
  - If you don’t know whether a part succeeded, you can re-upload that part number; S3 will use the last uploaded part for that number.
- **Only call `complete` when you have ETags for all parts**.
- **If the app restarts mid-upload:**
  - Load your persisted upload state.
  - Upload any missing parts.
  - Call `complete`.

### Handling common multipart failure modes

- **Network error / timeout while uploading a part:**
  - Retry the same part URL with exponential backoff.
  - Continue uploading remaining parts; you don’t need strict sequential uploads.

- **Missing `ETag` in browser:**
  - Fix bucket CORS to expose `ETag`.
  - Without ETags you cannot complete multipart.

- **403/expired presigned URLs (often after long pause):**
  - Presigned part URLs expire.
  - Kyora does not currently expose an endpoint to “refresh part URLs for an existing assetId/uploadId”.
  - The correct recovery is to **restart the upload**:
    1. Call `POST /v1/businesses/{businessDescriptor}/assets/uploads` again to get a new `assetId` and new URLs
    2. Upload again
    3. Persist references using the new `assetId`
  - The old asset will remain unreferenced and should be cleaned up by GC.

- **`complete` call fails:**
  - If it’s a transient 5xx/network error, retry `complete`.
  - If it fails due to bad/missing ETags, re-upload the missing parts, then call `complete` again.

### Local provider failure handling (dev)

- Local uploads are a single `POST` of full bytes.
- Retry the whole request on transient failures.
- Only persist references after you get a successful response from the local upload endpoint.

## Thumbnail behavior (important for UX)

- Backend decides whether a file "needs thumbnail" based on file category (images and videos).
- If a thumbnail is needed, the upload descriptor includes a nested `thumbnail` descriptor.
- The client is responsible for:
  - generating a thumbnail (canvas for images; extracted frame for video),
  - uploading thumbnail bytes with the thumbnail descriptor,
  - uploading the thumbnail using the returned `method`/`url`/`headers`.
- **Thumbnail uploads today:**
  - Local provider: single `POST` to internal URL (same pattern as main local upload).
  - S3 provider: single pre-signed `PUT` to a URL + required headers (not multipart).
- If you don't store `thumbnailUrl` in `AssetReference`, your UI will still work, but may load heavier media.

## See Also

- **forms.instructions.md** — FileUploadField and ImageUploadField components (high-level integration)
- **ky.instructions.md** — HTTP client patterns for upload API calls
- **backend-core.instructions.md** — Backend architecture for asset storage and GC
- **ui-implementation.instructions.md** — Progress indicators, error states for uploads

---

## Common integration pitfalls

- Wrong `contentType` or extension mismatch → “file type not allowed”.
- Browser uploads without S3 CORS exposing `ETag` → cannot complete multipart.
- Forgetting to strip quotes from `ETag` → completion fails.
- Saving references without `assetId` → GC cannot track ownership.
- Uploading under the wrong business descriptor → assets are scoped incorrectly.

## Notes for AI agents making changes

- Don’t invent new “asset states” or a separate “pending upload” table.
- If a new domain needs media, the correct pattern is:
  - store `AssetReference` (or a list) as JSONB on the domain model,
  - validate limits via request binding,
  - reuse the existing upload endpoints (don’t create purpose-specific upload routes).
