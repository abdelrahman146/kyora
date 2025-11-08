package list

type ListResponse[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"totalCount,omitempty"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages,omitempty"`
	HasMore    bool  `json:"hasMore"`
}

func NewListResponse[T any](items []T, page, pageSize int, totalCount int64, hasMore bool) *ListResponse[T] {
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))
	return &ListResponse[T]{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasMore:    hasMore,
	}
}
