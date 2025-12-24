package country_test

import (
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/utils/country"
	"github.com/stretchr/testify/require"
)

func TestCountries_ReturnsSortedCopy(t *testing.T) {
	t.Parallel()

	list1 := country.Countries()
	require.NotEmpty(t, list1)

	// Sorted ascending by name
	for i := 1; i < len(list1); i++ {
		require.LessOrEqual(t, list1[i-1].Name, list1[i].Name)
	}

	// Returned slice is a copy (callers can't mutate shared data)
	originalFirst := list1[0].Name
	list1[0].Name = "ZZZ"
	list2 := country.Countries()
	require.NotEmpty(t, list2)
	require.Equal(t, originalFirst, list2[0].Name)
}
