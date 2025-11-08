package billing

import (
	"context"
	"fmt"

	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/shopspring/decimal"
	stripelib "github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/price"
	"github.com/stripe/stripe-go/v83/product"
)

// SyncPlansToStripe ensures all local plans exist in Stripe (products + prices)
func (s *Service) SyncPlansToStripe(ctx context.Context) error {
	plans, err := s.storage.plan.FindMany(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch plans: %w", err)
	}
	logger.FromContext(ctx).Info("Starting plan sync to Stripe", "planCount", len(plans))
	for _, p := range plans {
		if err := s.syncSinglePlanToStripe(ctx, p); err != nil {
			logger.FromContext(ctx).Error("Failed to sync plan", "error", err, "planId", p.ID, "descriptor", p.Descriptor)
			continue
		}
	}
	logger.FromContext(ctx).Info("Completed plan sync to Stripe")
	return nil
}

// ensurePlanSynced guarantees a single plan has a valid StripePrice ID before use.
// It is invoked prior to subscription / checkout operations.
func (s *Service) ensurePlanSynced(ctx context.Context, p *Plan) error {
	if p == nil {
		return fmt.Errorf("plan is nil")
	}
	// Check cache first
	if p.StripePlanID != "" {
		if cached, ok := s.cache.get(p.ID); ok && cached == p.StripePlanID {
			return nil
		}
	}
	if p.StripePlanID == "" {
		// Full sync path
		if err := s.syncSinglePlanToStripe(ctx, p); err != nil {
			return err
		}
		if p.StripePlanID == "" { // still empty — treat as error
			return fmt.Errorf("plan %s not synced to Stripe (missing price id)", p.ID)
		}
		s.cache.set(p.ID, p.StripePlanID)
		return nil
	}
	// Validate existing price still conforms. If mismatch create new one.
	prod, err := s.findOrCreateProduct(ctx, p)
	if err != nil {
		return err
	}
	needNew, err := s.validateExistingPrice(ctx, p, prod.ID)
	if err != nil {
		return err
	}
	if needNew {
		newPrice, err := s.createPrice(ctx, p, prod.ID)
		if err != nil {
			return err
		}
		p.StripePlanID = newPrice.ID
		if err := s.storage.plan.UpdateOne(ctx, p); err != nil {
			return fmt.Errorf("failed updating plan price id: %w", err)
		}
		s.cache.set(p.ID, p.StripePlanID)
	} else {
		s.cache.set(p.ID, p.StripePlanID)
	}
	return nil
}

// syncSinglePlanToStripe handles syncing a single plan with proper error handling and conflict resolution
func (s *Service) syncSinglePlanToStripe(ctx context.Context, p *Plan) error {
	prod, err := s.findOrCreateProduct(ctx, p)
	if err != nil {
		return fmt.Errorf("failed to find or create product: %w", err)
	}
	needNewPrice, err := s.validateExistingPrice(ctx, p, prod.ID)
	if err != nil {
		return fmt.Errorf("failed to validate existing price: %w", err)
	}
	if needNewPrice {
		newPrice, err := s.createPrice(ctx, p, prod.ID)
		if err != nil {
			return fmt.Errorf("failed to create new price: %w", err)
		}
		p.StripePlanID = newPrice.ID
		if err := s.storage.plan.UpdateOne(ctx, p); err != nil {
			logger.FromContext(ctx).Error("Failed to update plan with new Stripe price ID", "error", err, "plan_id", p.ID, "price_id", newPrice.ID)
			return fmt.Errorf("failed to update plan: %w", err)
		}
		logger.FromContext(ctx).Info("Created/updated price for plan", "plan_id", p.ID, "price_id", newPrice.ID, "amount", p.Price)
	}
	return nil
}

// findOrCreateProduct finds existing product by metadata or creates new one
func (s *Service) findOrCreateProduct(ctx context.Context, p *Plan) (*stripelib.Product, error) {
	if p.StripePlanID != "" {
		if pr, err := price.Get(p.StripePlanID, nil); err == nil && pr != nil && pr.Product != nil {
			if prod, err := product.Get(pr.Product.ID, nil); err == nil {
				if prod.Metadata != nil {
					if kyoraID, exists := prod.Metadata["kyora_plan_id"]; exists && kyoraID == p.ID {
						return prod, nil
					}
					if descriptor, exists := prod.Metadata["descriptor"]; exists && descriptor == p.Descriptor {
						return prod, nil
					}
				}
				return s.updateProductMetadata(ctx, prod, p)
			}
		}
	}
	existingProd, err := s.findProductByMetadata(p)
	if err != nil {
		return nil, fmt.Errorf("failed to search for existing product: %w", err)
	}
	if existingProd != nil {
		return existingProd, nil
	}
	return s.createProduct(ctx, p)
}

