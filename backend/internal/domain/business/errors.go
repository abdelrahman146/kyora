package business

import (
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
)

// ErrBusinessNotFound indicates that a business could not be found in the authenticated workspace.
func ErrBusinessNotFound(businessDescriptor string, err error) error {
	return problem.NotFound("business not found").WithError(err).With("businessDescriptor", businessDescriptor).WithCode("business.not_found")
}

// ErrBusinessDescriptorAlreadyTaken indicates the requested descriptor is already used within the workspace.
func ErrBusinessDescriptorAlreadyTaken(descriptor string, err error) error {
	return problem.Conflict("business descriptor is already taken").WithError(err).With("descriptor", descriptor).WithCode("business.descriptor_taken")
}

// ErrInvalidBusinessDescriptor indicates the descriptor doesn't match allowed format.
func ErrInvalidBusinessDescriptor(descriptor string) error {
	return problem.BadRequest("invalid descriptor").With("descriptor", descriptor).With("hint", "use lowercase letters, numbers, and hyphens").WithCode("business.descriptor_invalid")
}

// ErrBusinessRateLimited indicates the client exceeded rate limits.
func ErrBusinessRateLimited() error {
	return problem.TooManyRequests("too many requests").WithCode("business.rate_limited")
}

func ErrShippingZoneNotFound(zoneID string, err error) error {
	return problem.NotFound("shipping zone not found").WithError(err).With("zoneId", zoneID).WithCode("business.shipping_zone_not_found")
}

func ErrShippingZoneNameAlreadyTaken(name string, err error) error {
	return problem.Conflict("shipping zone name is already taken").WithError(err).With("name", name).WithCode("business.shipping_zone_name_taken")
}

// Service validation errors

func ErrDescriptorRequired() error {
	return problem.BadRequest("descriptor is required").With("field", "descriptor").WithCode("business.descriptor_required")
}

func ErrStorefrontIdRequired() error {
	return problem.BadRequest("storefrontId is required").WithCode("business.storefront_id_required")
}

func ErrBusinessInputRequired() error {
	return problem.BadRequest("input is required").WithCode("business.input_required")
}

func ErrInvalidCountryCode() error {
	return problem.BadRequest("invalid countryCode").With("field", "countryCode").WithCode("business.invalid_country_code")
}

func ErrInvalidCurrency() error {
	return problem.BadRequest("invalid currency").With("field", "currency").WithCode("business.invalid_currency")
}

func ErrCountriesRequired() error {
	return problem.BadRequest("countries is required").With("field", "countries").WithCode("business.countries_required")
}

func ErrBusinessDataRequired() error {
	return problem.BadRequest("business is required").WithCode("business.data_required")
}

func ErrZoneIdRequired() error {
	return problem.BadRequest("zoneId is required").WithCode("business.zone_id_required")
}

func ErrShippingZoneRequestRequired() error {
	return problem.BadRequest("request is required").WithCode("business.shipping_zone_request_required")
}

func ErrShippingZoneNameRequired() error {
	return problem.BadRequest("name is required").With("field", "name").WithCode("business.shipping_zone_name_required")
}

func ErrShippingZoneNameTooLong() error {
	return problem.BadRequest("name is too long").With("field", "name").WithCode("business.shipping_zone_name_too_long")
}

func ErrShippingCostNegative() error {
	return problem.BadRequest("shippingCost cannot be negative").With("field", "shippingCost").WithCode("business.shipping_cost_negative")
}

func ErrFreeShippingThresholdNegative() error {
	return problem.BadRequest("freeShippingThreshold cannot be negative").With("field", "freeShippingThreshold").WithCode("business.free_shipping_threshold_negative")
}

func ErrShippingZoneNameEmpty() error {
	return problem.BadRequest("name cannot be empty").With("field", "name").WithCode("business.shipping_zone_name_empty")
}
