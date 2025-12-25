package billing

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
)

type Storage struct {
	cache         *cache.Cache
	db            *database.Database
	plan          *database.Repository[Plan]
	subscription  *database.Repository[Subscription]
	invoiceRecord *database.Repository[InvoiceRecord]
	event         *database.Repository[StripeEvent]
}

func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
	s := &Storage{
		cache:         cache,
		db:            db,
		plan:          database.NewRepository[Plan](db),
		subscription:  database.NewRepository[Subscription](db),
		invoiceRecord: database.NewRepository[InvoiceRecord](db),
		event:         database.NewRepository[StripeEvent](db),
	}

	return s
}

func (s *Storage) FindInvoiceRecordByWorkspaceAndStripeID(ctx context.Context, workspaceID, stripeInvoiceID string) (*InvoiceRecord, error) {
	return s.invoiceRecord.FindOne(
		ctx,
		s.invoiceRecord.ScopeEquals(InvoiceRecordSchema.WorkspaceID, workspaceID),
		s.invoiceRecord.ScopeEquals(InvoiceRecordSchema.StripeInvoiceID, stripeInvoiceID),
	)
}

func (s *Storage) UpsertInvoiceRecord(ctx context.Context, workspaceID, stripeInvoiceID, hostedURL, invoicePDF string) error {
	rec, err := s.FindInvoiceRecordByWorkspaceAndStripeID(ctx, workspaceID, stripeInvoiceID)
	if err != nil && !database.IsRecordNotFound(err) {
		return err
	}
	if rec == nil {
		rec = &InvoiceRecord{WorkspaceID: workspaceID, StripeInvoiceID: stripeInvoiceID, HostedInvoiceURL: hostedURL, InvoicePDF: invoicePDF}
		return s.invoiceRecord.CreateOne(ctx, rec)
	}
	rec.HostedInvoiceURL = hostedURL
	rec.InvoicePDF = invoicePDF
	return s.invoiceRecord.UpdateOne(ctx, rec)
}

// SyncPlans upserts all defined plans into the database
// This can be called manually via CLI command or automatically on startup
func (s *Storage) SyncPlans(ctx context.Context) error {
	logger.FromContext(ctx).Info("Syncing plans to database", "planCount", len(plans))

	for _, template := range plans {
		plan := template // copy
		existing, err := s.plan.FindOne(ctx, s.plan.ScopeEquals(PlanSchema.Descriptor, plan.Descriptor))
		if err != nil && !database.IsRecordNotFound(err) {
			logger.FromContext(ctx).Error("failed to fetch existing plan", "error", err, "descriptor", plan.Descriptor)
			continue
		}
		if existing == nil {
			// Keep StripePlanID NULL until Stripe sync populates it.
			if err := s.plan.CreateOne(ctx, &plan); err != nil {
				logger.FromContext(ctx).Error("failed to create plan", "error", err, "descriptor", plan.Descriptor)
				continue
			}
			logger.FromContext(ctx).Info("Plan created", "descriptor", plan.Descriptor, "id", plan.ID)
		} else {
			plan.ID = existing.ID
			// Never overwrite an existing Stripe price ID during a DB-only sync.
			if plan.StripePlanID == nil {
				plan.StripePlanID = existing.StripePlanID
			}
			if err := s.plan.UpdateOne(ctx, &plan); err != nil {
				logger.FromContext(ctx).Error("failed to update plan", "error", err, "descriptor", plan.Descriptor)
				continue
			}
			logger.FromContext(ctx).Info("Plan updated", "descriptor", plan.Descriptor, "id", plan.ID)
		}
	}

	logger.FromContext(ctx).Info("Plan sync completed")
	return nil
}

// CountBusinessesByWorkspace returns count of businesses for a workspace
func (s *Storage) CountBusinessesByWorkspace(ctx context.Context, workspaceID string) (int64, error) {
	var count int64
	err := s.db.Conn(ctx).Table("businesses").Where("workspace_id = ?", workspaceID).Count(&count).Error
	return count, err
}

// CountMonthlyOrdersByWorkspace counts orders within date range across all businesses in a workspace
func (s *Storage) CountMonthlyOrdersByWorkspace(ctx context.Context, workspaceID string, from, to time.Time) (int64, error) {
	var count int64
	q := s.db.Conn(ctx).Table("orders as o").
		Joins("join businesses b on b.id = o.business_id").
		Where("b.workspace_id = ?", workspaceID)
	if !from.IsZero() && !to.IsZero() {
		q = q.Where("o.ordered_at BETWEEN ? AND ?", from, to)
	} else if !from.IsZero() {
		q = q.Where("o.ordered_at >= ?", from)
	} else if !to.IsZero() {
		q = q.Where("o.ordered_at <= ?", to)
	}
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
