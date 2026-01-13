package analytics

import "github.com/abdelrahman146/kyora/internal/platform/types/problem"

func ErrInvalidDateFormat(field string, err error) error {
	return problem.BadRequest("invalid "+field+" date format, use YYYY-MM-DD").
		WithError(err).
		With("field", field).
		WithCode("analytics.invalid_date_format")
}

func ErrInvalidQueryParams(err error) error {
	return problem.BadRequest("invalid query parameters").
		WithError(err).
		WithCode("analytics.invalid_query_params")
}

func ErrInvalidDateRange(from, to string) error {
	return problem.BadRequest("to must be on or after from").
		With("from", from).
		With("to", to).
		WithCode("analytics.invalid_date_range")
}

func ErrAnalyticsQueryFailed(err error) error {
	return problem.InternalError().
		WithError(err).
		WithCode("analytics.query_failed")
}
