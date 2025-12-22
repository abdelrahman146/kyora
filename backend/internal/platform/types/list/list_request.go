package list

import (
	"fmt"
	"reflect"

	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
)

type ListRequest struct {
	page       int
	pageSize   int
	orderBy    []string // e.g., ["name", "-createdAt"]
	searchTerm string
}

func (r *ListRequest) Offset() int {
	if r.page < 1 {
		return 0
	}
	return (r.page - 1) * r.pageSize
}

func (r *ListRequest) Limit() int {
	if r.pageSize <= 0 {
		return 30
	}
	if r.pageSize > 100 {
		return 100
	}
	return r.pageSize
}

func (r *ListRequest) OrderBy() []string {
	if len(r.orderBy) == 0 {
		return []string{"-createdAt"}
	}
	return r.orderBy
}

func (r *ListRequest) SearchTerm() string {
	return r.searchTerm
}

func (r *ListRequest) Page() int {
	return r.page
}

func (r *ListRequest) PageSize() int {
	return r.pageSize
}

func (r *ListRequest) ParsedOrderBy(schemaDef any) []string {
	var result []string
	for _, order := range r.OrderBy() {
		// Determine the sort direction
		direction := "ASC"
		if order[0] == '-' {
			direction = "DESC"
			order = order[1:]
		}

		// Look up the field in the schema
		fields := ParseArrayToSchema([]string{order}, schemaDef)
		if len(fields) == 0 {
			continue
		}

		// Map the field to its sort expression
		for _, field := range fields {
			result = append(result, fmt.Sprintf("%s %s", field.Column(), direction))
		}
	}
	return result
}

func (r *ListRequest) ParsedOrderByWithDefault(schemaDef any, defaultOrder []string) []string {
	if len(r.OrderBy()) == 0 {
		return defaultOrder
	}
	return r.ParsedOrderBy(schemaDef)
}

func NewListRequest(page, pageSize int, orderBy []string, searchTerm string) *ListRequest {
	return &ListRequest{
		page:       page,
		pageSize:   pageSize,
		orderBy:    orderBy,
		searchTerm: searchTerm,
	}
}

func ParseArrayToSchema(arr []string, schemaDef any) (fields []schema.Field) {
	// Create a map for efficient O(1) lookups
	// Maps JSON field name (string) -> schema.Field
	jsonToField := make(map[string]schema.Field)

	// Use reflection to inspect the schemaDef
	val := reflect.ValueOf(schemaDef)

	// Handle cases where a pointer to the struct is passed
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	// We can only inspect structs
	if val.Kind() != reflect.Struct {
		// Return the empty 'fields' slice.
		// You could also return an error here.
		return
	}

	// --- 1. Build the lookup map ---
	// Iterate over all fields in the struct *once*
	numStructFields := val.NumField()
	for i := range numStructFields {
		fieldVal := val.Field(i)

		// Check if the struct field is actually a 'Field' type
		// *** In your code, change 'Field' to 'schema.Field' ***
		if f, ok := fieldVal.Interface().(schema.Field); ok {
			// It is, so add it to our map
			jsonToField[f.JSONField()] = f
		}
	}

	// --- 2. Look up the requested fields ---
	// 'fields' is already declared as the named return value
	for _, jsonName := range arr {
		// Look up the name in our map
		if f, ok := jsonToField[jsonName]; ok {
			// Found it, add to the result slice
			fields = append(fields, f)
		}
		// If 'ok' is false, the field name wasn't in the schema,
		// so we simply ignore it.
	}

	// Return the (potentially empty) slice of found fields
	return
}
