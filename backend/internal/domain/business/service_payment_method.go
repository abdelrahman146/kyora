package business

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/throttle"
	"github.com/shopspring/decimal"
)

// BusinessPaymentMethodView is the effective view (global catalog + per-business overrides).
type BusinessPaymentMethodView struct {
	Descriptor        PaymentMethodDescriptor `json:"descriptor"`
	Name              string                  `json:"name"`
	LogoURL           string                  `json:"logoUrl"`
	Enabled           bool                    `json:"enabled"`
	FeePercent        decimal.Decimal         `json:"feePercent"`
	FeeFixed          decimal.Decimal         `json:"feeFixed"`
	DefaultFeePercent decimal.Decimal         `json:"defaultFeePercent"`
	DefaultFeeFixed   decimal.Decimal         `json:"defaultFeeFixed"`
}

func (s *Service) ListPaymentMethods(ctx context.Context, actor *account.User, biz *Business) ([]BusinessPaymentMethodView, error) {
	if err := actor.Role.HasPermission(role.ActionView, role.ResourceBusiness); err != nil {
		return nil, err
	}
	rows, err := s.storage.ListBusinessPaymentMethods(ctx, biz.ID)
	if err != nil {
		return nil, err
	}
	byDesc := make(map[PaymentMethodDescriptor]*BusinessPaymentMethod, len(rows))
	for _, r := range rows {
		byDesc[r.Descriptor] = r
	}

	defs := GlobalPaymentMethods()
	out := make([]BusinessPaymentMethodView, 0, len(defs))
	for _, d := range defs {
		v := BusinessPaymentMethodView{
			Descriptor:        d.Descriptor,
			Name:              d.Name,
			LogoURL:           d.LogoURL,
			Enabled:           d.DefaultEnabled,
			FeePercent:        d.DefaultFeePercent,
			FeeFixed:          d.DefaultFeeFixed,
			DefaultFeePercent: d.DefaultFeePercent,
			DefaultFeeFixed:   d.DefaultFeeFixed,
		}
		if r, ok := byDesc[d.Descriptor]; ok {
			v.Enabled = r.Enabled
			v.FeePercent = r.FeePercent
			v.FeeFixed = r.FeeFixed
		}
		out = append(out, v)
	}
	return out, nil
}

func (s *Service) UpdatePaymentMethod(ctx context.Context, actor *account.User, biz *Business, descriptor string, req *UpdateBusinessPaymentMethodRequest) (*BusinessPaymentMethodView, error) {
	if err := actor.Role.HasPermission(role.ActionManage, role.ResourceBusiness); err != nil {
		return nil, err
	}
	if !throttle.Allow(s.storage.cache, fmt.Sprintf("rl:payment_method:update:%s:%s", biz.ID, actor.ID), time.Minute, 120, 250*time.Millisecond) {
		return nil, ErrBusinessRateLimited()
	}
	if req == nil {
		return nil, problem.BadRequest("request is required")
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	desc := PaymentMethodDescriptor(strings.TrimSpace(strings.ToLower(descriptor)))
	def, ok := FindPaymentMethodDefinition(desc)
	if !ok {
		return nil, problem.BadRequest("invalid payment method descriptor").With("descriptor", descriptor)
	}

	// Load existing override (if any)
	existing, err := s.storage.GetBusinessPaymentMethodByDescriptor(ctx, biz.ID, desc)
	if err != nil {
		if !database.IsRecordNotFound(err) {
			return nil, err
		}
		existing = nil
	}

	pm := existing
	if pm == nil {
		pm = &BusinessPaymentMethod{
			BusinessID: biz.ID,
			Descriptor: desc,
			Enabled:    def.DefaultEnabled,
			FeePercent: def.DefaultFeePercent,
			FeeFixed:   def.DefaultFeeFixed,
		}
	}

	if req.Enabled != nil {
		pm.Enabled = *req.Enabled
	}
	if req.FeePercent.Valid {
		pm.FeePercent = req.FeePercent.Decimal
	}
	if req.FeeFixed.Valid {
		pm.FeeFixed = req.FeeFixed.Decimal
	}

	// Persist create/update.
	if pm.ID == "" {
		if err := s.storage.payment.CreateOne(ctx, pm); err != nil {
			return nil, err
		}
	} else {
		if err := s.storage.payment.UpdateOne(ctx, pm); err != nil {
			return nil, err
		}
	}

	view := &BusinessPaymentMethodView{
		Descriptor:        def.Descriptor,
		Name:              def.Name,
		LogoURL:           def.LogoURL,
		Enabled:           pm.Enabled,
		FeePercent:        pm.FeePercent,
		FeeFixed:          pm.FeeFixed,
		DefaultFeePercent: def.DefaultFeePercent,
		DefaultFeeFixed:   def.DefaultFeeFixed,
	}
	return view, nil
}

// GetEffectivePaymentMethodFee returns the effective fee configuration for a business.
// It is intended for internal automation (e.g., creating transaction fee expenses) and does not enforce role permissions.
func (s *Service) GetEffectivePaymentMethodFee(ctx context.Context, businessID string, descriptor PaymentMethodDescriptor) (enabled bool, feePercent decimal.Decimal, feeFixed decimal.Decimal, err error) {
	def, ok := FindPaymentMethodDefinition(descriptor)
	if !ok {
		return false, decimal.Zero, decimal.Zero, problem.BadRequest("invalid payment method descriptor").With("descriptor", descriptor)
	}
	// Defaults.
	enabled = def.DefaultEnabled
	feePercent = def.DefaultFeePercent
	feeFixed = def.DefaultFeeFixed

	pm, err := s.storage.GetBusinessPaymentMethodByDescriptor(ctx, businessID, descriptor)
	if err == nil && pm != nil {
		enabled = pm.Enabled
		feePercent = pm.FeePercent
		feeFixed = pm.FeeFixed
	}
	return enabled, feePercent, feeFixed, nil
}