func (s *Service) findProductByMetadata(p *Plan) (*stripelib.Product, error) {
	params := &stripelib.ProductListParams{Active: stripelib.Bool(true)}
	params.Limit = stripelib.Int64(100)
	iter := product.List(params)
	for iter.Next() {
		prod := iter.Product()
		if prod.Metadata != nil {
			if kyoraID, exists := prod.Metadata["kyora_plan_id"]; exists && kyoraID == p.ID {
				return prod, nil
			}
			if descriptor, exists := prod.Metadata["descriptor"]; exists && descriptor == p.Descriptor {
				return prod, nil
			}
		}
	}
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}
	return nil, nil
}

func (s *Service) createProduct(ctx context.Context, p *Plan) (*stripelib.Product, error) {
	idempotencyKey := fmt.Sprintf("product_%s", p.ID)
	params := &stripelib.ProductParams{
		Name:        stripelib.String(p.Name),
		Description: stripelib.String(p.Description),
		Metadata: map[string]string{
			"kyora_plan_id": p.ID,
			"descriptor":    p.Descriptor,
		},
	}
	params.SetIdempotencyKey(idempotencyKey)
	prod, err := withStripeRetry(ctx, 3, func() (*stripelib.Product, error) { return product.New(params) })
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe product: %w", err)
	}
	logger.FromContext(ctx).Info("Created new Stripe product", "planId", p.ID, "productId", prod.ID, "name", p.Name)
	return prod, nil
}

func (s *Service) updateProductMetadata(ctx context.Context, prod *stripelib.Product, p *Plan) (*stripelib.Product, error) {
	params := &stripelib.ProductParams{
		Name:        stripelib.String(p.Name),
		Description: stripelib.String(p.Description),
		Metadata: map[string]string{
			"kyora_plan_id": p.ID,
			"descriptor":    p.Descriptor,
		},
	}
	updatedProd, err := withStripeRetry(ctx, 3, func() (*stripelib.Product, error) { return product.Update(prod.ID, params) })
	if err != nil {
		return nil, fmt.Errorf("failed to update product metadata: %w", err)
	}
	logger.FromContext(ctx).Info("Updated product metadata", "planId", p.ID, "productId", prod.ID)
	return updatedProd, nil
}

func (s *Service) validateExistingPrice(ctx context.Context, p *Plan, productID string) (bool, error) {
	if p.StripePlanID == "" {
		return true, nil
	}
	existingPrice, err := price.Get(p.StripePlanID, nil)
	if err != nil {
		logger.FromContext(ctx).Warn("Failed to fetch existing price, will create new one", "priceId", p.StripePlanID, "error", err)
		return true, nil
	}
	interval := "month"
	if p.BillingCycle == BillingCycleYearly {
		interval = "year"
	}
	expectedAmount := p.Price.Mul(decimal.NewFromInt(100)).IntPart()
	if string(existingPrice.Currency) != p.Currency || existingPrice.Recurring == nil || string(existingPrice.Recurring.Interval) != interval || existingPrice.UnitAmount != expectedAmount || existingPrice.Product == nil || existingPrice.Product.ID != productID {
		return true, nil
	}
	return false, nil
}

func (s *Service) createPrice(ctx context.Context, p *Plan, productID string) (*stripelib.Price, error) {
	interval := "month"
	if p.BillingCycle == BillingCycleYearly {
		interval = "year"
	}
	unitAmount := p.Price.Mul(decimal.NewFromInt(100)).IntPart()
	idempotencyKey := fmt.Sprintf("price_%s_%s_%d", p.ID, interval, unitAmount)
	params := &stripelib.PriceParams{
		Currency:   stripelib.String(p.Currency),
		UnitAmount: stripelib.Int64(unitAmount),
		Recurring:  &stripelib.PriceRecurringParams{Interval: stripelib.String(interval)},
		Product:    stripelib.String(productID),
		Metadata: map[string]string{
			"kyora_plan_id": p.ID,
			"descriptor":    p.Descriptor,
		},
	}
	params.SetIdempotencyKey(idempotencyKey)
	newPrice, err := withStripeRetry(ctx, 3, func() (*stripelib.Price, error) { return price.New(params) })
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe price: %w", err)
	}
	return newPrice, nil
}
