package asset

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
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
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/abdelrahman146/kyora/internal/platform/utils/throttle"
	"github.com/spf13/viper"
)

type Service struct {
	storage *Storage
	atomic  *database.AtomicProcess
	blob    blob.Provider

	localDir string
}

type NewServiceParams struct {
	Storage  *Storage
	Atomic   *database.AtomicProcess
	Blob     blob.Provider
	LocalDir string
}

func NewService(storage *Storage, atomic *database.AtomicProcess, provider blob.Provider) *Service {
	return NewServiceWithParams(NewServiceParams{Storage: storage, Atomic: atomic, Blob: provider})
}

func NewServiceWithParams(p NewServiceParams) *Service {
	localDir := strings.TrimSpace(p.LocalDir)
	if localDir == "" {
		localDir = filepath.Join("tmp", "assets")
	}
	return &Service{storage: p.Storage, atomic: p.Atomic, blob: p.Blob, localDir: localDir}
}

func (s *Service) CreateUpload(ctx context.Context, actor *account.User, biz *business.Business, purpose Purpose, req *CreateUploadRequest) (*CreateUploadResponse, error) {
	if actor == nil {
		return nil, problem.Unauthorized("unauthorized")
	}
	if biz == nil {
		return nil, problem.BadRequest("business is required")
	}
	if req == nil {
		return nil, problem.BadRequest("request is required")
	}

	// Abuse protection: per-business+actor token bucket.
	if !throttle.Allow(s.storage.Cache(), fmt.Sprintf("rl:asset:upload:%s:%s:%s", biz.ID, actor.ID, purpose), time.Minute, 60, 250*time.Millisecond) {
		return nil, ErrRateLimited()
	}

	fileName := strings.TrimSpace(req.FileName)
	contentType := strings.TrimSpace(req.ContentType)
	size := req.SizeBytes
	if fileName == "" {
		return nil, problem.BadRequest("fileName is required").With("field", "fileName")
	}
	if contentType == "" {
		return nil, problem.BadRequest("contentType is required").With("field", "contentType")
	}
	if size <= 0 {
		return nil, problem.BadRequest("sizeBytes must be > 0").With("field", "sizeBytes")
	}

	maxBytes := viper.GetInt64(config.UploadsMaxBytes)
	if maxBytes > 0 && size > maxBytes {
		return nil, problem.BadRequest("file is too large").With("field", "sizeBytes").With("maxBytes", maxBytes)
	}

	// Only allow image uploads for current use-cases.
	if !isAllowedImageContentType(contentType) {
		return nil, ErrUploadNotAllowed("contentType", "only image uploads are supported")
	}

	idem, err := NormalizeIdempotencyKey(req.IdempotencyKey)
	if err != nil {
		return nil, err
	}

	visibility := VisibilityPublic
	fingerprint := RequestFingerprint(purpose, visibility, fileName, contentType, size)

	// Idempotency fast-path.
	if idem != "" {
		if existing, gerr := s.storage.FindByBusinessAndIdempotencyKey(ctx, biz.ID, idem); gerr == nil && existing != nil {
			if existing.RequestHash != fingerprint {
				return nil, ErrIdempotencyConflict()
			}
			resp := &CreateUploadResponse{AssetID: existing.ID, PublicURL: existing.PublicURL, Status: existing.Status}
			if existing.Status == StatusPending {
				u, uerr := s.buildUploadDescriptor(ctx, biz, existing)
				if uerr != nil {
					return nil, uerr
				}
				resp.Upload = u
				if existing.UploadExpiresAt != nil {
					resp.ExpiresAt = existing.UploadExpiresAt
				}
			}
			return resp, nil
		}
	}

	out := &CreateUploadResponse{}
	err = s.atomic.Exec(ctx, func(tctx context.Context) error {
		a := &Asset{
			WorkspaceID:     biz.WorkspaceID,
			BusinessID:      biz.ID,
			CreatedByUserID: actor.ID,
			Purpose:         purpose,
			Visibility:      visibility,
			Status:          StatusPending,
			ContentType:     contentType,
			SizeBytes:       size,
			IdempotencyKey:  idem,
			RequestHash:     fingerprint,
		}
		// Pre-generate ID so we can use it in the object key.
		if a.ID == "" {
			a.ID = id.KsuidWithPrefix(AssetPrefix)
		}
		a.ObjectKey = s.buildObjectKey(biz.ID, a.ID, fileName)

		if err := s.storage.Create(tctx, a); err != nil {
			if database.IsUniqueViolation(err) && idem != "" {
				// Another request won the race. Re-fetch and return.
				existing, gerr := s.storage.FindByBusinessAndIdempotencyKey(tctx, biz.ID, idem)
				if gerr != nil {
					return gerr
				}
				if existing.RequestHash != fingerprint {
					return ErrIdempotencyConflict()
				}
				out.AssetID = existing.ID
				out.PublicURL = existing.PublicURL
				out.Status = existing.Status
				if existing.Status == StatusPending {
					u, uerr := s.buildUploadDescriptor(tctx, biz, existing)
					if uerr != nil {
						return uerr
					}
					out.Upload = u
					out.ExpiresAt = existing.UploadExpiresAt
				}
				return nil
			}
			return err
		}

		u, uerr := s.buildUploadDescriptor(tctx, biz, a)
		if uerr != nil {
			return uerr
		}
		a.PublicURL = s.buildPublicURL(biz, a)
		if u == nil {
			return problem.InternalError().With("asset", "upload_descriptor_missing")
		}

		// Persist expiry for idempotent replay.
		if u.URL != "" {
			exp := time.Now().UTC().Add(10 * time.Minute)
			a.UploadExpiresAt = &exp
			_ = s.storage.Update(tctx, a)
			out.ExpiresAt = &exp
		}

		out.AssetID = a.ID
		out.Upload = u
		out.PublicURL = a.PublicURL
		out.Status = a.Status
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) CompleteUpload(ctx context.Context, actor *account.User, biz *business.Business, assetID string, expectedPurpose Purpose) (*CompleteUploadResponse, error) {
	if actor == nil {
		return nil, problem.Unauthorized("unauthorized")
	}
	if biz == nil {
		return nil, problem.BadRequest("business is required")
	}
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return nil, problem.BadRequest("assetId is required")
	}

	// Abuse protection: completion calls are cheap but can be spammed.
	if !throttle.Allow(s.storage.Cache(), fmt.Sprintf("rl:asset:complete:%s:%s", biz.ID, actor.ID), time.Minute, 120, 150*time.Millisecond) {
		return nil, ErrRateLimited()
	}

	a, err := s.storage.GetByID(ctx, biz.ID, assetID)
	if err != nil || a == nil {
		return nil, ErrAssetNotFound(assetID, err)
	}
	if a.Purpose != expectedPurpose {
		return nil, problem.Forbidden("invalid asset purpose").With("assetPurpose", a.Purpose).With("expectedPurpose", expectedPurpose)
	}
	if a.Status == StatusReady {
		return &CompleteUploadResponse{AssetID: a.ID, PublicURL: a.PublicURL, Status: a.Status}, nil
	}

	provider := strings.ToLower(strings.TrimSpace(viper.GetString(config.StorageProvider)))

	// Local uploads: content must have been PUT to us already.
	if provider == "local" {
		if strings.TrimSpace(a.LocalFilePath) == "" {
			return nil, ErrUploadNotReady()
		}
		st, statErr := os.Stat(a.LocalFilePath)
		if statErr != nil {
			return nil, ErrUploadNotReady()
		}
		if st.Size() != a.SizeBytes {
			return nil, problem.Conflict("uploaded size mismatch").With("expected", a.SizeBytes).With("actual", st.Size())
		}

		return s.markReady(ctx, a, s.buildPublicURL(biz, a))
	}

	if s.blob == nil {
		return nil, problem.InternalError().With("blob", "not_configured")
	}

	info, err := s.blob.Head(ctx, a.ObjectKey)
	if err != nil {
		return nil, err
	}
	if info.SizeBytes != a.SizeBytes {
		return nil, problem.Conflict("uploaded size mismatch").With("expected", a.SizeBytes).With("actual", info.SizeBytes)
	}
	if strings.TrimSpace(info.ContentType) != "" && strings.TrimSpace(info.ContentType) != strings.TrimSpace(a.ContentType) {
		return nil, problem.Conflict("uploaded contentType mismatch").With("expected", a.ContentType).With("actual", info.ContentType)
	}

	publicURL := s.buildPublicURL(biz, a)
	if publicURL == "" {
		return nil, problem.InternalError().With("asset", "public_url_unavailable")
	}

	return s.markReady(ctx, a, publicURL)
}

// DeleteAsset removes an uploaded asset if it belongs to the business, the actor is authorized,
// and the asset is not currently referenced by any business, product, or variant.
func (s *Service) DeleteAsset(ctx context.Context, actor *account.User, biz *business.Business, assetID string) error {
	if actor == nil {
		return problem.Unauthorized("unauthorized")
	}
	if biz == nil {
		return problem.BadRequest("business is required")
	}

	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return problem.BadRequest("assetId is required")
	}

	if !throttle.Allow(s.storage.Cache(), fmt.Sprintf("rl:asset:delete:%s:%s", biz.ID, actor.ID), time.Minute, 120, 150*time.Millisecond) {
		return ErrRateLimited()
	}

	a, err := s.storage.GetByID(ctx, biz.ID, assetID)
	if err != nil || a == nil {
		return ErrAssetNotFound(assetID, err)
	}
	if a.WorkspaceID != biz.WorkspaceID {
		return problem.Forbidden("asset does not belong to workspace")
	}

	if err := s.authorizeDelete(actor, a); err != nil {
		return err
	}

	if a.Status == StatusReady {
		referenced, rerr := s.storage.IsReferenced(ctx, a)
		if rerr != nil {
			return rerr
		}
		if referenced {
			return ErrAssetInUse(a.ID)
		}
	}

	_, derr := s.deleteAsset(ctx, a)
	return derr
}

