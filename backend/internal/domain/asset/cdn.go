package asset

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// GenerateCDNURL converts a storage provider URL to a CDN URL if cdn_base_url is configured.
// If no CDN is configured, it returns the original storage URL.
func GenerateCDNURL(storageURL string) string {
	cdnBaseURL := viper.GetString("storage.cdn_base_url")
	if cdnBaseURL == "" {
		return storageURL
	}

	// Normalize CDN base URL (remove trailing slash)
	cdnBaseURL = strings.TrimSuffix(cdnBaseURL, "/")

	// Extract the path from storage URL
	// Storage URLs are typically: https://bucket.endpoint.com/path/to/file
	// or https://storage.example.com/path/to/file
	// We need to extract everything after the domain

	// Find the third slash (after https://)
	parts := strings.SplitN(storageURL, "/", 4)
	if len(parts) < 4 {
		// Malformed URL, return as-is
		return storageURL
	}

	objectPath := parts[3] // Everything after the domain

	// Construct CDN URL
	return fmt.Sprintf("%s/%s", cdnBaseURL, objectPath)
}

// GenerateCDNURLFromObjectKey constructs a CDN URL directly from an object key.
func GenerateCDNURLFromObjectKey(objectKey string) string {
	cdnBaseURL := viper.GetString("storage.cdn_base_url")
	if cdnBaseURL == "" {
		return ""
	}

	cdnBaseURL = strings.TrimSuffix(cdnBaseURL, "/")
	return fmt.Sprintf("%s/%s", cdnBaseURL, objectKey)
}
