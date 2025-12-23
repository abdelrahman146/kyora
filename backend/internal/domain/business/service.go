package business

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
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

var businessDescriptorRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,62}$`)

func normalizeBusinessDescriptor(v string) (string, error) {
	v = strings.TrimSpace(strings.ToLower(v))
	if v == "" {
		return "", problem.BadRequest("descriptor is required").With("field", "descriptor")
	}
	if !businessDescriptorRegex.MatchString(v) {
		return "", ErrInvalidBusinessDescriptor(v)
	}
	return v, nil
}

func normalizeCountryCode(v string) string {
	return strings.TrimSpace(strings.ToUpper(v))
}

func normalizeCurrency(v string) string {
	return strings.TrimSpace(strings.ToUpper(v))
}

func (s *Service) GetBusinessByID(ctx context.Context, actor *account.User, id string) (*Business, error) {
	workspaceID := actor.WorkspaceID
	if err := actor.Role.HasPermission(role.ActionView, role.ResourceBusiness); err != nil {
		return nil, err
	}
	return s.storage.business.FindOne(ctx,
		s.storage.business.ScopeWorkspaceID(workspaceID),
		s.storage.business.ScopeID(id),
	)
}

func (s *Service) GetBusinessByDescriptor(ctx context.Context, actor *account.User, descriptor string) (*Business, error) {
	if err := actor.Role.HasPermission(role.ActionView, role.ResourceBusiness); err != nil {
		return nil, err
	}
	norm, err := normalizeBusinessDescriptor(descriptor)
	if err != nil {
		return nil, err
	}
	return s.storage.business.FindOne(ctx,
		s.storage.business.ScopeWorkspaceID(actor.WorkspaceID),
		s.storage.business.ScopeEquals(BusinessSchema.Descriptor, norm),
	)
}

// GetBusinessByDescriptorForWorkspace returns a business by descriptor scoped to a workspace.
// It intentionally does not enforce role permissions; callers should enforce authorization separately.
func (s *Service) GetBusinessByDescriptorForWorkspace(ctx context.Context, workspaceID string, descriptor string) (*Business, error) {
	norm, err := normalizeBusinessDescriptor(descriptor)
	if err != nil {
		return nil, err
	}
	return s.storage.business.FindOne(ctx,
		s.storage.business.ScopeWorkspaceID(workspaceID),
		s.storage.business.ScopeEquals(BusinessSchema.Descriptor, norm),
	)
}

func (s *Service) ListBusinesses(ctx context.Context, actor *account.User) ([]*Business, error) {
	if err := actor.Role.HasPermission(role.ActionView, role.ResourceBusiness); err != nil {
		return nil, err
	}
	return s.storage.business.FindMany(ctx, s.storage.business.ScopeWorkspaceID(actor.WorkspaceID))
}

func (s *Service) CreateBusiness(ctx context.Context, actor *account.User, input *CreateBusinessInput) (*Business, error) {
	if err := actor.Role.HasPermission(role.ActionManage, role.ResourceBusiness); err != nil {
		return nil, err
	}
	if input == nil {
		return nil, problem.BadRequest("input is required")
	}
	normDescriptor, err := normalizeBusinessDescriptor(input.Descriptor)
	if err != nil {
		return nil, err
	}
	country := normalizeCountryCode(input.CountryCode)
	currency := normalizeCurrency(input.Currency)
	if len(country) != 2 {
		return nil, problem.BadRequest("invalid countryCode").With("field", "countryCode")
	}
	if len(currency) != 3 {
		return nil, problem.BadRequest("invalid currency").With("field", "currency")
	}
	available, err := s.IsBusinessDescriptorAvailable(ctx, actor, normDescriptor)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, ErrBusinessDescriptorAlreadyTaken(normDescriptor, nil)
	}

	business := &Business{
		WorkspaceID: actor.WorkspaceID,
		Descriptor:  normDescriptor,
		Name:        input.Name,
		CountryCode: country,
		VatRate:     input.VatRate,
		Currency:    currency,
	}
	if !input.SafetyBuffer.IsZero() {
		business.SafetyBuffer = input.SafetyBuffer
	}
	if !input.EstablishedAt.IsZero() {
		business.EstablishedAt = input.EstablishedAt.Time
	}
	err = s.storage.business.CreateOne(ctx, business)
	if err != nil {
		return nil, err
	}
	return business, nil
}

func (s *Service) ArchiveBusiness(ctx context.Context, actor *account.User, id string) error {
	if err := actor.Role.HasPermission(role.ActionManage, role.ResourceBusiness); err != nil {
		return err
	}
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
	if err := actor.Role.HasPermission(role.ActionManage, role.ResourceBusiness); err != nil {
		return err
	}
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
	if err := actor.Role.HasPermission(role.ActionManage, role.ResourceBusiness); err != nil {
		return nil, err
	}
	if input == nil {
		return nil, problem.BadRequest("input is required")
	}
	business, err := s.GetBusinessByID(ctx, actor, id)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		business.Name = strings.TrimSpace(*input.Name)
	}
	if input.Descriptor != nil {
		norm, err := normalizeBusinessDescriptor(*input.Descriptor)
		if err != nil {
			return nil, err
		}
		// If descriptor changed, enforce uniqueness inside workspace.
		if norm != business.Descriptor {
			available, err := s.IsBusinessDescriptorAvailable(ctx, actor, norm)
			if err != nil {
				return nil, err
			}
			if !available {
				return nil, ErrBusinessDescriptorAlreadyTaken(norm, nil)
			}
		}
		business.Descriptor = norm
	}
	if input.CountryCode != nil {
		cc := normalizeCountryCode(*input.CountryCode)
		if len(cc) != 2 {
			return nil, problem.BadRequest("invalid countryCode").With("field", "countryCode")
		}
		business.CountryCode = cc
	}
	if input.Currency != nil {
		cur := normalizeCurrency(*input.Currency)
		if len(cur) != 3 {
			return nil, problem.BadRequest("invalid currency").With("field", "currency")
		}
		business.Currency = cur
	}
	if input.VatRate.Valid {
		business.VatRate = transformer.FromNullDecimal(input.VatRate)
	}
	if input.SafetyBuffer.Valid {
		business.SafetyBuffer = transformer.FromNullDecimal(input.SafetyBuffer)
	}
	if input.EstablishedAt != nil {
		business.EstablishedAt = input.EstablishedAt.Time
	}
	err = s.storage.business.UpdateOne(ctx, business)
	if err != nil {
		return nil, err
	}
	return business, nil
}

func (s *Service) DeleteBusiness(ctx context.Context, actor *account.User, id string) error {
	if err := actor.Role.HasPermission(role.ActionManage, role.ResourceBusiness); err != nil {
		return err
	}
	business, err := s.GetBusinessByID(ctx, actor, id)
	if err != nil {
		return err
	}
	return s.storage.business.DeleteOne(ctx, business)
}

func (s *Service) IsBusinessDescriptorAvailable(ctx context.Context, actor *account.User, descriptor string) (bool, error) {
	if err := actor.Role.HasPermission(role.ActionView, role.ResourceBusiness); err != nil {
		return false, err
	}
	norm, err := normalizeBusinessDescriptor(descriptor)
	if err != nil {
		return false, err
	}
	business, err := s.storage.business.FindOne(ctx,
		s.storage.business.ScopeWorkspaceID(actor.WorkspaceID),
		s.storage.business.ScopeEquals(BusinessSchema.Descriptor, norm),
	)
	if err != nil {
		if database.IsRecordNotFound(err) {
			return true, nil
		}
		return false, err
	}
	return business == nil, nil
}

func (s *Service) CountBusinesses(ctx context.Context, actor *account.User) (int64, error) {
	if err := actor.Role.HasPermission(role.ActionView, role.ResourceBusiness); err != nil {
		return 0, err
	}
	return s.storage.business.Count(ctx, s.storage.business.ScopeWorkspaceID(actor.WorkspaceID))
}

func (s *Service) CountActiveBusinesses(ctx context.Context, actor *account.User) (int64, error) {
	if err := actor.Role.HasPermission(role.ActionView, role.ResourceBusiness); err != nil {
		return 0, err
	}
	return s.storage.business.Count(ctx,
		s.storage.business.ScopeWorkspaceID(actor.WorkspaceID),
		s.storage.business.ScopeIsNull(BusinessSchema.ArchivedAt),
	)
}

func (s *Service) MaxBusinessesEnforceFunc(ctx context.Context, actor *account.User, businessID string) (int64, error) {
	return s.CountActiveBusinesses(ctx, actor)
}