func (s *Service) authorizeDelete(actor *account.User, a *Asset) error {
	switch a.Purpose {
	case PurposeBusinessLogo:
		return actor.Role.HasPermission(role.ActionManage, role.ResourceBusiness)
	case PurposeProductPhoto, PurposeVariantPhoto:
		return actor.Role.HasPermission(role.ActionManage, role.ResourceInventory)
	default:
		return problem.Forbidden("unsupported asset purpose").With("purpose", a.Purpose)
	}
}

type deleteOutcome struct {
	deletedFile bool
	deletedBlob bool
}

func (s *Service) deleteAsset(ctx context.Context, a *Asset) (*deleteOutcome, error) {
	out := &deleteOutcome{}
	if a == nil {
		return out, nil
	}

	if a.LocalFilePath != "" {
		if err := os.Remove(a.LocalFilePath); err == nil {
			out.deletedFile = true
		} else if !os.IsNotExist(err) {
			return out, err
		}
	}

	if s.blob != nil && a.ObjectKey != "" {
		if err := s.blob.Delete(ctx, a.ObjectKey); err != nil {
			return out, err
		}
		out.deletedBlob = true
	}

	if err := s.storage.Delete(ctx, a); err != nil {
		return out, err
	}
	return out, nil
}

func (s *Service) StoreLocalContent(ctx context.Context, actor *account.User, biz *business.Business, assetID string, expectedPurpose Purpose, contentType string, r io.Reader) (*Asset, error) {
	if actor == nil {
		return nil, problem.Unauthorized("unauthorized")
	}
	if biz == nil {
		return nil, problem.BadRequest("business is required")
	}
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return nil, problem.BadRequest("assetId is required")
	}
	provider := strings.ToLower(strings.TrimSpace(viper.GetString(config.StorageProvider)))
	if provider != "local" {
		return nil, problem.BadRequest("local upload endpoint is disabled")
	}

	a, err := s.storage.GetByID(ctx, biz.ID, assetID)
	if err != nil || a == nil {
		return nil, ErrAssetNotFound(assetID, err)
	}
	if a.Purpose != expectedPurpose {
		return nil, problem.Forbidden("invalid asset purpose").With("assetPurpose", a.Purpose).With("expectedPurpose", expectedPurpose)
	}
	if a.Status == StatusReady {
		return a, nil
	}

	if !isAllowedImageContentType(a.ContentType) {
		return nil, ErrUploadNotAllowed("contentType", "only image uploads are supported")
	}

	// Ensure the client isn't smuggling a different Content-Type.
	ct := strings.TrimSpace(contentType)
	if ct != "" && strings.TrimSpace(a.ContentType) != "" && ct != a.ContentType {
		return nil, problem.Conflict("contentType mismatch").With("expected", a.ContentType).With("actual", ct)
	}

	maxBytes := viper.GetInt64(config.UploadsMaxBytes)
	limit := a.SizeBytes
	if maxBytes > 0 && limit > maxBytes {
		limit = maxBytes
	}
	if limit <= 0 {
		limit = maxBytes
	}
	if limit <= 0 {
		limit = 5 * 1024 * 1024
	}

	// Read with a hard limit to prevent abuse.
	data, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return nil, problem.InternalError().WithError(err)
	}
	if int64(len(data)) != a.SizeBytes {
		return nil, problem.BadRequest("uploaded size must match sizeBytes").With("expected", a.SizeBytes).With("actual", len(data))
	}

	// Persist to disk.
	if err := os.MkdirAll(s.localDir, 0o755); err != nil {
		return nil, problem.InternalError().WithError(err)
	}
	path := filepath.Join(s.localDir, a.ID)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return nil, problem.InternalError().WithError(err)
	}

	a.LocalFilePath = path
	if err := s.storage.Update(ctx, a); err != nil {
		return nil, problem.InternalError().WithError(err)
	}

	return a, nil
}

