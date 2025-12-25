package blob

import (
	"strings"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/spf13/viper"
)

// FromConfig returns a blob Provider configured from Viper settings.
//
// If storage.provider is "local" (or empty), it returns (nil, nil) to indicate
// no external blob storage is configured.
func FromConfig() (Provider, error) {
	storageProvider := strings.ToLower(strings.TrimSpace(viper.GetString(config.StorageProvider)))
	if storageProvider == "" || storageProvider == "local" {
		return nil, nil
	}

	p, err := NewS3CompatibleProvider(S3CompatibleConfig{
		Bucket:          viper.GetString(config.StorageBucket),
		Region:          viper.GetString(config.StorageRegion),
		Endpoint:        viper.GetString(config.StorageEndpoint),
		AccessKeyID:     viper.GetString(config.StorageAccessKeyID),
		SecretAccessKey: viper.GetString(config.StorageSecretAccessKey),
		PublicBaseURL:   viper.GetString(config.StoragePublicBaseURL),
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}
