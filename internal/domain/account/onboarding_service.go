package account

import (
	"context"
	"fmt"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/utils"
)

type OnboardingService struct {
	userRepo         *userRepository
	organizationRepo *organizationRepository
	atomicProcess    *db.AtomicProcess
	storeProvisioner StoreProvisioner
}

type StoreProvisioner interface {
	ProvisionInitialStore(ctx context.Context, organizationID string, req *CreateInitialStoreRequest) error
}

type CreateInitialStoreRequest struct {
	Name        string
	CountryCode string
	Currency    string
}

func NewOnboardingService(userRepo *userRepository, organizationRepo *organizationRepository, atomicProcess *db.AtomicProcess, storeProvisioner StoreProvisioner) *OnboardingService {
	return &OnboardingService{
		userRepo:         userRepo,
		organizationRepo: organizationRepo,
		atomicProcess:    atomicProcess,
		storeProvisioner: storeProvisioner,
	}
}

func (s *OnboardingService) OnboardNewOrganization(ctx context.Context, orgReq *CreateOrganizationRequest, userReq *CreateUserRequest, storeReq *CreateInitialStoreRequest) (*User, error) {
	if s.storeProvisioner == nil {
		return nil, utils.Problem.InternalError().WithError(fmt.Errorf("store provisioner is not configured"))
	}
	if storeReq == nil {
		return nil, utils.Problem.BadRequest("missing store information")
	}

	var createdUser *User
	err := s.atomicProcess.Exec(ctx, func(txCtx context.Context) error {
		org, err := s.createOrganization(txCtx, orgReq)
		if err != nil {
			return err
		}
		createdUser, err = s.createUser(txCtx, org.ID, userReq)
		if err != nil {
			return err
		}
		return s.storeProvisioner.ProvisionInitialStore(txCtx, org.ID, storeReq)
	})

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *OnboardingService) createOrganization(ctx context.Context, orgReq *CreateOrganizationRequest) (*Organization, error) {
	org := &Organization{
		Name: orgReq.Name,
	}
	if err := s.organizationRepo.createOne(ctx, org); err != nil {
		return nil, err
	}
	return org, nil
}

func (s *OnboardingService) createUser(ctx context.Context, organizationID string, userReq *CreateUserRequest) (*User, error) {
	passwordHash, err := utils.Hash.Make(userReq.Password)
	if err != nil {
		return nil, utils.Problem.InternalError().WithError(err)
	}
	if existingUser, _ := s.userRepo.findOne(ctx, s.userRepo.scopeEmail(userReq.Email)); existingUser != nil {
		return nil, utils.Problem.Conflict("User with the given email already exists").With("email", userReq.Email)
	}
	user := &User{
		FirstName:      userReq.FirstName,
		LastName:       userReq.LastName,
		Email:          userReq.Email,
		PasswordHash:   passwordHash,
		OrganizationID: organizationID,
	}
	if err := s.userRepo.createOne(ctx, user); err != nil {
		return nil, err
	}
	createdUser, err := s.userRepo.findByID(ctx, user.ID, db.WithPreload(OrganizationStruct))
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (s *OnboardingService) IsOrganizationSlugAvailable(ctx context.Context, slug string) (bool, error) {
	existingOrg, err := s.organizationRepo.findOne(ctx, s.organizationRepo.scopeSlug(slug))
	if err != nil {
		if db.IsRecordNotFound(err) {
			return true, nil
		}
		return false, err
	}
	return existingOrg == nil, nil
}

func (s *OnboardingService) IsEmailAvailable(ctx context.Context, email string) (bool, error) {
	existingUser, err := s.userRepo.findOne(ctx, s.userRepo.scopeEmail(email))
	if err != nil {
		if db.IsRecordNotFound(err) {
			return true, nil
		}
		return false, err
	}
	return existingUser == nil, nil
}
