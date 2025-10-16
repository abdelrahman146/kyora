package types

type ListRequest struct {
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
	OrderBy  []string `json:"orderBy"` // e.g., ["name", "-createdAt"]
	Search   string   `json:"search"`
}

type ListResponse[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"totalCount"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
	HasMore    bool  `json:"hasMore"`
}

func NewListResponse[T any](items []T, page, pageSize int, totalCount int64) *ListResponse[T] {
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))
	hasMore := page < totalPages
	return &ListResponse[T]{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasMore:    hasMore,
	}
}
