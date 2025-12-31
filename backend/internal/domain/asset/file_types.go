package asset

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// FileCategory represents the category of a file based on its content type.
type FileCategory string

const (
	FileCategoryImage      FileCategory = "image"
	FileCategoryVideo      FileCategory = "video"
	FileCategoryAudio      FileCategory = "audio"
	FileCategoryDocument   FileCategory = "document"
	FileCategoryCompressed FileCategory = "compressed"
	FileCategoryOther      FileCategory = "other"
)

// defaultAllowedExtensions provides sensible defaults for file type validation.
var defaultAllowedExtensions = map[FileCategory][]string{
	FileCategoryImage:      {"jpg", "jpeg", "png", "webp", "gif", "heic", "heif"},
	FileCategoryVideo:      {"mp4", "mov", "avi", "mkv", "webm"},
	FileCategoryAudio:      {"mp3", "wav", "ogg", "m4a", "aac"},
	FileCategoryDocument:   {"pdf", "doc", "docx", "txt", "rtf", "odt"},
	FileCategoryCompressed: {"zip", "tar", "gz", "rar", "7z"},
}

// defaultMaxSizeBytes provides sensible defaults for per-category size limits.
var defaultMaxSizeBytes = map[FileCategory]int64{
	FileCategoryImage:      10 * 1024 * 1024,  // 10 MB
	FileCategoryVideo:      100 * 1024 * 1024, // 100 MB
	FileCategoryAudio:      20 * 1024 * 1024,  // 20 MB
	FileCategoryDocument:   10 * 1024 * 1024,  // 10 MB
	FileCategoryCompressed: 50 * 1024 * 1024,  // 50 MB
	FileCategoryOther:      5 * 1024 * 1024,   // 5 MB
}

// FileTypeValidator validates file types and sizes based on configuration.
type FileTypeValidator struct {
	allowedExtensions map[FileCategory][]string
	maxSizeBytes      map[FileCategory]int64
}

// NewFileTypeValidator creates a new validator with config-driven rules.
func NewFileTypeValidator() *FileTypeValidator {
	validator := &FileTypeValidator{
		allowedExtensions: defaultAllowedExtensions,
		maxSizeBytes:      defaultMaxSizeBytes,
	}

	// Override with config if present
	if viper.IsSet("uploads.allowed_extensions") {
		configMap := make(map[string][]string)
		if err := viper.UnmarshalKey("uploads.allowed_extensions", &configMap); err == nil {
			validator.allowedExtensions = make(map[FileCategory][]string)
			for category, extensions := range configMap {
				validator.allowedExtensions[FileCategory(category)] = extensions
			}
		}
	}

	if viper.IsSet("uploads.max_size_bytes") {
		configMap := make(map[string]int64)
		if err := viper.UnmarshalKey("uploads.max_size_bytes", &configMap); err == nil {
			validator.maxSizeBytes = make(map[FileCategory]int64)
			for category, size := range configMap {
				validator.maxSizeBytes[FileCategory(category)] = size
			}
		}
	}

	return validator
}

// GetCategory returns the file category based on content type.
func (v *FileTypeValidator) GetCategory(contentType string) FileCategory {
	contentType = strings.ToLower(contentType)

	if strings.HasPrefix(contentType, "image/") {
		return FileCategoryImage
	}
	if strings.HasPrefix(contentType, "video/") {
		return FileCategoryVideo
	}
	if strings.HasPrefix(contentType, "audio/") {
		return FileCategoryAudio
	}
	if strings.Contains(contentType, "pdf") || strings.Contains(contentType, "document") ||
		strings.Contains(contentType, "text") || strings.Contains(contentType, "msword") ||
		strings.Contains(contentType, "wordprocessingml") {
		return FileCategoryDocument
	}
	if strings.Contains(contentType, "zip") || strings.Contains(contentType, "compressed") ||
		strings.Contains(contentType, "tar") || strings.Contains(contentType, "gzip") ||
		strings.Contains(contentType, "x-rar") || strings.Contains(contentType, "x-7z") {
		return FileCategoryCompressed
	}

	return FileCategoryOther
}

// IsAllowed checks if a file is allowed based on extension and content type.
func (v *FileTypeValidator) IsAllowed(fileName, contentType string) bool {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), ".")
	if ext == "" {
		return false
	}

	category := v.GetCategory(contentType)
	allowed, ok := v.allowedExtensions[category]
	if !ok {
		return false
	}

	for _, allowedExt := range allowed {
		if ext == strings.ToLower(allowedExt) {
			return true
		}
	}

	return false
}

// GetMaxSize returns the maximum allowed file size for a category.
func (v *FileTypeValidator) GetMaxSize(category FileCategory) int64 {
	if size, ok := v.maxSizeBytes[category]; ok {
		return size
	}
	return defaultMaxSizeBytes[FileCategoryOther]
}

// NeedsThumbnail returns true if the file category should generate thumbnails.
func (v *FileTypeValidator) NeedsThumbnail(category FileCategory) bool {
	return category == FileCategoryImage || category == FileCategoryVideo
}
