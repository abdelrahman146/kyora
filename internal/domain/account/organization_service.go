package account

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
)

type OrganizationService struct {
	orgRepo *organizationRepository
}

func NewOrganizationService(orgRepo *organizationRepository) *OrganizationService {
	return &OrganizationService{orgRepo: orgRepo}
}

func (s *OrganizationService) UpdateOrganization(ctx context.Context, orgID string, orgReq *UpdateOrganizationRequest) (*Organization, error) {
	if _, err := s.orgRepo.findOne(ctx, s.orgRepo.scopeID(orgID)); err != nil {
		return nil, err
	}
	org := &Organization{
		Name: orgReq.Name,
	}
	if err := s.orgRepo.patchOne(ctx, org, s.orgRepo.scopeID(orgID), db.WithReturning(&org)); err != nil {
		return nil, err
	}
	return org, nil
}

func (s *OrganizationService) GetOrganizationByID(ctx context.Context, id string) (*Organization, error) {
	org, err := s.orgRepo.findOne(ctx, db.WithScopes(s.orgRepo.scopeID(id)))
	if err != nil {
		return nil, err
	}
	return org, nil
}
