package storefront

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"gorm.io/gorm"
)

const (
	StorefrontRequestStruct = "StorefrontRequest"
	StorefrontRequestTable  = "storefront_requests"
	StorefrontRequestPrefix = "sreq"
)

var StorefrontRequestSchema = struct {
	BusinessID     schema.Field
	IdempotencyKey schema.Field
}{
	BusinessID:     schema.NewField("business_id", "businessId"),
	IdempotencyKey: schema.NewField("idempotency_key", "idempotencyKey"),
}

// StorefrontRequest is an idempotency record for public storefront operations.
//
// It ensures that repeated POSTs with the same idempotency key for the same business
// either return the same created order or fail with a conflict when the payload differs.
//
// Only minimal data is stored: request hash + created order ID.
// The full request payload is not persisted to avoid storing PII.
//
//nolint:lll // gorm tags are long by nature
type StorefrontRequest struct {
	gorm.Model

	ID             string `gorm:"column:id;primaryKey;type:text"`
	BusinessID     string `gorm:"column:business_id;type:text;not null;index:idx_storefront_req_business_key,unique"`
	IdempotencyKey string `gorm:"column:idempotency_key;type:text;not null;index:idx_storefront_req_business_key,unique"`
	RequestHash    Hash   `gorm:"column:request_hash;type:text;not null"`
	OrderID        string `gorm:"column:order_id;type:text"`
}

func (StorefrontRequest) TableName() string { return StorefrontRequestTable }

func (r *StorefrontRequest) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = id.KsuidWithPrefix(StorefrontRequestPrefix)
	}
	return nil
}

// Hash is a hex-encoded SHA-256 digest stored as text.
// It implements sql.Scanner / driver.Valuer via Scan/Value.
//
//nolint:revive // exported for json + db interoperability
type Hash [32]byte

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *Hash) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	if len(decoded) != len(h[:]) {
		return fmt.Errorf("invalid hash length: %d", len(decoded))
	}
	copy(h[:], decoded)
	return nil
}

func (h Hash) Value() (driver.Value, error) {
	return h.String(), nil
}

func (h *Hash) Scan(value any) error {
	if value == nil {
		return errors.New("hash is NULL")
	}
	s, ok := value.(string)
	if !ok {
		b, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("unsupported hash type: %T", value)
		}
		s = string(b)
	}
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	if len(decoded) != len(h[:]) {
		return fmt.Errorf("invalid hash length: %d", len(decoded))
	}
	copy(h[:], decoded)
	return nil
}
