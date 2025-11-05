package business

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
)

type Service struct {
	storage         *Storage
	atomicProcessor atomic.AtomicProcessor
	bus             *bus.Bus
}

func NewService(storage *Storage, atomicProcessor atomic.AtomicProcessor, bus *bus.Bus) *Service {
	return &Service{
		storage:         storage,
		atomicProcessor: atomicProcessor,
		bus:             bus,
	}
}

func (s *Service) GetBusinessByID(ctx context.Context, actor *account.User, id string) (*Business, error) {
	workspaceID := actor.WorkspaceID
	if err := actor.Role.HasPermission(role.ActionView, role.ResourceBusiness); err != nil {
		return nil, err
	}
	return s.storage.business.FindOne(ctx,
		s.storage.business.ScopeWorkspaceID(workspaceID),
		s.storage.business.ScopeID(id),
		s.storage.business.WithPreload(account.WorkspaceStruct),
	)
}

func (s *Service) GetBusinessByDescriptor(ctx context.Context, actor *account.User, descriptor string) (*Business, error) {
	return s.storage.business.FindOne(ctx,
		s.storage.business.ScopeWorkspaceID(actor.WorkspaceID),
		s.storage.business.ScopeEquals(BusinessSchema.Descriptor, descriptor),
		s.storage.business.WithPreload(account.WorkspaceStruct),
	)
}

func (s *Service) ListBusinesses(ctx context.Context, actor *account.User) ([]*Business, error) {
	return s.storage.business.FindMany(ctx, s.storage.business.ScopeWorkspaceID(actor.WorkspaceID))
}

func (s *Service) CreateBusiness(ctx context.Context, actor *account.User, input *CreateBusinessInput) (*Business, error) {
	business := &Business{
		WorkspaceID: actor.WorkspaceID,
		Descriptor:  input.Descriptor,
		Name:        input.Name,
		CountryCode: input.CountryCode,
		VatRate:     input.VatRate,
		Currency:    input.Currency,
	}
	if !input.SafetyBuffer.IsZero() {
		business.SafetyBuffer = input.SafetyBuffer
	}
	if !input.EstablishedAt.IsZero() {
		business.EstablishedAt = input.EstablishedAt
	}
	err := s.storage.business.CreateOne(ctx, business)
	if err != nil {
		return nil, err
	}
	return business, nil
}

func (s *Service) ArchiveBusiness(ctx context.Context, actor *account.User, id string) error {
	business, err := s.GetBusinessByID(ctx, actor, id)
	if err != nil {
		return err
	}
	now := business.ArchivedAt
	if now != nil {
		// already archived
		return nil
	}
	now = new(time.Time)
	*now = time.Now()
	business.ArchivedAt = now
	return s.storage.business.UpdateOne(ctx, business)
}

func (s *Service) UnarchiveBusiness(ctx context.Context, actor *account.User, id string) error {
	business, err := s.GetBusinessByID(ctx, actor, id)
	if err != nil {
		return err
	}
	if business.ArchivedAt == nil {
		// not archived
		return nil
	}
	business.ArchivedAt = nil
	return s.storage.business.UpdateOne(ctx, business)
}

func (s *Service) UpdateBusiness(ctx context.Context, actor *account.User, id string, input *UpdateBusinessInput) (*Business, error) {
	business, err := s.GetBusinessByID(ctx, actor, id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		business.Name = input.Name
	}
	if input.Descriptor != "" {
		business.Descriptor = input.Descriptor
	}
	if input.CountryCode != "" {
		business.CountryCode = input.CountryCode
	}
	if input.Currency != "" {
		business.Currency = input.Currency
	}
	if input.VatRate.Valid {
		business.VatRate = transformer.FromNullDecimal(input.VatRate)
	}
	if input.SafetyBuffer.Valid {
		business.SafetyBuffer = transformer.FromNullDecimal(input.SafetyBuffer)
	}
	if !input.EstablishedAt.IsZero() {
		business.EstablishedAt = input.EstablishedAt
	}
	err = s.storage.business.UpdateOne(ctx, business)
	if err != nil {
		return nil, err
	}
	return business, nil
}

func (s *Service) DeleteBusiness(ctx context.Context, actor *account.User, id string) error {
	business, err := s.GetBusinessByID(ctx, actor, id)
	if err != nil {
		return err
	}
	return s.storage.business.DeleteOne(ctx, business)
}

func (s *Service) IsBusinessDescriptorAvailable(ctx context.Context, actor *account.User, descriptor string) (bool, error) {
	business, err := s.storage.business.FindOne(ctx,
		s.storage.business.ScopeWorkspaceID(actor.WorkspaceID),
		s.storage.business.ScopeEquals(BusinessSchema.Descriptor, descriptor),
	)
	if err != nil {
		return false, err
	}
	return business == nil, nil
}

func (s *Service) CountBusinesses(ctx context.Context, actor *account.User) (int64, error) {
	return s.storage.business.Count(ctx, s.storage.business.ScopeWorkspaceID(actor.WorkspaceID))
}

func (s *Service) CountActiveBusinesses(ctx context.Context, actor *account.User) (int64, error) {
	return s.storage.business.Count(ctx,
		s.storage.business.ScopeWorkspaceID(actor.WorkspaceID),
		s.storage.business.ScopeIsNull(BusinessSchema.ArchivedAt),
	)
}

func (s *Service) MaxBusinessesEnforceFunc(ctx context.Context, actor *account.User, businessID string) (int64, error) {
	return s.CountActiveBusinesses(ctx, actor)
}
