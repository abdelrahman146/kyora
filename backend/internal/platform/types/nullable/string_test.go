package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/types/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullableString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    nullable.String
		expected string
	}{
		{
			name:     "valid string",
			input:    nullable.NewString("test value"),
			expected: `"test value"`,
		},
		{
			name:     "empty string becomes null",
			input:    nullable.NewString(""),
			expected: `null`,
		},
		{
			name:     "invalid (null) value",
			input:    nullable.String{},
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))
		})
	}
}

func TestNullableString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected nullable.String
	}{
		{
			name:     "valid string",
			input:    `"test value"`,
			expected: nullable.NewString("test value"),
		},
		{
			name:     "null value",
			input:    `null`,
			expected: nullable.String{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result nullable.String
			err := json.Unmarshal([]byte(tt.input), &result)
			require.NoError(t, err)
			assert.Equal(t, tt.expected.Valid, result.Valid)
			assert.Equal(t, tt.expected.String, result.String)
		})
	}
}

func TestNullableString_InStruct(t *testing.T) {
	type TestStruct struct {
		Name   string          `json:"name"`
		Street nullable.String `json:"street"`
		ZipCode nullable.String `json:"zipCode"`
	}

	t.Run("marshal with valid and null values", func(t *testing.T) {
		s := TestStruct{
			Name:    "Test",
			Street:  nullable.NewString("123 Main St"),
			ZipCode: nullable.String{}, // null
		}

		data, err := json.Marshal(s)
		require.NoError(t, err)
		assert.JSONEq(t, `{"name":"Test","street":"123 Main St","zipCode":null}`, string(data))
	})

	t.Run("unmarshal with valid and null values", func(t *testing.T) {
		input := `{"name":"Test","street":"123 Main St","zipCode":null}`
		
		var result TestStruct
		err := json.Unmarshal([]byte(input), &result)
		require.NoError(t, err)
		
		assert.Equal(t, "Test", result.Name)
		assert.True(t, result.Street.Valid)
		assert.Equal(t, "123 Main St", result.Street.String)
		assert.False(t, result.ZipCode.Valid)
	})
}

func TestNullableString_Helpers(t *testing.T) {
	t.Run("Ptr returns pointer for valid value", func(t *testing.T) {
		ns := nullable.NewString("test")
		ptr := ns.Ptr()
		require.NotNil(t, ptr)
		assert.Equal(t, "test", *ptr)
	})

	t.Run("Ptr returns nil for invalid value", func(t *testing.T) {
		ns := nullable.String{}
		ptr := ns.Ptr()
		assert.Nil(t, ptr)
	})

	t.Run("ValueOrDefault returns value for valid", func(t *testing.T) {
		ns := nullable.NewString("test")
		assert.Equal(t, "test", ns.ValueOrDefault("default"))
	})

	t.Run("ValueOrDefault returns default for invalid", func(t *testing.T) {
		ns := nullable.String{}
		assert.Equal(t, "default", ns.ValueOrDefault("default"))
	})

	t.Run("NewStringFromPtr with nil", func(t *testing.T) {
		ns := nullable.NewStringFromPtr(nil)
		assert.False(t, ns.Valid)
	})

	t.Run("NewStringFromPtr with value", func(t *testing.T) {
		val := "test"
		ns := nullable.NewStringFromPtr(&val)
		assert.True(t, ns.Valid)
		assert.Equal(t, "test", ns.String)
	})
}
