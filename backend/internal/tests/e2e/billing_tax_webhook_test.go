package e2e_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

type BillingTaxWebhookSuite struct {
	suite.Suite
	helper *BillingTestHelper
}

func (s *BillingTaxWebhookSuite) SetupSuite() {
	s.helper = NewBillingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *BillingTaxWebhookSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database,
		"plans",
		"subscriptions",
		"billing_invoice_records",
		"stripe_events",
		"users",
		"workspaces",
	)
	s.NoError(err)
}

func (s *BillingTaxWebhookSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database,
		"plans",
		"subscriptions",
		"billing_invoice_records",
		"stripe_events",
		"users",
		"workspaces",
	)
	s.NoError(err)
}

func (s *BillingTaxWebhookSuite) TestCalculateTax_Auth_AndValidation() {
	ctx := s.T().Context()

	_, _, adminToken := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("admin"), role.RoleAdmin)
	_, _, memberToken := s.helper.CreateTestUser(ctx, s.helper.UniqueEmail("member"), role.RoleUser)

	// Unauthorized
	respUnauthed, err := s.helper.Client().Post("/v1/billing/tax/calculate", map[string]interface{}{"amount": 1000, "currency": "usd"})
	s.NoError(err)
	defer respUnauthed.Body.Close()
	s.Equal(http.StatusUnauthorized, respUnauthed.StatusCode)

	// Validation
	respBad, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/tax/calculate", map[string]interface{}{"amount": 0, "currency": "usd"}, adminToken)
	s.NoError(err)
	defer respBad.Body.Close()
	s.Equal(http.StatusBadRequest, respBad.StatusCode)

	// Member can view billing, so should not be forbidden.
	resp, err := s.helper.Client().AuthenticatedRequest("POST", "/v1/billing/tax/calculate", map[string]interface{}{"amount": 1000, "currency": "usd"}, memberToken)
	s.NoError(err)
	defer resp.Body.Close()

	// stripe-mock may not implement tax calculations; assert non-flaky contract:
	// - if it works: 200 + has an id
	// - if unsupported: 500 (stripe operation failed)
	if resp.StatusCode == http.StatusOK {
		var calc map[string]interface{}
		s.NoError(testutils.DecodeJSON(resp, &calc))
		s.Contains(calc, "id")
	} else {
		s.Equal(http.StatusInternalServerError, resp.StatusCode)
	}
}

func stripeTestSignatureHeader(secret string, payload []byte, timestamp int64) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strconv.FormatInt(timestamp, 10)))
	mac.Write([]byte("."))
	mac.Write(payload)
	sig := hex.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("t=%d,v1=%s", timestamp, sig)
}

func (s *BillingTaxWebhookSuite) TestStripeWebhook_Signature_AndIdempotency() {
	payload := []byte(`{"id":"evt_test_123","type":"kyora.unhandled","data":{"object":{"foo":"bar"}}}`)

	// Missing signature
	respMissing, err := s.helper.Client().PostRaw("/webhooks/stripe", payload, map[string]string{"Content-Type": "application/json"})
	s.NoError(err)
	defer respMissing.Body.Close()
	s.Equal(http.StatusBadRequest, respMissing.StatusCode)

	// Invalid signature
	respInvalid, err := s.helper.Client().PostRaw("/webhooks/stripe", payload, map[string]string{
		"Content-Type":     "application/json",
		"Stripe-Signature": "t=1,v1=deadbeef",
	})
	s.NoError(err)
	defer respInvalid.Body.Close()
	s.Equal(http.StatusBadRequest, respInvalid.StatusCode)

	// Valid signature + idempotency (send twice)
	ts := time.Now().Unix()
	sig := stripeTestSignatureHeader("whsec_test", payload, ts)

	respOk1, err := s.helper.Client().PostRaw("/webhooks/stripe", payload, map[string]string{
		"Content-Type":     "application/json",
		"Stripe-Signature": sig,
	})
	s.NoError(err)
	defer respOk1.Body.Close()
	s.Equal(http.StatusOK, respOk1.StatusCode)

	respOk2, err := s.helper.Client().PostRaw("/webhooks/stripe", payload, map[string]string{
		"Content-Type":     "application/json",
		"Stripe-Signature": sig,
	})
	s.NoError(err)
	defer respOk2.Body.Close()
	s.Equal(http.StatusOK, respOk2.StatusCode)
}

func TestBillingTaxWebhookSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(BillingTaxWebhookSuite))
}
