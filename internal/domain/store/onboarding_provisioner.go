package store

import (
	"context"
	"strings"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/utils"
)

type onboardingStoreProvisioner struct {
	service *StoreService
}

func NewOnboardingStoreProvisioner(service *StoreService) account.StoreProvisioner {
	return &onboardingStoreProvisioner{service: service}
}

func (p *onboardingStoreProvisioner) ProvisionInitialStore(ctx context.Context, organizationID string, req *account.CreateInitialStoreRequest) error {
	if req == nil {
		return utils.Problem.BadRequest("missing store information")
	}
	name := strings.TrimSpace(req.Name)
	if name == "" || req.CountryCode == "" || req.Currency == "" {
		return utils.Problem.BadRequest("store name, country code and currency are required")
	}

	storeReq := &CreateStoreRequest{
		Name:        name,
		CountryCode: req.CountryCode,
		Currency:    req.Currency,
	}
	_, err := p.service.CreateStore(ctx, organizationID, storeReq)
	return err
}
