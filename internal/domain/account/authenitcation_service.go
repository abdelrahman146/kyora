package account

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/utils"
)

type AuthenticationService struct {
	userRepo *UserRepository
}

func NewAuthenticationService(userRepo *UserRepository) *AuthenticationService {
	return &AuthenticationService{userRepo: userRepo}
}

func (s *AuthenticationService) Authenticate(ctx context.Context, email, password string) (*User, string, error) {
	user, err := s.userRepo.FindOne(ctx, s.userRepo.ScopeEmail(email), db.WithPreload(OrganizationStruct))
	if err != nil {
		return nil, "", db.HandleDBError(err)
	}
	if user == nil || !utils.Hash.Validate(password, user.PasswordHash) {
		return nil, "", utils.Problem.Unauthorized("Invalid email or password")
	}
	jwt, err := utils.JWT.GenerateToken(user.ID, user.OrganizationID)
	if err != nil {
		return nil, "", utils.Problem.InternalError().WithError(err)
	}
	return user, jwt, nil
}

func (s *AuthenticationService) GetUserByID(ctx context.Context, id string) (*User, error) {
	user, err := s.userRepo.FindByID(ctx, id, db.WithPreload(OrganizationStruct))
	if err != nil {
		return nil, db.HandleDBError(err)
	}
	return user, nil
}