func (s *Service) GetPublicAsset(ctx context.Context, assetID string) (*Asset, error) {
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return nil, problem.BadRequest("assetId is required")
	}
	// Public assets are safe to resolve without auth; still enforce visibility + status.
	// We query by ID without a business scope.
	a, err := s.storage.asset.FindByID(ctx, assetID)
	if err != nil || a == nil {
		return nil, ErrAssetNotFound(assetID, err)
	}
	if a.Visibility != VisibilityPublic {
		return nil, problem.NotFound("asset not found")
	}
	if a.Status != StatusReady {
		return nil, problem.NotFound("asset not found")
	}
	return a, nil
}

func (s *Service) buildObjectKey(businessID, assetID, fileName string) string {
	name := sanitizeFilename(fileName)
	return fmt.Sprintf("business/%s/assets/%s/%s", businessID, assetID, name)
}

func (s *Service) buildPublicURL(biz *business.Business, a *Asset) string {
	if a == nil {
		return ""
	}
	provider := strings.ToLower(strings.TrimSpace(viper.GetString(config.StorageProvider)))
	if provider == "local" {
		base := strings.TrimRight(strings.TrimSpace(viper.GetString(config.HTTPBaseURL)), "/")
		if base == "" {
			base = "http://localhost:8080"
		}
		return fmt.Sprintf("%s/v1/public/assets/%s", base, a.ID)
	}
	if s.blob == nil {
		return ""
	}
	if url, ok := s.blob.PublicURL(a.ObjectKey); ok {
		return url
	}
	return ""
}

