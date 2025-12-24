package database_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestIsRetryableTxError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "pgx serialization failure",
			err:  &pgconn.PgError{Code: "40001"},
			want: true,
		},
		{
			name: "pgx deadlock",
			err:  &pgconn.PgError{Code: "40P01"},
			want: true,
		},
		{
			name: "pgx lock not available",
			err:  &pgconn.PgError{Code: "55P03"},
			want: true,
		},
		{
			name: "pq serialization failure",
			err:  &pq.Error{Code: pq.ErrorCode("40001")},
			want: true,
		},
		{
			name: "wrapped pgx error",
			err:  fmt.Errorf("wrapped: %w", &pgconn.PgError{Code: "40P01"}),
			want: true,
		},
		{
			name: "wrapped pq error",
			err:  fmt.Errorf("wrapped: %w", &pq.Error{Code: pq.ErrorCode("55P03")}),
			want: true,
		},
		{
			name: "fallback sqlstate string",
			err:  errors.New("pq: could not serialize access due to concurrent update (SQLSTATE 40001)"),
			want: true,
		},
		{
			name: "non-retryable",
			err:  &pgconn.PgError{Code: "23505"},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, database.IsRetryableTxError(tt.err))
		})
	}
}
