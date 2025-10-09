package account

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
)

type UserService struct {
	userRepo *UserRepository
}

func NewUserService(userRepo *UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) UpdateUser(ctx context.Context, userID string, userReq *UpdateUserRequest) (*User, error) {
	if _, err := s.userRepo.FindOne(ctx, s.userRepo.ScopeID(userID)); err != nil {
		return nil, db.HandleDBError(err)
	}
	user := &User{}
	if userReq.FirstName != "" {
		user.FirstName = userReq.FirstName
	}
	if userReq.LastName != "" {
		user.LastName = userReq.LastName
	}
	if err := s.userRepo.PatchOne(ctx, user, s.userRepo.ScopeID(userID), db.WithReturning(&user)); err != nil {
		return nil, db.HandleDBError(err)
	}
	return user, nil
}

func (s *UserService) GetOrganizationUsers(ctx context.Context, orgID string) ([]*User, error) {
	users, err := s.userRepo.List(ctx, s.userRepo.ScopeOrganizationID(orgID))
	if err != nil {
		return nil, db.HandleDBError(err)
	}
	return users, nil
}
