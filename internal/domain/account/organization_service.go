package account

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
)

type OrganizationService struct {
	orgRepo *OrganizationRepository
}

func NewOrganizationService(orgRepo *OrganizationRepository) *OrganizationService {
	return &OrganizationService{orgRepo: orgRepo}
}

func (s *OrganizationService) UpdateOrganization(ctx context.Context, orgID string, orgReq *UpdateOrganizationRequest) (*Organization, error) {
	if _, err := s.orgRepo.FindOne(ctx, s.orgRepo.ScopeID(orgID)); err != nil {
		return nil, db.HandleDBError(err)
	}
	org := &Organization{
		Name: orgReq.Name,
	}
	if err := s.orgRepo.PatchOne(ctx, org, s.orgRepo.ScopeID(orgID), db.WithReturning(&org)); err != nil {
		return nil, db.HandleDBError(err)
	}
	return org, nil
}

func (s *OrganizationService) GetOrganizationByID(ctx context.Context, id string) (*Organization, error) {
	org, err := s.orgRepo.FindOne(ctx, db.WithScopes(s.orgRepo.ScopeID(id)))
	if err != nil {
		return nil, db.HandleDBError(err)
	}
	return org, nil
}
