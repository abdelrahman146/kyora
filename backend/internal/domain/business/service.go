package business

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/throttle"
	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
	"github.com/shopspring/decimal"
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

// GetBusinessByStorefrontPublicID returns a business by its public storefront ID.
// This is used by unauthenticated storefront endpoints.
// It only returns enabled, non-archived businesses.
func (s *Service) GetBusinessByStorefrontPublicID(ctx context.Context, storefrontPublicID string) (*Business, error) {
	id := strings.TrimSpace(storefrontPublicID)
	if id == "" {
		return nil, problem.BadRequest("storefrontId is required")
	}
	biz, err := s.storage.business.FindOne(ctx,
		s.storage.business.ScopeEquals(BusinessSchema.StorefrontPublicID, id),
		s.storage.business.ScopeEquals(BusinessSchema.StorefrontEnabled, true),
		s.storage.business.ScopeIsNull(BusinessSchema.ArchivedAt),
	)
	if err != nil {
		return nil, ErrBusinessNotFound(id, err)
	}
	return biz, nil
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

	var created *Business
	err = s.atomicProcessor.Exec(ctx, func(tctx context.Context) error {
		biz := &Business{
			WorkspaceID:       actor.WorkspaceID,
			Descriptor:        normDescriptor,
			Name:              input.Name,
			Brand:             strings.TrimSpace(input.Brand),
			CountryCode:       country,
			VatRate:           input.VatRate,
			Currency:          currency,
			StorefrontEnabled: input.StorefrontEnabled,
			StorefrontTheme:   input.StorefrontTheme,
			SupportEmail:      strings.TrimSpace(input.SupportEmail),
			PhoneNumber:       strings.TrimSpace(input.PhoneNumber),
			WhatsappNumber:    strings.TrimSpace(input.WhatsappNumber),
			Address:           strings.TrimSpace(input.Address),
			WebsiteURL:        strings.TrimSpace(input.WebsiteURL),
			InstagramURL:      strings.TrimSpace(input.InstagramURL),
			FacebookURL:       strings.TrimSpace(input.FacebookURL),
			TikTokURL:         strings.TrimSpace(input.TikTokURL),
			XURL:              strings.TrimSpace(input.XURL),
			SnapchatURL:       strings.TrimSpace(input.SnapchatURL),
		}
		if !input.SafetyBuffer.IsZero() {
			biz.SafetyBuffer = input.SafetyBuffer
		}
		if !input.EstablishedAt.IsZero() {
			biz.EstablishedAt = input.EstablishedAt.Time
		}
		if err := s.storage.business.CreateOne(tctx, biz); err != nil {
			return err
		}
		// Always create a default shipping zone for ease-of-use.
		// Name=country, countries=[country], cost=free, currency=business currency.
		zone := &ShippingZone{
			BusinessID:            biz.ID,
			Name:                  biz.CountryCode,
			Countries:             CountryCodeList{biz.CountryCode},
			Currency:              biz.Currency,
			ShippingCost:          decimal.Zero,
			FreeShippingThreshold: decimal.Zero,
		}
		if err := s.storage.CreateShippingZone(tctx, zone); err != nil {
			return err
		}
		created = biz
		return nil
	}, atomic.WithIsolationLevel(atomic.LevelSerializable), atomic.WithRetries(3))
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Service) normalizeZoneCountries(in []string) (CountryCodeList, error) {
	if len(in) == 0 {
		return nil, problem.BadRequest("countries is required").With("field", "countries")
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, c := range in {
		cc := normalizeCountryCode(c)
		if len(cc) != 2 {
			return nil, problem.BadRequest("invalid country code").With("field", "countries")
		}
		if _, ok := seen[cc]; ok {
			continue
		}
		seen[cc] = struct{}{}
		out = append(out, cc)
	}
	if len(out) == 0 {
		return nil, problem.BadRequest("countries is required").With("field", "countries")
	}
	return CountryCodeList(out), nil
}

func (s *Service) ListShippingZones(ctx context.Context, actor *account.User, biz *Business) ([]*ShippingZone, error) {
	if err := actor.Role.HasPermission(role.ActionView, role.ResourceBusiness); err != nil {
		return nil, err
	}
	return s.storage.ListShippingZones(ctx, biz.ID)
}

func (s *Service) GetShippingZoneByID(ctx context.Context, actor *account.User, biz *Business, zoneID string) (*ShippingZone, error) {
	if err := actor.Role.HasPermission(role.ActionView, role.ResourceBusiness); err != nil {
		return nil, err
	}
	if zoneID == "" {
		return nil, problem.BadRequest("zoneId is required")
	}
	zone, err := s.storage.GetShippingZoneByID(ctx, biz.ID, zoneID)
	if err != nil {
		return nil, ErrShippingZoneNotFound(zoneID, err)
	}
	return zone, nil
}

func (s *Service) CreateShippingZone(ctx context.Context, actor *account.User, biz *Business, req *CreateShippingZoneRequest) (*ShippingZone, error) {
	if err := actor.Role.HasPermission(role.ActionManage, role.ResourceBusiness); err != nil {
		return nil, err
	}
	if !throttle.Allow(s.storage.cache, fmt.Sprintf("rl:shipping_zone:create:%s:%s", biz.ID, actor.ID), time.Minute, 60, 1*time.Second) {
		return nil, ErrBusinessRateLimited()
	}
	if req == nil {
		return nil, problem.BadRequest("request is required")
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, problem.BadRequest("name is required").With("field", "name")
	}
	if len(name) > 80 {
		return nil, problem.BadRequest("name is too long").With("field", "name")
	}
	if req.ShippingCost.LessThan(decimal.Zero) {
		return nil, problem.BadRequest("shippingCost cannot be negative").With("field", "shippingCost")
	}
	if req.FreeShippingThreshold.LessThan(decimal.Zero) {
		return nil, problem.BadRequest("freeShippingThreshold cannot be negative").With("field", "freeShippingThreshold")
	}
	countries, err := s.normalizeZoneCountries(req.Countries)
	if err != nil {
		return nil, err
	}
	zone := &ShippingZone{
		BusinessID:            biz.ID,
		Name:                  name,
		Countries:             countries,
		Currency:              biz.Currency,
		ShippingCost:          req.ShippingCost,
		FreeShippingThreshold: req.FreeShippingThreshold,
	}
	if err := s.storage.CreateShippingZone(ctx, zone); err != nil {
		if database.IsUniqueViolation(err) {
			return nil, ErrShippingZoneNameAlreadyTaken(name, err)
		}
		return nil, err
	}
	return zone, nil
}

func (s *Service) UpdateShippingZone(ctx context.Context, actor *account.User, biz *Business, zoneID string, req *UpdateShippingZoneRequest) (*ShippingZone, error) {
	if err := actor.Role.HasPermission(role.ActionManage, role.ResourceBusiness); err != nil {
		return nil, err
	}
	if !throttle.Allow(s.storage.cache, fmt.Sprintf("rl:shipping_zone:update:%s:%s", biz.ID, actor.ID), time.Minute, 120, 1*time.Second) {
		return nil, ErrBusinessRateLimited()
	}
	if req == nil {
		return nil, problem.BadRequest("request is required")
	}
	zone, err := s.storage.GetShippingZoneByID(ctx, biz.ID, zoneID)
	if err != nil {
		return nil, ErrShippingZoneNotFound(zoneID, err)
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, problem.BadRequest("name cannot be empty").With("field", "name")
		}
		if len(name) > 80 {
			return nil, problem.BadRequest("name is too long").With("field", "name")
		}
		zone.Name = name
	}
	if req.Countries != nil {
		countries, err := s.normalizeZoneCountries(req.Countries)
		if err != nil {
			return nil, err
		}
		zone.Countries = countries
	}
	if req.ShippingCost.Valid {
		if req.ShippingCost.Decimal.LessThan(decimal.Zero) {
			return nil, problem.BadRequest("shippingCost cannot be negative").With("field", "shippingCost")
		}
		zone.ShippingCost = transformer.FromNullDecimal(req.ShippingCost)
	}
	if req.FreeShippingThreshold.Valid {
		if req.FreeShippingThreshold.Decimal.LessThan(decimal.Zero) {
			return nil, problem.BadRequest("freeShippingThreshold cannot be negative").With("field", "freeShippingThreshold")
		}
		zone.FreeShippingThreshold = transformer.FromNullDecimal(req.FreeShippingThreshold)
	}
	zone.Currency = biz.Currency
	if err := s.storage.UpdateShippingZone(ctx, zone); err != nil {
		if database.IsUniqueViolation(err) {
			return nil, ErrShippingZoneNameAlreadyTaken(zone.Name, err)
		}
		return nil, err
	}
	return zone, nil
}

