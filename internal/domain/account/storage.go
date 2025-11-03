package account

import (
	"encoding/json"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/spf13/viper"
)

const (
	resetPasswordTokenPrefix = "pwreset:"
	verifyEmailTokenPrefix   = "emailverify:"
)

type Storage struct {
	cache     *cache.Cache
	workspace *database.Repository[Workspace]
	user      *database.Repository[User]
}

func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
	return &Storage{
		cache:     cache,
		workspace: database.NewRepository[Workspace](db),
		user:      database.NewRepository[User](db),
	}
}

type PasswordResetPayload struct {
	UserID      string    `json:"userId"`
	WorkspaceID string    `json:"workspaceId"`
	Email       string    `json:"email"`
	ExpAt       time.Time `json:"expAt"`
}

func (s *Storage) CreatePasswordResetToken(payload *PasswordResetPayload) (string, time.Time, error) {
	key := resetPasswordTokenPrefix + payload.UserID
	ttl := viper.GetInt32(config.PasswordResetTokenExpirySeconds)
	payload.ExpAt = time.Now().Add(time.Duration(ttl) * time.Second)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", time.Time{}, err
	}
	err = s.cache.SetX(key, payloadBytes, ttl)
	if err != nil {
		return "", time.Time{}, err
	}
	return key, payload.ExpAt, nil
}

func (s *Storage) GetPasswordResetToken(token string) (*PasswordResetPayload, error) {
	var payload PasswordResetPayload
	data, err := s.cache.Get(resetPasswordTokenPrefix + token)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}

func (s *Storage) ConsumePasswordResetToken(token string) error {
	return s.cache.Delete(resetPasswordTokenPrefix + token)
}

type VerifyEmailPayload struct {
	UserID      string    `json:"userId"`
	WorkspaceID string    `json:"workspaceId"`
	Email       string    `json:"email"`
	ExpAt       time.Time `json:"expAt"`
}

func (s *Storage) CreateVerifyEmailToken(payload *VerifyEmailPayload) (string, time.Time, error) {
	key := verifyEmailTokenPrefix + payload.UserID
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", time.Time{}, err
	}
	ttl := viper.GetInt32(config.VerifyEmailTokenExpirySeconds)
	payload.ExpAt = time.Now().Add(time.Duration(ttl) * time.Second)
	err = s.cache.SetX(key, payloadBytes, ttl)
	if err != nil {
		return "", time.Time{}, err
	}
	return key, payload.ExpAt, nil
}

func (s *Storage) GetVerifyEmailToken(token string) (*VerifyEmailPayload, error) {
	var payload VerifyEmailPayload
	data, err := s.cache.Get(verifyEmailTokenPrefix + token)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}

func (s *Storage) ConsumeVerifyEmailToken(token string) error {
	return s.cache.Delete(verifyEmailTokenPrefix + token)
}
