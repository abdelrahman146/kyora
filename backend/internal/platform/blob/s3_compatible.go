package blob

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3CompatibleProvider implements Provider for any S3-compatible storage
// (AWS S3, DigitalOcean Spaces, MinIO, etc.).
//
// It assumes bucket access is configured externally (bucket policy / ACL).
// If PublicBaseURL is provided, PublicURL will be derived by concatenation.
// Otherwise, PublicURL returns ok=false.
//
// Note: For DigitalOcean Spaces, use endpoint like https://nyc3.digitaloceanspaces.com
// and region like nyc3.
type S3CompatibleProvider struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string

	publicBaseURL string
}

type S3CompatibleConfig struct {
	Bucket          string
	Region          string
	Endpoint        string // optional; for Spaces or custom S3
	AccessKeyID     string
	SecretAccessKey string

	// PublicBaseURL, when set, is used to build public URLs: <PublicBaseURL>/<key>
	// Example: https://my-bucket.nyc3.digitaloceanspaces.com
	PublicBaseURL string
}

func NewS3CompatibleProvider(cfg S3CompatibleConfig) (*S3CompatibleProvider, error) {
	if strings.TrimSpace(cfg.Bucket) == "" {
		return nil, errors.New("blob s3: bucket is required")
	}
	if strings.TrimSpace(cfg.Region) == "" {
		return nil, errors.New("blob s3: region is required")
	}
	if strings.TrimSpace(cfg.AccessKeyID) == "" {
		return nil, errors.New("blob s3: accessKeyId is required")
	}
	if strings.TrimSpace(cfg.SecretAccessKey) == "" {
		return nil, errors.New("blob s3: secretAccessKey is required")
	}

	awsCfg := aws.Config{
		Region:      cfg.Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
	}

	if strings.TrimSpace(cfg.Endpoint) != "" {
		endpoint := strings.TrimRight(strings.TrimSpace(cfg.Endpoint), "/")
		awsCfg.BaseEndpoint = aws.String(endpoint)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		// For Spaces / many S3-compatible providers, path-style is typically required.
		// DO Spaces supports virtual-hosted too, but path-style is the safer default across S3-compatible providers.
		o.UsePathStyle = true
	})

	return &S3CompatibleProvider{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucket:        strings.TrimSpace(cfg.Bucket),
		publicBaseURL: strings.TrimRight(strings.TrimSpace(cfg.PublicBaseURL), "/"),
	}, nil
}

func (p *S3CompatibleProvider) PresignPut(ctx context.Context, in PresignPutInput) (*PresignPutOutput, error) {
	if p == nil || p.presignClient == nil {
		return nil, ErrProviderNotConfigured()
	}
	if strings.TrimSpace(in.Key) == "" {
		return nil, errors.New("blob s3: key is required")
	}
	if strings.TrimSpace(in.ContentType) == "" {
		return nil, errors.New("blob s3: contentType is required")
	}
	if in.ExpiresIn <= 0 {
		in.ExpiresIn = 10 * time.Minute
	}

	req := &s3.PutObjectInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(in.Key),
		ContentType: aws.String(in.ContentType),
		// We deliberately do NOT set ACL here; many S3-compatible providers require
		// bucket policy instead. If you want public objects, configure bucket policy or CDN.
	}

	presigned, err := p.presignClient.PresignPutObject(ctx, req, func(po *s3.PresignOptions) {
		po.Expires = in.ExpiresIn
	})
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Content-Type": in.ContentType,
	}

	return &PresignPutOutput{
		Method:    "PUT",
		URL:       presigned.URL,
		Headers:   headers,
		ExpiresAt: time.Now().UTC().Add(in.ExpiresIn),
	}, nil
}

func (p *S3CompatibleProvider) Head(ctx context.Context, key string) (*ObjectInfo, error) {
	if p == nil || p.client == nil {
		return nil, ErrProviderNotConfigured()
	}
	if strings.TrimSpace(key) == "" {
		return nil, errors.New("blob s3: key is required")
	}

	out, err := p.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(p.bucket), Key: aws.String(key)})
	if err != nil {
		var nsk *types.NotFound
		if errors.As(err, &nsk) {
			return nil, ErrObjectNotFound(key)
		}
		return nil, err
	}

	etag := ""
	if out.ETag != nil {
		etag = strings.Trim(*out.ETag, "\"")
	}

	lm := time.Time{}
	if out.LastModified != nil {
		lm = *out.LastModified
	}

	ct := ""
	if out.ContentType != nil {
		ct = *out.ContentType
	}

	sz := int64(0)
	if out.ContentLength != nil {
		sz = *out.ContentLength
	}

	return &ObjectInfo{
		Key:          key,
		SizeBytes:    sz,
		ContentType:  ct,
		ETag:         etag,
		LastModified: lm,
	}, nil
}

func (p *S3CompatibleProvider) Delete(ctx context.Context, key string) error {
	if p == nil || p.client == nil {
		return ErrProviderNotConfigured()
	}
	if strings.TrimSpace(key) == "" {
		return nil
	}
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(p.bucket), Key: aws.String(key)})
	// S3 delete is idempotent.
	return err
}

func (p *S3CompatibleProvider) PublicURL(key string) (string, bool) {
	if p == nil {
		return "", false
	}
	base := strings.TrimSpace(p.publicBaseURL)
	if base == "" {
		return "", false
	}
	u, err := url.Parse(base)
	if err != nil {
		return "", false
	}
	// Ensure no double slashes.
	path := strings.TrimLeft(key, "/")
	u.Path = strings.TrimRight(u.Path, "/") + "/" + path
	return u.String(), true
}