func (s *Service) DeleteShippingZone(ctx context.Context, actor *account.User, biz *Business, zoneID string) error {
	if err := actor.Role.HasPermission(role.ActionManage, role.ResourceBusiness); err != nil {
		return err
	}
	if !throttle.Allow(s.storage.cache, fmt.Sprintf("rl:shipping_zone:delete:%s:%s", biz.ID, actor.ID), time.Minute, 60, 1*time.Second) {
		return ErrBusinessRateLimited()
	}
	zone, err := s.storage.GetShippingZoneByID(ctx, biz.ID, zoneID)
	if err != nil {
		return ErrShippingZoneNotFound(zoneID, err)
	}
	return s.storage.DeleteShippingZone(ctx, zone)
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
	if input.Brand != nil {
		business.Brand = strings.TrimSpace(*input.Brand)
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
	if input.StorefrontEnabled != nil {
		business.StorefrontEnabled = *input.StorefrontEnabled
	}
	if input.StorefrontTheme != nil {
		business.StorefrontTheme = *input.StorefrontTheme
	}
	if input.SupportEmail != nil {
		business.SupportEmail = strings.TrimSpace(*input.SupportEmail)
	}
	if input.PhoneNumber != nil {
		business.PhoneNumber = strings.TrimSpace(*input.PhoneNumber)
	}
	if input.WhatsappNumber != nil {
		business.WhatsappNumber = strings.TrimSpace(*input.WhatsappNumber)
	}
	if input.Address != nil {
		business.Address = strings.TrimSpace(*input.Address)
	}
	if input.WebsiteURL != nil {
		business.WebsiteURL = strings.TrimSpace(*input.WebsiteURL)
	}
	if input.InstagramURL != nil {
		business.InstagramURL = strings.TrimSpace(*input.InstagramURL)
	}
	if input.FacebookURL != nil {
		business.FacebookURL = strings.TrimSpace(*input.FacebookURL)
	}
	if input.TikTokURL != nil {
		business.TikTokURL = strings.TrimSpace(*input.TikTokURL)
	}
	if input.XURL != nil {
		business.XURL = strings.TrimSpace(*input.XURL)
	}
	if input.SnapchatURL != nil {
		business.SnapchatURL = strings.TrimSpace(*input.SnapchatURL)
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
