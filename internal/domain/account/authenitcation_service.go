package account

import (
	"context"
	"encoding/json"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/spf13/viper"
)

type AuthenticationService struct {
	userRepo *UserRepository
	cache    *db.Memcache
}

func NewAuthenticationService(userRepo *UserRepository, cache *db.Memcache) *AuthenticationService {
	return &AuthenticationService{userRepo: userRepo, cache: cache}
}

func (s *AuthenticationService) Authenticate(ctx context.Context, email, password string) (*User, string, error) {
	user, err := s.userRepo.FindOne(ctx, s.userRepo.ScopeEmail(email), db.WithPreload(OrganizationStruct))
	if err != nil {
		return nil, "", err
	}
	if user == nil || !utils.Hash.Validate(password, user.PasswordHash) {
		return nil, "", utils.Problem.Unauthorized("Invalid email or password")
	}
	jwt, err := utils.JWT.GenerateToken(user.ID, user.OrganizationID)
	if err != nil {
		return nil, "", err
	}
	return user, jwt, nil
}

func (s *AuthenticationService) GetUserByID(ctx context.Context, id string) (*User, error) {
	user, err := s.userRepo.FindByID(ctx, id, db.WithPreload(OrganizationStruct))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthenticationService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.userRepo.FindOne(ctx, s.userRepo.ScopeEmail(email), db.WithPreload(OrganizationStruct))
	if err != nil {
		return nil, err
	}
	return user, nil
}

type resetPayload struct {
	UserID         string    `json:"userId"`
	OrganizationID string    `json:"organizationId"`
	Email          string    `json:"email"`
	ExpAt          time.Time `json:"expAt"`
}

const resetPrefix = "pwreset:"

func (s *AuthenticationService) CreateResetToken(ctx context.Context, email string) (string, error) {

	user, err := s.userRepo.FindOne(ctx, s.userRepo.ScopeEmail(email))
	if err != nil {
		return "", err
	}
	token, err := utils.ID.RandomString(40)
	if err != nil {
		return "", err
	}
	passwordResetTtlSeconds := viper.GetInt64("password_reset_ttl_seconds")
	ttl := time.Duration(passwordResetTtlSeconds) * time.Second
	payload := resetPayload{UserID: user.ID, OrganizationID: user.OrganizationID, Email: user.Email, ExpAt: time.Now().Add(ttl)}
	b, _ := json.Marshal(payload)
	exp := int32(ttl.Seconds())
	if exp <= 0 {
		exp = 900 // default 15m
	}
	if err := s.cache.Set(resetPrefix+token, b, exp); err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthenticationService) ValidateResetToken(ctx context.Context, token string) (*User, error) {
	b, err := s.cache.Get(resetPrefix + token)
	if err != nil || len(b) == 0 {
		return nil, utils.Problem.NotFound("Invalid or expired token")
	}
	var p resetPayload
	_ = json.Unmarshal(b, &p)
	if time.Now().After(p.ExpAt) {
		_ = s.cache.Delete(resetPrefix + token)
		return nil, utils.Problem.NotFound("Token expired")
	}
	user, err := s.userRepo.FindByID(ctx, p.UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthenticationService) ConsumeResetToken(ctx context.Context, token string, newPassword string) error {
	user, err := s.ValidateResetToken(ctx, token)
	if err != nil {
		return err
	}
	hash, err := utils.Hash.Make(newPassword)
	if err != nil {
		return utils.Problem.InternalError().WithError(err)
	}
	updates := &User{PasswordHash: hash}
	if err := s.userRepo.PatchOne(ctx, updates, s.userRepo.ScopeID(user.ID)); err != nil {
		return err
	}
	_ = s.cache.Delete(resetPrefix + token)
	return nil
}
