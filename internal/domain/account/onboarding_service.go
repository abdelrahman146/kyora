package account

import (
	"context"
	"fmt"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/utils"
)

type OnboardingService struct {
	userRepo         *UserRepository
	organizationRepo *OrganizationRepository
	atomicProcess    *db.AtomicProcess
}

func NewOnboardingService(userRepo *UserRepository, organizationRepo *OrganizationRepository, atomicProcess *db.AtomicProcess) *OnboardingService {
	return &OnboardingService{
		userRepo:         userRepo,
		organizationRepo: organizationRepo,
		atomicProcess:    atomicProcess,
	}
}

func (s *OnboardingService) OnboardNewOrganization(ctx context.Context, orgReq *CreateOrganizationRequest, userReq *CreateUserRequest) (*User, error) {
	var createdUser *User

	err := s.atomicProcess.Exec(ctx, func(ctx context.Context) error {
		if existingOrg, _ := s.organizationRepo.FindOne(ctx, s.organizationRepo.ScopeSlug(orgReq.Slug)); existingOrg != nil {
			return utils.Problem.Conflict("Organization with the given slug already exists").WithError(fmt.Errorf("organization with slug %q already exists", orgReq.Slug)).With("slug", orgReq.Slug)
		}
		org := &Organization{
			Slug: orgReq.Slug,
			Name: orgReq.Name,
		}
		if err := s.organizationRepo.CreateOne(ctx, org); err != nil {
			return err
		}
		passwordHash, err := utils.Hash.Make(userReq.Password)
		if err != nil {
			return utils.Problem.InternalError().WithError(err)
		}
		if existingUser, _ := s.userRepo.FindOne(ctx, s.userRepo.ScopeEmail(userReq.Email)); existingUser != nil {
			return utils.Problem.Conflict("User with the given email already exists").With("email", userReq.Email)
		}
		user := &User{
			FirstName:      userReq.FirstName,
			LastName:       userReq.LastName,
			Email:          userReq.Email,
			PasswordHash:   passwordHash,
			OrganizationID: org.ID,
		}
		if err := s.userRepo.CreateOne(ctx, user); err != nil {
			return err
		}
		if createdUser, err = s.userRepo.FindByID(ctx, user.ID, db.WithPreload(OrganizationStruct)); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *OnboardingService) IsOrganizationSlugAvailable(ctx context.Context, slug string) (bool, error) {
	existingOrg, err := s.organizationRepo.FindOne(ctx, s.organizationRepo.ScopeSlug(slug))
	if err != nil {
		if db.IsRecordNotFound(err) {
			return true, nil
		}
		return false, err
	}
	return existingOrg == nil, nil
}

func (s *OnboardingService) IsEmailAvailable(ctx context.Context, email string) (bool, error) {
	existingUser, err := s.userRepo.FindOne(ctx, s.userRepo.ScopeEmail(email))
	if err != nil {
		if db.IsRecordNotFound(err) {
			return true, nil
		}
		return false, err
	}
	return existingUser == nil, nil
}