func (s *Service) buildUploadDescriptor(ctx context.Context, biz *business.Business, a *Asset) (*UploadDescriptor, error) {
	provider := strings.ToLower(strings.TrimSpace(viper.GetString(config.StorageProvider)))
	base := strings.TrimRight(strings.TrimSpace(viper.GetString(config.HTTPBaseURL)), "/")
	if base == "" {
		base = "http://localhost:8080"
	}

	if provider == "local" {
		// Local dev: upload directly to the API.
		return &UploadDescriptor{
			Method:  http.MethodPut,
			URL:     fmt.Sprintf("%s/v1/businesses/%s/assets/uploads/%s/content/%s", base, biz.Descriptor, a.ID, a.Purpose),
			Headers: map[string]string{"Content-Type": a.ContentType},
		}, nil
	}

	if s.blob == nil {
		return nil, problem.InternalError().With("blob", "not_configured")
	}

	expires := 10 * time.Minute
	out, err := s.blob.PresignPut(ctx, blob.PresignPutInput{
		Key:         a.ObjectKey,
		ContentType: a.ContentType,
		SizeBytes:   a.SizeBytes,
		ExpiresIn:   expires,
	})
	if err != nil {
		return nil, err
	}

	return &UploadDescriptor{Method: out.Method, URL: out.URL, Headers: out.Headers}, nil
}

func (s *Service) markReady(ctx context.Context, a *Asset, publicURL string) (*CompleteUploadResponse, error) {
	now := time.Now().UTC()
	a.Status = StatusReady
	a.CompletedAt = &now
	a.PublicURL = publicURL
	if err := s.storage.Update(ctx, a); err != nil {
		return nil, problem.InternalError().WithError(err)
	}
	return &CompleteUploadResponse{AssetID: a.ID, PublicURL: a.PublicURL, Status: a.Status}, nil
}

func isAllowedImageContentType(ct string) bool {
	ct = strings.TrimSpace(strings.ToLower(ct))
	if ct == "" {
		return false
	}
	// Normalize common "image/jpg".
	if ct == "image/jpg" {
		ct = "image/jpeg"
	}
	switch ct {
	case "image/png", "image/jpeg", "image/webp", "image/gif":
		return true
	default:
		return false
	}
}

func sniffContentType(headerValue string) string {
	if headerValue == "" {
		return ""
	}
	mt, _, err := mime.ParseMediaType(headerValue)
	if err != nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(mt))
}

var errMissingLocalFile = errors.New("local file missing")
