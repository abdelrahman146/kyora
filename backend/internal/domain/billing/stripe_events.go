package billing

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"gorm.io/gorm"
)

// StripeEvent persists processed webhook events for idempotency and replay protection.
// Only the event ID and type plus processed timestamp are stored.
const (
	StripeEventTable  = "stripe_events"
	StripeEventPrefix = "stev"
)

type StripeEvent struct {
	gorm.Model
	ID          string    `json:"id" gorm:"column:id;primaryKey;type:text"`
	EventID     string    `json:"eventId" gorm:"column:event_id;type:text;uniqueIndex;not null"`
	Type        string    `json:"type" gorm:"column:type;type:text;not null"`
	ProcessedAt time.Time `json:"processedAt" gorm:"column:processed_at;type:timestamp;not null"`
}

func (e *StripeEvent) TableName() string { return StripeEventTable }

func (e *StripeEvent) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == "" {
		e.ID = id.KsuidWithPrefix(StripeEventPrefix)
	}
	if e.ProcessedAt.IsZero() {
		e.ProcessedAt = time.Now()
	}
	return nil
}

var StripeEventSchema = struct {
	ID      schema.Field
	EventID schema.Field
	Type    schema.Field
}{
	ID:      schema.NewField("id", "id"),
	EventID: schema.NewField("event_id", "eventId"),
	Type:    schema.NewField("type", "type"),
}
