package business

import "github.com/shopspring/decimal"

// PaymentMethodDescriptor is the canonical identifier for a payment method.
// It is stable and used in URLs and persisted overrides.
type PaymentMethodDescriptor string

const (
	PaymentMethodCashOnDelivery PaymentMethodDescriptor = "cash_on_delivery"
	PaymentMethodBankTransfer   PaymentMethodDescriptor = "bank_transfer"
	PaymentMethodCreditCard     PaymentMethodDescriptor = "credit_card"
	PaymentMethodTamara         PaymentMethodDescriptor = "tamara"
	PaymentMethodTabby          PaymentMethodDescriptor = "tabby"
	PaymentMethodPayPal         PaymentMethodDescriptor = "paypal"
)

// PaymentMethodDefinition is the global payment method catalog entry.
// Businesses can enable/disable entries and override fee values per business.
type PaymentMethodDefinition struct {
	Descriptor        PaymentMethodDescriptor `json:"descriptor"`
	Name              string                  `json:"name"`
	LogoURL           string                  `json:"logoUrl"`
	DefaultFeePercent decimal.Decimal         `json:"defaultFeePercent"`
	DefaultFeeFixed   decimal.Decimal         `json:"defaultFeeFixed"`
	DefaultEnabled    bool                    `json:"defaultEnabled"`
}

// GlobalPaymentMethods returns the global catalog of payment methods.
// Keep this list small and static to avoid extra DB tables.
func GlobalPaymentMethods() []PaymentMethodDefinition {
	return []PaymentMethodDefinition{
		{
			Descriptor:        PaymentMethodCashOnDelivery,
			Name:              "Cash on delivery",
			LogoURL:           "",
			DefaultFeePercent: decimal.Zero,
			DefaultFeeFixed:   decimal.Zero,
			DefaultEnabled:    true,
		},
		{
			Descriptor:        PaymentMethodBankTransfer,
			Name:              "Bank transfer",
			LogoURL:           "",
			DefaultFeePercent: decimal.Zero,
			DefaultFeeFixed:   decimal.Zero,
			DefaultEnabled:    true,
		},
		{
			Descriptor:        PaymentMethodCreditCard,
			Name:              "Credit card",
			LogoURL:           "",
			DefaultFeePercent: decimal.Zero,
			DefaultFeeFixed:   decimal.Zero,
			DefaultEnabled:    false,
		},
		{
			Descriptor:        PaymentMethodTamara,
			Name:              "Tamara",
			LogoURL:           "",
			DefaultFeePercent: decimal.Zero,
			DefaultFeeFixed:   decimal.Zero,
			DefaultEnabled:    false,
		},
		{
			Descriptor:        PaymentMethodTabby,
			Name:              "Tabby",
			LogoURL:           "",
			DefaultFeePercent: decimal.Zero,
			DefaultFeeFixed:   decimal.Zero,
			DefaultEnabled:    false,
		},
		{
			Descriptor:        PaymentMethodPayPal,
			Name:              "PayPal",
			LogoURL:           "",
			DefaultFeePercent: decimal.Zero,
			DefaultFeeFixed:   decimal.Zero,
			DefaultEnabled:    false,
		},
	}
}

// FindPaymentMethodDefinition looks up a payment method in the global catalog.
func FindPaymentMethodDefinition(descriptor PaymentMethodDescriptor) (*PaymentMethodDefinition, bool) {
	for _, m := range GlobalPaymentMethods() {
		if m.Descriptor == descriptor {
			mm := m
			return &mm, true
		}
	}
	return nil, false
}
